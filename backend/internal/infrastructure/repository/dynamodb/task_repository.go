package dynamodb

import (
	"context"
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

// TaskRepositoryImpl はDynamoDBを使用したTaskRepositoryの実装
type TaskRepositoryImpl struct {
	client    *dynamodb.Client
	tableName string
	logger    logger.Logger
}

// NewTaskRepository は新しいTaskRepositoryImplインスタンスを作成する
func NewTaskRepository(client *dynamodb.Client, cfg *config.Config, logger logger.Logger) repository.TaskRepository {
	return &TaskRepositoryImpl{
		client:    client,
		tableName: cfg.DynamoUnifiedTable, // 統一テーブルを使用
		logger:    logger,
	}
}

// Create はタスクを作成する
func (r *TaskRepositoryImpl) Create(ctx context.Context, task *model.Task) error {
	pk := UserPartitionKey(task.UserID.String())
	sk := TaskSortKey(task.ID.String())

	item := map[string]types.AttributeValue{
		"PK":           &types.AttributeValueMemberS{Value: pk},
		"SK":           &types.AttributeValueMemberS{Value: sk},
		"user_id":      &types.AttributeValueMemberS{Value: task.UserID.String()},
		"task_id":      &types.AttributeValueMemberS{Value: task.ID.String()},
		"detail":       &types.AttributeValueMemberS{Value: task.Detail},
		"is_completed": &types.AttributeValueMemberBOOL{Value: task.IsCompleted},
		"created_at":   &types.AttributeValueMemberS{Value: task.CreatedAt.Format(time.RFC3339)},
		"updated_at":   &types.AttributeValueMemberS{Value: task.UpdatedAt.Format(time.RFC3339)},
	}

	input := &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
	}

	_, err := r.client.PutItem(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB PutItem エラー: %v", err)
		return appErrors.NewDynamoDBOperationError("create_task", err)
	}

	r.logger.Infof("タスク作成成功: ID=%s, Detail=%s", task.ID.String(), task.Detail)
	return nil
}

// GetByID はIDによってタスクを取得する
func (r *TaskRepositoryImpl) GetByID(ctx context.Context, id, userID uuid.UUID) (*model.Task, error) {
	pk := UserPartitionKey(userID.String())
	sk := TaskSortKey(id.String())

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
		return nil, appErrors.NewDynamoDBOperationError("get_task_by_id", err)
	}

	if result.Item == nil {
		r.logger.Debugf("タスクが見つかりません: ID=%s, UserID=%s", id.String(), userID.String())
		return nil, appErrors.ErrRecordNotFound
	}

	task, err := r.itemToTask(result.Item)
	if err != nil {
		r.logger.Errorf("タスク変換エラー: %v", err)
		return nil, appErrors.NewDynamoDBOperationError("conversion", err)
	}

	r.logger.Debugf("タスク取得成功: ID=%s", id.String())
	return task, nil
}

// GetAllByUserID はユーザーIDに紐づくすべてのタスクを取得する
func (r *TaskRepositoryImpl) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Task, error) {
	pk := UserPartitionKey(userID.String())

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk_prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":        &types.AttributeValueMemberS{Value: pk},
			":sk_prefix": &types.AttributeValueMemberS{Value: TaskQueryPrefix()},
		},
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB Query エラー (GetAllByUserID): %v", err)
		return nil, appErrors.NewDynamoDBOperationError("get_tasks_by_user_id", err)
	}

	tasks := make([]*model.Task, 0, len(result.Items))
	for _, item := range result.Items {
		task, err := r.itemToTask(item)
		if err != nil {
			r.logger.Warnf("タスク変換エラー（スキップ）: %v", err)
			continue
		}
		tasks = append(tasks, task)
	}

	r.logger.Infof("ユーザータスク取得成功: UserID=%s, count=%d", userID.String(), len(tasks))
	return tasks, nil
}

