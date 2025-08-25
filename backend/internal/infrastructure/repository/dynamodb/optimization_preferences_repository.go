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
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	roundVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/round"
	sessionVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/session"
	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// OptimizationPreferencesRepositoryImpl はDynamoDBを使用したOptimizationPreferencesRepositoryの実装
type OptimizationPreferencesRepositoryImpl struct {
	client    *dynamodb.Client
	tableName string
	logger    logger.Logger
}

// NewOptimizationPreferencesRepository は新しいOptimizationPreferencesRepositoryImplインスタンスを作成する
func NewOptimizationPreferencesRepository(client *dynamodb.Client, cfg *config.Config, logger logger.Logger) repository.OptimizationPreferencesRepository {
	return &OptimizationPreferencesRepositoryImpl{
		client:    client,
		tableName: cfg.DynamoUnifiedTable,
		logger:    logger,
	}
}

// Get は最適化設定を取得する
func (r *OptimizationPreferencesRepositoryImpl) Get(ctx context.Context, userID userVO.UserID) (*entity.OptimizationPreferences, error) {
	pk := UserPartitionKey(userID.String())
	sk := OptimizationPreferencesSortKey()

	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB GetItem エラー (Get): %v", err)
		return nil, appErrors.NewDynamoDBOperationError("get_optimization_preferences", err)
	}

	if result.Item == nil {
		r.logger.Debugf("最適化設定が見つかりません: %s", userID.String())
		return nil, appErrors.ErrRecordNotFound
	}

	preferences, err := r.itemToOptimizationPreferences(result.Item)
	if err != nil {
		r.logger.Errorf("最適化設定変換エラー: %v", err)
		return nil, appErrors.NewDynamoDBOperationError("conversion", err)
	}

	r.logger.Debugf("最適化設定取得成功: UserID=%s", userID.String())
	return preferences, nil
}

// Create は新しい最適化設定を作成する
func (r *OptimizationPreferencesRepositoryImpl) Create(ctx context.Context, preferences *entity.OptimizationPreferences) error {
	pk := UserPartitionKey(preferences.UserID.String())
	sk := OptimizationPreferencesSortKey()

	item := map[string]types.AttributeValue{
		"PK":                 &types.AttributeValueMemberS{Value: pk},
		"SK":                 &types.AttributeValueMemberS{Value: sk},
		"user_id":            &types.AttributeValueMemberS{Value: preferences.UserID.String()},
		"round_work_time":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", preferences.RoundWorkTime.Minutes())},
		"round_break_time":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", preferences.RoundBreakTime.Minutes())},
		"session_rounds":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", preferences.SessionRounds.Count())},
		"session_break_time": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", preferences.SessionBreakTime.Minutes())},
		"entity_type":        &types.AttributeValueMemberS{Value: "optimization_preferences"},
		"created_at":         &types.AttributeValueMemberS{Value: preferences.CreatedAt.Format(time.RFC3339)},
		"updated_at":         &types.AttributeValueMemberS{Value: preferences.UpdatedAt.Format(time.RFC3339)},
	}

	input := &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	}

	_, err := r.client.PutItem(ctx, input)
	if err != nil {
		var conditionalCheckFailedException *types.ConditionalCheckFailedException
		if errors.As(err, &conditionalCheckFailedException) {
			r.logger.Errorf("最適化設定作成条件チェック失敗（既存）: %v", err)
			return appErrors.NewBadRequestError("最適化設定が既に存在します")
		}

		r.logger.Errorf("DynamoDB PutItem エラー (Create): %v", err)
		return appErrors.NewDynamoDBOperationError("create_optimization_preferences", err)
	}

	r.logger.Infof("最適化設定作成成功: UserID=%s", preferences.UserID.String())
	return nil
}

// Update は最適化設定を更新する
func (r *OptimizationPreferencesRepositoryImpl) Update(ctx context.Context, preferences *entity.OptimizationPreferences) error {
	pk := UserPartitionKey(preferences.UserID.String())
	sk := OptimizationPreferencesSortKey()

	preferences.UpdatedAt = time.Now()

	item := map[string]types.AttributeValue{
		"PK":                 &types.AttributeValueMemberS{Value: pk},
		"SK":                 &types.AttributeValueMemberS{Value: sk},
		"user_id":            &types.AttributeValueMemberS{Value: preferences.UserID.String()},
		"round_work_time":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", preferences.RoundWorkTime.Minutes())},
		"round_break_time":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", preferences.RoundBreakTime.Minutes())},
		"session_rounds":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", preferences.SessionRounds.Count())},
		"session_break_time": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", preferences.SessionBreakTime.Minutes())},
		"entity_type":        &types.AttributeValueMemberS{Value: "optimization_preferences"},
		"created_at":         &types.AttributeValueMemberS{Value: preferences.CreatedAt.Format(time.RFC3339)},
		"updated_at":         &types.AttributeValueMemberS{Value: preferences.UpdatedAt.Format(time.RFC3339)},
	}

	input := &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_exists(PK)"),
	}

	_, err := r.client.PutItem(ctx, input)
	if err != nil {
		var conditionalCheckFailedException *types.ConditionalCheckFailedException
		if errors.As(err, &conditionalCheckFailedException) {
			r.logger.Errorf("最適化設定更新条件チェック失敗（存在しない）: %v", err)
			return appErrors.ErrRecordNotFound
		}

		r.logger.Errorf("DynamoDB PutItem エラー (Update): %v", err)
		return appErrors.NewDynamoDBOperationError("update_optimization_preferences", err)
	}

	r.logger.Infof("最適化設定更新成功: UserID=%s", preferences.UserID.String())
	return nil
}

