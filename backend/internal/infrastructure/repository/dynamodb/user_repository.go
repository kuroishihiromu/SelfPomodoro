package dynamodb

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// UserRepositoryImpl はDynamoDBを使用したUserRepositoryの実装
type UserRepositoryImpl struct {
	client    *dynamodb.Client
	tableName string
	logger    logger.Logger
}

// NewUserRepository は新しいUserRepositoryImplインスタンスを作成する
func NewUserRepository(client *dynamodb.Client, cfg *config.Config, logger logger.Logger) repository.UserRepository {
	return &UserRepositoryImpl{
		client:    client,
		tableName: cfg.DynamoUnifiedTable, // 統一テーブルを使用
		logger:    logger,
	}
}

// GetByID はIDによってユーザーを取得する
func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	pk := UserPartitionKey(id.String())
	sk := "PROFILE"

	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB GetItem エラー (GetByID): %v", err)
		return nil, appErrors.NewDynamoDBOperationError("get_user_by_id", err)
	}

	if result.Item == nil {
		r.logger.Debugf("ユーザーが見つかりません: %s", id.String())
		return nil, appErrors.ErrRecordNotFound
	}

	user, err := r.itemToUser(result.Item)
	if err != nil {
		r.logger.Errorf("ユーザー変換エラー: %v", err)
		return nil, appErrors.NewDynamoDBOperationError("conversion", err)
	}

	r.logger.Debugf("ユーザー取得成功: ID=%s", id.String())
	return user, nil
}

// GetByEmail はメールアドレスによってユーザーを取得する
func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	// GSI1を使用してメールアドレス検索
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("email = :email"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
		Limit: aws.Int32(1),
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB Query エラー (GetByEmail): %v", err)
		return nil, appErrors.NewDynamoDBOperationError("get_user_by_email", err)
	}

	if len(result.Items) == 0 {
		r.logger.Debugf("メールアドレスでユーザーが見つかりません: %s", email)
		return nil, appErrors.ErrRecordNotFound
	}

	user, err := r.itemToUser(result.Items[0])
	if err != nil {
		r.logger.Errorf("ユーザー変換エラー: %v", err)
		return nil, appErrors.NewDynamoDBOperationError("conversion", err)
	}

	r.logger.Debugf("メールアドレスでユーザー取得成功: %s", email)
	return user, nil
}

// Create は新しいユーザーを作成する
func (r *UserRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	if !user.IsValidForCreation() {
		return appErrors.NewBadRequestError("ユーザー作成に必要な情報が不足しています")
	}

	pk := UserPartitionKey(user.ID.String())
	sk := "PROFILE"

	item := map[string]types.AttributeValue{
		"PK":         &types.AttributeValueMemberS{Value: pk},
		"SK":         &types.AttributeValueMemberS{Value: sk},
		"user_id":    &types.AttributeValueMemberS{Value: user.ID.String()},
		"name":       &types.AttributeValueMemberS{Value: user.Name},
		"email":      &types.AttributeValueMemberS{Value: user.Email},
		"provider":   &types.AttributeValueMemberS{Value: user.Provider},
		"created_at": &types.AttributeValueMemberS{Value: user.CreatedAt.Format(time.RFC3339)},
		"updated_at": &types.AttributeValueMemberS{Value: user.UpdatedAt.Format(time.RFC3339)},
	}

	// provider_idは任意フィールド
	if user.ProviderID != nil {
		item["provider_id"] = &types.AttributeValueMemberS{Value: *user.ProviderID}
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
			r.logger.Errorf("ユーザー作成条件チェック失敗（既存）: %v", err)
			return appErrors.NewBadRequestError("ユーザーが既に存在します")
		}

		r.logger.Errorf("DynamoDB PutItem エラー: %v", err)
		return appErrors.NewDynamoDBOperationError("create_user", err)
	}

	r.logger.Infof("ユーザー作成成功: ID=%s, Name=%s", user.ID.String(), user.Name)
	return nil
}

// Update はユーザー情報を更新する
func (r *UserRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	pk := UserPartitionKey(user.ID.String())
	sk := "PROFILE"

	user.UpdatedAt = time.Now()

	item := map[string]types.AttributeValue{
		"PK":         &types.AttributeValueMemberS{Value: pk},
		"SK":         &types.AttributeValueMemberS{Value: sk},
		"user_id":    &types.AttributeValueMemberS{Value: user.ID.String()},
		"name":       &types.AttributeValueMemberS{Value: user.Name},
		"email":      &types.AttributeValueMemberS{Value: user.Email},
		"provider":   &types.AttributeValueMemberS{Value: user.Provider},
		"created_at": &types.AttributeValueMemberS{Value: user.CreatedAt.Format(time.RFC3339)},
		"updated_at": &types.AttributeValueMemberS{Value: user.UpdatedAt.Format(time.RFC3339)},
	}

	if user.ProviderID != nil {
		item["provider_id"] = &types.AttributeValueMemberS{Value: *user.ProviderID}
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
			r.logger.Errorf("ユーザー更新条件チェック失敗（存在しない）: %v", err)
			return appErrors.ErrRecordNotFound
		}

		r.logger.Errorf("DynamoDB PutItem エラー: %v", err)
		return appErrors.NewDynamoDBOperationError("update_user", err)
	}

	r.logger.Infof("ユーザー更新成功: ID=%s", user.ID.String())
	return nil
}

