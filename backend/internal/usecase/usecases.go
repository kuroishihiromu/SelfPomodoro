package usecase

import (
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/sqs"
)

// UseCases はすべてのユースケースをまとめた構造体
type UseCases struct {
	Auth       AuthUseCase
	Task       TaskUseCase
	Session    SessionUseCase
	Round      RoundUseCase
	Statistics StatisticsUsecase
	UserConfig UserConfigUseCase
	// TODO: 他のユースケースを追加する場合はここにフィールドを追加
}

// NewUseCases はすべてのユースケースを初期化する
func NewUseCases(taskRepo repository.TaskRepository, sessionRepo repository.SessionRepository, roundRepo repository.RoundRepository, statisticsRepo repository.StatisticsRepository, userConfigRepo repository.UserConfigRepository, sqsClient *sqs.SQSClient, logger logger.Logger) *UseCases {
	return &UseCases{
		Auth:       NewAuthUseCase(logger),
		Task:       NewTaskUseCase(taskRepo, logger),
		Session:    NewSessionUseCase(sessionRepo, roundRepo, userConfigRepo, sqsClient, logger),
		Round:      NewRoundUseCase(roundRepo, sessionRepo, userConfigRepo, sqsClient, logger),
		Statistics: NewStatisticsUsecase(statisticsRepo, logger),
		UserConfig: NewUserConfigUseCase(userConfigRepo, logger),
		// TODO: 他のユースケースを初期化する場合はここに追加
	}
}
