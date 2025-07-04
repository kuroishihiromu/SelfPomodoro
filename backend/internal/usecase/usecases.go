package usecase

import (
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/sqs"
	dynamo "github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/dynamodb"
)

// UseCases はすべてのユースケースをまとめた構造体（AuthUseCase統合版）
type UseCases struct {
	Auth         AuthUseCase // 認証UseCase追加
	User         UserUseCase
	Task         TaskUseCase
	Session      SessionUseCase
	Round        RoundUseCase
	Statistics   StatisticsUsecase
	UserConfig   UserConfigUseCase
	Onboarding   OnboardingUseCase
	Optimization OptimizationUseCase
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
	optimizationRepo repository.OptimizationRepository,
	authRepo repository.AuthRepository, // AuthRepository追加
	sqsClient *sqs.SQSClient,
	cfg *config.Config, // Config追加（AuthUseCase用）
	logger logger.Logger,
) *UseCases {
	// 統計集約サービスの作成（DynamoDB統一後）
	var statsAggregator StatisticsAggregationService
	// DynamoDB統計リポジトリの型チェック
	if dynamoStatsRepo, ok := statisticsRepo.(*dynamo.StatisticsRepositoryImpl); ok {
		statsAggregator = NewStatisticsAggregationService(dynamoStatsRepo, logger)
		logger.Info("DynamoDB統計集約サービスを有効化しました")
	} else {
		// フォールバック：統計集約サービスを無効化
		statsAggregator = &noOpStatisticsAggregationService{logger: logger}
		logger.Warn("統計リポジトリがDynamoDB実装ではないため統計集約サービスを無効化しました")
	}

	return &UseCases{
		Auth:         NewAuthUseCase(authRepo, userRepo, cfg, logger), // AuthUseCase初期化
		User:         NewUserUseCase(userRepo, logger),
		Task:         NewTaskUseCase(taskRepo, logger),
		Session:      NewSessionUseCase(sessionRepo, roundRepo, userConfigRepo, sqsClient, logger),
		Round:        NewRoundUseCase(roundRepo, sessionRepo, userConfigRepo, sqsClient, statsAggregator, logger),
		Statistics:   NewStatisticsUsecase(statisticsRepo, logger),
		UserConfig:   NewUserConfigUseCase(userConfigRepo, logger),
		Onboarding:   NewOnboardingUseCase(userRepo, userConfigRepo, sampleDataRepo, logger),
		Optimization: NewOptimizationUseCase(optimizationRepo, userConfigRepo, logger),
	}
}
