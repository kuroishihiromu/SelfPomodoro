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
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// StatisticsRepositoryImpl はDynamoDBを使用したStatisticsRepositoryの実装
type StatisticsRepositoryImpl struct {
	client    *dynamodb.Client
	tableName string
	logger    logger.Logger
}

// NewStatisticsRepository は新しいStatisticsRepositoryImplインスタンスを作成する
func NewStatisticsRepository(client *dynamodb.Client, cfg *config.Config, logger logger.Logger) repository.StatisticsRepository {
	return &StatisticsRepositoryImpl{
		client:    client,
		tableName: cfg.DynamoUnifiedTable, // 統一テーブルを使用
		logger:    logger,
	}
}

// GetFocusTrend は指定期間内の日別集中度統計を取得する
func (r *StatisticsRepositoryImpl) GetFocusTrend(ctx context.Context, userID uuid.UUID, period *model.StatisticsPeriod) ([]*model.FocusTrendItem, error) {
	pk := UserPartitionKey(userID.String())
	
	// クエリで期間内の日別統計を取得
	startDate := DateFromTime(period.StartDate)
	endDate := DateFromTime(period.EndDate)
	startSK, endSK := DailyStatsQueryRange(startDate, endDate)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND SK BETWEEN :start_sk AND :end_sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":       &types.AttributeValueMemberS{Value: pk},
			":start_sk": &types.AttributeValueMemberS{Value: startSK},
			":end_sk":   &types.AttributeValueMemberS{Value: endSK},
		},
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB Query エラー (FocusTrend): %v", err)
		return nil, appErrors.NewDynamoDBOperationError("get_focus_trend", err)
	}

	// DynamoDBアイテムをドメインモデルに変換
	var trendItems []*model.FocusTrendItem
	for _, item := range result.Items {
		dailyStats, err := r.itemToAggregatedDailyStats(item)
		if err != nil {
			r.logger.Warnf("日別統計変換エラー: %v", err)
			continue
		}
		trendItems = append(trendItems, dailyStats.ToFocusTrendItem())
	}

	// データがない日付を補完
	trendItems = r.fillMissingDates(trendItems, period)

	r.logger.Infof("集中度トレンド取得成功: period=%s～%s, items=%d",
		period.StartDate.Format("2006-01-02"), period.EndDate.Format("2006-01-02"), len(trendItems))

	return trendItems, nil
}

// GetFocusHeatmap は指定期間内の時間帯別集中度統計を取得する
func (r *StatisticsRepositoryImpl) GetFocusHeatmap(ctx context.Context, userID uuid.UUID, period *model.StatisticsPeriod) ([]*model.FocusHeatmapItem, error) {
	pk := UserPartitionKey(userID.String())
	
	// クエリで期間内の時間別統計を取得
	startDate := DateFromTime(period.StartDate)
	endDate := DateFromTime(period.EndDate)
	startSK, endSK := HourlyStatsQueryRange(startDate, endDate)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND SK BETWEEN :start_sk AND :end_sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":       &types.AttributeValueMemberS{Value: pk},
			":start_sk": &types.AttributeValueMemberS{Value: startSK},
			":end_sk":   &types.AttributeValueMemberS{Value: endSK},
		},
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB Query エラー (FocusHeatmap): %v", err)
		return nil, appErrors.NewDynamoDBOperationError("get_focus_heatmap", err)
	}

	// DynamoDBアイテムをドメインモデルに変換
	var heatmapItems []*model.FocusHeatmapItem
	for _, item := range result.Items {
		hourlyStats, err := r.itemToAggregatedHourlyStats(item)
		if err != nil {
			r.logger.Warnf("時間別統計変換エラー: %v", err)
			continue
		}
		heatmapItems = append(heatmapItems, hourlyStats.ToFocusHeatmapItem())
	}

	r.logger.Infof("集中度ヒートマップ取得成功: period=%s～%s, items=%d",
		period.StartDate.Format("2006-01-02"), period.EndDate.Format("2006-01-02"), len(heatmapItems))

	return heatmapItems, nil
}

