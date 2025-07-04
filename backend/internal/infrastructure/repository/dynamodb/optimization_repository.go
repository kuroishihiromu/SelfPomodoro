package dynamodb

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// OptimizationRepositoryImpl は最適化データリポジトリのDynamoDB実装
type OptimizationRepositoryImpl struct {
	client    *dynamodb.Client
	tableName string
	logger    logger.Logger
}

// NewOptimizationRepository は新しい最適化リポジトリを作成する
func NewOptimizationRepository(client *dynamodb.Client, tableName string, logger logger.Logger) repository.OptimizationRepository {
	return &OptimizationRepositoryImpl{
		client:    client,
		tableName: tableName,
		logger:    logger,
	}
}

// SaveRoundOptimizationLog はラウンド最適化ログを保存する
func (r *OptimizationRepositoryImpl) SaveRoundOptimizationLog(ctx context.Context, log *model.RoundOptimizationLog) error {
	r.logger.Infof("ラウンド最適化ログ保存開始: UserID=%s", log.UserID[:8]+"...")

	// 統合テーブル用アイテム作成
	// PK: USER#{user_id}, SK: OPTIMIZATION_ROUND#{timestamp}
	item := map[string]types.AttributeValue{
		"PK":          &types.AttributeValueMemberS{Value: UserPartitionKey(log.UserID)},
		"SK":          &types.AttributeValueMemberS{Value: OptimizationRoundSortKey(log.Timestamp)},
		"user_id":     &types.AttributeValueMemberS{Value: log.UserID},
		"timestamp":   &types.AttributeValueMemberS{Value: log.Timestamp},
		"work_time":   &types.AttributeValueMemberN{Value: strconv.Itoa(log.WorkTime)},
		"break_time":  &types.AttributeValueMemberN{Value: strconv.Itoa(log.BreakTime)},
		"focus_score": &types.AttributeValueMemberN{Value: strconv.Itoa(log.FocusScore)},
		"created_at":  &types.AttributeValueMemberS{Value: log.CreatedAt},
	}

	_, err := r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	if err != nil {
		r.logger.Errorf("ラウンド最適化ログ保存失敗: %v", err)
		return appErrors.NewDynamoDBOperationError("put_round_optimization_log", err)
	}

	r.logger.Infof("ラウンド最適化ログ保存完了: UserID=%s", log.UserID[:8]+"...")
	return nil
}

// GetRoundOptimizationHistory はラウンド最適化履歴を取得する
func (r *OptimizationRepositoryImpl) GetRoundOptimizationHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*model.RoundOptimizationLog, error) {
	r.logger.Infof("ラウンド最適化履歴取得開始: UserID=%s, Limit=%d", userID.String()[:8]+"...", limit)

	// クエリ実行
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk_prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":        &types.AttributeValueMemberS{Value: UserPartitionKey(userID.String())},
			":sk_prefix": &types.AttributeValueMemberS{Value: "OPTIMIZATION_LOG#ROUND#"},
		},
		ScanIndexForward: aws.Bool(false), // 新しい順
		Limit:            aws.Int32(int32(limit)),
	})

	if err != nil {
		r.logger.Errorf("ラウンド最適化履歴取得失敗: %v", err)
		return nil, appErrors.NewDynamoDBOperationError("query_round_optimization_history", err)
	}

	// アイテムを変換
	logs := make([]*model.RoundOptimizationLog, 0, len(result.Items))
	for _, item := range result.Items {
		log, err := r.convertToRoundOptimizationLog(item)
		if err != nil {
			r.logger.Warnf("アイテム変換エラー: %v", err)
			continue
		}
		logs = append(logs, log)
	}


	r.logger.Infof("ラウンド最適化履歴取得完了: UserID=%s, 件数=%d", userID.String()[:8]+"...", len(logs))
	return logs, nil
}

// GetLatestRoundOptimizationResult は最新のラウンド最適化結果を取得する
func (r *OptimizationRepositoryImpl) GetLatestRoundOptimizationResult(ctx context.Context, userID uuid.UUID) (*model.RoundOptimizationLog, error) {
	logs, err := r.GetRoundOptimizationHistory(ctx, userID, 1)
	if err != nil {
		return nil, err
	}

	if len(logs) == 0 {
		return nil, appErrors.NewNotFoundError("ラウンド最適化結果")
	}

	return logs[0], nil
}

