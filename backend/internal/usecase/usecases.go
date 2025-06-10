package usecase

import (
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/sqs"
)

// UseCases はすべてのユースケースをまとめた構造体（AuthUseCase統合版）
type UseCases struct {
	Auth       AuthUseCase // 認証UseCase追加
	User       UserUseCase
	Task       TaskUseCase
	Session    SessionUseCase
	Round      RoundUseCase
	Statistics StatisticsUsecase
	UserConfig UserConfigUseCase
	Onboarding OnboardingUseCase
}

// NewUseCases はすべてのユースケースを初期化する（AuthUseCase統合版）
func NewUseCases(
	userRepo repository.UserRepository,
	taskRepo repository.TaskRepository,
	sessionRepo repository.SessionRepository,
	roundRepo repository.RoundRepository,
	statisticsRepo repository.StatisticsRepository,
	userConfigRepo repository.UserConfigRepository,
	sampleDataRepo repository.SampleOptimizationDataRepository,
	authRepo repository.AuthRepository, // AuthRepository追加
	sqsClient *sqs.SQSClient,
	cfg *config.Config, // Config追加（AuthUseCase用）
	logger logger.Logger,
) *UseCases {
	return &UseCases{
		Auth:       NewAuthUseCase(authRepo, userRepo, cfg, logger), // AuthUseCase初期化
		User:       NewUserUseCase(userRepo, logger),
		Task:       NewTaskUseCase(taskRepo, logger),
		Session:    NewSessionUseCase(sessionRepo, roundRepo, userConfigRepo, sqsClient, logger),
		Round:      NewRoundUseCase(roundRepo, sessionRepo, userConfigRepo, sqsClient, logger),
		Statistics: NewStatisticsUsecase(statisticsRepo, logger),
		UserConfig: NewUserConfigUseCase(userConfigRepo, logger),
		Onboarding: NewOnboardingUseCase(userRepo, userConfigRepo, sampleDataRepo, logger),
	}
}
