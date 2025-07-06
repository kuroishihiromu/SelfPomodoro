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

// RoundRepositoryImpl はDynamoDBを使用したRoundRepositoryの実装
type RoundRepositoryImpl struct {
	client    *dynamodb.Client
	tableName string
	logger    logger.Logger
}

// NewRoundRepository は新しいRoundRepositoryImplインスタンスを作成する
func NewRoundRepository(client *dynamodb.Client, cfg *config.Config, logger logger.Logger) repository.RoundRepository {
	return &RoundRepositoryImpl{
		client:    client,
		tableName: cfg.DynamoUnifiedTable, // 統一テーブルを使用
		logger:    logger,
	}
}

// Create はラウンドを作成する（UserIDを直接指定）
func (r *RoundRepositoryImpl) Create(ctx context.Context, round *entity.Round, userID uuid.UUID) error {
	date := round.StartTime.Format("2006-01-02")
	userIDStr := userID.String()

	r.logger.Infof("ラウンド作成開始: RoundID=%s, SessionID=%s, UserID=%s",
		round.ID.String(), round.SessionID.String(), userIDStr)

	pk := UserPartitionKey(userIDStr)
	sk := RoundSortKey(date, round.SessionID.String(), round.RoundOrder)

	r.logger.Infof("DynamoDB キー生成: PK=%s, SK=%s", pk, sk)

	// TTL設定: 30日後に自動削除（完了時のみ作成されるため）
	ttl := round.CreatedAt.Add(30 * 24 * time.Hour).Unix()

	item := map[string]types.AttributeValue{
		"PK":          &types.AttributeValueMemberS{Value: pk},
		"SK":          &types.AttributeValueMemberS{Value: sk},
		"user_id":     &types.AttributeValueMemberS{Value: userIDStr},
		"round_id":    &types.AttributeValueMemberS{Value: round.ID.String()},
		"session_id":  &types.AttributeValueMemberS{Value: round.SessionID.String()},
		"date":        &types.AttributeValueMemberS{Value: date},
		"round_order": &types.AttributeValueMemberN{Value: fmt.Sprintf("%03d", round.RoundOrder)},
		"start_time":  &types.AttributeValueMemberS{Value: round.StartTime.Format(time.RFC3339)},
		"ttl":         &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", ttl)},
		"created_at":  &types.AttributeValueMemberS{Value: round.CreatedAt.Format(time.RFC3339)},
		"updated_at":  &types.AttributeValueMemberS{Value: round.UpdatedAt.Format(time.RFC3339)},
	}

	// Optional fields
	if round.EndTime != nil {
		item["end_time"] = &types.AttributeValueMemberS{Value: round.EndTime.Format(time.RFC3339)}
	}
	if round.WorkTime != nil {
		item["work_time"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", *round.WorkTime)}
	}
	if round.BreakTime != nil {
		item["break_time"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", *round.BreakTime)}
	}
	if round.FocusScore != nil {
		item["focus_score"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", *round.FocusScore)}
	}

	input := &dynamodb.PutItemInput{
		TableName:           aws.String(r.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
	}

	_, err := r.client.PutItem(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB PutItem エラー: %v", err)
		// デバッグ: 詳細なエラー情報をログ出力
		r.logger.Errorf("PutItem詳細: TableName=%s, PK=%s, SK=%s", r.tableName, pk, sk)
		return appErrors.NewDynamoDBOperationError("create_round", err)
	}

	r.logger.Infof("ラウンド作成成功: ID=%s, SessionID=%s, Order=%d, UserID=%s, PK=%s, SK=%s",
		round.ID.String(), round.SessionID.String(), round.RoundOrder, userID, pk, sk)
	return nil
}

// GetByID はIDによってラウンドを取得する（GSI最適化版）
func (r *RoundRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Round, error) {
	r.logger.Infof("ラウンド取得開始: RoundID=%s", id.String())

	// 🎯 RoundIdIndex GSIを使用して効率的に検索
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("RoundIdIndex"), // GSI使用
		KeyConditionExpression: aws.String("round_id = :round_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":round_id": &types.AttributeValueMemberS{Value: id.String()},
		},
		Limit: aws.Int32(1), // 最初の1件のみ
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		r.logger.Errorf("RoundIdIndex GSI Query エラー: %v", err)
		// GSIが利用できない場合はScanにフォールバック
		return r.getByIDFallback(ctx, id)
	}

	r.logger.Infof("RoundIdIndex GSI Query結果: %d件のレコードが見つかりました", len(result.Items))

	if len(result.Items) == 0 {
		r.logger.Debugf("GSI: ラウンドが見つかりません: ID=%s", id.String())
		return nil, appErrors.ErrRecordNotFound
	}

	// GSI結果からラウンドを変換
	round, err := r.itemToRound(result.Items[0])
	if err != nil {
		r.logger.Errorf("GSI結果のラウンド変換エラー: %v", err)
		return nil, appErrors.NewDynamoDBOperationError("round_conversion", err)
	}

	r.logger.Infof("GSI経由でラウンド取得成功: ID=%s", id.String())
	return round, nil
}

// getByIDFallback はGSIが利用できない場合のフォールバック（従来のScan方式）
func (r *RoundRepositoryImpl) getByIDFallback(ctx context.Context, id uuid.UUID) (*entity.Round, error) {
	r.logger.Warnf("GSI利用不可、Scanフォールバック使用: RoundID=%s", id.String())

	input := &dynamodb.ScanInput{
		TableName:        aws.String(r.tableName),
		FilterExpression: aws.String("round_id = :round_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":round_id": &types.AttributeValueMemberS{Value: id.String()},
		},
		Limit: aws.Int32(1),
	}

	result, err := r.client.Scan(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB Scan エラー (Fallback): %v", err)
		return nil, appErrors.NewDynamoDBOperationError("get_round_by_id_fallback", err)
	}

	if len(result.Items) == 0 {
		r.logger.Debugf("Scan: ラウンドが見つかりません: ID=%s", id.String())
		return nil, appErrors.ErrRecordNotFound
	}

	round, err := r.itemToRound(result.Items[0])
	if err != nil {
		r.logger.Errorf("Scan結果のラウンド変換エラー: %v", err)
		return nil, appErrors.NewDynamoDBOperationError("round_conversion", err)
	}

	r.logger.Infof("Scan経由でラウンド取得成功: ID=%s", id.String())
	return round, nil
}

// GetByIDWithUserID はIDとユーザーIDによってラウンドを取得する
func (r *RoundRepositoryImpl) GetByIDWithUserID(ctx context.Context, id, userID uuid.UUID) (*entity.Round, error) {
	pk := UserPartitionKey(userID.String())

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk_prefix)"),
		FilterExpression:       aws.String("round_id = :round_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":        &types.AttributeValueMemberS{Value: pk},
			":sk_prefix": &types.AttributeValueMemberS{Value: RoundQueryPrefix("")},
			":round_id":  &types.AttributeValueMemberS{Value: id.String()},
		},
		Limit: aws.Int32(1),
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB Query エラー (GetByIDWithUserID): %v", err)
		return nil, appErrors.NewDynamoDBOperationError("get_round_by_id_with_user_id", err)
	}

	if len(result.Items) == 0 {
		r.logger.Debugf("ラウンドが見つかりません: ID=%s, UserID=%s", id.String(), userID.String())
		return nil, appErrors.ErrRecordNotFound
	}

	round, err := r.itemToRound(result.Items[0])
	if err != nil {
		r.logger.Errorf("ラウンド変換エラー: %v", err)
		return nil, appErrors.NewDynamoDBOperationError("round_conversion", err)
	}

	r.logger.Debugf("ラウンド取得成功: ID=%s", id.String())
	return round, nil
}

// GetBySessionIDWithUserID はセッションIDとユーザーIDに紐づくすべてのラウンドを取得する
func (r *RoundRepositoryImpl) GetBySessionIDWithUserID(ctx context.Context, sessionID, userID uuid.UUID) ([]*entity.Round, error) {
	pk := UserPartitionKey(userID.String())

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk_prefix)"),
		FilterExpression:       aws.String("session_id = :session_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":         &types.AttributeValueMemberS{Value: pk},
			":sk_prefix":  &types.AttributeValueMemberS{Value: RoundQueryPrefix("")},
			":session_id": &types.AttributeValueMemberS{Value: sessionID.String()},
		},
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB Query エラー (GetBySessionIDWithUserID): %v", err)
		return nil, appErrors.NewDynamoDBOperationError("get_rounds_by_session_id_with_user_id", err)
	}

	rounds := make([]*entity.Round, 0, len(result.Items))
	for _, item := range result.Items {
		round, err := r.itemToRound(item)
		if err != nil {
			r.logger.Warnf("ラウンド変換エラー（スキップ）: %v", err)
			continue
		}
		rounds = append(rounds, round)
	}

	r.logger.Infof("セッションラウンド取得成功: SessionID=%s, UserID=%s, count=%d",
		sessionID.String(), userID.String(), len(rounds))
	return rounds, nil
}

