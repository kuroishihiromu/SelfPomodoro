package usecase

import (
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// UseCases はすべてのユースケースをまとめた構造体
type UseCases struct {
	Auth    AuthUseCase
	Task    TaskUseCase
	Session SessionUseCase
	// TODO: 他のユースケースを追加する場合はここにフィールドを追加
}

// NewUseCases はすべてのユースケースを初期化する
func NewUseCases(taskRepo repository.TaskRepository, sessionRepo repository.SessionRepository, logger logger.Logger) *UseCases {
	return &UseCases{
		Auth:    NewAuthUseCase(logger),
		Task:    NewTaskUseCase(taskRepo, logger),
		Session: NewSessionUseCase(sessionRepo, nil, logger),
		// TODO: 他のユースケースを初期化する場合はここに追加
	}
}
