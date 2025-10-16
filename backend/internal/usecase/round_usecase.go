package usecase

import (
	"context"
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/domain/service"
	roundVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/round"
	sessionVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/session"
	statisticsVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/statistics"
	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/usecase/dto"
	"github.com/tsunakit99/selfpomodoro/internal/usecase/mapper"
)

// RoundUseCase はラウンドに関するユースケースを定義するインターフェース
type RoundUseCase interface {
	// StartRound は新しいラウンドを開始する
	StartRound(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, req *dto.RoundCreateRequest) (*dto.RoundResponse, error)

	// CompleteRound はラウンドを完了する(SQSメッセージ送信付き)
	CompleteRound(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *dto.RoundCompleteRequest) (*dto.RoundResponse, error)
}

// roundUseCase はラウンドに関するユースケースの実装（新エラーハンドリング完全対応版）
type roundUseCase struct {
	sessionRepo                 repository.SessionRepository
	optimizationPreferencesRepo repository.OptimizationPreferencesRepository
	messagingService            service.MessagingService
	statisticsRepo              repository.StatisticsRepository
	roundMapper                 *mapper.RoundMapper
	logger                      logger.Logger
}

// NewRoundUseCase は新しいラウンドユースケースを作成する
func NewRoundUseCase(sessionRepo repository.SessionRepository, optimizationPreferencesRepo repository.OptimizationPreferencesRepository, messagingService service.MessagingService, statisticsRepo repository.StatisticsRepository, logger logger.Logger) RoundUseCase {
	return &roundUseCase{
		sessionRepo:                 sessionRepo,
		optimizationPreferencesRepo: optimizationPreferencesRepo,
		messagingService:            messagingService,
		statisticsRepo:              statisticsRepo,
		roundMapper:                 mapper.NewRoundMapper(),
		logger:                      logger,
	}
}

// StartRound は新しいラウンドを開始する（元の設計に戻す：開始時DB作成）
func (uc *roundUseCase) StartRound(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, req *dto.RoundCreateRequest) (*dto.RoundResponse, error) {
	// UUIDをValue Objectに変換
	sessionIDVO, err := sessionVO.NewSessionIDFromString(sessionID.String())
	if err != nil {
		return nil, appErrors.NewInternalError(err)
	}
	userIDVO, err := userVO.NewUserIDFromString(userID.String())
	if err != nil {
		return nil, appErrors.NewInternalError(err)
	}

	// セッションが存在するか確認
	session, err := uc.sessionRepo.GetSession(ctx, sessionIDVO, userIDVO)
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
	rounds, err := uc.sessionRepo.GetRoundsBySession(ctx, session.ID, userIDVO)
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

	// OptimizationPreferencesから直接session_roundsを取得
	optimizationPreferences := uc.getOptimizationPreferencesWithFallback(ctx, userIDVO)
	maxRounds := optimizationPreferences.GetSessionRoundsOrDefault()

	// ✅ DDD Aggregate: Sessionドメインロジックでラウンド追加
	round, err := session.AddRound(maxRounds.Count())
	if err != nil {
		uc.logger.Errorf("ラウンド追加エラー: %v", err)
		return nil, appErrors.NewBadRequestError(err.Error())
	}

	// データベースにラウンドを保存（開始時作成）
	if err = uc.sessionRepo.AddRoundToSession(ctx, session.ID, userIDVO, round); err != nil {
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
		round.ID.String(), sessionID.String(), round.RoundOrder.Order())

	return uc.roundMapper.ToRoundResponse(round), nil
}

