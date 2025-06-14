package config

import (
	"os"
	"strconv"
)

// Config はLambda環境での設定を保持する構造体
type Config struct {
	// データベース設定
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSL_MODE"`

	// Cognito設定
	CognitoUserPoolID string `mapstructure:"COGNITO_USER_POOL_ID"`
	CognitoClientID   string `mapstructure:"COGNITO_CLIENT_ID"`

	// DynamoDB設定
	DynamoRegion                   string `mapstructure:"DYNAMO_REGION"`
	DynamoUserConfigTable          string `mapstructure:"DYNAMO_USER_CONFIG_TABLE"`
	DynamoRoundOptimizationTable   string `mapstructure:"DYNAMO_ROUND_OPTIMIZATION_TABLE"`
	DynamoSessionOptimizationTable string `mapstructure:"DYNAMO_SESSION_OPTIMIZATION_TABLE"`

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

// Load はLambda環境変数から設定を読み込む
func Load() (*Config, error) {
	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

	return &Config{
		// データベース設定
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     port,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "pomodoro"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "require"),

		// Cognito設定
		CognitoUserPoolID: getEnv("COGNITO_USER_POOL_ID", ""),
		CognitoClientID:   getEnv("COGNITO_CLIENT_ID", ""),

		// DynamoDB設定
		DynamoUserConfigTable:          getEnv("DYNAMO_USER_CONFIG_TABLE", "selfpomodoro_user_configs_dev"),
		DynamoRoundOptimizationTable:   getEnv("DYNAMO_ROUND_OPTIMIZATION_TABLE", "selfpomodoro_round_optimization_logs_dev"),
		DynamoSessionOptimizationTable: getEnv("DYNAMO_SESSION_OPTIMIZATION_TABLE", "selfpomodoro_session_optimization_logs_dev"),

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
