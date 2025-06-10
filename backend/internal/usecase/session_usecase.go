package usecase

import (
	"context"

	"errors"

	"github.com/google/uuid"
	domainErrors "github.com/tsunakit99/selfpomodoro/internal/domain/errors"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	sessionerrors "github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/postgres/errors"
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

// sessionUseCase はSessionUseCaseインターフェースの実装
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

// StartSession は新しいセッションを開始する（UserConfig デフォルト値フォールバック付き）
func (uc *sessionUseCase) StartSession(ctx context.Context, userID uuid.UUID) (*model.SessionResponse, error) {
	// UserConfig取得（存在前提・フォールバック付き）
	var userConfig *model.UserConfig
	if uc.userConfigRepo != nil {
		config, err := uc.userConfigRepo.GetUserConfig(ctx, userID)
		if err != nil {
			uc.logger.Warnf("UserConfig取得失敗、デフォルト値でセッション開始: %v", err)
			// フォールバック：デフォルト値のUserConfigを作成
			userConfig = model.NewUserConfig(userID)
		} else {
			userConfig = config
		}
	} else {
		uc.logger.Warn("UserConfigRepository が初期化されていません - デフォルト値使用")
		// フォールバック：デフォルト値のUserConfigを作成
		userConfig = model.NewUserConfig(userID)
	}

	// ユーザー設定確認ログ（フォールバック対応）
	uc.logger.Infof("セッション開始 - ユーザー設定確認完了: work=%d分, break=%d分, rounds=%d",
		userConfig.RoundWorkTime, userConfig.RoundBreakTime, userConfig.SessionRounds)

	// セッションを作成
	session := model.NewSession(userID)

	// DBにセッションを保存
	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		uc.logger.Errorf("セッション開始エラー: %v", err)
		if errors.Is(err, sessionerrors.ErrSessionInProgress) {
			return nil, domainErrors.NewSessionInProgressError()
		}
		return nil, domainErrors.NewInternalServerError()
	}

	uc.logger.Infof("セッション開始成功: %s", session.ID.String())
	return session.ToResponse(), nil
}

// GetSession はセッションを取得する
func (uc *sessionUseCase) GetSession(ctx context.Context, id, userID uuid.UUID) (*model.SessionResponse, error) {
	session, err := uc.sessionRepo.GetByID(ctx, id, userID)
	if err != nil {
		uc.logger.Errorf("セッション取得エラー: %v", err)
		if errors.Is(err, sessionerrors.ErrSessionNotFound) {
			return nil, domainErrors.NewSessionNotFoundError()
		}
		return nil, domainErrors.NewInternalServerError()
	}

	// セッションをレスポンス用に変換
	return session.ToResponse(), nil
}

// GetAllSessions はユーザの全セッションを取得する
func (uc *sessionUseCase) GetAllSessions(ctx context.Context, userID uuid.UUID) (*model.SessionsResponse, error) {
	sessions, err := uc.sessionRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		uc.logger.Errorf("セッション一覧取得エラー: %v", err)
		return nil, domainErrors.NewInternalServerError()
	}

	// セッションをレスポンス用に変換
	sessionResponses := make([]*model.SessionResponse, len(sessions))
	for i, session := range sessions {
		sessionResponses[i] = session.ToResponse()
	}

	return &model.SessionsResponse{Sessions: sessionResponses}, nil
}

// CompleteSession はセッションを完了する（SQSメッセージ送信付き）
func (uc *sessionUseCase) CompleteSession(ctx context.Context, id, userID uuid.UUID) (*model.SessionResponse, error) {
	_, err := uc.sessionRepo.GetByID(ctx, id, userID)
	if err != nil {
		uc.logger.Errorf("セッション取得エラー: %v", err)
		if errors.Is(err, sessionerrors.ErrSessionNotFound) {
			return nil, domainErrors.NewSessionNotFoundError()
		}
		return nil, domainErrors.NewInternalServerError()
	}

	// ラウンドの統計情報を計算
	averageFocus, totalWorkMin, roundCount, breakTime, err := uc.roundRepo.CalculateSessionStats(ctx, id)
	if err != nil {
		uc.logger.Errorf("セッション統計情報取得エラー: %v", err)
		return nil, domainErrors.NewInternalServerError()
	}

	// セッションを完了する
	err = uc.sessionRepo.Complete(ctx, id, userID, averageFocus, totalWorkMin, roundCount, breakTime)
	if err != nil {
		uc.logger.Errorf("セッション完了エラー: %v", err)
		if errors.Is(err, sessionerrors.ErrSessionAlreadyEnded) {
			return nil, domainErrors.NewSessionAlreadyEndedError()
		}
		return nil, domainErrors.NewInternalServerError()
	}

	// ✅ 完了したラウンドが存在する場合のみSQS送信
	if roundCount > 0 {
		uc.logger.Infof("セッション完了 - SQS最適化メッセージ送信開始: SessionID=%s, RoundCount=%d, AvgFocus=%.1f, TotalWork=%dmin",
			id.String(), roundCount, averageFocus, totalWorkMin)

		// 同期でSQS送信
		uc.sendSessionOptimizationMessage(ctx, userID, id, averageFocus, totalWorkMin)
	} else {
		uc.logger.Warn("完了したラウンドが存在しないため、SQS最適化メッセージは送信しません")
	}

	// 完了したセッションを取得
	updatedSession, err := uc.sessionRepo.GetByID(ctx, id, userID)
	if err != nil {
		uc.logger.Errorf("完了後のセッション取得エラー: %v", err)
		return nil, domainErrors.NewInternalServerError()
	}

	uc.logger.Infof("セッション完了成功: %s (平均集中度: %.1f, 総作業時間: %d分, ラウンド数: %d)",
		id.String(), averageFocus, totalWorkMin, roundCount)

	// TODO: セッション最適化Lambda実行
	// TODO: 最適化ログの記録

	// セッションをレスポンス用に変換
	return updatedSession.ToResponse(), nil
}

// DeleteSession はセッションを削除する
func (uc *sessionUseCase) DeleteSession(ctx context.Context, id, userID uuid.UUID) error {
	if err := uc.sessionRepo.Delete(ctx, id, userID); err != nil {
		uc.logger.Errorf("セッション削除エラー: %v", err)
		if errors.Is(err, sessionerrors.ErrSessionNotFound) {
			return domainErrors.NewSessionNotFoundError()
		}
		return domainErrors.NewInternalServerError()
	}
	uc.logger.Infof("セッション削除成功: %s", id.String())
	return nil
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