// GetWeeklyStats は指定期間内の週別統計を取得する
func (r *StatisticsRepositoryImpl) GetWeeklyStats(ctx context.Context, userID uuid.UUID, period *model.StatisticsPeriod) ([]*model.FocusTrendItem, error) {
	pk := UserPartitionKey(userID.String())
	
	// 期間の週境界を計算
	startWeek, _ := model.GetWeekBoundaries(period.StartDate)
	endWeek, _ := model.GetWeekBoundaries(period.EndDate)
	
	// クエリで期間内の週別統計を取得
	startSK, endSK := WeeklyStatsQueryRange(startWeek, endWeek)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND SK BETWEEN :start_sk AND :end_sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":       &types.AttributeValueMemberS{Value: pk},
			":start_sk": &types.AttributeValueMemberS{Value: startSK},
			":end_sk":   &types.AttributeValueMemberS{Value: endSK},
		},
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB Query エラー (GetWeeklyStats): %v", err)
		return nil, appErrors.NewDynamoDBOperationError("get_weekly_stats", err)
	}

	// DynamoDBアイテムをドメインモデルに変換
	var trendItems []*model.FocusTrendItem
	for _, item := range result.Items {
		weeklyStats, err := r.itemToAggregatedWeeklyStats(item)
		if err != nil {
			r.logger.Warnf("週別統計変換エラー: %v", err)
			continue
		}
		trendItems = append(trendItems, weeklyStats.ToFocusTrendItem())
	}

	// データがない週を補完
	trendItems = r.fillMissingWeeks(trendItems, period)

	r.logger.Infof("週別統計取得成功: period=%s～%s, items=%d",
		period.StartDate.Format("2006-01-02"), period.EndDate.Format("2006-01-02"), len(trendItems))

	return trendItems, nil
}

// GetAvgFocusScoreByDate は指定日の平均集中度を取得する
func (r *StatisticsRepositoryImpl) GetAvgFocusScoreByDate(ctx context.Context, userID uuid.UUID, date time.Time) (float64, error) {
	pk := UserPartitionKey(userID.String())
	sk := DailyStatsSortKey(DateFromTime(date))

	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB GetItem エラー (AvgFocusScoreByDate): %v", err)
		return 0, appErrors.NewDynamoDBOperationError("get_avg_focus_score_by_date", err)
	}

	if result.Item == nil {
		r.logger.Debugf("指定日の統計データが見つかりません: %s", date.Format("2006-01-02"))
		return 0.0, nil // データがない場合は0を返す
	}

	dailyStats, err := r.itemToAggregatedDailyStats(result.Item)
	if err != nil {
		r.logger.Errorf("日別統計変換エラー: %v", err)
		return 0, appErrors.NewDynamoDBOperationError("conversion", err)
	}

	r.logger.Debugf("指定日の平均集中度取得成功: date=%s, avgFocus=%.1f",
		date.Format("2006-01-02"), dailyStats.AvgFocusScore)

	return dailyStats.AvgFocusScore, nil
}

// GetAvgFocusScoreByHour は指定日時の平均集中度を取得する
func (r *StatisticsRepositoryImpl) GetAvgFocusScoreByHour(ctx context.Context, userID uuid.UUID, date time.Time, hour int) (float64, error) {
	pk := UserPartitionKey(userID.String())
	sk := HourlyStatsSortKey(DateFromTime(date), hour)

	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		r.logger.Errorf("DynamoDB GetItem エラー (AvgFocusScoreByHour): %v", err)
		return 0, appErrors.NewDynamoDBOperationError("get_avg_focus_score_by_hour", err)
	}

	if result.Item == nil {
		r.logger.Debugf("指定日時の統計データが見つかりません: %s %02d:00", date.Format("2006-01-02"), hour)
		return 0.0, nil // データがない場合は0を返す
	}

	hourlyStats, err := r.itemToAggregatedHourlyStats(result.Item)
	if err != nil {
		r.logger.Errorf("時間別統計変換エラー: %v", err)
		return 0, appErrors.NewDynamoDBOperationError("conversion", err)
	}

	r.logger.Debugf("指定日時の平均集中度取得成功: date=%s %02d:00, avgFocus=%.1f",
		date.Format("2006-01-02"), hour, hourlyStats.AvgFocusScore)

	return hourlyStats.AvgFocusScore, nil
}

// Helper methods