// Complete はラウンドを完了する（インターフェース準拠）
func (r *RoundRepositoryImpl) Complete(ctx context.Context, id uuid.UUID, focusScore *int, worktime, breaktime int) error {
	// ラウンドを取得
	round, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// ドメインロジックで完了処理
	err = round.CompleteWith(focusScore, worktime, breaktime)
	if err != nil {
		return err
	}

	// GSIから取得した既存レコードを直接更新
	err = r.updateRoundByGSI(ctx, round)
	if err != nil {
		return err
	}

	r.logger.Infof("ラウンド完了成功: ID=%s, focusScore=%v", id.String(), focusScore)
	return nil
}

// Helper methods

// itemToRound はDynamoDBアイテムをRoundモデルに変換する
func (r *RoundRepositoryImpl) itemToRound(item map[string]types.AttributeValue) (*entity.Round, error) {
	round := &entity.Round{}

	// round_id
	if roundIDAttr, exists := item["round_id"]; exists {
		if s, ok := roundIDAttr.(*types.AttributeValueMemberS); ok {
			if id, err := uuid.Parse(s.Value); err == nil {
				round.ID = id
			}
		}
	}

	// session_id
	if sessionIDAttr, exists := item["session_id"]; exists {
		if s, ok := sessionIDAttr.(*types.AttributeValueMemberS); ok {
			if id, err := uuid.Parse(s.Value); err == nil {
				round.SessionID = id
			}
		}
	}

	// round_order
	if roundOrderAttr, exists := item["round_order"]; exists {
		if n, ok := roundOrderAttr.(*types.AttributeValueMemberN); ok {
			if roundOrder, err := parseInt(n.Value); err == nil {
				round.RoundOrder = roundOrder
			}
		}
	}

	// start_time
	if startTimeAttr, exists := item["start_time"]; exists {
		if s, ok := startTimeAttr.(*types.AttributeValueMemberS); ok {
			if startTime, err := time.Parse(time.RFC3339, s.Value); err == nil {
				round.StartTime = startTime
			}
		}
	}

	// end_time
	if endTimeAttr, exists := item["end_time"]; exists {
		if s, ok := endTimeAttr.(*types.AttributeValueMemberS); ok {
			if endTime, err := time.Parse(time.RFC3339, s.Value); err == nil {
				round.EndTime = &endTime
			}
		}
	}

	// work_time
	if workTimeAttr, exists := item["work_time"]; exists {
		if n, ok := workTimeAttr.(*types.AttributeValueMemberN); ok {
			if workTime, err := parseInt(n.Value); err == nil {
				round.WorkTime = &workTime
			}
		}
	}

	// break_time
	if breakTimeAttr, exists := item["break_time"]; exists {
		if n, ok := breakTimeAttr.(*types.AttributeValueMemberN); ok {
			if breakTime, err := parseInt(n.Value); err == nil {
				round.BreakTime = &breakTime
			}
		}
	}

	// focus_score
	if focusScoreAttr, exists := item["focus_score"]; exists {
		if n, ok := focusScoreAttr.(*types.AttributeValueMemberN); ok {
			if focusScore, err := parseInt(n.Value); err == nil {
				round.FocusScore = &focusScore
			}
		}
	}

	// created_at
	if createdAttr, exists := item["created_at"]; exists {
		if s, ok := createdAttr.(*types.AttributeValueMemberS); ok {
			if createdAt, err := time.Parse(time.RFC3339, s.Value); err == nil {
				round.CreatedAt = createdAt
			}
		}
	}

	// updated_at
	if updatedAttr, exists := item["updated_at"]; exists {
		if s, ok := updatedAttr.(*types.AttributeValueMemberS); ok {
			if updatedAt, err := time.Parse(time.RFC3339, s.Value); err == nil {
				round.UpdatedAt = updatedAt
			}
		}
	}

	return round, nil
}

