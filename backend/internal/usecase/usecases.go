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
	Round   RoundUseCase
	// TODO: 他のユースケースを追加する場合はここにフィールドを追加
}

// NewUseCases はすべてのユースケースを初期化する
func NewUseCases(taskRepo repository.TaskRepository, sessionRepo repository.SessionRepository, roundRepo repository.RoundRepository, logger logger.Logger) *UseCases {
	return &UseCases{
		Auth:    NewAuthUseCase(logger),
		Task:    NewTaskUseCase(taskRepo, logger),
		Session: NewSessionUseCase(sessionRepo, roundRepo, logger),
		Round:   NewRoundUseCase(roundRepo, sessionRepo, logger),
		// TODO: 他のユースケースを初期化する場合はここに追加
	}
}