// itemToAggregatedDailyStats はDynamoDBアイテムを日別統計に変換する
func (r *StatisticsRepositoryImpl) itemToAggregatedDailyStats(item map[string]types.AttributeValue) (*model.AggregatedDailyStats, error) {
	stats := &model.AggregatedDailyStats{}

	// user_id
	if userIDAttr, exists := item["user_id"]; exists {
		if s, ok := userIDAttr.(*types.AttributeValueMemberS); ok {
			stats.UserID = s.Value
		}
	}

	// date
	if dateAttr, exists := item["date"]; exists {
		if s, ok := dateAttr.(*types.AttributeValueMemberS); ok {
			stats.Date = s.Value
		}
	}

	// total_rounds
	if totalRoundsAttr, exists := item["total_rounds"]; exists {
		if n, ok := totalRoundsAttr.(*types.AttributeValueMemberN); ok {
			if rounds, err := strconv.Atoi(n.Value); err == nil {
				stats.TotalRounds = rounds
			}
		}
	}

	// avg_focus_score
	if avgFocusAttr, exists := item["avg_focus_score"]; exists {
		if n, ok := avgFocusAttr.(*types.AttributeValueMemberN); ok {
			if score, err := strconv.ParseFloat(n.Value, 64); err == nil {
				stats.AvgFocusScore = score
			}
		}
	}

	// total_work_min
	if totalWorkAttr, exists := item["total_work_min"]; exists {
		if n, ok := totalWorkAttr.(*types.AttributeValueMemberN); ok {
			if workMin, err := strconv.Atoi(n.Value); err == nil {
				stats.TotalWorkMin = workMin
			}
		}
	}

	// total_break_min
	if totalBreakAttr, exists := item["total_break_min"]; exists {
		if n, ok := totalBreakAttr.(*types.AttributeValueMemberN); ok {
			if breakMin, err := strconv.Atoi(n.Value); err == nil {
				stats.TotalBreakMin = breakMin
			}
		}
	}

	// session_count
	if sessionCountAttr, exists := item["session_count"]; exists {
		if n, ok := sessionCountAttr.(*types.AttributeValueMemberN); ok {
			if count, err := strconv.Atoi(n.Value); err == nil {
				stats.SessionCount = count
			}
		}
	}

	// updated_at
	if updatedAttr, exists := item["updated_at"]; exists {
		if s, ok := updatedAttr.(*types.AttributeValueMemberS); ok {
			if updatedAt, err := time.Parse(time.RFC3339, s.Value); err == nil {
				stats.UpdatedAt = updatedAt
			}
		}
	}

	return stats, nil
}

// itemToAggregatedHourlyStats はDynamoDBアイテムを時間別統計に変換する
func (r *StatisticsRepositoryImpl) itemToAggregatedHourlyStats(item map[string]types.AttributeValue) (*model.AggregatedHourlyStats, error) {
	stats := &model.AggregatedHourlyStats{}

	// user_id
	if userIDAttr, exists := item["user_id"]; exists {
		if s, ok := userIDAttr.(*types.AttributeValueMemberS); ok {
			stats.UserID = s.Value
		}
	}

	// date
	if dateAttr, exists := item["date"]; exists {
		if s, ok := dateAttr.(*types.AttributeValueMemberS); ok {
			stats.Date = s.Value
		}
	}

	// hour
	if hourAttr, exists := item["hour"]; exists {
		if n, ok := hourAttr.(*types.AttributeValueMemberN); ok {
			if hour, err := strconv.Atoi(n.Value); err == nil {
				stats.Hour = hour
			}
		}
	}

	// total_rounds
	if totalRoundsAttr, exists := item["total_rounds"]; exists {
		if n, ok := totalRoundsAttr.(*types.AttributeValueMemberN); ok {
			if rounds, err := strconv.Atoi(n.Value); err == nil {
				stats.TotalRounds = rounds
			}
		}
	}

	// avg_focus_score
	if avgFocusAttr, exists := item["avg_focus_score"]; exists {
		if n, ok := avgFocusAttr.(*types.AttributeValueMemberN); ok {
			if score, err := strconv.ParseFloat(n.Value, 64); err == nil {
				stats.AvgFocusScore = score
			}
		}
	}

	// total_work_min
	if totalWorkAttr, exists := item["total_work_min"]; exists {
		if n, ok := totalWorkAttr.(*types.AttributeValueMemberN); ok {
			if workMin, err := strconv.Atoi(n.Value); err == nil {
				stats.TotalWorkMin = workMin
			}
		}
	}

	// total_break_min
	if totalBreakAttr, exists := item["total_break_min"]; exists {
		if n, ok := totalBreakAttr.(*types.AttributeValueMemberN); ok {
			if breakMin, err := strconv.Atoi(n.Value); err == nil {
				stats.TotalBreakMin = breakMin
			}
		}
	}

	// updated_at
	if updatedAttr, exists := item["updated_at"]; exists {
		if s, ok := updatedAttr.(*types.AttributeValueMemberS); ok {
			if updatedAt, err := time.Parse(time.RFC3339, s.Value); err == nil {
				stats.UpdatedAt = updatedAt
			}
		}
	}

	return stats, nil
}