// SaveSessionOptimizationLog はセッション最適化ログを保存する
func (r *OptimizationRepositoryImpl) SaveSessionOptimizationLog(ctx context.Context, log *model.SessionOptimizationLog) error {
	r.logger.Infof("セッション最適化ログ保存開始: UserID=%s", log.UserID[:8]+"...")

	// 統合テーブル用アイテム作成
	// PK: USER#{user_id}, SK: OPTIMIZATION_SESSION#{timestamp}
	item := map[string]types.AttributeValue{
		"PK":               &types.AttributeValueMemberS{Value: UserPartitionKey(log.UserID)},
		"SK":               &types.AttributeValueMemberS{Value: OptimizationSessionSortKey(log.Timestamp)},
		"user_id":          &types.AttributeValueMemberS{Value: log.UserID},
		"timestamp":        &types.AttributeValueMemberS{Value: log.Timestamp},
		"round_count":      &types.AttributeValueMemberN{Value: strconv.Itoa(log.RoundCount)},
		"break_time":       &types.AttributeValueMemberN{Value: strconv.Itoa(log.BreakTime)},
		"avg_focus_score":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", log.AvgFocusScore)},
		"total_work_time":  &types.AttributeValueMemberN{Value: strconv.Itoa(log.TotalWorkTime)},
		"created_at":       &types.AttributeValueMemberS{Value: log.CreatedAt},
	}

	_, err := r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	if err != nil {
		r.logger.Errorf("セッション最適化ログ保存失敗: %v", err)
		return appErrors.NewDynamoDBOperationError("put_session_optimization_log", err)
	}

	r.logger.Infof("セッション最適化ログ保存完了: UserID=%s", log.UserID[:8]+"...")
	return nil
}

// GetSessionOptimizationHistory はセッション最適化履歴を取得する
func (r *OptimizationRepositoryImpl) GetSessionOptimizationHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*model.SessionOptimizationLog, error) {
	r.logger.Infof("セッション最適化履歴取得開始: UserID=%s, Limit=%d", userID.String()[:8]+"...", limit)

	// クエリ実行
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk_prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":        &types.AttributeValueMemberS{Value: UserPartitionKey(userID.String())},
			":sk_prefix": &types.AttributeValueMemberS{Value: "OPTIMIZATION_LOG#SESSION#"},
		},
		ScanIndexForward: aws.Bool(false), // 新しい順
		Limit:            aws.Int32(int32(limit)),
	})

	if err != nil {
		r.logger.Errorf("セッション最適化履歴取得失敗: %v", err)
		return nil, appErrors.NewDynamoDBOperationError("query_session_optimization_history", err)
	}

	// アイテムを変換
	logs := make([]*model.SessionOptimizationLog, 0, len(result.Items))
	for _, item := range result.Items {
		log, err := r.convertToSessionOptimizationLog(item)
		if err != nil {
			r.logger.Warnf("アイテム変換エラー: %v", err)
			continue
		}
		logs = append(logs, log)
	}


	r.logger.Infof("セッション最適化履歴取得完了: UserID=%s, 件数=%d", userID.String()[:8]+"...", len(logs))
	return logs, nil
}

// GetLatestSessionOptimizationResult は最新のセッション最適化結果を取得する
func (r *OptimizationRepositoryImpl) GetLatestSessionOptimizationResult(ctx context.Context, userID uuid.UUID) (*model.SessionOptimizationLog, error) {
	logs, err := r.GetSessionOptimizationHistory(ctx, userID, 1)
	if err != nil {
		return nil, err
	}

	if len(logs) == 0 {
		return nil, appErrors.NewNotFoundError("セッション最適化結果")
	}

	return logs[0], nil
}

// GetOptimizationEffectiveness は最適化の効果性を取得する
func (r *OptimizationRepositoryImpl) GetOptimizationEffectiveness(ctx context.Context, userID uuid.UUID, since time.Time) (*model.OptimizationEffectiveness, error) {
	r.logger.Infof("最適化効果性取得開始: UserID=%s, Since=%s", userID.String()[:8]+"...", since.Format("2006-01-02"))

	// 期間内のラウンド最適化データ取得
	roundLogs, err := r.GetRoundOptimizationHistory(ctx, userID, 100)
	if err != nil {
		return nil, err
	}

	// 期間内のセッション最適化データ取得
	sessionLogs, err := r.GetSessionOptimizationHistory(ctx, userID, 100)
	if err != nil {
		return nil, err
	}

	// 効果性計算
	effectiveness := model.NewOptimizationEffectiveness(userID, since, time.Now())
	effectiveness.RoundOptimizations = len(roundLogs)
	effectiveness.SessionOptimizations = len(sessionLogs)

	// 平均集中度改善の計算（簡単な実装）
	if len(roundLogs) > 0 {
		totalImprovement := 0.0
		for _, log := range roundLogs {
			totalImprovement += float64(log.FocusScore)
		}
		effectiveness.AvgFocusImprovement = totalImprovement / float64(len(roundLogs))
	}

	// 最後の最適化時刻設定
	if len(roundLogs) > 0 {
		latestTime, _ := time.Parse(time.RFC3339, roundLogs[0].Timestamp)
		effectiveness.LastOptimizedAt = latestTime
	}

	r.logger.Infof("最適化効果性取得完了: UserID=%s", userID.String()[:8]+"...")
	return effectiveness, nil
}