// Delete はユーザーを削除する
func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	pk := UserPartitionKey(id.String())
	sk := "PROFILE"

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
			r.logger.Errorf("ユーザー削除条件チェック失敗（存在しない）: %v", err)
			return appErrors.ErrRecordNotFound
		}

		r.logger.Errorf("DynamoDB DeleteItem エラー: %v", err)
		return appErrors.NewDynamoDBOperationError("delete_user", err)
	}

	r.logger.Infof("ユーザー削除成功: ID=%s", id.String())
	return nil
}

// UpdateProfile はユーザープロフィールを更新する
func (r *UserRepositoryImpl) UpdateProfile(ctx context.Context, id uuid.UUID, name, email string) (*entity.User, error) {
	// 現在のユーザー情報を取得
	user, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// プロフィール更新
	user.UpdateProfile(name, email)

	// データベースに保存
	err = r.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	r.logger.Infof("ユーザープロフィール更新成功: ID=%s", id.String())
	return user, nil
}

// ExistsByID はユーザーの存在確認を行う
func (r *UserRepositoryImpl) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	pk := UserPartitionKey(id.String())
	sk := "PROFILE"

	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		ProjectionExpression: aws.String("PK"),
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB GetItem エラー (ExistsByID): %v", err)
		return false, appErrors.NewDynamoDBOperationError("exists_user_by_id", err)
	}

	exists := result.Item != nil
	r.logger.Debugf("ユーザー存在確認: ID=%s, exists=%v", id.String(), exists)
	return exists, nil
}

// GetUsersByProvider はプロバイダー別にユーザーを取得する
func (r *UserRepositoryImpl) GetUsersByProvider(ctx context.Context, provider string, limit, offset int) ([]*entity.User, error) {
	// GSI2を使用してプロバイダー検索
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("GSI2"),
		KeyConditionExpression: aws.String("provider = :provider"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":provider": &types.AttributeValueMemberS{Value: provider},
		},
		Limit: aws.Int32(int32(limit)),
	}

	// オフセット処理（簡易実装）
	// 実際のプロダクションではLastEvaluatedKeyを使用した方が効率的
	if offset > 0 {
		// スキャンでオフセット分読み飛ばし（非効率だが簡易実装）
		input.Limit = aws.Int32(int32(offset + limit))
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB Query エラー (GetUsersByProvider): %v", err)
		return nil, appErrors.NewDynamoDBOperationError("get_users_by_provider", err)
	}

	users := make([]*entity.User, 0, len(result.Items))
	startIndex := offset
	if startIndex > len(result.Items) {
		startIndex = len(result.Items)
	}

	for i := startIndex; i < len(result.Items) && len(users) < limit; i++ {
		user, err := r.itemToUser(result.Items[i])
		if err != nil {
			r.logger.Warnf("ユーザー変換エラー（スキップ）: %v", err)
			continue
		}
		users = append(users, user)
	}

	r.logger.Infof("プロバイダー別ユーザー取得成功: provider=%s, count=%d", provider, len(users))
	return users, nil
}

// Helper methods

// itemToUser はDynamoDBアイテムをUserモデルに変換する
func (r *UserRepositoryImpl) itemToUser(item map[string]types.AttributeValue) (*entity.User, error) {
	user := &entity.User{}

	// user_id
	if userIDAttr, exists := item["user_id"]; exists {
		if s, ok := userIDAttr.(*types.AttributeValueMemberS); ok {
			if id, err := uuid.Parse(s.Value); err == nil {
				user.ID = id
			}
		}
	}

	// name
	if nameAttr, exists := item["name"]; exists {
		if s, ok := nameAttr.(*types.AttributeValueMemberS); ok {
			user.Name = s.Value
		}
	}

	// email
	if emailAttr, exists := item["email"]; exists {
		if s, ok := emailAttr.(*types.AttributeValueMemberS); ok {
			user.Email = s.Value
		}
	}

	// provider
	if providerAttr, exists := item["provider"]; exists {
		if s, ok := providerAttr.(*types.AttributeValueMemberS); ok {
			user.Provider = s.Value
		}
	}

	// provider_id (optional)
	if providerIDAttr, exists := item["provider_id"]; exists {
		if s, ok := providerIDAttr.(*types.AttributeValueMemberS); ok {
			user.ProviderID = &s.Value
		}
	}

	// created_at
	if createdAttr, exists := item["created_at"]; exists {
		if s, ok := createdAttr.(*types.AttributeValueMemberS); ok {
			if createdAt, err := time.Parse(time.RFC3339, s.Value); err == nil {
				user.CreatedAt = createdAt
			}
		}
	}

	// updated_at
	if updatedAttr, exists := item["updated_at"]; exists {
		if s, ok := updatedAttr.(*types.AttributeValueMemberS); ok {
			if updatedAt, err := time.Parse(time.RFC3339, s.Value); err == nil {
				user.UpdatedAt = updatedAt
			}
		}
	}

	return user, nil
}

