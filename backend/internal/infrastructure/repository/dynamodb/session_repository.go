package dynamodb

import (
	"context"
	"fmt"
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

// SessionRepositoryImpl はDynamoDBを使用したSessionRepositoryの実装
type SessionRepositoryImpl struct {
	client    *dynamodb.Client
	tableName string
	logger    logger.Logger
}

// NewSessionRepository は新しいSessionRepositoryImplインスタンスを作成する
func NewSessionRepository(client *dynamodb.Client, cfg *config.Config, logger logger.Logger) repository.SessionRepository {
	return &SessionRepositoryImpl{
		client:    client,
		tableName: cfg.DynamoUnifiedTable, // 統一テーブルを使用
		logger:    logger,
	}
}

// Create はセッションを作成する
func (r *SessionRepositoryImpl) Create(ctx context.Context, session *entity.Session) error {
	date := session.StartTime.Format("2006-01-02")
	pk := UserPartitionKey(session.UserID.String())
	sk := SessionSortKey(date, session.ID.String())

	// TTL設定: 30日後に自動削除（完了時のみ作成されるため）
	ttl := session.CreatedAt.Add(30 * 24 * time.Hour).Unix()

	item := map[string]types.AttributeValue{
		"PK":         &types.AttributeValueMemberS{Value: pk},
		"SK":         &types.AttributeValueMemberS{Value: sk},
		"user_id":    &types.AttributeValueMemberS{Value: session.UserID.String()},
		"session_id": &types.AttributeValueMemberS{Value: session.ID.String()},
		"date":       &types.AttributeValueMemberS{Value: date},
		"start_time": &types.AttributeValueMemberS{Value: session.StartTime.Format(time.RFC3339)},
		"ttl":        &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", ttl)},
		"created_at": &types.AttributeValueMemberS{Value: session.CreatedAt.Format(time.RFC3339)},
		"updated_at": &types.AttributeValueMemberS{Value: session.UpdatedAt.Format(time.RFC3339)},
	}

	// Optional fields
	if session.EndTime != nil {
		item["end_time"] = &types.AttributeValueMemberS{Value: session.EndTime.Format(time.RFC3339)}
	}
	if session.AverageFocus != nil {
		item["average_focus"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", *session.AverageFocus)}
	}
	if session.TotalWorkMin != nil {
		item["total_work_min"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", *session.TotalWorkMin)}
	}
	if session.RoundCount != nil {
		item["round_count"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", *session.RoundCount)}
	}
	if session.BreakTime != nil {
		item["break_time"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", *session.BreakTime)}
	}

	input := &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
	}

	_, err := r.client.PutItem(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB PutItem エラー: %v", err)
		return appErrors.NewDynamoDBOperationError("create_session", err)
	}

	r.logger.Infof("セッション作成成功: ID=%s", session.ID.String())
	return nil
}

// GetByID はIDによってセッションを取得する（GSI最適化版）
func (r *SessionRepositoryImpl) GetByID(ctx context.Context, id, userID uuid.UUID) (*entity.Session, error) {
	r.logger.Infof("セッション取得開始: SessionID=%s, UserID=%s", id.String(), userID.String())

	// 🎯 SessionIdIndex GSIを使用して効率的に検索
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("SessionIdIndex"), // GSI使用
		KeyConditionExpression: aws.String("session_id = :session_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":session_id": &types.AttributeValueMemberS{Value: id.String()},
		},
		Limit: aws.Int32(1), // 最初の1件のみ
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		r.logger.Errorf("SessionIdIndex GSI Query エラー: %v", err)
		return nil, appErrors.NewDynamoDBOperationError("get_session_by_id_gsi", err)
	}

	r.logger.Infof("SessionIdIndex GSI Query結果: %d件のレコードが見つかりました", len(result.Items))

	if len(result.Items) == 0 {
		r.logger.Debugf("GSI: セッションが見つかりません: ID=%s", id.String())
		return nil, appErrors.ErrRecordNotFound
	}

	// GSI結果からセッションを変換
	session, err := r.itemToSession(result.Items[0])
	if err != nil {
		r.logger.Errorf("GSI結果のセッション変換エラー: %v", err)
		return nil, appErrors.NewDynamoDBOperationError("session_conversion", err)
	}

	// ユーザーIDの整合性チェック
	if session.UserID != userID {
		r.logger.Warnf("セッションのユーザーID不一致: SessionID=%s, Expected=%s, Actual=%s",
			id.String(), userID.String(), session.UserID.String())
		return nil, appErrors.ErrRecordNotFound
	}

	r.logger.Infof("GSI経由でセッション取得成功: ID=%s, UserID=%s", id.String(), userID.String())
	return session, nil
}

// GetAllByUserID はユーザーIDに紐づくすべてのセッションを取得する
func (r *SessionRepositoryImpl) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Session, error) {
	pk := UserPartitionKey(userID.String())

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk_prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":        &types.AttributeValueMemberS{Value: pk},
			":sk_prefix": &types.AttributeValueMemberS{Value: SessionQueryPrefix("")},
		},
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB Query エラー (GetAllByUserID): %v", err)
		return nil, appErrors.NewDynamoDBOperationError("get_sessions_by_user_id", err)
	}

	sessions := make([]*entity.Session, 0, len(result.Items))
	for _, item := range result.Items {
		session, err := r.itemToSession(item)
		if err != nil {
			r.logger.Warnf("セッション変換エラー（スキップ）: %v", err)
			continue
		}
		sessions = append(sessions, session)
	}

	r.logger.Infof("ユーザーセッション取得成功: UserID=%s, count=%d", userID.String(), len(sessions))
	return sessions, nil
}