// fillMissingDates は指定された期間内の日付を補完する
func (r *StatisticsRepositoryImpl) fillMissingDates(items []*model.FocusTrendItem, period *model.StatisticsPeriod) []*model.FocusTrendItem {
	// 既存データの日付をマップに格納
	dateMap := make(map[string]bool)
	for _, item := range items {
		dateMap[item.Date] = true
	}

	// 日付のフォーマット
	dateFormat := "2006-01-02"

	// 全期間の日付を作成
	var result []*model.FocusTrendItem
	result = append(result, items...)

	// 開始日から終了日まで1日ずつ増やしてチェック
	for d := period.StartDate; !d.After(period.EndDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format(dateFormat)
		if !dateMap[dateStr] {
			// データが存在しない場合は0値を設定
			result = append(result, &model.FocusTrendItem{
				Date:       dateStr,
				FocusScore: 0,
			})
		}
	}

	// 日付でソート
	result = r.sortByDate(result)

	return result
}

// fillMissingWeeks は指定された期間内の週を補完する
func (r *StatisticsRepositoryImpl) fillMissingWeeks(items []*model.FocusTrendItem, period *model.StatisticsPeriod) []*model.FocusTrendItem {
	// 既存データの週開始日をマップに格納
	weekMap := make(map[string]bool)
	for _, item := range items {
		weekMap[item.Date] = true
	}

	// 全期間の週を作成
	var result []*model.FocusTrendItem
	result = append(result, items...)

	// 期間内の全ての週を生成
	currentDate := period.StartDate
	for !currentDate.After(period.EndDate) {
		weekStart, _ := model.GetWeekBoundaries(currentDate)
		
		// この週のデータが存在しない場合は0値を設定
		if !weekMap[weekStart] {
			result = append(result, &model.FocusTrendItem{
				Date:       weekStart,
				FocusScore: 0,
			})
		}
		
		// 次の週へ
		currentDate = currentDate.AddDate(0, 0, 7)
	}

	// 週開始日でソート
	result = r.sortByDateWeeks(result)

	return result
}

// sortByDateWeeks は週開始日でソート
func (r *StatisticsRepositoryImpl) sortByDateWeeks(items []*model.FocusTrendItem) []*model.FocusTrendItem {
	// シンプルなソート実装
	for i := 0; i < len(items)-1; i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i].Date > items[j].Date {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	return items
}

// sortByDate は日付でソート
func (r *StatisticsRepositoryImpl) sortByDate(items []*model.FocusTrendItem) []*model.FocusTrendItem {
	// シンプルなソート実装
	for i := 0; i < len(items)-1; i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i].Date > items[j].Date {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	return items
}

// 統計データの更新・管理メソッド群

// UpdateDailyStats は日別統計を更新する
func (r *StatisticsRepositoryImpl) UpdateDailyStats(ctx context.Context, userID string, date string, round *model.Round) error {
	pk := UserPartitionKey(userID)
	sk := DailyStatsSortKey(date)

	// 現在の統計を取得
	existingStats, err := r.getDailyStats(ctx, pk, sk)
	if err != nil {
		// 新規作成
		existingStats = model.NewAggregatedDailyStats(userID, date)
	}

	// ラウンドデータで更新
	existingStats.UpdateWithRound(round)

	// DynamoDBに保存
	return r.putDailyStats(ctx, pk, sk, existingStats)
}

