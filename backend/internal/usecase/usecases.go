package usecase

import (
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/messaging/sqs"
)

// UseCases はすべてのユースケースをまとめた構造体（API Gateway Authorizer対応版）
type UseCases struct {
	User         UserUseCase
	Task         TaskUseCase
	Session      SessionUseCase
	Round        RoundUseCase
	Statistics   StatisticsUsecase
	OptimizationPreferences OptimizationPreferencesUseCase
	Onboarding   OnboardingUseCase
}

// NewUseCases はすべてのユースケースを初期化する（API Gateway Authorizer対応版）
func NewUseCases(
	userRepo repository.UserRepository,
	taskRepo repository.TaskRepository,
	sessionRepo repository.SessionRepository,
	statisticsRepo repository.StatisticsRepository,
	optimizationPreferencesRepo repository.OptimizationPreferencesRepository,
	sqsClient *sqs.SQSClient,
	cfg *config.Config,
	logger logger.Logger,
) *UseCases {

	// SQSメッセージングサービスを作成
	messagingService := sqs.NewSQSMessagingService(sqsClient, logger)

	return &UseCases{
		User:         NewUserUseCase(userRepo, logger),
		Task:         NewTaskUseCase(taskRepo, logger),
		Session:      NewSessionUseCase(sessionRepo, optimizationPreferencesRepo, messagingService, logger),
		Round:        NewRoundUseCase(sessionRepo, optimizationPreferencesRepo, messagingService, statisticsRepo, logger),
		Statistics:   NewStatisticsUsecase(statisticsRepo, logger),
		OptimizationPreferences: NewOptimizationPreferencesUseCase(optimizationPreferencesRepo, logger),
		Onboarding:   NewOnboardingUseCase(userRepo, optimizationPreferencesRepo, sessionRepo, logger),
	}
}
