package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// UserConfigRepositoryImpl はDynamoDBを使用したUserConfigRepositoryの実装（新エラーハンドリング対応版）
type UserConfigRepositoryImpl struct {
	client    *dynamodb.Client
	tableName string
	logger    logger.Logger
}

// NewUserConfigRepository は新しいUserConfigRepositoryImplインスタンスを作成する
func NewUserConfigRepository(client *dynamodb.Client, cfg *config.Config, logger logger.Logger) repository.UserConfigRepository {
	return &UserConfigRepositoryImpl{
		client:    client,
		tableName: cfg.DynamoUserConfigTable,
		logger:    logger,
	}
}

// GetUserConfig はユーザーIDからユーザー設定を取得する（新エラーハンドリング対応版）
func (r *UserConfigRepositoryImpl) GetUserConfig(ctx context.Context, userID uuid.UUID) (*model.UserConfig, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: userID.String()},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB GetItem エラー: %v", err)
		return nil, appErrors.NewDynamoDBOperationError("get_user_config", err)
	}

	if result.Item == nil {
		r.logger.Debugf("ユーザー設定が見つかりません: %s", userID.String())
		return nil, appErrors.ErrDynamoDBItemNotFound // Infrastructure Error
	}

	// 手動でDynamoDBアイテムから構造体に変換
	config := &model.UserConfig{}

	// user_id
	if userIDAttr, exists := result.Item["user_id"]; exists {
		if s, ok := userIDAttr.(*types.AttributeValueMemberS); ok {
			config.UserID = s.Value
		}
	}

	// round_work_time
	if workTimeAttr, exists := result.Item["round_work_time"]; exists {
		if n, ok := workTimeAttr.(*types.AttributeValueMemberN); ok {
			if workTime, err := strconv.Atoi(n.Value); err == nil {
				config.RoundWorkTime = workTime
			}
		}
	}

	// round_break_time
	if breakTimeAttr, exists := result.Item["round_break_time"]; exists {
		if n, ok := breakTimeAttr.(*types.AttributeValueMemberN); ok {
			if breakTime, err := strconv.Atoi(n.Value); err == nil {
				config.RoundBreakTime = breakTime
			}
		}
	}

	// session_rounds
	if roundsAttr, exists := result.Item["session_rounds"]; exists {
		if n, ok := roundsAttr.(*types.AttributeValueMemberN); ok {
			if rounds, err := strconv.Atoi(n.Value); err == nil {
				config.SessionRounds = rounds
			}
		}
	}

	// session_break_time
	if sessionBreakAttr, exists := result.Item["session_break_time"]; exists {
		if n, ok := sessionBreakAttr.(*types.AttributeValueMemberN); ok {
			if sessionBreak, err := strconv.Atoi(n.Value); err == nil {
				config.SessionBreakTime = sessionBreak
			}
		}
	}

	// created_at
	if createdAttr, exists := result.Item["created_at"]; exists {
		if s, ok := createdAttr.(*types.AttributeValueMemberS); ok {
			if createdAt, err := time.Parse(time.RFC3339, s.Value); err == nil {
				config.CreatedAt = createdAt
			}
		}
	}

	// updated_at
	if updatedAttr, exists := result.Item["updated_at"]; exists {
		if s, ok := updatedAttr.(*types.AttributeValueMemberS); ok {
			if updatedAt, err := time.Parse(time.RFC3339, s.Value); err == nil {
				config.UpdatedAt = updatedAt
			}
		}
	}

	r.logger.Infof("DynamoDB手動デシリアライズ完了: workTime=%d, breakTime=%d",
		config.RoundWorkTime, config.RoundBreakTime)

	return config, nil
}

