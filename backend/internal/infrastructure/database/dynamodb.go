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

	// Lambda環境での最適化: 統合テーブル存在確認
	if cfg.Environment != "production" {
		// 開発環境でのみ統合テーブル存在確認を実行
		err := checkTableExists(client, cfg.DynamoUnifiedTable)
		if err != nil {
			logger.Warnf("DynamoDB統合テーブル %s が存在しないか、アクセスできません: %v", cfg.DynamoUnifiedTable, err)
		} else {
			logger.Infof("DynamoDB統合テーブル %s が利用可能", cfg.DynamoUnifiedTable)
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

// GetTableName は統合テーブル名を取得する
func (d *DynamoDB) GetTableName() string {
	return d.Config.DynamoUnifiedTable
}

// Close はリソースをクリーンアップする（Lambda用）
func (d *DynamoDB) Close() error {
	// DynamoDBクライアントには明示的なCloseメソッドはないが、
	// 将来的な拡張のためにメソッドを定義
	d.logger.Debug("DynamoDB接続をクリーンアップ")
	return nil
}