// updateRoundByGSI はGSIを使用してラウンドを更新する
func (r *RoundRepositoryImpl) updateRoundByGSI(ctx context.Context, round *entity.Round) error {
	// GSIを使用して既存のレコードのPK/SKを取得
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("RoundIdIndex"),
		KeyConditionExpression: aws.String("round_id = :round_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":round_id": &types.AttributeValueMemberS{Value: round.ID.String()},
		},
		ProjectionExpression: aws.String("PK, SK, user_id"),
		Limit:                aws.Int32(1),
	}

	result, err := r.client.Query(ctx, input)
	if err != nil || len(result.Items) == 0 {
		r.logger.Errorf("ラウンド更新用GSI Query失敗: %v", err)
		return appErrors.NewDynamoDBOperationError("gsi_query_for_update", err)
	}

	// 既存レコードからPK/SKを取得
	existingItem := result.Items[0]
	pkAttr, ok := existingItem["PK"]
	if !ok {
		return appErrors.NewDynamoDBOperationError("pk_not_found", nil)
	}
	skAttr, ok := existingItem["SK"]
	if !ok {
		return appErrors.NewDynamoDBOperationError("sk_not_found", nil)
	}

	pk := pkAttr.(*types.AttributeValueMemberS).Value
	sk := skAttr.(*types.AttributeValueMemberS).Value

	// UpdateExpressionを使用して効率的に更新
	round.UpdatedAt = time.Now()

	updateExpression := "SET end_time = :end_time, work_time = :work_time, break_time = :break_time, updated_at = :updated_at"
	expressionAttributeValues := map[string]types.AttributeValue{
		":end_time":   &types.AttributeValueMemberS{Value: round.EndTime.Format(time.RFC3339)},
		":work_time":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", *round.WorkTime)},
		":break_time": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", *round.BreakTime)},
		":updated_at": &types.AttributeValueMemberS{Value: round.UpdatedAt.Format(time.RFC3339)},
	}

	if round.FocusScore != nil {
		updateExpression += ", focus_score = :focus_score"
		expressionAttributeValues[":focus_score"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", *round.FocusScore)}
	}

	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
		ConditionExpression:       aws.String("attribute_exists(PK) AND attribute_exists(SK)"),
	}

	_, err = r.client.UpdateItem(ctx, updateInput)
	if err != nil {
		r.logger.Errorf("DynamoDB UpdateItem エラー: %v", err)
		return appErrors.NewDynamoDBOperationError("update_round_gsi", err)
	}

	r.logger.Infof("ラウンドGSI更新成功: ID=%s", round.ID.String())
	return nil
}
