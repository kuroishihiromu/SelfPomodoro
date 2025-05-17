package repository

import (
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/database"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/postgres"
)

// RepositoryFactory はすべてのリポジトリを管理するファクトリ
type RepositoryFactory struct {
	Task    repository.TaskRepository
	Session repository.SessionRepository
	// TODO: 他のリポジトリを追加する場合はここにフィールドを追加
}

// NewRepositoryFactory はすべてのリポジトリを初期化する
func NewRepositoryFactory(postgresDB *database.PostgresDB, dynamoDB *database.DynamoDB, logger logger.Logger) *RepositoryFactory {
	// PostgresDBを使用してリポジトリを初期化
	taskRepo := postgres.NewTaskRepository(postgresDB, logger)
	sessionRepo := postgres.NewSessionRepository(postgresDB, logger)

	// TODO: DynamoDBを使用してリポジトリを初期化する場合はここに追加

	return &RepositoryFactory{
		Task:    taskRepo,
		Session: sessionRepo,
	}
}
