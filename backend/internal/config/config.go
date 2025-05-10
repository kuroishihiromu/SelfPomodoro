package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// 設定を一時的に保持するための内部構造体
type configValues struct {
	// サーバー設定
	ServerPort         int `mapstructure:"SERVER_PORT"`
	ServerReadTimeout  int `mapstructure:"SERVER_READ_TIMEOUT"`
	ServerWriteTimeout int `mapstructure:"SERVER_WRITE_TIMEOUT"`
	ServerIdleTimeout  int `mapstructure:"SERVER_IDLE_TIMEOUT"`

	// データベース設定
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSL_MODE"`

	// DynamoDB設定
	DynamoRegion                   string `mapstructure:"DYNAMO_REGION"`
	DynamoRoundOptimizationTable   string `mapstructure:"DYNAMO_ROUND_OPTIMIZATION_TABLE"`
	DynamoSessionOptimizationTable string `mapstructure:"DYNAMO_SESSION_OPTIMIZATION_TABLE"`
	DynamoUserConfigTable          string `mapstructure:"DYNAMO_USER_CONFIG_TABLE"`

	// 認証設定
	CognitoRegion        string `mapstructure:"COGNITO_REGION"`
	CognitoUserPoolID    string `mapstructure:"COGNITO_USER_POOL_ID"`
	CognitoClientID      string `mapstructure:"COGNITO_CLIENT_ID"`
	JWTSecret            string `mapstructure:"JWT_SECRET"`
	JWTExpiration        int    `mapstructure:"JWT_EXPIRATION"`
	JWTRefreshExpiration int    `mapstructure:"JWT_REFRESH_EXPIRATION"`

	// SNS/SQS設定
	SNSRegion          string `mapstructure:"SNS_REGION"`
	SNSRoundTopicArn   string `mapstructure:"SNS_ROUND_TOPIC_ARN"`
	SNSSessionTopicArn string `mapstructure:"SNS_SESSION_TOPIC_ARN"`

	// ロギング設定
	LogLevel string `mapstructure:"LOG_LEVEL"`

	// 環境設定
	Environment string `mapstructure:"ENVIRONMENT"`
}

// Config はアプリケーションの設定を保持する構造体
type Config struct {
	// サーバー設定
	ServerPort         int           `mapstructure:"SERVER_PORT"`
	ServerReadTimeout  time.Duration `mapstructure:"SERVER_READ_TIMEOUT"`
	ServerWriteTimeout time.Duration `mapstructure:"SERVER_WRITE_TIMEOUT"`
	ServerIdleTimeout  time.Duration `mapstructure:"SERVER_IDLE_TIMEOUT"`

	// データベース設定
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSL_MODE"`

	// DynamoDB設定
	DynamoRegion                   string `mapstructure:"DYNAMO_REGION"`
	DynamoRoundOptimizationTable   string `mapstructure:"DYNAMO_ROUND_OPTIMIZATION_TABLE"`
	DynamoSessionOptimizationTable string `mapstructure:"DYNAMO_SESSION_OPTIMIZATION_TABLE"`
	DynamoUserConfigTable          string `mapstructure:"DYNAMO_USER_CONFIG_TABLE"`

	// 認証設定
	CognitoRegion        string        `mapstructure:"COGNITO_REGION"`
	CognitoUserPoolID    string        `mapstructure:"COGNITO_USER_POOL_ID"`
	CognitoClientID      string        `mapstructure:"COGNITO_CLIENT_ID"`
	JWTSecret            string        `mapstructure:"JWT_SECRET"`
	JWTExpiration        time.Duration `mapstructure:"JWT_EXPIRATION"`
	JWTRefreshExpiration time.Duration `mapstructure:"JWT_REFRESH_EXPIRATION"`

	// SNS/SQS設定
	SNSRegion          string `mapstructure:"SNS_REGION"`
	SNSRoundTopicArn   string `mapstructure:"SNS_ROUND_TOPIC_ARN"`
	SNSSessionTopicArn string `mapstructure:"SNS_SESSION_TOPIC_ARN"`

	// ロギング設定
	LogLevel string `mapstructure:"LOG_LEVEL"`

	// 環境設定
	Environment string `mapstructure:"ENVIRONMENT"`
}