// Update はタスクの詳細を更新する
func (r *TaskRepositoryImpl) Update(ctx context.Context, task *model.Task) error {
	pk := UserPartitionKey(task.UserID.String())
	sk := TaskSortKey(task.ID.String())

	task.UpdatedAt = time.Now()

	item := map[string]types.AttributeValue{
		"PK":           &types.AttributeValueMemberS{Value: pk},
		"SK":           &types.AttributeValueMemberS{Value: sk},
		"user_id":      &types.AttributeValueMemberS{Value: task.UserID.String()},
		"task_id":      &types.AttributeValueMemberS{Value: task.ID.String()},
		"detail":       &types.AttributeValueMemberS{Value: task.Detail},
		"is_completed": &types.AttributeValueMemberBOOL{Value: task.IsCompleted},
		"created_at":   &types.AttributeValueMemberS{Value: task.CreatedAt.Format(time.RFC3339)},
		"updated_at":   &types.AttributeValueMemberS{Value: task.UpdatedAt.Format(time.RFC3339)},
	}

	input := &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_exists(PK) AND attribute_exists(SK)"),
	}

	_, err := r.client.PutItem(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB PutItem エラー: %v", err)
		return appErrors.NewDynamoDBOperationError("update_task", err)
	}

	r.logger.Infof("タスク更新成功: ID=%s", task.ID.String())
	return nil
}

// ToggleCompletion はタスクの完了状態を切り替える
func (r *TaskRepositoryImpl) ToggleCompletion(ctx context.Context, id, userID uuid.UUID) error {
	// 現在のタスクを取得
	task, err := r.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	// 完了状態を切り替え
	task.ToggleCompletion()

	// 更新処理
	err = r.Update(ctx, task)
	if err != nil {
		return err
	}

	r.logger.Infof("タスク完了状態切り替え成功: ID=%s, IsCompleted=%v", id.String(), task.IsCompleted)
	return nil
}

// Delete はタスクを削除する
func (r *TaskRepositoryImpl) Delete(ctx context.Context, id, userID uuid.UUID) error {
	pk := UserPartitionKey(userID.String())
	sk := TaskSortKey(id.String())

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		ConditionExpression: aws.String("attribute_exists(PK) AND attribute_exists(SK)"),
	}

	_, err := r.client.DeleteItem(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB DeleteItem エラー: %v", err)
		return appErrors.NewDynamoDBOperationError("delete_task", err)
	}

	r.logger.Infof("タスク削除成功: ID=%s", id.String())
	return nil
}

// Helper methods

// itemToTask はDynamoDBアイテムをTaskモデルに変換する
func (r *TaskRepositoryImpl) itemToTask(item map[string]types.AttributeValue) (*model.Task, error) {
	task := &model.Task{}

	// user_id
	if userIDAttr, exists := item["user_id"]; exists {
		if s, ok := userIDAttr.(*types.AttributeValueMemberS); ok {
			if id, err := uuid.Parse(s.Value); err == nil {
				task.UserID = id
			}
		}
	}

	// task_id
	if taskIDAttr, exists := item["task_id"]; exists {
		if s, ok := taskIDAttr.(*types.AttributeValueMemberS); ok {
			if id, err := uuid.Parse(s.Value); err == nil {
				task.ID = id
			}
		}
	}

	// detail
	if detailAttr, exists := item["detail"]; exists {
		if s, ok := detailAttr.(*types.AttributeValueMemberS); ok {
			task.Detail = s.Value
		}
	}

	// is_completed
	if isCompletedAttr, exists := item["is_completed"]; exists {
		if b, ok := isCompletedAttr.(*types.AttributeValueMemberBOOL); ok {
			task.IsCompleted = b.Value
		}
	}

	// created_at
	if createdAttr, exists := item["created_at"]; exists {
		if s, ok := createdAttr.(*types.AttributeValueMemberS); ok {
			if createdAt, err := time.Parse(time.RFC3339, s.Value); err == nil {
				task.CreatedAt = createdAt
			}
		}
	}

	// updated_at
	if updatedAttr, exists := item["updated_at"]; exists {
		if s, ok := updatedAttr.(*types.AttributeValueMemberS); ok {
			if updatedAt, err := time.Parse(time.RFC3339, s.Value); err == nil {
				task.UpdatedAt = updatedAt
			}
		}
	}

	return task, nil
}