// CompleteRound はラウンドを完了する（元の設計に戻す：既存ラウンドの更新）
func (uc *roundUseCase) CompleteRound(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *dto.RoundCompleteRequest) (*dto.RoundResponse, error) {
	// 入力値検証
	if id == uuid.Nil {
		uc.logger.Error("CompleteRound: 無効なRoundIDが指定されました")
		return nil, appErrors.NewValidationError("ラウンドIDが無効です")
	}
	if userID == uuid.Nil {
		uc.logger.Error("CompleteRound: 無効なUserIDが指定されました")
		return nil, appErrors.NewValidationError("ユーザーIDが無効です")
	}
	if req == nil {
		uc.logger.Error("CompleteRound: リクエストが無効です")
		return nil, appErrors.NewValidationError("リクエストが無効です")
	}

	// UUIDをValue Objectに変換
	roundIDVO, err := roundVO.NewRoundIDFromString(id.String())
	if err != nil {
		uc.logger.Errorf("CompleteRound: RoundID Value Object変換エラー: %v", err)
		return nil, appErrors.NewInternalError(err)
	}
	userIDVO, err := userVO.NewUserIDFromString(userID.String())
	if err != nil {
		uc.logger.Errorf("CompleteRound: UserID Value Object変換エラー: %v", err)
		return nil, appErrors.NewInternalError(err)
	}

	// SessionRepository経由でラウンドを取得
	round, err := uc.sessionRepo.GetRoundByID(ctx, roundIDVO, userIDVO)
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

	// ✅ ドメインロジック活用：OptimizationPreferences統合（デフォルト値フォールバック）
	optimizationPreferences := uc.getOptimizationPreferencesWithFallback(ctx, userIDVO)
	workTime := optimizationPreferences.GetWorkTimeOrDefault()
	breakTime := optimizationPreferences.GetBreakTimeOrDefault()

	// ✅ ドメインロジック活用：ラウンド完了処理
	if err := round.CompleteWith(req.FocusScore, workTime.Minutes(), breakTime.Minutes()); err != nil {
		uc.logger.Errorf("ラウンド完了ドメインエラー: %v", err)
		return nil, appErrors.NewBadRequestError(err.Error())
	}

	// SessionRepository経由でラウンドを完了
	err = uc.sessionRepo.CompleteRound(ctx, round.SessionID, userIDVO, roundIDVO, req.FocusScore, workTime.Minutes(), breakTime.Minutes())
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

		// SessionIDをUUIDに変換
		sessionUUID, err := uuid.Parse(round.SessionID.String())
		if err != nil {
			uc.logger.Errorf("SessionID変換エラー: %v", err)
		} else {
			uc.sendRoundOptimizationMessage(ctx, userID, id, focusScore, sessionUUID, workTime.Minutes())
		}
	} else {
		uc.logger.Info("集中度スコア未設定または最小値未満のため、SQS最適化メッセージは送信しません")
	}

	// 統計データの更新（StatisticsRepository経由）
	if uc.statisticsRepo != nil {
		// EndTimeが*time.Timeの場合の処理
		endTime := round.EndTime
		if endTime == nil {
			endTime = &round.StartTime // フォールバック
		}
		if err := uc.statisticsRepo.UpdateStatisticsWithRound(ctx, userIDVO, round, *endTime); err != nil {
			// 統計更新エラーはログ出力のみで、ラウンド完了処理は継続
			uc.logger.Warnf("統計データ更新エラー（処理続行）: %v", err)
		}
	}

	focusScoreLog := "未設定"
	if req.FocusScore != nil {
		focusScoreLog = strconv.Itoa(*req.FocusScore)
	}

	uc.logger.Infof("ラウンド完了成功: ID=%s, FocusScore=%s, WorkTime=%d分, BreakTime=%d分",
		id.String(), focusScoreLog, workTime.Minutes(), breakTime.Minutes())

	return uc.roundMapper.ToRoundResponse(round), nil
}

// ✅ ドメインロジック活用：OptimizationPreferences安全取得（フォールバック）
func (uc *roundUseCase) getOptimizationPreferencesWithFallback(ctx context.Context, userID userVO.UserID) *entity.OptimizationPreferences {
	if uc.optimizationPreferencesRepo == nil {
		uc.logger.Warn("OptimizationPreferencesRepository が nil です - デフォルト設定使用")
		defaultPrefs, _ := entity.NewOptimizationPreferences(userID)
		return defaultPrefs
	}

	optimizationPreferences, err := uc.optimizationPreferencesRepo.GetOrCreateDefault(ctx, userID)
	if err != nil {
		uc.logger.Warnf("OptimizationPreferences取得エラー、デフォルト設定を使用: %v", err)
		defaultPrefs, _ := entity.NewOptimizationPreferences(userID)
		return defaultPrefs
	}

	uc.logger.Infof("OptimizationPreferences取得成功: workTime=%d, breakTime=%d, sessionRounds=%d",
		optimizationPreferences.GetWorkTimeOrDefault().Minutes(), optimizationPreferences.GetBreakTimeOrDefault().Minutes(), optimizationPreferences.GetSessionRoundsOrDefault().Count())
	return optimizationPreferences
}

// sendRoundOptimizationMessage はラウンド最適化メッセージを送信する（MessagingService統合版）
func (uc *roundUseCase) sendRoundOptimizationMessage(ctx context.Context, userID, roundID uuid.UUID, focusScore int, sessionID uuid.UUID, workTime int) {
	if uc.messagingService == nil {
		uc.logger.Warn("MessagingServiceが初期化されていません。最適化メッセージは送信されません。")
		return
	}

	// UUIDをValue Objectに変換
	userIDVO, err := userVO.NewUserIDFromString(userID.String())
	if err != nil {
		uc.logger.Errorf("UserID変換エラー: %v", err)
		return
	}

	sessionIDVO, err := sessionVO.NewSessionIDFromString(sessionID.String())
	if err != nil {
		uc.logger.Errorf("SessionID変換エラー: %v", err)
		return
	}

	roundIDVO, err := roundVO.NewRoundIDFromString(roundID.String())
	if err != nil {
		uc.logger.Errorf("RoundID変換エラー: %v", err)
		return
	}

	// Value Objectsに変換
	focusScoreVO, err := statisticsVO.NewFocusScore(float64(focusScore))
	if err != nil {
		uc.logger.Errorf("FocusScore変換エラー: %v", err)
		return
	}

	workTimeVO, err := statisticsVO.NewWorkTime(workTime)
	if err != nil {
		uc.logger.Errorf("WorkTime変換エラー: %v", err)
		return
	}

	// ドメインサービス経由でメッセージ送信
	roundData := service.RoundOptimizationData{
		UserID:     userIDVO,
		SessionID:  sessionIDVO,
		RoundID:    roundIDVO,
		FocusScore: focusScoreVO,
		WorkTime:   workTimeVO,
	}

	if err := uc.messagingService.SendRoundOptimizationMessage(ctx, roundData); err != nil {
		uc.logger.Errorf("ラウンド最適化メッセージ送信エラー: %v", err)
	}
}
