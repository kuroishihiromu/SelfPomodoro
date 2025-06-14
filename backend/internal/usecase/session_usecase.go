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

// SessionUseCase はセッションに関するユースケースを定義するインターフェース
type SessionUseCase interface {
	// StartSession は新しいセッションを開始する
	StartSession(ctx context.Context, userID uuid.UUID) (*model.SessionResponse, error)

	// GetSession はセッションを取得する
	GetSession(ctx context.Context, id, userID uuid.UUID) (*model.SessionResponse, error)

	// GetAllSessions はユーザの全セッションを取得する
	GetAllSessions(ctx context.Context, userID uuid.UUID) (*model.SessionsResponse, error)

	// CompleteSession はセッションを完了する（SQSメッセージ送信付き）
	CompleteSession(ctx context.Context, id, userID uuid.UUID) (*model.SessionResponse, error)

	// DeleteSession はセッションを削除する
	DeleteSession(ctx context.Context, id, userID uuid.UUID) error
}

// sessionUseCase はSessionUseCaseインターフェースの実装（新エラーハンドリング対応版）
type sessionUseCase struct {
	sessionRepo    repository.SessionRepository
	roundRepo      repository.RoundRepository
	userConfigRepo repository.UserConfigRepository
	sqsClient      *sqs.SQSClient
	logger         logger.Logger
}

// NewSessionUseCase は新しいSessionUseCaseインスタンスを作成する
func NewSessionUseCase(
	sessionRepo repository.SessionRepository,
	roundRepo repository.RoundRepository,
	userConfigRepo repository.UserConfigRepository,
	sqsClient *sqs.SQSClient,
	logger logger.Logger,
) SessionUseCase {
	return &sessionUseCase{
		sessionRepo:    sessionRepo,
		roundRepo:      roundRepo,
		userConfigRepo: userConfigRepo,
		sqsClient:      sqsClient,
		logger:         logger,
	}
}

