package container

import (
	"context"
	"fmt"
	"sync"

	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/database"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/messaging/sqs"
	"github.com/tsunakit99/selfpomodoro/internal/usecase"
)

// Container はDependency Injection コンテナのインターフェース
type Container interface {
	// Core Dependencies
	GetUseCases() *usecase.UseCases
	GetLogger() logger.Logger
	GetConfig() *config.Config

	// Lifecycle
	Initialize(ctx context.Context) error
	Cleanup() error
	HealthCheck(ctx context.Context) error
}

// LambdaContainer はLambda環境最適化されたDIコンテナ
type LambdaContainer struct {
	// Dependencies
	useCases *usecase.UseCases
	logger   logger.Logger
	config   *config.Config

	// Infrastructure Services
	infraServices *InfrastructureServices

	// State management
	initialized bool
	mu          sync.RWMutex
}

// NewLambdaContainer は新しいコンテナを作成
func NewLambdaContainer() Container {
	return &LambdaContainer{}
}

// Initialize はコンテナを初期化（Lambda init phase用）
func (c *LambdaContainer) Initialize(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized {
		return nil // 既に初期化済み
	}

	// 1. Config読み込み
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("設定読み込みエラー: %w", err)
	}
	c.config = cfg

	// 2. Logger初期化
	appLogger, err := logger.NewLogger(cfg.LogLevel, cfg.Environment)
	if err != nil {
		return fmt.Errorf("ロガー初期化エラー: %w", err)
	}
	c.logger = appLogger

	// 3. Infrastructure Services初期化
	infraServices, err := c.initializeInfrastructure(cfg, appLogger)
	if err != nil {
		return fmt.Errorf("インフラサービス初期化エラー: %w", err)
	}
	c.infraServices = infraServices // 保存

	// 4. UseCases初期化（API Gateway Authorizer対応版）
	c.useCases = usecase.NewUseCases(
		infraServices.Repositories.User,
		infraServices.Repositories.Task,
		infraServices.Repositories.Session,
		infraServices.Repositories.Statistics,
		infraServices.Repositories.OptimizationPreferences,
		infraServices.SQSClient,
		cfg,
		appLogger,
	)

	c.initialized = true
	appLogger.Info("Lambda Container初期化完了")
	return nil
}

// initializeInfrastructure はインフラストラクチャサービスを初期化（DynamoDB完全移行版）
func (c *LambdaContainer) initializeInfrastructure(cfg *config.Config, logger logger.Logger) (*InfrastructureServices, error) {
	// 1. DynamoDB接続（必須）
	dynamoDB, err := database.NewDynamoDB(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("DynamoDB接続エラー: %w", err)
	}

	// 2. SQS接続（オプショナル）
	var sqsClient *sqs.SQSClient
	if cfg.SQSRoundOptimizationURL != "" {
		sqsClient, err = sqs.NewSQSClient(cfg, logger)
		if err != nil {
			logger.Warnf("SQS接続失敗、続行します: %v", err)
			sqsClient = nil
		}
	}

	// 3. Repository Factory初期化（DynamoDBのみ）
	repositoryFactory := repository.NewRepositoryFactory(dynamoDB, cfg, logger)

	return &InfrastructureServices{
		Repositories: repositoryFactory,
		SQSClient:    sqsClient,
		DynamoDB:     dynamoDB,
	}, nil
}

// Getter methods
func (c *LambdaContainer) GetUseCases() *usecase.UseCases {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.useCases
}

func (c *LambdaContainer) GetLogger() logger.Logger {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.logger
}

func (c *LambdaContainer) GetConfig() *config.Config {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}

// Cleanup はリソースのクリーンアップを行う
func (c *LambdaContainer) Cleanup() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.initialized {
		return nil
	}

	var errors []error

	// DynamoDBクローズ
	if c.infraServices != nil && c.infraServices.DynamoDB != nil {
		if err := c.infraServices.DynamoDB.Close(); err != nil {
			errors = append(errors, fmt.Errorf("DynamoDB close error: %w", err))
		}
	}

	// SQSクローズ
	if c.infraServices != nil && c.infraServices.SQSClient != nil {
		if err := c.infraServices.SQSClient.Close(); err != nil {
			errors = append(errors, fmt.Errorf("SQS close error: %w", err))
		}
	}

	if len(errors) > 0 {
		c.logger.Errorf("クリーンアップエラー: %v", errors)
		return errors[0] // 最初のエラーを返す
	}

	c.logger.Info("リソースクリーンアップ完了")
	c.initialized = false
	return nil
}

// HealthCheck は全サービスの接続確認
func (c *LambdaContainer) HealthCheck(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.initialized {
		return fmt.Errorf("コンテナが初期化されていません")
	}

	// API Gateway Authorizerを使用するため、認証ヘルスチェックは不要

	// SQSのヘルスチェック
	if c.infraServices.SQSClient != nil {
		if err := c.infraServices.SQSClient.HealthCheck(ctx); err != nil {
			return fmt.Errorf("SQS接続確認失敗: %w", err)
		}
	}

	c.logger.Info("全サービス接続確認成功")
	return nil
}

// InfrastructureServices はインフラ層のサービス群（DynamoDB完全移行版）
type InfrastructureServices struct {
	Repositories *repository.RepositoryFactory
	SQSClient    *sqs.SQSClient
	DynamoDB     *database.DynamoDB
}
