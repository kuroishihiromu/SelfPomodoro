package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
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

// StartSession は新しいセッションを開始する（UserConfig統合版）
func (uc *sessionUseCase) StartSession(ctx context.Context, userID uuid.UUID) (*model.SessionResponse, error) {
	// セッション開始前にユーザー設定を取得・初期化
	userConfig, err := uc.ensureUserConfig(ctx, userID)
	if err != nil {
		uc.logger.Errorf("ユーザー設定確認エラー: %v", err)
		// 設定取得に失敗してもセッションは開始する（フォールバック）
		uc.logger.Warn("ユーザー設定取得に失敗しましたが、セッションを開始します")
	} else {
		uc.logger.Infof("セッション開始 - ユーザー設定確認完了: work=%d分, break=%d分, rounds=%d",
			userConfig.RoundWorkTime, userConfig.RoundBreakTime, userConfig.SessionRounds)
	}

	// セッションを作成
	session := model.NewSession(userID)

	// DBにセッションを保存
	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		uc.logger.Errorf("セッション開始エラー: %v", err)
		return nil, err
	}

	uc.logger.Infof("セッション開始成功: %s", session.ID.String())

	// セッションをレスポンス用に変換
	return session.ToResponse(), nil
}

// GetSession はセッションを取得する
func (uc *sessionUseCase) GetSession(ctx context.Context, id, userID uuid.UUID) (*model.SessionResponse, error) {
	session, err := uc.sessionRepo.GetByID(ctx, id, userID)
	if err != nil {
		uc.logger.Errorf("セッション取得エラー: %v", err)
		return nil, err
	}

	// セッションをレスポンス用に変換
	return session.ToResponse(), nil
}

// GetAllSessions はユーザの全セッションを取得する
func (uc *sessionUseCase) GetAllSessions(ctx context.Context, userID uuid.UUID) (*model.SessionsResponse, error) {
	sessions, err := uc.sessionRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		uc.logger.Errorf("セッション取得エラー: %v", err)
		return nil, err
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
	// セッションを取得
	_, err := uc.sessionRepo.GetByID(ctx, id, userID)
	if err != nil {
		uc.logger.Errorf("セッション取得エラー: %v", err)
		return nil, err
	}

	// ラウンドの統計情報を計算
	averageFocus, totalWorkMin, roundCount, breakTime, err := uc.roundRepo.CalculateSessionStats(ctx, id)
	if err != nil {
		uc.logger.Errorf("セッション統計情報取得エラー: %v", err)
		return nil, err
	}

	// セッションを完了する
	err = uc.sessionRepo.Complete(ctx, id, userID, averageFocus, totalWorkMin, roundCount, breakTime)
	if err != nil {
		uc.logger.Errorf("セッション完了エラー: %v", err)
		return nil, err
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
		uc.logger.Errorf("セッション取得エラー: %v", err)
		return nil, err
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
	// セッションを削除
	if err := uc.sessionRepo.Delete(ctx, id, userID); err != nil {
		uc.logger.Errorf("セッション削除エラー: %v", err)
		return err
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

// ensureUserConfig はユーザー設定を確認し、存在しない場合は作成する
func (uc *sessionUseCase) ensureUserConfig(ctx context.Context, userID uuid.UUID) (*model.UserConfig, error) {
	if uc.userConfigRepo == nil {
		uc.logger.Warn("UserConfigRepositoryが利用できません")
		return nil, nil
	}

	// GetOrCreateを使用して設定を取得または作成
	userConfig, err := uc.userConfigRepo.GetOrCreateUserConfig(ctx, userID)
	if err != nil {
		uc.logger.Errorf("ユーザー設定の取得または作成に失敗: %v", err)
		return nil, err
	}

	uc.logger.Infof("ユーザー設定確認完了: %s (work=%d分, break=%d分, rounds=%d)",
		userID.String(), userConfig.RoundWorkTime, userConfig.RoundBreakTime, userConfig.SessionRounds)

	return userConfig, nil
}
