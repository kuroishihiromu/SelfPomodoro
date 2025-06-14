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

// NewDynamoDB は新しいDynamoDBクライアントを作成する（Lambda最適化版）
func NewDynamoDB(cfg *config.Config, logger logger.Logger) (*DynamoDB, error) {
	// AWS設定の読み込み
	logger.Infof("DynamoDBクライアントを初期化: リージョン=%s", cfg.AWSRegion)

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.AWSRegion),
	)
	if err != nil {
		return nil, fmt.Errorf("AWS設定読み込みエラー: %w", err)
	}

	// DynamoDBクライアントの作成
	client := dynamodb.NewFromConfig(awsCfg)

	// Lambda環境での最適化: テーブル存在確認は軽量化
	if cfg.Environment != "production" {
		// 開発環境でのみテーブル存在確認を実行
		tables := []string{
			cfg.DynamoUserConfigTable,
			cfg.DynamoRoundOptimizationTable,
			cfg.DynamoSessionOptimizationTable,
		}

		for _, table := range tables {
			err := checkTableExists(client, table)
			if err != nil {
				logger.Warnf("DynamoDBテーブル %s が存在しないか、アクセスできません: %v", table, err)
			} else {
				logger.Infof("DynamoDBテーブル %s が利用可能", table)
			}
		}
	}

	logger.Info("DynamoDB接続成功")

	return &DynamoDB{
		Client: client,
		Config: cfg,
		logger: logger,
	}, nil
}

// checkTableExists はテーブルの存在確認を行う（軽量版）
func checkTableExists(client *dynamodb.Client, tableName string) error {
	_, err := client.DescribeTable(context.Background(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	return err
}

// GetTableName はテーブル名を取得する
func (d *DynamoDB) GetTableName(tableType string) string {
	switch tableType {
	case "user_config":
		return d.Config.DynamoUserConfigTable
	case "round_optimization":
		return d.Config.DynamoRoundOptimizationTable
	case "session_optimization":
		return d.Config.DynamoSessionOptimizationTable
	default:
		d.logger.Warnf("未知のテーブルタイプ: %s", tableType)
		return ""
	}
}

// Close はリソースをクリーンアップする（Lambda用）
func (d *DynamoDB) Close() error {
	// DynamoDBクライアントには明示的なCloseメソッドはないが、
	// 将来的な拡張のためにメソッドを定義
	d.logger.Debug("DynamoDB接続をクリーンアップ")
	return nil
}