// UpdateHourlyStats は時間別統計を更新する
func (r *StatisticsRepositoryImpl) UpdateHourlyStats(ctx context.Context, userID string, date string, hour int, round *model.Round) error {
	pk := UserPartitionKey(userID)
	sk := HourlyStatsSortKey(date, hour)

	// 現在の統計を取得
	existingStats, err := r.getHourlyStats(ctx, pk, sk)
	if err != nil {
		// 新規作成
		existingStats = model.NewAggregatedHourlyStats(userID, date, hour)
	}

	// ラウンドデータで更新
	existingStats.UpdateWithRound(round)

	// DynamoDBに保存
	return r.putHourlyStats(ctx, pk, sk, existingStats)
}

// UpdateWeeklyStats は週別統計を更新する
func (r *StatisticsRepositoryImpl) UpdateWeeklyStats(ctx context.Context, userID string, date string, round *model.Round) error {
	// 週の境界を計算
	roundDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return fmt.Errorf("invalid date format: %s", date)
	}
	weekStart, weekEnd := model.GetWeekBoundaries(roundDate)
	
	pk := UserPartitionKey(userID)
	sk := WeeklyStatsSortKey(weekStart)

	// 現在の統計を取得
	existingStats, err := r.getWeeklyStats(ctx, pk, sk)
	if err != nil {
		// 新規作成
		existingStats = model.NewAggregatedWeeklyStats(userID, weekStart, weekEnd)
	}

	// ラウンドデータで更新
	existingStats.UpdateWithRound(round)

	// DynamoDBに保存
	return r.putWeeklyStats(ctx, pk, sk, existingStats)
}

// getDailyStats は日別統計を取得する
func (r *StatisticsRepositoryImpl) getDailyStats(ctx context.Context, pk, sk string) (*model.AggregatedDailyStats, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, fmt.Errorf("daily stats not found")
	}

	return r.itemToAggregatedDailyStats(result.Item)
}

// getHourlyStats は時間別統計を取得する
func (r *StatisticsRepositoryImpl) getHourlyStats(ctx context.Context, pk, sk string) (*model.AggregatedHourlyStats, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, fmt.Errorf("hourly stats not found")
	}

	return r.itemToAggregatedHourlyStats(result.Item)
}