// DeleteAllUserData はユーザーに関連するすべてのデータを削除する（統合テーブル対応）
func (r *UserRepositoryImpl) DeleteAllUserData(ctx context.Context, userID uuid.UUID) error {
	pk := UserPartitionKey(userID.String())

	r.logger.Infof("ユーザーデータ包括削除開始: UserID=%s", userID.String()[:8]+"...")

	// 1. 該当PKのすべてのアイテムを取得
	items, err := r.queryAllUserItems(ctx, pk)
	if err != nil {
		r.logger.Errorf("ユーザーアイテム取得エラー: %v", err)
		return appErrors.NewDynamoDBOperationError("query_user_items", err)
	}

	if len(items) == 0 {
		r.logger.Debugf("削除対象アイテムなし: UserID=%s", userID.String()[:8]+"...")
		return nil
	}

	// 2. バッチ削除実行
	err = r.batchDeleteItems(ctx, items)
	if err != nil {
		r.logger.Errorf("バッチ削除エラー: %v", err)
		return appErrors.NewDynamoDBOperationError("batch_delete_user_data", err)
	}

	r.logger.Infof("ユーザーデータ包括削除完了: UserID=%s, 削除件数=%d",
		userID.String()[:8]+"...", len(items))

	return nil
}

// queryAllUserItems は指定PKのすべてのアイテムを取得する
func (r *UserRepositoryImpl) queryAllUserItems(ctx context.Context, pk string) ([]map[string]types.AttributeValue, error) {
	var allItems []map[string]types.AttributeValue
	var lastEvaluatedKey map[string]types.AttributeValue

	for {
		input := &dynamodb.QueryInput{
			TableName:              aws.String(r.tableName),
			KeyConditionExpression: aws.String("PK = :pk"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: pk},
			},
			ProjectionExpression: aws.String("PK, SK"), // キーのみ取得
		}

		if lastEvaluatedKey != nil {
			input.ExclusiveStartKey = lastEvaluatedKey
		}

		result, err := r.client.Query(ctx, input)
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, result.Items...)

		// ページネーション処理
		if result.LastEvaluatedKey == nil {
			break
		}
		lastEvaluatedKey = result.LastEvaluatedKey
	}

	r.logger.Debugf("ユーザーアイテム取得完了: PK=%s, 件数=%d", pk, len(allItems))
	return allItems, nil
}

// batchDeleteItems はアイテムをバッチ削除する
func (r *UserRepositoryImpl) batchDeleteItems(ctx context.Context, items []map[string]types.AttributeValue) error {
	const batchSize = 25 // DynamoDB BatchWriteItemの上限

	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}

		batch := items[i:end]
		err := r.batchDeleteBatch(ctx, batch)
		if err != nil {
			r.logger.Errorf("バッチ削除失敗: batch %d-%d, error: %v", i, end-1, err)
			return err
		}

		r.logger.Debugf("バッチ削除成功: %d-%d件", i+1, end)
	}

	return nil
}

// batchDeleteBatch は単一バッチの削除を実行する
func (r *UserRepositoryImpl) batchDeleteBatch(ctx context.Context, items []map[string]types.AttributeValue) error {
	writeRequests := make([]types.WriteRequest, 0, len(items))

	for _, item := range items {
		// PK, SKでキーを構成
		key := map[string]types.AttributeValue{
			"PK": item["PK"],
			"SK": item["SK"],
		}

		writeRequests = append(writeRequests, types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{
				Key: key,
			},
		})
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			r.tableName: writeRequests,
		},
	}

	_, err := r.client.BatchWriteItem(ctx, input)
	return err
}

// CountUserItems はユーザーに関連するアイテム数を取得する（テスト用）
func (r *UserRepositoryImpl) CountUserItems(ctx context.Context, userID uuid.UUID) (int, error) {
	pk := UserPartitionKey(userID.String())

	items, err := r.queryAllUserItems(ctx, pk)
	if err != nil {
		r.logger.Errorf("ユーザーアイテム数取得エラー: %v", err)
		return 0, appErrors.NewDynamoDBOperationError("count_user_items", err)
	}

	return len(items), nil
}
