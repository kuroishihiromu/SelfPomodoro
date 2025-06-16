package repository

import (
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/database"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/auth"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/dynamodb"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/postgres"
)

// RepositoryFactory はすべてのリポジトリを管理するファクトリ
type RepositoryFactory struct {
	Auth                   repository.AuthRepository
	User                   repository.UserRepository
	Task                   repository.TaskRepository
	Session                repository.SessionRepository
	Round                  repository.RoundRepository
	Statistics             repository.StatisticsRepository
	UserConfig             repository.UserConfigRepository
	SampleOptimizationData repository.SampleOptimizationDataRepository
	// TODO: 他のリポジトリを追加する場合はここにフィールドを追加
}

// NewRepositoryFactory はすべてのリポジトリを初期化する
func NewRepositoryFactory(postgresDB *database.PostgresDB, dynamoDB *database.DynamoDB, cfg *config.Config, logger logger.Logger) *RepositoryFactory {
	// PostgresDBを使用してリポジトリを初期化
	authRepo := auth.NewCognitoAuthRepository(cfg, logger)
	userRepo := postgres.NewUserRepository(postgresDB, logger)
	taskRepo := postgres.NewTaskRepository(postgresDB, logger)
	sessionRepo := postgres.NewSessionRepository(postgresDB, logger)
	roundRepo := postgres.NewRoundRepository(postgresDB, logger)
	statisticsRepo := postgres.NewStatisticsRepository(postgresDB, logger)

	// DynamoDBを使用してリポジトリを初期化
	var userConfigRepo repository.UserConfigRepository
	var sampleOptimizationDataRepo repository.SampleOptimizationDataRepository
	if dynamoDB != nil {
		userConfigRepo = dynamodb.NewUserConfigRepository(dynamoDB.Client, dynamoDB.Config, logger)
		sampleOptimizationDataRepo = dynamodb.NewSampleOptimizationDataRepository(dynamoDB.Client, dynamoDB.Config, logger)
	} else {
		// DynamoDBが利用できない場合はnilを設定(エラーハンドリングは各usecaseで行う)
		logger.Warn("DynamoDBが利用できないため、UserConfigRepositoryはnilになります")
		userConfigRepo = nil
		sampleOptimizationDataRepo = nil
	}

	// TODO: DynamoDBを使用してリポジトリを初期化する場合はここに追加

	return &RepositoryFactory{
		Auth:                   authRepo,
		User:                   userRepo,
		Task:                   taskRepo,
		Session:                sessionRepo,
		Round:                  roundRepo,
		Statistics:             statisticsRepo,
		UserConfig:             userConfigRepo,
		SampleOptimizationData: sampleOptimizationDataRepo,
	}
}