// CountOptimizationLogs は最適化ログ数をカウントする
func (r *OptimizationRepositoryImpl) CountOptimizationLogs(ctx context.Context, userID uuid.UUID) (roundCount, sessionCount int, err error) {
	r.logger.Infof("最適化ログカウント開始: UserID=%s", userID.String()[:8]+"...")

	// ラウンド最適化ログカウント
	roundResult, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk_prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":        &types.AttributeValueMemberS{Value: UserPartitionKey(userID.String())},
			":sk_prefix": &types.AttributeValueMemberS{Value: "OPTIMIZATION_LOG#ROUND#"},
		},
		Select: types.SelectCount,
	})

	if err != nil {
		r.logger.Errorf("ラウンド最適化ログカウント失敗: %v", err)
		return 0, 0, appErrors.NewDynamoDBOperationError("count_round_optimization_logs", err)
	}

	// セッション最適化ログカウント
	sessionResult, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk_prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":        &types.AttributeValueMemberS{Value: UserPartitionKey(userID.String())},
			":sk_prefix": &types.AttributeValueMemberS{Value: "OPTIMIZATION_LOG#SESSION#"},
		},
		Select: types.SelectCount,
	})

	if err != nil {
		r.logger.Errorf("セッション最適化ログカウント失敗: %v", err)
		return 0, 0, appErrors.NewDynamoDBOperationError("count_session_optimization_logs", err)
	}

	roundCount = int(roundResult.Count)
	sessionCount = int(sessionResult.Count)

	r.logger.Infof("最適化ログカウント完了: UserID=%s, Round=%d, Session=%d", 
		userID.String()[:8]+"...", roundCount, sessionCount)
	
	return roundCount, sessionCount, nil
}

// convertToRoundOptimizationLog はDynamoDBアイテムをRoundOptimizationLogに変換する
func (r *OptimizationRepositoryImpl) convertToRoundOptimizationLog(item map[string]types.AttributeValue) (*model.RoundOptimizationLog, error) {
	log := &model.RoundOptimizationLog{}

	if val, ok := item["user_id"]; ok && val.(*types.AttributeValueMemberS) != nil {
		log.UserID = val.(*types.AttributeValueMemberS).Value
	}

	if val, ok := item["timestamp"]; ok && val.(*types.AttributeValueMemberS) != nil {
		log.Timestamp = val.(*types.AttributeValueMemberS).Value
	}

	if val, ok := item["work_time"]; ok && val.(*types.AttributeValueMemberN) != nil {
		if workTime, err := strconv.Atoi(val.(*types.AttributeValueMemberN).Value); err == nil {
			log.WorkTime = workTime
		}
	}

	if val, ok := item["break_time"]; ok && val.(*types.AttributeValueMemberN) != nil {
		if breakTime, err := strconv.Atoi(val.(*types.AttributeValueMemberN).Value); err == nil {
			log.BreakTime = breakTime
		}
	}

	if val, ok := item["focus_score"]; ok && val.(*types.AttributeValueMemberN) != nil {
		if focusScore, err := strconv.Atoi(val.(*types.AttributeValueMemberN).Value); err == nil {
			log.FocusScore = focusScore
		}
	}

	if val, ok := item["created_at"]; ok && val.(*types.AttributeValueMemberS) != nil {
		log.CreatedAt = val.(*types.AttributeValueMemberS).Value
	}

	return log, nil
}

// convertToSessionOptimizationLog はDynamoDBアイテムをSessionOptimizationLogに変換する
func (r *OptimizationRepositoryImpl) convertToSessionOptimizationLog(item map[string]types.AttributeValue) (*model.SessionOptimizationLog, error) {
	log := &model.SessionOptimizationLog{}

	if val, ok := item["user_id"]; ok && val.(*types.AttributeValueMemberS) != nil {
		log.UserID = val.(*types.AttributeValueMemberS).Value
	}

	if val, ok := item["timestamp"]; ok && val.(*types.AttributeValueMemberS) != nil {
		log.Timestamp = val.(*types.AttributeValueMemberS).Value
	}

	if val, ok := item["round_count"]; ok && val.(*types.AttributeValueMemberN) != nil {
		if roundCount, err := strconv.Atoi(val.(*types.AttributeValueMemberN).Value); err == nil {
			log.RoundCount = roundCount
		}
	}

	if val, ok := item["break_time"]; ok && val.(*types.AttributeValueMemberN) != nil {
		if breakTime, err := strconv.Atoi(val.(*types.AttributeValueMemberN).Value); err == nil {
			log.BreakTime = breakTime
		}
	}

	if val, ok := item["avg_focus_score"]; ok && val.(*types.AttributeValueMemberN) != nil {
		if avgFocusScore, err := strconv.ParseFloat(val.(*types.AttributeValueMemberN).Value, 64); err == nil {
			log.AvgFocusScore = avgFocusScore
		}
	}

	if val, ok := item["total_work_time"]; ok && val.(*types.AttributeValueMemberN) != nil {
		if totalWorkTime, err := strconv.Atoi(val.(*types.AttributeValueMemberN).Value); err == nil {
			log.TotalWorkTime = totalWorkTime
		}
	}

	if val, ok := item["created_at"]; ok && val.(*types.AttributeValueMemberS) != nil {
		log.CreatedAt = val.(*types.AttributeValueMemberS).Value
	}

	return log, nil
}