// StartSession は新しいセッションを開始する（新エラーハンドリング対応版）
func (uc *sessionUseCase) StartSession(ctx context.Context, userID uuid.UUID) (*model.SessionResponse, error) {
	// ✅ ドメインロジック活用：UserConfig安全取得（デフォルト値フォールバック）
	userConfig := uc.getUserConfigWithFallback(ctx, userID)

	// ユーザー設定確認ログ
	uc.logger.Infof("セッション開始 - ユーザー設定確認完了: work=%d分, break=%d分, rounds=%d",
		userConfig.GetWorkTimeOrDefault(), userConfig.GetBreakTimeOrDefault(), userConfig.GetSessionRoundsOrDefault())

	// ✅ ドメインファクトリー使用
	session := model.NewSession(userID)

	// DBにセッションを保存
	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		uc.logger.Errorf("セッション開始エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrUniqueConstraint) {
			return nil, appErrors.NewSessionInProgressError()
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	uc.logger.Infof("セッション開始成功: %s", session.ID.String())
	return session.ToResponse(), nil
}

// GetSession はセッションを取得する（新エラーハンドリング対応版）
func (uc *sessionUseCase) GetSession(ctx context.Context, id, userID uuid.UUID) (*model.SessionResponse, error) {
	session, err := uc.sessionRepo.GetByID(ctx, id, userID)
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

	return session.ToResponse(), nil
}

// GetAllSessions はユーザの全セッションを取得する（新エラーハンドリング対応版）
func (uc *sessionUseCase) GetAllSessions(ctx context.Context, userID uuid.UUID) (*model.SessionsResponse, error) {
	sessions, err := uc.sessionRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		uc.logger.Errorf("セッション一覧取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// セッションをレスポンス用に変換
	sessionResponses := make([]*model.SessionResponse, len(sessions))
	for i, session := range sessions {
		sessionResponses[i] = session.ToResponse()
	}

	return &model.SessionsResponse{Sessions: sessionResponses}, nil
}

// CompleteSession はセッションを完了する（新エラーハンドリング対応版）
func (uc *sessionUseCase) CompleteSession(ctx context.Context, id, userID uuid.UUID) (*model.SessionResponse, error) {
	// セッション取得
	session, err := uc.sessionRepo.GetByID(ctx, id, userID)
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

	// ✅ ドメインロジック活用：状態チェック
	if session.IsCompleted() {
		uc.logger.Errorf("セッションは既に完了しています: %s", id.String())
		return nil, appErrors.NewSessionAlreadyEndedError()
	}

	// ✅ ドメインロジック活用：ラウンド取得と統計計算
	rounds, err := uc.roundRepo.GetAllBySessionID(ctx, id)
	if err != nil {
		uc.logger.Errorf("セッションラウンド取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// ✅ ドメインロジック活用：セッション統計計算（Repository依存なし！）
	stats := session.CalculateStatistics(rounds)
	uc.logger.Infof("セッション統計計算完了: 平均集中度=%.1f, 総作業時間=%d分, ラウンド数=%d, 休憩時間=%d分",
		stats.AverageFocus, stats.TotalWorkMin, stats.RoundCount, stats.TotalBreakTime)

	// ✅ ドメインロジック活用：セッション完了処理
	session.CompleteWithStatistics(stats)

	// データベース更新
	err = uc.sessionRepo.Complete(ctx, id, userID, stats.AverageFocus, stats.TotalWorkMin, stats.RoundCount, stats.TotalBreakTime)
	if err != nil {
		uc.logger.Errorf("セッション完了永続化エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return nil, appErrors.NewSessionNotFoundError()
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// ✅ ドメインロジック活用：最適化メッセージ送信判定
	if session.ShouldSendOptimizationMessage() {
		avgFocus, totalWork, _ := session.GetOptimizationMessageData()
		uc.logger.Infof("セッション完了 - SQS最適化メッセージ送信開始: SessionID=%s, RoundCount=%d, AvgFocus=%.1f, TotalWork=%dmin",
			id.String(), stats.RoundCount, avgFocus, totalWork)

		uc.sendSessionOptimizationMessage(ctx, userID, id, avgFocus, totalWork)
	} else {
		uc.logger.Warn("完了したラウンドが存在しないため、SQS最適化メッセージは送信しません")
	}

	// 完了したセッションを取得
	updatedSession, err := uc.sessionRepo.GetByID(ctx, id, userID)
	if err != nil {
		uc.logger.Errorf("完了後のセッション取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return nil, appErrors.NewSessionNotFoundError()
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// ✅ ドメインロジック活用：セッション品質評価ログ
	quality := updatedSession.GetSessionQuality()
	efficiency := updatedSession.GetEfficiency()
	uc.logger.Infof("セッション完了成功: SessionID=%s, 品質=%s, 効率=%.1f%%, 平均集中度=%.1f, 総作業時間=%d分, ラウンド数=%d",
		id.String(), quality, efficiency, stats.AverageFocus, stats.TotalWorkMin, stats.RoundCount)

	return updatedSession.ToResponse(), nil
}

// DeleteSession はセッションを削除する（新エラーハンドリング対応版）
func (uc *sessionUseCase) DeleteSession(ctx context.Context, id, userID uuid.UUID) error {
	if err := uc.sessionRepo.Delete(ctx, id, userID); err != nil {
		uc.logger.Errorf("セッション削除エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return appErrors.NewSessionNotFoundError()
		}
		if appErrors.IsDatabaseError(err) {
			return appErrors.NewInternalError(err)
		}

		return appErrors.NewInternalError(err)
	}
	uc.logger.Infof("セッション削除成功: %s", id.String())
	return nil
}

// ✅ ドメインロジック活用：UserConfig安全取得（フォールバック）
func (uc *sessionUseCase) getUserConfigWithFallback(ctx context.Context, userID uuid.UUID) *model.UserConfig {
	if uc.userConfigRepo == nil {
		uc.logger.Warn("UserConfigRepository が nil です - デフォルト設定使用")
		return model.NewDefaultUserConfig(userID)
	}

	userConfig, err := uc.userConfigRepo.GetUserConfig(ctx, userID)
	if err != nil {
		uc.logger.Warnf("UserConfig取得エラー、デフォルト設定を使用: %v", err)
		return model.NewDefaultUserConfig(userID)
	}

	return userConfig
}

// sendSessionOptimizationMessage はセッション最適化メッセージをSQSに送信する（同期）
func (uc *sessionUseCase) sendSessionOptimizationMessage(ctx context.Context, userID, sessionID uuid.UUID, avgFocusScore float64, totalWorkTime int) {
	if uc.sqsClient == nil {
		uc.logger.Warn("SQSクライアントが初期化されていません。最適化メッセージは送信されません。")
		return
	}

	// 計算結果を使用してメッセージを作成
	message := model.NewSessionOptimizationMessage(userID, sessionID, avgFocusScore, totalWorkTime)

	// SQS送信実行
	err := uc.sqsClient.SendSessionOptimizationMessage(ctx, message)
	if err != nil {
		uc.logger.Errorf("セッション最適化メッセージ送信エラー: %v", err)
		return
	}

	uc.logger.Infof("セッション最適化メッセージ送信成功: %s", message.ToLogString())
}
