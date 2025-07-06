package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/sqs"
)

// RoundUseCase はラウンドに関するユースケースを定義するインターフェース
type RoundUseCase interface {
	// StartRound は新しいラウンドを開始する
	StartRound(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, req *entity.RoundCreateRequest) (*entity.RoundResponse, error)

	// CompleteRound はラウンドを完了する(SQSメッセージ送信付き)
	CompleteRound(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *entity.RoundCompleteRequest) (*entity.RoundResponse, error)
}

// roundUseCase はラウンドに関するユースケースの実装（新エラーハンドリング完全対応版）
type roundUseCase struct {
	roundRepo       repository.RoundRepository
	sessionRepo     repository.SessionRepository
	userConfigRepo  repository.UserConfigRepository
	sqsClient       *sqs.SQSClient
	statsAggregator StatisticsAggregationService
	logger          logger.Logger
}

// NewRoundUseCase は新しいラウンドユースケースを作成する
func NewRoundUseCase(roundrepo repository.RoundRepository, sessionRepo repository.SessionRepository, userConfigRepo repository.UserConfigRepository, sqsClient *sqs.SQSClient, statsAggregator StatisticsAggregationService, logger logger.Logger) RoundUseCase {
	return &roundUseCase{
		roundRepo:       roundrepo,
		sessionRepo:     sessionRepo,
		userConfigRepo:  userConfigRepo,
		sqsClient:       sqsClient,
		statsAggregator: statsAggregator,
		logger:          logger,
	}
}