// CreateUserConfig は新しいユーザー設定を作成する（新エラーハンドリング対応版）
func (r *UserConfigRepositoryImpl) CreateUserConfig(ctx context.Context, config *model.UserConfig) error {
	r.logger.Infof("CreateUserConfig 入力データ: UserID=%s, WorkTime=%d, BreakTime=%d",
		config.UserID, config.RoundWorkTime, config.RoundBreakTime)

	// 手動でDynamoDBアイテムを作成
	item := map[string]types.AttributeValue{
		"user_id":            &types.AttributeValueMemberS{Value: config.UserID},
		"round_work_time":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", config.RoundWorkTime)},
		"round_break_time":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", config.RoundBreakTime)},
		"session_rounds":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", config.SessionRounds)},
		"session_break_time": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", config.SessionBreakTime)},
		"created_at":         &types.AttributeValueMemberS{Value: config.CreatedAt.Format(time.RFC3339)},
		"updated_at":         &types.AttributeValueMemberS{Value: config.UpdatedAt.Format(time.RFC3339)},
	}

	r.logger.Infof("手動マップ結果 user_id: %s", config.UserID)

	input := &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(user_id)"),
	}

	_, err := r.client.PutItem(ctx, input)
	if err != nil {
		// DynamoDB固有のエラーハンドリング
		var conditionalCheckFailedException *types.ConditionalCheckFailedException
		if errors.As(err, &conditionalCheckFailedException) {
			r.logger.Errorf("ユーザー設定作成条件チェック失敗（既存）: %v", err)
			return appErrors.NewDynamoDBConditionError(err)
		}

		// その他のDynamoDBエラー
		r.logger.Errorf("DynamoDB PutItem エラー: %v", err)
		return appErrors.NewDynamoDBOperationError("create_user_config", err)
	}

	r.logger.Infof("ユーザー設定作成成功: %s", config.UserID)
	return nil
}

// UpdateUserConfig はユーザー設定を更新する（新エラーハンドリング対応版）
func (r *UserConfigRepositoryImpl) UpdateUserConfig(ctx context.Context, config *model.UserConfig) error {
	// 更新時刻を設定
	config.UpdatedAt = time.Now()

	// 手動でDynamoDBアイテムを作成
	item := map[string]types.AttributeValue{
		"user_id":            &types.AttributeValueMemberS{Value: config.UserID},
		"round_work_time":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", config.RoundWorkTime)},
		"round_break_time":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", config.RoundBreakTime)},
		"session_rounds":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", config.SessionRounds)},
		"session_break_time": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", config.SessionBreakTime)},
		"created_at":         &types.AttributeValueMemberS{Value: config.CreatedAt.Format(time.RFC3339)},
		"updated_at":         &types.AttributeValueMemberS{Value: config.UpdatedAt.Format(time.RFC3339)},
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
		// アイテムが存在する場合のみ更新を許可
		ConditionExpression: aws.String("attribute_exists(user_id)"),
	}

	_, err := r.client.PutItem(ctx, input)
	if err != nil {
		// DynamoDB固有のエラーハンドリング
		var conditionalCheckFailedException *types.ConditionalCheckFailedException
		if errors.As(err, &conditionalCheckFailedException) {
			r.logger.Errorf("ユーザー設定更新条件チェック失敗（存在しない）: %v", err)
			return appErrors.ErrDynamoDBItemNotFound // Infrastructure Error（存在しない）
		}

		// その他のDynamoDBエラー
		r.logger.Errorf("DynamoDB PutItem エラー: %v", err)
		return appErrors.NewDynamoDBOperationError("update_user_config", err)
	}

	r.logger.Infof("ユーザー設定更新成功: %s", config.UserID)
	return nil
}

// DeleteUserConfig はユーザー設定を削除する（新エラーハンドリング対応版）
func (r *UserConfigRepositoryImpl) DeleteUserConfig(ctx context.Context, userID uuid.UUID) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: userID.String()},
		},
		// アイテムが存在する場合のみ削除を許可
		ConditionExpression: aws.String("attribute_exists(user_id)"),
	}

	_, err := r.client.DeleteItem(ctx, input)
	if err != nil {
		// DynamoDB固有のエラーハンドリング
		var conditionalCheckFailedException *types.ConditionalCheckFailedException
		if errors.As(err, &conditionalCheckFailedException) {
			r.logger.Errorf("ユーザー設定削除条件チェック失敗（存在しない）: %v", err)
			return appErrors.ErrDynamoDBItemNotFound // Infrastructure Error（存在しない）
		}

		// その他のDynamoDBエラー
		r.logger.Errorf("DynamoDB DeleteItem エラー: %v", err)
		return appErrors.NewDynamoDBOperationError("delete_user_config", err)
	}

	r.logger.Infof("ユーザー設定削除成功: %s", userID.String())
	return nil
}
