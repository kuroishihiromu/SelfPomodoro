package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/sqs"
)

// RoundUseCase はラウンドに関するユースケースを定義するインターフェース
type RoundUseCase interface {
	// StartRound は新しいラウンドを開始する
	StartRound(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, req *model.RoundCreateRequest) (*model.RoundResponse, error)

	// GetRound は指定されたIDのラウンドを取得する
	GetRound(ctx context.Context, id uuid.UUID) (*model.RoundResponse, error)

	// GetAllRoundsBySessionID は指定されたセッションIDのラウンドを取得する
	GetAllRoundsBySessionID(ctx context.Context, sessionID uuid.UUID) (*model.RoundsResponse, error)

	// CompleteRound はラウンドを完了する(SQSメッセージ送信付き)
	CompleteRound(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *model.RoundCompleteRequest) (*model.RoundResponse, error)

	// AbortRound はラウンドを中止する（SQSメッセージは送信しない）
	AbortRound(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.RoundResponse, error)
}

// roundUseCase はラウンドに関するユースケースの実装（新エラーハンドリング完全対応版）
type roundUseCase struct {
	roundRepo      repository.RoundRepository
	sessionRepo    repository.SessionRepository
	userConfigRepo repository.UserConfigRepository
	sqsClient      *sqs.SQSClient
	logger         logger.Logger
}

// NewRoundUseCase は新しいラウンドユースケースを作成する
func NewRoundUseCase(roundrepo repository.RoundRepository, sessionRepo repository.SessionRepository, userConfigRepo repository.UserConfigRepository, sqsClient *sqs.SQSClient, logger logger.Logger) RoundUseCase {
	return &roundUseCase{
		roundRepo:      roundrepo,
		sessionRepo:    sessionRepo,
		userConfigRepo: userConfigRepo,
		sqsClient:      sqsClient,
		logger:         logger,
	}
}

// StartRound は新しいラウンドを開始する（新エラーハンドリング完全対応版）
func (uc *roundUseCase) StartRound(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, req *model.RoundCreateRequest) (*model.RoundResponse, error) {
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

	// 最後のラウンドを取得して、完了しているか確認
	lastRound, err := uc.roundRepo.GetLastRoundBySessionID(ctx, session.ID)
	if err != nil && !errors.Is(err, appErrors.ErrRecordNotFound) {
		uc.logger.Errorf("最後のラウンド取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// ✅ ドメインロジック活用：進行中ラウンドチェック
	if lastRound != nil && lastRound.IsInProgress() {
		uc.logger.Errorf("進行中のラウンドがあります: %v", lastRound.ID)
		return nil, appErrors.NewRoundInProgressError()
	}

	// 新しいラウンド順序を計算
	var newRoundOrder int
	if lastRound == nil {
		newRoundOrder = 1 // 最初のラウンド
	} else {
		newRoundOrder = lastRound.RoundOrder + 1 // 次のラウンド
	}

	// ✅ ドメインファクトリー使用
	round := model.NewRound(session.ID, newRoundOrder)

	// ラウンドをデータベースに保存
	if err = uc.roundRepo.Create(ctx, round); err != nil {
		uc.logger.Errorf("ラウンド作成エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrUniqueConstraint) {
			return nil, appErrors.NewRoundInProgressError() // 制約違反は進行中ラウンドとして扱う
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	uc.logger.Infof("ラウンド開始成功: ID=%s, SessionID=%s, Order=%d",
		round.ID.String(), sessionID.String(), newRoundOrder)

	return round.ToResponse(), nil
}

// GetRound は指定されたIDのラウンドを取得する（新エラーハンドリング完全対応版）
func (uc *roundUseCase) GetRound(ctx context.Context, id uuid.UUID) (*model.RoundResponse, error) {
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
	return round.ToResponse(), nil
}

// GetAllRoundsBySessionID はセッションIDに紐づくすべてのラウンドを取得する（新エラーハンドリング完全対応版）
func (uc *roundUseCase) GetAllRoundsBySessionID(ctx context.Context, sessionID uuid.UUID) (*model.RoundsResponse, error) {
	rounds, err := uc.roundRepo.GetAllBySessionID(ctx, sessionID)
	if err != nil {
		uc.logger.Errorf("ラウンド一覧取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// ラウンドをAPIレスポンス形式に変換
	roundResponses := make([]*model.RoundResponse, len(rounds))
	for i, round := range rounds {
		roundResponses[i] = round.ToResponse()
	}

	return &model.RoundsResponse{
		Rounds: roundResponses,
	}, nil
}

// CompleteRound はラウンドを完了する（新エラーハンドリング完全対応版）
func (uc *roundUseCase) CompleteRound(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *model.RoundCompleteRequest) (*model.RoundResponse, error) {
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

	uc.logger.Infof("ラウンド完了成功: ID=%s, FocusScore=%v, WorkTime=%d分, BreakTime=%d分",
		id.String(), req.FocusScore, workTime, breakTime)

	return completedRound.ToResponse(), nil
}

// AbortRound はラウンドを中止する（新エラーハンドリング完全対応版）
func (uc *roundUseCase) AbortRound(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.RoundResponse, error) {
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

	// セッションを取得してユーザーIDでアクセス権限を確認
	_, err = uc.sessionRepo.GetByID(ctx, round.SessionID, userID)
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

	// ✅ ドメインロジック活用：中止可能性チェック
	if err := round.CanBeAborted(); err != nil {
		uc.logger.Errorf("ラウンド中止不可: %v", err)
		return nil, appErrors.NewRoundAlreadyEndedError()
	}

	// ✅ ドメインロジック活用：ラウンド中止処理
	if err := round.Abort(); err != nil {
		uc.logger.Errorf("ラウンド中止ドメインエラー: %v", err)
		return nil, appErrors.NewInternalError(err)
	}

	// データベース更新
	if err := uc.roundRepo.AbortRound(ctx, id); err != nil {
		uc.logger.Errorf("ラウンド中止永続化エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return nil, appErrors.NewRoundAlreadyEndedError()
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// 中止時はSQS送信なし（ビジネスルール）
	uc.logger.Infof("ラウンド中止完了: %s (SQS送信なし)", id.String())

	// 中止したラウンドを取得して返す
	abortedRound, err := uc.roundRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Errorf("中止ラウンド取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return nil, appErrors.NewRoundNotFoundError()
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	return abortedRound.ToResponse(), nil
}

// ✅ ドメインロジック活用：UserConfig安全取得（フォールバック）
func (uc *roundUseCase) getUserConfigWithFallback(ctx context.Context, userID uuid.UUID) *model.UserConfig {
	if uc.userConfigRepo == nil {
		uc.logger.Warn("UserConfigRepository が nil です - デフォルト設定使用")
		return model.NewDefaultUserConfig(userID)
	}

	userConfig, err := uc.userConfigRepo.GetUserConfig(ctx, userID)
	if err != nil {
		uc.logger.Warnf("UserConfig取得エラー、デフォルト設定を使用: %v", err)
		return model.NewDefaultUserConfig(userID)
	}

	uc.logger.Infof("UserConfig取得成功: workTime=%d, breakTime=%d",
		userConfig.GetWorkTimeOrDefault(), userConfig.GetBreakTimeOrDefault())
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