// Update はセッションを更新する
func (r *SessionRepositoryImpl) Update(ctx context.Context, session *entity.Session) error {
	date := session.StartTime.Format("2006-01-02")
	pk := UserPartitionKey(session.UserID.String())
	sk := SessionSortKey(date, session.ID.String())

	session.UpdatedAt = time.Now()

	item := map[string]types.AttributeValue{
		"PK":         &types.AttributeValueMemberS{Value: pk},
		"SK":         &types.AttributeValueMemberS{Value: sk},
		"user_id":    &types.AttributeValueMemberS{Value: session.UserID.String()},
		"session_id": &types.AttributeValueMemberS{Value: session.ID.String()},
		"date":       &types.AttributeValueMemberS{Value: date},
		"start_time": &types.AttributeValueMemberS{Value: session.StartTime.Format(time.RFC3339)},
		"created_at": &types.AttributeValueMemberS{Value: session.CreatedAt.Format(time.RFC3339)},
		"updated_at": &types.AttributeValueMemberS{Value: session.UpdatedAt.Format(time.RFC3339)},
	}

	// Optional fields
	if session.EndTime != nil {
		item["end_time"] = &types.AttributeValueMemberS{Value: session.EndTime.Format(time.RFC3339)}
	}
	if session.AverageFocus != nil {
		item["average_focus"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", *session.AverageFocus)}
	}
	if session.TotalWorkMin != nil {
		item["total_work_min"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", *session.TotalWorkMin)}
	}
	if session.RoundCount != nil {
		item["round_count"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", *session.RoundCount)}
	}
	if session.BreakTime != nil {
		item["break_time"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", *session.BreakTime)}
	}

	input := &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_exists(PK) AND attribute_exists(SK)"),
	}

	_, err := r.client.PutItem(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB PutItem エラー: %v", err)
		return appErrors.NewDynamoDBOperationError("update_session", err)
	}

	r.logger.Infof("セッション更新成功: ID=%s", session.ID.String())
	return nil
}

// Complete はセッションを完了する(終了時刻、平均集中度、総作業時間を設定)
func (r *SessionRepositoryImpl) Complete(ctx context.Context, id, userID uuid.UUID, averageFocus float64, totalWorkMin, roundCount, breakTime int) error {
	// セッションを取得
	session, err := r.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	// セッションを完了状態に更新
	now := time.Now()
	session.EndTime = &now
	session.AverageFocus = &averageFocus
	session.TotalWorkMin = &totalWorkMin
	session.RoundCount = &roundCount
	session.BreakTime = &breakTime
	session.UpdatedAt = now

	// データベースに保存
	err = r.Update(ctx, session)
	if err != nil {
		return err
	}

	r.logger.Infof("セッション完了成功: ID=%s, avgFocus=%.2f, workMin=%d, rounds=%d",
		id.String(), averageFocus, totalWorkMin, roundCount)

	return nil
}

// Delete はセッションを削除する
func (r *SessionRepositoryImpl) Delete(ctx context.Context, id, userID uuid.UUID) error {
	// まずセッションを取得して日付を特定
	session, err := r.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	date := session.StartTime.Format("2006-01-02")
	pk := UserPartitionKey(userID.String())
	sk := SessionSortKey(date, id.String())

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		ConditionExpression: aws.String("attribute_exists(PK) AND attribute_exists(SK)"),
	}

	_, err = r.client.DeleteItem(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB DeleteItem エラー: %v", err)
		return appErrors.NewDynamoDBOperationError("delete_session", err)
	}

	r.logger.Infof("セッション削除成功: ID=%s", id.String())
	return nil
}

// GetUserIDBySessionID はセッションIDからユーザーIDを効率的に取得する（GSI使用）
func (r *SessionRepositoryImpl) GetUserIDBySessionID(ctx context.Context, sessionID uuid.UUID) (uuid.UUID, error) {
	r.logger.Infof("GSIでSessionIDからUserID検索開始: SessionID=%s", sessionID.String())

	// 🎯 GSIを使用してsession_idで効率的に検索
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("SessionIdIndex"), // GSI使用
		KeyConditionExpression: aws.String("session_id = :session_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":session_id": &types.AttributeValueMemberS{Value: sessionID.String()},
		},
		ProjectionExpression: aws.String("user_id, PK, SK"),
		Limit:                aws.Int32(1), // 最初の1件のみ
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		r.logger.Errorf("GSI Query エラー: %v", err)
		return uuid.Nil, appErrors.NewDynamoDBOperationError("gsi_session_id_query", err)
	}

	r.logger.Infof("GSI Query結果: %d件のレコードが見つかりました", len(result.Items))

	if len(result.Items) == 0 {
		r.logger.Warnf("GSI: SessionIDに対応するセッションが見つかりません: %s", sessionID.String())
		return uuid.Nil, appErrors.ErrRecordNotFound
	}

	// 最初のレコードからuser_idを取得
	item := result.Items[0]
	if userIDAttr, exists := item["user_id"]; exists {
		if s, ok := userIDAttr.(*types.AttributeValueMemberS); ok {
			if userID, err := uuid.Parse(s.Value); err == nil {
				r.logger.Infof("GSI SessionID->UserID変換成功: %s -> %s", sessionID.String(), userID.String())
				return userID, nil
			} else {
				r.logger.Errorf("UserIDパースエラー: %s, エラー: %v", s.Value, err)
			}
		}
	}

	r.logger.Errorf("GSI結果のuser_id属性が無効: SessionID=%s", sessionID.String())
	return uuid.Nil, appErrors.NewBadRequestError("invalid user_id attribute in GSI result")
}

// Helper methods

// itemToSession はDynamoDBアイテムをSessionモデルに変換する
func (r *SessionRepositoryImpl) itemToSession(item map[string]types.AttributeValue) (*entity.Session, error) {
	session := &entity.Session{}

	// user_id
	if userIDAttr, exists := item["user_id"]; exists {
		if s, ok := userIDAttr.(*types.AttributeValueMemberS); ok {
			if id, err := uuid.Parse(s.Value); err == nil {
				session.UserID = id
			}
		}
	}

	// session_id
	if sessionIDAttr, exists := item["session_id"]; exists {
		if s, ok := sessionIDAttr.(*types.AttributeValueMemberS); ok {
			if id, err := uuid.Parse(s.Value); err == nil {
				session.ID = id
			}
		}
	}

	// start_time
	if startTimeAttr, exists := item["start_time"]; exists {
		if s, ok := startTimeAttr.(*types.AttributeValueMemberS); ok {
			if startTime, err := time.Parse(time.RFC3339, s.Value); err == nil {
				session.StartTime = startTime
			}
		}
	}

	// end_time
	if endTimeAttr, exists := item["end_time"]; exists {
		if s, ok := endTimeAttr.(*types.AttributeValueMemberS); ok {
			if endTime, err := time.Parse(time.RFC3339, s.Value); err == nil {
				session.EndTime = &endTime
			}
		}
	}

	// average_focus
	if avgFocusAttr, exists := item["average_focus"]; exists {
		if n, ok := avgFocusAttr.(*types.AttributeValueMemberN); ok {
			if avgFocus, err := parseFloat64(n.Value); err == nil {
				session.AverageFocus = &avgFocus
			}
		}
	}

	// total_work_min
	if totalWorkAttr, exists := item["total_work_min"]; exists {
		if n, ok := totalWorkAttr.(*types.AttributeValueMemberN); ok {
			if totalWork, err := parseInt(n.Value); err == nil {
				session.TotalWorkMin = &totalWork
			}
		}
	}

	// round_count
	if roundCountAttr, exists := item["round_count"]; exists {
		if n, ok := roundCountAttr.(*types.AttributeValueMemberN); ok {
			if roundCount, err := parseInt(n.Value); err == nil {
				session.RoundCount = &roundCount
			}
		}
	}

	// break_time
	if breakTimeAttr, exists := item["break_time"]; exists {
		if n, ok := breakTimeAttr.(*types.AttributeValueMemberN); ok {
			if breakTime, err := parseInt(n.Value); err == nil {
				session.BreakTime = &breakTime
			}
		}
	}

	// created_at
	if createdAttr, exists := item["created_at"]; exists {
		if s, ok := createdAttr.(*types.AttributeValueMemberS); ok {
			if createdAt, err := time.Parse(time.RFC3339, s.Value); err == nil {
				session.CreatedAt = createdAt
			}
		}
	}

	// updated_at
	if updatedAttr, exists := item["updated_at"]; exists {
		if s, ok := updatedAttr.(*types.AttributeValueMemberS); ok {
			if updatedAt, err := time.Parse(time.RFC3339, s.Value); err == nil {
				session.UpdatedAt = updatedAt
			}
		}
	}

	return session, nil
}