// putDailyStats は日別統計を保存する
func (r *StatisticsRepositoryImpl) putDailyStats(ctx context.Context, pk, sk string, stats *model.AggregatedDailyStats) error {
	item := map[string]types.AttributeValue{
		"PK":               &types.AttributeValueMemberS{Value: pk},
		"SK":               &types.AttributeValueMemberS{Value: sk},
		"user_id":          &types.AttributeValueMemberS{Value: stats.UserID},
		"date":             &types.AttributeValueMemberS{Value: stats.Date},
		"total_rounds":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", stats.TotalRounds)},
		"avg_focus_score":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", stats.AvgFocusScore)},
		"total_work_min":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", stats.TotalWorkMin)},
		"total_break_min":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", stats.TotalBreakMin)},
		"session_count":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", stats.SessionCount)},
		"updated_at":       &types.AttributeValueMemberS{Value: stats.UpdatedAt.Format(time.RFC3339)},
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err := r.client.PutItem(ctx, input)
	return err
}

// putHourlyStats は時間別統計を保存する
func (r *StatisticsRepositoryImpl) putHourlyStats(ctx context.Context, pk, sk string, stats *model.AggregatedHourlyStats) error {
	item := map[string]types.AttributeValue{
		"PK":               &types.AttributeValueMemberS{Value: pk},
		"SK":               &types.AttributeValueMemberS{Value: sk},
		"user_id":          &types.AttributeValueMemberS{Value: stats.UserID},
		"date":             &types.AttributeValueMemberS{Value: stats.Date},
		"hour":             &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", stats.Hour)},
		"total_rounds":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", stats.TotalRounds)},
		"avg_focus_score":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", stats.AvgFocusScore)},
		"total_work_min":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", stats.TotalWorkMin)},
		"total_break_min":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", stats.TotalBreakMin)},
		"updated_at":       &types.AttributeValueMemberS{Value: stats.UpdatedAt.Format(time.RFC3339)},
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err := r.client.PutItem(ctx, input)
	return err
}

// getWeeklyStats は週別統計を取得する
func (r *StatisticsRepositoryImpl) getWeeklyStats(ctx context.Context, pk, sk string) (*model.AggregatedWeeklyStats, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, fmt.Errorf("weekly stats not found")
	}

	return r.itemToAggregatedWeeklyStats(result.Item)
}

// putWeeklyStats は週別統計を保存する
func (r *StatisticsRepositoryImpl) putWeeklyStats(ctx context.Context, pk, sk string, stats *model.AggregatedWeeklyStats) error {
	item := map[string]types.AttributeValue{
		"PK":               &types.AttributeValueMemberS{Value: pk},
		"SK":               &types.AttributeValueMemberS{Value: sk},
		"user_id":          &types.AttributeValueMemberS{Value: stats.UserID},
		"week_start":       &types.AttributeValueMemberS{Value: stats.WeekStart},
		"week_end":         &types.AttributeValueMemberS{Value: stats.WeekEnd},
		"total_rounds":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", stats.TotalRounds)},
		"avg_focus_score":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", stats.AvgFocusScore)},
		"total_work_min":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", stats.TotalWorkMin)},
		"total_break_min":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", stats.TotalBreakMin)},
		"session_count":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", stats.SessionCount)},
		"days_active":      &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", stats.DaysActive)},
		"updated_at":       &types.AttributeValueMemberS{Value: stats.UpdatedAt.Format(time.RFC3339)},
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err := r.client.PutItem(ctx, input)
	return err
}

// itemToAggregatedWeeklyStats はDynamoDBアイテムを週別統計に変換する
func (r *StatisticsRepositoryImpl) itemToAggregatedWeeklyStats(item map[string]types.AttributeValue) (*model.AggregatedWeeklyStats, error) {
	stats := &model.AggregatedWeeklyStats{}

	// user_id
	if userIDAttr, exists := item["user_id"]; exists {
		if s, ok := userIDAttr.(*types.AttributeValueMemberS); ok {
			stats.UserID = s.Value
		}
	}

	// week_start
	if weekStartAttr, exists := item["week_start"]; exists {
		if s, ok := weekStartAttr.(*types.AttributeValueMemberS); ok {
			stats.WeekStart = s.Value
		}
	}

	// week_end
	if weekEndAttr, exists := item["week_end"]; exists {
		if s, ok := weekEndAttr.(*types.AttributeValueMemberS); ok {
			stats.WeekEnd = s.Value
		}
	}

	// total_rounds
	if totalRoundsAttr, exists := item["total_rounds"]; exists {
		if n, ok := totalRoundsAttr.(*types.AttributeValueMemberN); ok {
			if rounds, err := strconv.Atoi(n.Value); err == nil {
				stats.TotalRounds = rounds
			}
		}
	}

	// avg_focus_score
	if avgFocusAttr, exists := item["avg_focus_score"]; exists {
		if n, ok := avgFocusAttr.(*types.AttributeValueMemberN); ok {
			if score, err := strconv.ParseFloat(n.Value, 64); err == nil {
				stats.AvgFocusScore = score
			}
		}
	}

	// total_work_min
	if totalWorkAttr, exists := item["total_work_min"]; exists {
		if n, ok := totalWorkAttr.(*types.AttributeValueMemberN); ok {
			if workMin, err := strconv.Atoi(n.Value); err == nil {
				stats.TotalWorkMin = workMin
			}
		}
	}

	// total_break_min
	if totalBreakAttr, exists := item["total_break_min"]; exists {
		if n, ok := totalBreakAttr.(*types.AttributeValueMemberN); ok {
			if breakMin, err := strconv.Atoi(n.Value); err == nil {
				stats.TotalBreakMin = breakMin
			}
		}
	}

	// session_count
	if sessionCountAttr, exists := item["session_count"]; exists {
		if n, ok := sessionCountAttr.(*types.AttributeValueMemberN); ok {
			if count, err := strconv.Atoi(n.Value); err == nil {
				stats.SessionCount = count
			}
		}
	}

	// days_active
	if daysActiveAttr, exists := item["days_active"]; exists {
		if n, ok := daysActiveAttr.(*types.AttributeValueMemberN); ok {
			if days, err := strconv.Atoi(n.Value); err == nil {
				stats.DaysActive = days
			}
		}
	}

	// updated_at
	if updatedAttr, exists := item["updated_at"]; exists {
		if s, ok := updatedAttr.(*types.AttributeValueMemberS); ok {
			if updatedAt, err := time.Parse(time.RFC3339, s.Value); err == nil {
				stats.UpdatedAt = updatedAt
			}
		}
	}

	return stats, nil
}