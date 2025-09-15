package repository

import (
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/database"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/dynamodb"
)

// RepositoryFactory はすべてのリポジトリを管理するファクトリ（DDD準拠・API Gateway Authorizer対応版）
type RepositoryFactory struct {
	User                   repository.UserRepository
	Task                   repository.TaskRepository
	Session                repository.SessionRepository
	Statistics             repository.StatisticsRepository
	OptimizationPreferences repository.OptimizationPreferencesRepository
	// DDD原則: 1集約 = 1リポジトリ
	// 削除済み: Auth (API Gateway Authorizer使用), Round (SessionAggregateに統合), 
	// UserConfig (OptimizationPreferencesに統合), SampleOptimizationData (DDDルール違反), 
	// Optimization (DDDルール違反)
}

// NewRepositoryFactory はすべてのリポジトリを初期化する（API Gateway Authorizer対応版）
func NewRepositoryFactory(dynamoDB *database.DynamoDB, cfg *config.Config, logger logger.Logger) *RepositoryFactory {
	if dynamoDB == nil {
		logger.Fatal("DynamoDBが必須ですが初期化されていません")
		return nil
	}

	// DynamoDBを使用してDDD準拠のリポジトリを初期化
	userRepo := dynamodb.NewUserRepository(dynamoDB.Client, cfg, logger)
	taskRepo := dynamodb.NewTaskRepository(dynamoDB.Client, cfg, logger)
	sessionRepo := dynamodb.NewSessionRepository(dynamoDB.Client, cfg, logger)
	statisticsRepo := dynamodb.NewStatisticsRepository(dynamoDB.Client, cfg, logger)
	optimizationPreferencesRepo := dynamodb.NewOptimizationPreferencesRepository(dynamoDB.Client, cfg, logger)

	logger.Info("全リポジトリがDynamoDBで初期化されました（API Gateway Authorizer対応・DDD準拠）")

	return &RepositoryFactory{
		User:                   userRepo,
		Task:                   taskRepo,
		Session:                sessionRepo,
		Statistics:             statisticsRepo,
		OptimizationPreferences: optimizationPreferencesRepo,
	}
}