// Load は設定ファイルから設定を読み込む
func Load() (*Config, error) {
	// 環境設定
	viper.SetDefault("ENVIRONMENT", "development")

	// .env ファイルを読み込む
	envFile := ".env"
	if _, err := os.Stat(envFile); err == nil {
		viper.SetConfigFile(envFile)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("警告: %s を読み込めませんでした: %v\n", envFile, err)
		} else {
			fmt.Printf("設定ファイル読み込み: %s\n", envFile)
		}
	}

	// デフォルト値設定
	setDefaults()

	// 環境変数から設定を読み込む（.envファイルの内容を上書きする場合）
	viper.AutomaticEnv()

	// 設定を一時構造体に読み込む
	var values configValues
	if err := viper.Unmarshal(&values); err != nil {
		return nil, fmt.Errorf("設定読み込みエラー: %w", err)
	}

	// 最終的な設定構造体を作成
	config := &Config{
		// サーバー設定
		ServerPort:         values.ServerPort,
		ServerReadTimeout:  time.Duration(values.ServerReadTimeout) * time.Second,
		ServerWriteTimeout: time.Duration(values.ServerWriteTimeout) * time.Second,
		ServerIdleTimeout:  time.Duration(values.ServerIdleTimeout) * time.Second,

		// データベース設定
		DBHost:     values.DBHost,
		DBPort:     values.DBPort,
		DBUser:     values.DBUser,
		DBPassword: values.DBPassword,
		DBName:     values.DBName,
		DBSSLMode:  values.DBSSLMode,

		// DynamoDB設定
		DynamoRegion:                   values.DynamoRegion,
		DynamoRoundOptimizationTable:   values.DynamoRoundOptimizationTable,
		DynamoSessionOptimizationTable: values.DynamoSessionOptimizationTable,
		DynamoUserConfigTable:          values.DynamoUserConfigTable,

		// 認証設定
		CognitoRegion:        values.CognitoRegion,
		CognitoUserPoolID:    values.CognitoUserPoolID,
		CognitoClientID:      values.CognitoClientID,
		JWTSecret:            values.JWTSecret,
		JWTExpiration:        time.Duration(values.JWTExpiration) * time.Hour,
		JWTRefreshExpiration: time.Duration(values.JWTRefreshExpiration) * time.Hour,

		// SNS/SQS設定
		SNSRegion:          values.SNSRegion,
		SNSRoundTopicArn:   values.SNSRoundTopicArn,
		SNSSessionTopicArn: values.SNSSessionTopicArn,

		// ロギング設定
		LogLevel: values.LogLevel,

		// 環境設定
		Environment: values.Environment,
	}

	return config, nil
}

// setDefaults 設定のデフォルト値を設定
func setDefaults() {
	// サーバー
	viper.SetDefault("SERVER_PORT", 8080)
	viper.SetDefault("SERVER_READ_TIMEOUT", 10)  // 秒
	viper.SetDefault("SERVER_WRITE_TIMEOUT", 10) // 秒
	viper.SetDefault("SERVER_IDLE_TIMEOUT", 60)  // 秒

	// DB
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_NAME", "pomodoro")
	viper.SetDefault("DB_SSL_MODE", "disable")

	// DynamoDB
	viper.SetDefault("DYNAMO_REGION", "ap-northeast-1")
	viper.SetDefault("DYNAMO_ROUND_OPTIMIZATION_TABLE", "round_optimization_logs")
	viper.SetDefault("DYNAMO_SESSION_OPTIMIZATION_TABLE", "session_optimization_logs")
	viper.SetDefault("DYNAMO_USER_CONFIG_TABLE", "user_configs")

	// 認証
	viper.SetDefault("JWT_EXPIRATION", 24)          // 時間
	viper.SetDefault("JWT_REFRESH_EXPIRATION", 168) // 時間 (7日)

	// ロギング
	viper.SetDefault("LOG_LEVEL", "info")
}
