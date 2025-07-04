package repository

import (
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/database"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/auth"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/dynamodb"
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
	Optimization           repository.OptimizationRepository
	// TODO: 他のリポジトリを追加する場合はここにフィールドを追加
}

// NewRepositoryFactory はすべてのリポジトリを初期化する（DynamoDB完全移行版）
func NewRepositoryFactory(dynamoDB *database.DynamoDB, cfg *config.Config, logger logger.Logger) *RepositoryFactory {
	if dynamoDB == nil {
		logger.Fatal("DynamoDBが必須ですが初期化されていません")
		return nil
	}

	// 認証リポジトリ（Cognito使用）
	authRepo := auth.NewCognitoAuthRepository(cfg, logger)
	
	// DynamoDBを使用してすべてのリポジトリを初期化
	userRepo := dynamodb.NewUserRepository(dynamoDB.Client, cfg, logger)
	taskRepo := dynamodb.NewTaskRepository(dynamoDB.Client, cfg, logger)
	sessionRepo := dynamodb.NewSessionRepository(dynamoDB.Client, cfg, logger)
	// RoundRepository（依存関係削除済み）
	roundRepo := dynamodb.NewRoundRepository(dynamoDB.Client, cfg, logger)
	statisticsRepo := dynamodb.NewStatisticsRepository(dynamoDB.Client, cfg, logger)
	userConfigRepo := dynamodb.NewUserConfigRepository(dynamoDB.Client, cfg, logger)
	sampleOptimizationDataRepo := dynamodb.NewSampleOptimizationDataRepository(dynamoDB.Client, cfg, logger)
	optimizationRepo := dynamodb.NewOptimizationRepository(dynamoDB.Client, cfg.DynamoUnifiedTable, logger)

	logger.Info("全リポジトリがDynamoDBで初期化されました")

	return &RepositoryFactory{
		Auth:                   authRepo,
		User:                   userRepo,
		Task:                   taskRepo,
		Session:                sessionRepo,
		Round:                  roundRepo,
		Statistics:             statisticsRepo,
		UserConfig:             userConfigRepo,
		SampleOptimizationData: sampleOptimizationDataRepo,
		Optimization:           optimizationRepo,
	}
}
