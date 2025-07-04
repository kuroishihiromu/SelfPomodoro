package config

import (
	"os"
)

// Config はLambda環境での設定を保持する構造体（DynamoDB完全移行版）
type Config struct {
	// Cognito設定
	CognitoUserPoolID string `mapstructure:"COGNITO_USER_POOL_ID"`
	CognitoClientID   string `mapstructure:"COGNITO_CLIENT_ID"`

	// DynamoDB設定
	DynamoRegion      string `mapstructure:"DYNAMO_REGION"`
	DynamoUnifiedTable string `mapstructure:"DYNAMO_UNIFIED_TABLE"`

	// SQS設定
	SQSRoundOptimizationURL   string `mapstructure:"SQS_ROUND_OPTIMIZATION_URL"`
	SQSSessionOptimizationURL string `mapstructure:"SQS_SESSION_OPTIMIZATION_URL"`

	// AWS設定
	AWSRegion string `mapstructure:"AWS_REGION"`

	// ロギング設定
	LogLevel string `mapstructure:"LOG_LEVEL"`

	// 環境設定
	Environment string `mapstructure:"ENVIRONMENT"`
}

// Load はLambda環境変数から設定を読み込む（DynamoDB完全移行版）
func Load() (*Config, error) {
	return &Config{
		// Cognito設定
		CognitoUserPoolID: getEnv("COGNITO_USER_POOL_ID", ""),
		CognitoClientID:   getEnv("COGNITO_CLIENT_ID", ""),

		// DynamoDB設定
		DynamoUnifiedTable: getEnv("DYNAMO_UNIFIED_TABLE", "selfpomodoro_unified_table_dev"),

		// SQS設定
		SQSRoundOptimizationURL:   getEnv("SQS_ROUND_OPTIMIZATION_URL", ""),
		SQSSessionOptimizationURL: getEnv("SQS_SESSION_OPTIMIZATION_URL", ""),

		// AWS設定
		AWSRegion: getEnv("AWS_REGION", "ap-northeast-1"),

		// ロギング設定
		LogLevel: getEnv("LOG_LEVEL", "info"),

		// 環境設定
		Environment: getEnv("ENVIRONMENT", "development"),
	}, nil
}

// getEnv は環境変数を取得し、存在しない場合はデフォルト値を返す
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