// Delete は最適化設定を削除する
func (r *OptimizationPreferencesRepositoryImpl) Delete(ctx context.Context, userID userVO.UserID) error {
	pk := UserPartitionKey(userID.String())
	sk := OptimizationPreferencesSortKey()

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		ConditionExpression: aws.String("attribute_exists(PK)"),
	}

	_, err := r.client.DeleteItem(ctx, input)
	if err != nil {
		var conditionalCheckFailedException *types.ConditionalCheckFailedException
		if errors.As(err, &conditionalCheckFailedException) {
			r.logger.Errorf("最適化設定削除条件チェック失敗（存在しない）: %v", err)
			return appErrors.ErrRecordNotFound
		}

		r.logger.Errorf("DynamoDB DeleteItem エラー (Delete): %v", err)
		return appErrors.NewDynamoDBOperationError("delete_optimization_preferences", err)
	}

	r.logger.Infof("最適化設定削除成功: UserID=%s", userID.String())
	return nil
}

// GetOrCreateDefault は設定を取得し、存在しない場合はデフォルト設定を作成して返す
func (r *OptimizationPreferencesRepositoryImpl) GetOrCreateDefault(ctx context.Context, userID userVO.UserID) (*entity.OptimizationPreferences, error) {
	// まず設定を取得を試行
	preferences, err := r.Get(ctx, userID)
	if err == nil {
		return preferences, nil
	}

	// 設定が存在しない場合、デフォルト設定を作成
	if errors.Is(err, appErrors.ErrRecordNotFound) {
		r.logger.Infof("最適化設定が見つからないため、デフォルト設定を作成します: UserID=%s", userID.String())

		// デフォルト設定を作成
		defaultPreferences, err := entity.NewOptimizationPreferences(userID)
		if err != nil {
			return nil, fmt.Errorf("デフォルト最適化設定の作成に失敗: %w", err)
		}

		// データベースに保存
		createErr := r.Create(ctx, defaultPreferences)
		if createErr != nil {
			// 既存エラーの場合は取得を再試行
			if errors.Is(createErr, appErrors.ErrUniqueConstraint) {
				r.logger.Debugf("最適化設定作成時の競合状態、再取得を試行: UserID=%s", userID.String())
				return r.Get(ctx, userID)
			}
			return nil, createErr
		}

		return defaultPreferences, nil
	}

	// その他のエラー
	return nil, err
}

// Helper methods

// itemToOptimizationPreferences はDynamoDBアイテムをOptimizationPreferencesモデルに変換する
func (r *OptimizationPreferencesRepositoryImpl) itemToOptimizationPreferences(item map[string]types.AttributeValue) (*entity.OptimizationPreferences, error) {
	preferences := &entity.OptimizationPreferences{}

	// user_id
	if userIDAttr, exists := item["user_id"]; exists {
		if s, ok := userIDAttr.(*types.AttributeValueMemberS); ok {
			if userID, err := userVO.NewUserIDFromString(s.Value); err == nil {
				preferences.UserID = userID
			}
		}
	}

	// round_work_time
	if workTimeAttr, exists := item["round_work_time"]; exists {
		if n, ok := workTimeAttr.(*types.AttributeValueMemberN); ok {
			if workTime, err := strconv.Atoi(n.Value); err == nil {
				if wt, err := roundVO.NewWorkTime(workTime); err == nil {
					preferences.RoundWorkTime = wt
				}
			}
		}
	}

	// round_break_time
	if breakTimeAttr, exists := item["round_break_time"]; exists {
		if n, ok := breakTimeAttr.(*types.AttributeValueMemberN); ok {
			if breakTime, err := strconv.Atoi(n.Value); err == nil {
				if bt, err := roundVO.NewBreakTime(breakTime); err == nil {
					preferences.RoundBreakTime = bt
				}
			}
		}
	}

	// session_rounds
	if roundsAttr, exists := item["session_rounds"]; exists {
		if n, ok := roundsAttr.(*types.AttributeValueMemberN); ok {
			if rounds, err := strconv.Atoi(n.Value); err == nil {
				if sr, err := sessionVO.NewSessionRounds(rounds); err == nil {
					preferences.SessionRounds = sr
				}
			}
		}
	}

	// session_break_time
	if sessionBreakTimeAttr, exists := item["session_break_time"]; exists {
		if n, ok := sessionBreakTimeAttr.(*types.AttributeValueMemberN); ok {
			if sessionBreakTime, err := strconv.Atoi(n.Value); err == nil {
				if sbt, err := sessionVO.NewBreakTime(sessionBreakTime); err == nil {
					preferences.SessionBreakTime = sbt
				}
			}
		}
	}

	// created_at
	if createdAttr, exists := item["created_at"]; exists {
		if s, ok := createdAttr.(*types.AttributeValueMemberS); ok {
			if createdAt, err := time.Parse(time.RFC3339, s.Value); err == nil {
				preferences.CreatedAt = createdAt
			}
		}
	}

	// updated_at
	if updatedAttr, exists := item["updated_at"]; exists {
		if s, ok := updatedAttr.(*types.AttributeValueMemberS); ok {
			if updatedAt, err := time.Parse(time.RFC3339, s.Value); err == nil {
				preferences.UpdatedAt = updatedAt
			}
		}
	}

	return preferences, nil
}
