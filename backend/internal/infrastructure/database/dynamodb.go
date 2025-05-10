package database

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// DynamoDB はDynamoDBクライアントを管理する
type DynamoDB struct {
	Client *dynamodb.Client
	Config *config.Config
	logger logger.Logger
}

// NewDynamoDB は新しいDynamoDBクライアントを作成する
func NewDynamoDB(cfg *config.Config, logger logger.Logger) (*DynamoDB, error) {
	// AWS設定の読み込み
	logger.Infof("DynamoDBクライアントを初期化: リージョン=%s", cfg.DynamoRegion)
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.DynamoRegion),
	)
	if err != nil {
		return nil, fmt.Errorf("AWS設定読み込みエラー: %w", err)
	}

	// DynamoDBクライアントの作成
	client := dynamodb.NewFromConfig(awsCfg)

	// テーブルの存在確認
	tables := []string{
		cfg.DynamoRoundOptimizationTable,
		cfg.DynamoSessionOptimizationTable,
		cfg.DynamoUserConfigTable,
	}

	for _, table := range tables {
		_, err := client.DescribeTable(context.Background(), &dynamodb.DescribeTableInput{
			TableName: aws.String(table),
		})
		if err != nil {
			// テーブルがない場合は警告のみ（開発環境ではテーブルが存在しない場合がある）
			logger.Warnf("DynamoDBテーブル %s が存在しないか、アクセスできません: %v", table, err)
		} else {
			logger.Infof("DynamoDBテーブル %s が利用可能", table)
		}
	}

	return &DynamoDB{
		Client: client,
		Config: cfg,
		logger: logger,
	}, nil
}

// GetTableName はテーブル名を取得する
func (d *DynamoDB) GetTableName(tableType string) string {
	switch tableType {
	case "round_optimization":
		return d.Config.DynamoRoundOptimizationTable
	case "session_optimization":
		return d.Config.DynamoSessionOptimizationTable
	case "user_config":
		return d.Config.DynamoUserConfigTable
	default:
		d.logger.Warnf("未知のテーブルタイプ: %s", tableType)
		return ""
	}
}