// StartRound は新しいラウンドを開始する（元の設計に戻す：開始時DB作成）
func (uc *roundUseCase) StartRound(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, req *entity.RoundCreateRequest) (*entity.RoundResponse, error) {
	// セッションが存在するか確認
	session, err := uc.sessionRepo.GetByID(ctx, sessionID, userID)
	if err != nil {
		uc.logger.Errorf("セッション取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return nil, appErrors.NewSessionNotFoundError()
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// ✅ DDD Aggregate: 既存ラウンドをSessionに読み込み
	rounds, err := uc.roundRepo.GetBySessionIDWithUserID(ctx, session.ID, userID)
	if err != nil && !errors.Is(err, appErrors.ErrRecordNotFound) {
		uc.logger.Errorf("ラウンド一覧取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// ✅ DDD Aggregate: SessionにRoundsを読み込み
	session.LoadRounds(rounds)

	// UserConfigから直接session_roundsを取得
	userConfig := uc.getUserConfigWithFallback(ctx, userID)
	maxRounds := userConfig.GetSessionRoundsOrDefault()

	// ✅ DDD Aggregate: Sessionドメインロジックでラウンド追加
	round, err := session.AddRound(maxRounds)
	if err != nil {
		uc.logger.Errorf("ラウンド追加エラー: %v", err)
		return nil, appErrors.NewBadRequestError(err.Error())
	}

	// データベースにラウンドを保存（開始時作成）
	if err = uc.roundRepo.Create(ctx, round, userID); err != nil {
		uc.logger.Errorf("ラウンド作成エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrUniqueConstraint) {
			return nil, appErrors.NewBadRequestError("同じラウンドが既に存在します")
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	uc.logger.Infof("ラウンド開始成功: ID=%s, SessionID=%s, Order=%d",
		round.ID.String(), sessionID.String(), round.RoundOrder)

	return round.ToResponse(), nil
}

// CompleteRound はラウンドを完了する（元の設計に戻す：既存ラウンドの更新）
func (uc *roundUseCase) CompleteRound(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *entity.RoundCompleteRequest) (*entity.RoundResponse, error) {
	// ラウンドの存在確認
	round, err := uc.roundRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Errorf("ラウンド取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return nil, appErrors.NewRoundNotFoundError()
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// ✅ ドメインロジック活用：完了可能性チェック
	if err := round.CanBeCompleted(); err != nil {
		uc.logger.Errorf("ラウンド完了不可: %v", err)
		return nil, appErrors.NewRoundAlreadyEndedError()
	}

	// ✅ ドメインロジック活用：UserConfig統合（デフォルト値フォールバック）
	userConfig := uc.getUserConfigWithFallback(ctx, userID)
	workTime := userConfig.GetWorkTimeOrDefault()
	breakTime := userConfig.GetBreakTimeOrDefault()

	// ✅ ドメインロジック活用：ラウンド完了処理
	if err := round.CompleteWith(req.FocusScore, workTime, breakTime); err != nil {
		uc.logger.Errorf("ラウンド完了ドメインエラー: %v", err)
		return nil, appErrors.NewBadRequestError(err.Error())
	}

	// データベース更新（ドメインオブジェクトの状態をそのまま永続化）
	err = uc.roundRepo.Complete(ctx, id, req.FocusScore, workTime, breakTime)
	if err != nil {
		uc.logger.Errorf("ラウンド完了永続化エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return nil, appErrors.NewRoundAlreadyEndedError()
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// ✅ ドメインロジック活用：最適化メッセージ送信判定
	if round.ShouldSendOptimizationMessage() {
		focusScore, _ := round.GetOptimizationMessageData()
		uc.logger.Infof("ラウンド完了 - SQS最適化メッセージ送信開始: RoundID=%s, FocusScore=%d",
			id.String(), focusScore)

		uc.sendRoundOptimizationMessage(ctx, userID, id, focusScore)
	} else {
		uc.logger.Info("集中度スコア未設定または最小値未満のため、SQS最適化メッセージは送信しません")
	}

	// 完了したラウンドを取得して返す
	completedRound, err := uc.roundRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Errorf("完了ラウンド取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return nil, appErrors.NewRoundNotFoundError()
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// 統計データの事前集約処理
	if uc.statsAggregator != nil {
		if err := uc.statsAggregator.UpdateStatisticsOnRoundComplete(ctx, userID, completedRound); err != nil {
			// 統計更新エラーはログ出力のみで、ラウンド完了処理は継続
			uc.logger.Warnf("統計データ更新エラー（処理続行）: %v", err)
		}
	}

	uc.logger.Infof("ラウンド完了成功: ID=%s, FocusScore=%v, WorkTime=%d分, BreakTime=%d分",
		id.String(), req.FocusScore, workTime, breakTime)

	return completedRound.ToResponse(), nil
}

// ✅ ドメインロジック活用：UserConfig安全取得（フォールバック）
func (uc *roundUseCase) getUserConfigWithFallback(ctx context.Context, userID uuid.UUID) *entity.UserConfig {
	if uc.userConfigRepo == nil {
		uc.logger.Warn("UserConfigRepository が nil です - デフォルト設定使用")
		return entity.NewDefaultUserConfig(userID)
	}

	userConfig, err := uc.userConfigRepo.GetUserConfig(ctx, userID)
	if err != nil {
		uc.logger.Warnf("UserConfig取得エラー、デフォルト設定を使用: %v", err)
		return entity.NewDefaultUserConfig(userID)
	}

	uc.logger.Infof("UserConfig取得成功: workTime=%d, breakTime=%d, sessionRounds=%d",
		userConfig.GetWorkTimeOrDefault(), userConfig.GetBreakTimeOrDefault(), userConfig.GetSessionRoundsOrDefault())
	return userConfig
}

// sendRoundOptimizationMessage はラウンド最適化メッセージをSQSに送信する（同期）
func (uc *roundUseCase) sendRoundOptimizationMessage(ctx context.Context, userID, roundID uuid.UUID, focusScore int) {
	if uc.sqsClient == nil {
		uc.logger.Warn("SQSクライアントが初期化されていません。最適化メッセージは送信されません。")
		return
	}

	// 最小限のメッセージを作成
	message := model.NewRoundOptimizationMessage(userID, roundID, focusScore)

	// SQS送信実行
	err := uc.sqsClient.SendRoundOptimizationMessage(ctx, message)
	if err != nil {
		uc.logger.Errorf("ラウンド最適化メッセージ送信エラー: %v", err)
		return
	}

	uc.logger.Infof("ラウンド最適化メッセージ送信成功: %s", message.ToLogString())
}
