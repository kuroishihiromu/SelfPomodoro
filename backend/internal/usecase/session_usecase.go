package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// SessionUseCase はセッションに関するユースケースを定義するインターフェース
type SessionUseCase interface {
	// StartSession は新しいセッションを開始する
	StartSession(ctx context.Context, userID uuid.UUID) (*model.SessionResponse, error)

	// GetSession はセッションを取得する
	GetSession(ctx context.Context, id, userID uuid.UUID) (*model.SessionResponse, error)

	// GetAllSessions はユーザの全セッションを取得する
	GetAllSessions(ctx context.Context, userID uuid.UUID) (*model.SessionsResponse, error)

	// CompleteSession はセッションを完了する
	CompleteSession(ctx context.Context, id, userID uuid.UUID) (*model.SessionResponse, error)

	// DeleteSession はセッションを削除する
	DeleteSession(ctx context.Context, id, userID uuid.UUID) error
}

// sessionUseCase はSessionUseCaseインターフェースの実装
type sessionUseCase struct {
	sessionRepo repository.SessionRepository
	roundRepo   repository.RoundRepository
	logger      logger.Logger
}

// NewSessionUseCase は新しいSessionUseCaseインスタンスを作成する
func NewSessionUseCase(sessionRepo repository.SessionRepository, roundRepo repository.RoundRepository, logger logger.Logger) SessionUseCase {
	return &sessionUseCase{
		sessionRepo: sessionRepo,
		roundRepo:   roundRepo,
		logger:      logger,
	}
}

// StartSession は新しいセッションを開始する
func (uc *sessionUseCase) StartSession(ctx context.Context, userID uuid.UUID) (*model.SessionResponse, error) {
	session := model.NewSession(userID)

	// DBにセッションを保存
	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		uc.logger.Errorf("セッション開始エラー: %v", err)
		return nil, err
	}

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

// CompleteSession はセッションを完了する
func (uc *sessionUseCase) CompleteSession(ctx context.Context, id, userID uuid.UUID) (*model.SessionResponse, error) {
	// セッションを取得
	_, err := uc.sessionRepo.GetByID(ctx, id, userID)
	if err != nil {
		uc.logger.Errorf("セッション取得エラー: %v", err)
		return nil, err
	}

	// TODO: ラウンドの平均集中度、総作業時間、ラウンド数、休憩時間を計算するロジックを実装

	// ここでは仮の値を設定
	roundCount := 0
	totalWorkMin := 0
	breakTime := 0
	averageFocus := 0.0

	// セッションを完了する
	err = uc.sessionRepo.Complete(ctx, id, userID, averageFocus, totalWorkMin, roundCount, breakTime)
	if err != nil {
		uc.logger.Errorf("セッション完了エラー: %v", err)
		return nil, err
	}

	updatedSession, err := uc.sessionRepo.GetByID(ctx, id, userID)
	if err != nil {
		uc.logger.Errorf("セッション取得エラー: %v", err)
		return nil, err
	}

	// TODO: 最適化Lambda実行
	// TODO: ログの設定と更新

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
	return nil
}
