package dynamodb

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	statisticsVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/statistics"
	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
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
		tableName: cfg.DynamoUnifiedTable,
		logger:    logger,
	}
}

// ===============================================
// 日別統計
// ===============================================

// GetDailyStatistics は指定された日の日別統計を取得する
func (r *StatisticsRepositoryImpl) GetDailyStatistics(ctx context.Context, userID userVO.UserID, date string) (*entity.DailyStatistics, error) {
	pk := UserPartitionKey(userID.String())
	sk := DailyStatsSortKey(date)

	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		r.logger.Errorf("DailyStatistics取得エラー: %v", err)
		return nil, appErrors.NewInternalError(err)
	}

	if result.Item == nil {
		return nil, appErrors.NewNotFoundError("DailyStatistics")
	}

	return r.itemToDailyStatistics(result.Item, userID)
}

// SaveDailyStatistics は日別統計を保存する
func (r *StatisticsRepositoryImpl) SaveDailyStatistics(ctx context.Context, stats *entity.DailyStatistics) error {
	pk := UserPartitionKey(stats.UserID().String())
	sk := DailyStatsSortKey(stats.DateString())

	item := map[string]types.AttributeValue{
		"PK":              &types.AttributeValueMemberS{Value: pk},
		"SK":              &types.AttributeValueMemberS{Value: sk},
		"user_id":         &types.AttributeValueMemberS{Value: stats.UserID().String()},
		"date":            &types.AttributeValueMemberS{Value: stats.DateString()},
		"total_rounds":    &types.AttributeValueMemberN{Value: strconv.Itoa(stats.TotalRoundsCount())},
		"avg_focus_score": &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", stats.AvgFocusScore().Score())},
		"total_work_min":  &types.AttributeValueMemberN{Value: strconv.Itoa(stats.TotalWorkMin().Minutes())},
		"total_break_min": &types.AttributeValueMemberN{Value: strconv.Itoa(stats.TotalBreakMin().Minutes())},
		"session_count":   &types.AttributeValueMemberN{Value: strconv.Itoa(stats.SessionCountValue())},
		"updated_at":      &types.AttributeValueMemberS{Value: stats.UpdatedAt().Format(time.RFC3339)},
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err := r.client.PutItem(ctx, input)
	if err != nil {
		r.logger.Errorf("DailyStatistics保存エラー: %v", err)
		return appErrors.NewInternalError(err)
	}

	return nil
}

// GetDailyStatisticsByPeriod は指定期間の日別統計を取得する
func (r *StatisticsRepositoryImpl) GetDailyStatisticsByPeriod(ctx context.Context, userID userVO.UserID, period statisticsVO.StatisticsPeriod) ([]*entity.DailyStatistics, error) {
	pk := UserPartitionKey(userID.String())
	startSK, endSK := DailyStatsQueryRange(
		period.StartDate().Format("2006-01-02"),
		period.EndDate().Format("2006-01-02"),
	)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND SK BETWEEN :startSK AND :endSK"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":      &types.AttributeValueMemberS{Value: pk},
			":startSK": &types.AttributeValueMemberS{Value: startSK},
			":endSK":   &types.AttributeValueMemberS{Value: endSK},
		},
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		r.logger.Errorf("DailyStatistics期間検索エラー: %v", err)
		return nil, appErrors.NewInternalError(err)
	}

	var stats []*entity.DailyStatistics
	for _, item := range result.Items {
		stat, err := r.itemToDailyStatistics(item, userID)
		if err != nil {
			r.logger.Warnf("DailyStatistics変換エラー（スキップ）: %v", err)
			continue
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

// ===============================================
// 時間別統計
// ===============================================

// GetHourlyStatistics は指定された日時の時間別統計を取得する
func (r *StatisticsRepositoryImpl) GetHourlyStatistics(ctx context.Context, userID userVO.UserID, date string, hour int) (*entity.HourlyStatistics, error) {
	pk := UserPartitionKey(userID.String())
	sk := HourlyStatsSortKey(date, hour)

	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		r.logger.Errorf("HourlyStatistics取得エラー: %v", err)
		return nil, appErrors.NewInternalError(err)
	}

	if result.Item == nil {
		return nil, appErrors.NewNotFoundError("HourlyStatistics")
	}

	return r.itemToHourlyStatistics(result.Item, userID)
}

// SaveHourlyStatistics は時間別統計を保存する
func (r *StatisticsRepositoryImpl) SaveHourlyStatistics(ctx context.Context, stats *entity.HourlyStatistics) error {
	pk := UserPartitionKey(stats.UserID().String())
	sk := HourlyStatsSortKey(stats.DateString(), stats.HourValue())

	item := map[string]types.AttributeValue{
		"PK":              &types.AttributeValueMemberS{Value: pk},
		"SK":              &types.AttributeValueMemberS{Value: sk},
		"user_id":         &types.AttributeValueMemberS{Value: stats.UserID().String()},
		"date":            &types.AttributeValueMemberS{Value: stats.DateString()},
		"hour":            &types.AttributeValueMemberN{Value: strconv.Itoa(stats.HourValue())},
		"total_rounds":    &types.AttributeValueMemberN{Value: strconv.Itoa(stats.TotalRoundsCount())},
		"avg_focus_score": &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", stats.AvgFocusScore().Score())},
		"total_work_min":  &types.AttributeValueMemberN{Value: strconv.Itoa(stats.TotalWorkMin().Minutes())},
		"total_break_min": &types.AttributeValueMemberN{Value: strconv.Itoa(stats.TotalBreakMin().Minutes())},
		"updated_at":      &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err := r.client.PutItem(ctx, input)
	if err != nil {
		r.logger.Errorf("HourlyStatistics保存エラー: %v", err)
		return appErrors.NewInternalError(err)
	}

	return nil
}

// GetHourlyStatisticsByPeriod は指定期間の時間別統計を取得する
func (r *StatisticsRepositoryImpl) GetHourlyStatisticsByPeriod(ctx context.Context, userID userVO.UserID, period statisticsVO.StatisticsPeriod) ([]*entity.HourlyStatistics, error) {
	pk := UserPartitionKey(userID.String())
	startSK, endSK := HourlyStatsQueryRange(
		period.StartDate().Format("2006-01-02"),
		period.EndDate().Format("2006-01-02"),
	)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND SK BETWEEN :startSK AND :endSK"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":      &types.AttributeValueMemberS{Value: pk},
			":startSK": &types.AttributeValueMemberS{Value: startSK},
			":endSK":   &types.AttributeValueMemberS{Value: endSK},
		},
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		r.logger.Errorf("HourlyStatistics期間検索エラー: %v", err)
		return nil, appErrors.NewInternalError(err)
	}

	var stats []*entity.HourlyStatistics
	for _, item := range result.Items {
		stat, err := r.itemToHourlyStatistics(item, userID)
		if err != nil {
			r.logger.Warnf("HourlyStatistics変換エラー（スキップ）: %v", err)
			continue
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

// ===============================================
// 週別統計
// ===============================================

// GetWeeklyStatistics は指定された週の週別統計を取得する
func (r *StatisticsRepositoryImpl) GetWeeklyStatistics(ctx context.Context, userID userVO.UserID, weekStart, weekEnd string) (*entity.WeeklyStatistics, error) {
	pk := UserPartitionKey(userID.String())
	sk := WeeklyStatsSortKey(weekStart)

	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	}

	result, err := r.client.GetItem(ctx, input)
	if err != nil {
		r.logger.Errorf("WeeklyStatistics取得エラー: %v", err)
		return nil, appErrors.NewInternalError(err)
	}

	if result.Item == nil {
		return nil, appErrors.NewNotFoundError("WeeklyStatistics")
	}

	return r.itemToWeeklyStatistics(result.Item, userID)
}

// SaveWeeklyStatistics は週別統計を保存する
func (r *StatisticsRepositoryImpl) SaveWeeklyStatistics(ctx context.Context, stats *entity.WeeklyStatistics) error {
	pk := UserPartitionKey(stats.UserID().String())
	sk := WeeklyStatsSortKey(stats.WeekStart())

	item := map[string]types.AttributeValue{
		"PK":              &types.AttributeValueMemberS{Value: pk},
		"SK":              &types.AttributeValueMemberS{Value: sk},
		"user_id":         &types.AttributeValueMemberS{Value: stats.UserID().String()},
		"week_start":      &types.AttributeValueMemberS{Value: stats.WeekStart()},
		"week_end":        &types.AttributeValueMemberS{Value: stats.WeekEnd()},
		"total_rounds":    &types.AttributeValueMemberN{Value: strconv.Itoa(stats.TotalRoundsCount())},
		"avg_focus_score": &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", stats.AvgFocusScore().Score())},
		"total_work_min":  &types.AttributeValueMemberN{Value: strconv.Itoa(stats.TotalWorkMin().Minutes())},
		"total_break_min": &types.AttributeValueMemberN{Value: strconv.Itoa(stats.TotalBreakMin().Minutes())},
		"session_count":   &types.AttributeValueMemberN{Value: strconv.Itoa(stats.SessionCountValue())},
		"days_active":     &types.AttributeValueMemberN{Value: strconv.Itoa(stats.DaysActiveValue())},
		"updated_at":      &types.AttributeValueMemberS{Value: stats.UpdatedAt().Format(time.RFC3339)},
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err := r.client.PutItem(ctx, input)
	if err != nil {
		r.logger.Errorf("WeeklyStatistics保存エラー: %v", err)
		return appErrors.NewInternalError(err)
	}

	return nil
}

// GetWeeklyStatisticsByPeriod は指定期間の週別統計を取得する
func (r *StatisticsRepositoryImpl) GetWeeklyStatisticsByPeriod(ctx context.Context, userID userVO.UserID, period statisticsVO.StatisticsPeriod) ([]*entity.WeeklyStatistics, error) {
	pk := UserPartitionKey(userID.String())
	startWeek := period.StartDate().Format("2006-01-02")
	endWeek := period.EndDate().Format("2006-01-02")
	startSK, endSK := WeeklyStatsQueryRange(startWeek, endWeek)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND SK BETWEEN :startSK AND :endSK"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":      &types.AttributeValueMemberS{Value: pk},
			":startSK": &types.AttributeValueMemberS{Value: startSK},
			":endSK":   &types.AttributeValueMemberS{Value: endSK},
		},
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		r.logger.Errorf("WeeklyStatistics期間検索エラー: %v", err)
		return nil, appErrors.NewInternalError(err)
	}

	var stats []*entity.WeeklyStatistics
	for _, item := range result.Items {
		stat, err := r.itemToWeeklyStatistics(item, userID)
		if err != nil {
			r.logger.Warnf("WeeklyStatistics変換エラー（スキップ）: %v", err)
			continue
		}
		stats = append(stats, stat)
	}

	return stats, nil
}

// ===============================================
// 集約用クエリは削除済み - UseCase層でMapper使用に変更
// ===============================================

// ===============================================
// 統計更新処理
// ===============================================

// UpdateStatisticsWithRound はラウンド完了時に統計を更新する
func (r *StatisticsRepositoryImpl) UpdateStatisticsWithRound(ctx context.Context, userID userVO.UserID, round *entity.Round, timestamp time.Time) error {
	date := timestamp.Format("2006-01-02")
	hour := timestamp.Hour()

	// 日別統計の更新
	if err := r.updateDailyStatisticsWithRound(ctx, userID, date, round); err != nil {
		r.logger.Errorf("日別統計更新エラー: %v", err)
		return err
	}

	// 時間別統計の更新
	if err := r.updateHourlyStatisticsWithRound(ctx, userID, date, hour, round); err != nil {
		r.logger.Errorf("時間別統計更新エラー: %v", err)
		return err
	}

	// 週別統計の更新
	if err := r.updateWeeklyStatisticsWithRound(ctx, userID, date, round); err != nil {
		r.logger.Errorf("週別統計更新エラー: %v", err)
		return err
	}

	return nil
}

// ===============================================
// プライベートヘルパーメソッド
// ===============================================

// updateDailyStatisticsWithRound は日別統計を更新する
func (r *StatisticsRepositoryImpl) updateDailyStatisticsWithRound(ctx context.Context, userID userVO.UserID, date string, round *entity.Round) error {
	// 既存の日別統計を取得（存在しない場合は新規作成）
	stats, err := r.GetDailyStatistics(ctx, userID, date)
	if err != nil {
		if appErrors.IsNotFoundError(err) {
			stats = entity.NewDailyStatistics(userID, date)
		} else {
			return err
		}
	}

	// ラウンドデータで統計を更新
	stats.UpdateWithRound(round)

	// 統計を保存
	return r.SaveDailyStatistics(ctx, stats)
}

// updateHourlyStatisticsWithRound は時間別統計を更新する
func (r *StatisticsRepositoryImpl) updateHourlyStatisticsWithRound(ctx context.Context, userID userVO.UserID, date string, hour int, round *entity.Round) error {
	// 既存の時間別統計を取得（存在しない場合は新規作成）
	stats, err := r.GetHourlyStatistics(ctx, userID, date, hour)
	if err != nil {
		if appErrors.IsNotFoundError(err) {
			stats = entity.NewHourlyStatistics(userID, date, hour)
		} else {
			return err
		}
	}

	// ラウンドデータで統計を更新
	stats.UpdateWithRound(round)

	// 統計を保存
	return r.SaveHourlyStatistics(ctx, stats)
}

// updateWeeklyStatisticsWithRound は週別統計を更新する
func (r *StatisticsRepositoryImpl) updateWeeklyStatisticsWithRound(ctx context.Context, userID userVO.UserID, date string, round *entity.Round) error {
	dateVO, err := statisticsVO.NewDate(date)
	if err != nil {
		return appErrors.NewInternalError(err)
	}

	weekPeriod := statisticsVO.NewWeekPeriodFromDate(dateVO)
	weekStart := weekPeriod.WeekStartString()
	weekEnd := weekPeriod.WeekEndString()

	stats, err := r.GetWeeklyStatistics(ctx, userID, weekStart, weekEnd)
	if err != nil {
		if appErrors.IsNotFoundError(err) {
			stats = entity.NewWeeklyStatistics(userID, weekStart, weekEnd)
		} else {
			return err
		}
	}

	stats.UpdateWithRound(round)

	return r.SaveWeeklyStatistics(ctx, stats)
}

// itemToDailyStatistics はDynamoDBアイテムを日別統計エンティティに変換する
func (r *StatisticsRepositoryImpl) itemToDailyStatistics(item map[string]types.AttributeValue, userID userVO.UserID) (*entity.DailyStatistics, error) {
	dateValue, ok := item["date"].(*types.AttributeValueMemberS)
	if !ok {
		return nil, fmt.Errorf("date属性が見つかりません")
	}

	// 基本の統計オブジェクトを作成
	date, err := statisticsVO.NewDate(dateValue.Value)
	if err != nil {
		return nil, fmt.Errorf("date値の変換エラー: %v", err)
	}

	// Value Objectsを作成
	totalRounds := statisticsVO.NewZeroTotalRounds()
	if totalRoundsValue, ok := item["total_rounds"].(*types.AttributeValueMemberN); ok {
		if count, err := strconv.Atoi(totalRoundsValue.Value); err == nil {
			if tr, err := statisticsVO.NewTotalRounds(count); err == nil {
				totalRounds = tr
			}
		}
	}

	avgFocusScore := statisticsVO.NewZeroAverageFocusScore()
	if avgFocusValue, ok := item["avg_focus_score"].(*types.AttributeValueMemberN); ok {
		if score, err := strconv.ParseFloat(avgFocusValue.Value, 64); err == nil {
			if afs, err := statisticsVO.NewAverageFocusScore(score); err == nil {
				avgFocusScore = afs
			}
		}
	}

	totalWorkMin := statisticsVO.NewZeroTotalWorkMinutes()
	if totalWorkValue, ok := item["total_work_min"].(*types.AttributeValueMemberN); ok {
		if minutes, err := strconv.Atoi(totalWorkValue.Value); err == nil {
			if twm, err := statisticsVO.NewTotalWorkMinutes(minutes); err == nil {
				totalWorkMin = twm
			}
		}
	}

	totalBreakMin := statisticsVO.NewZeroTotalBreakMinutes()
	if totalBreakValue, ok := item["total_break_min"].(*types.AttributeValueMemberN); ok {
		if minutes, err := strconv.Atoi(totalBreakValue.Value); err == nil {
			if tbm, err := statisticsVO.NewTotalBreakMinutes(minutes); err == nil {
				totalBreakMin = tbm
			}
		}
	}

	sessionCount := statisticsVO.NewZeroSessionCount()
	if sessionCountValue, ok := item["session_count"].(*types.AttributeValueMemberN); ok {
		if count, err := strconv.Atoi(sessionCountValue.Value); err == nil {
			if sc, err := statisticsVO.NewSessionCount(count); err == nil {
				sessionCount = sc
			}
		}
	}

	updatedAt := time.Now()
	if updatedAtValue, ok := item["updated_at"].(*types.AttributeValueMemberS); ok {
		if t, err := time.Parse(time.RFC3339, updatedAtValue.Value); err == nil {
			updatedAt = t
		}
	}

	return entity.NewDailyStatisticsWithValues(userID, date, totalRounds, avgFocusScore, totalWorkMin, totalBreakMin, sessionCount, updatedAt), nil
}

// itemToHourlyStatistics はDynamoDBアイテムを時間別統計エンティティに変換する
func (r *StatisticsRepositoryImpl) itemToHourlyStatistics(item map[string]types.AttributeValue, userID userVO.UserID) (*entity.HourlyStatistics, error) {
	dateValue, ok := item["date"].(*types.AttributeValueMemberS)
	if !ok {
		return nil, fmt.Errorf("date属性が見つかりません")
	}

	hourValue, ok := item["hour"].(*types.AttributeValueMemberN)
	if !ok {
		return nil, fmt.Errorf("hour属性が見つかりません")
	}

	hour, err := strconv.Atoi(hourValue.Value)
	if err != nil {
		return nil, fmt.Errorf("hour値の変換エラー: %v", err)
	}

	// Value Objectsを作成
	date, err := statisticsVO.NewDate(dateValue.Value)
	if err != nil {
		return nil, fmt.Errorf("date値の変換エラー: %v", err)
	}

	hourVO, err := statisticsVO.NewHour(hour)
	if err != nil {
		return nil, fmt.Errorf("hour値の変換エラー: %v", err)
	}

	totalRounds := statisticsVO.NewZeroTotalRounds()
	if totalRoundsValue, ok := item["total_rounds"].(*types.AttributeValueMemberN); ok {
		if count, err := strconv.Atoi(totalRoundsValue.Value); err == nil {
			if tr, err := statisticsVO.NewTotalRounds(count); err == nil {
				totalRounds = tr
			}
		}
	}

	avgFocusScore := statisticsVO.NewZeroAverageFocusScore()
	if avgFocusValue, ok := item["avg_focus_score"].(*types.AttributeValueMemberN); ok {
		if score, err := strconv.ParseFloat(avgFocusValue.Value, 64); err == nil {
			if afs, err := statisticsVO.NewAverageFocusScore(score); err == nil {
				avgFocusScore = afs
			}
		}
	}

	totalWorkMin := statisticsVO.NewZeroTotalWorkMinutes()
	if totalWorkValue, ok := item["total_work_min"].(*types.AttributeValueMemberN); ok {
		if minutes, err := strconv.Atoi(totalWorkValue.Value); err == nil {
			if twm, err := statisticsVO.NewTotalWorkMinutes(minutes); err == nil {
				totalWorkMin = twm
			}
		}
	}

	totalBreakMin := statisticsVO.NewZeroTotalBreakMinutes()
	if totalBreakValue, ok := item["total_break_min"].(*types.AttributeValueMemberN); ok {
		if minutes, err := strconv.Atoi(totalBreakValue.Value); err == nil {
			if tbm, err := statisticsVO.NewTotalBreakMinutes(minutes); err == nil {
				totalBreakMin = tbm
			}
		}
	}

	updatedAt := time.Now()
	if updatedAtValue, ok := item["updated_at"].(*types.AttributeValueMemberS); ok {
		if t, err := time.Parse(time.RFC3339, updatedAtValue.Value); err == nil {
			updatedAt = t
		}
	}

	return entity.NewHourlyStatisticsWithValues(userID, date, hourVO, totalRounds, avgFocusScore, totalWorkMin, totalBreakMin, updatedAt), nil
}

// itemToWeeklyStatistics はDynamoDBアイテムを週別統計エンティティに変換する
func (r *StatisticsRepositoryImpl) itemToWeeklyStatistics(item map[string]types.AttributeValue, userID userVO.UserID) (*entity.WeeklyStatistics, error) {
	weekStartValue, ok := item["week_start"].(*types.AttributeValueMemberS)
	if !ok {
		return nil, fmt.Errorf("week_start属性が見つかりません")
	}

	weekEndValue, ok := item["week_end"].(*types.AttributeValueMemberS)
	if !ok {
		return nil, fmt.Errorf("week_end属性が見つかりません")
	}

	// WeekPeriod Value Objectを作成
	weekStartDate, err := statisticsVO.NewDate(weekStartValue.Value)
	if err != nil {
		return nil, fmt.Errorf("week_start値の変換エラー: %v", err)
	}

	weekEndDate, err := statisticsVO.NewDate(weekEndValue.Value)
	if err != nil {
		return nil, fmt.Errorf("week_end値の変換エラー: %v", err)
	}

	weekPeriod, err := statisticsVO.NewWeekPeriod(weekStartDate, weekEndDate)
	if err != nil {
		return nil, fmt.Errorf("week_period値の変換エラー: %v", err)
	}

	// 他のValue Objectsを作成
	totalRounds := statisticsVO.NewZeroTotalRounds()
	if totalRoundsValue, ok := item["total_rounds"].(*types.AttributeValueMemberN); ok {
		if count, err := strconv.Atoi(totalRoundsValue.Value); err == nil {
			if tr, err := statisticsVO.NewTotalRounds(count); err == nil {
				totalRounds = tr
			}
		}
	}

	avgFocusScore := statisticsVO.NewZeroAverageFocusScore()
	if avgFocusValue, ok := item["avg_focus_score"].(*types.AttributeValueMemberN); ok {
		if score, err := strconv.ParseFloat(avgFocusValue.Value, 64); err == nil {
			if afs, err := statisticsVO.NewAverageFocusScore(score); err == nil {
				avgFocusScore = afs
			}
		}
	}

	totalWorkMin := statisticsVO.NewZeroTotalWorkMinutes()
	if totalWorkValue, ok := item["total_work_min"].(*types.AttributeValueMemberN); ok {
		if minutes, err := strconv.Atoi(totalWorkValue.Value); err == nil {
			if twm, err := statisticsVO.NewTotalWorkMinutes(minutes); err == nil {
				totalWorkMin = twm
			}
		}
	}

	totalBreakMin := statisticsVO.NewZeroTotalBreakMinutes()
	if totalBreakValue, ok := item["total_break_min"].(*types.AttributeValueMemberN); ok {
		if minutes, err := strconv.Atoi(totalBreakValue.Value); err == nil {
			if tbm, err := statisticsVO.NewTotalBreakMinutes(minutes); err == nil {
				totalBreakMin = tbm
			}
		}
	}

	sessionCount := statisticsVO.NewZeroSessionCount()
	if sessionCountValue, ok := item["session_count"].(*types.AttributeValueMemberN); ok {
		if count, err := strconv.Atoi(sessionCountValue.Value); err == nil {
			if sc, err := statisticsVO.NewSessionCount(count); err == nil {
				sessionCount = sc
			}
		}
	}

	daysActive := statisticsVO.NewZeroDaysActive()
	if daysActiveValue, ok := item["days_active"].(*types.AttributeValueMemberN); ok {
		if days, err := strconv.Atoi(daysActiveValue.Value); err == nil {
			if da, err := statisticsVO.NewDaysActive(days); err == nil {
				daysActive = da
			}
		}
	}

	updatedAt := time.Now()
	if updatedAtValue, ok := item["updated_at"].(*types.AttributeValueMemberS); ok {
		if t, err := time.Parse(time.RFC3339, updatedAtValue.Value); err == nil {
			updatedAt = t
		}
	}

	return entity.NewWeeklyStatisticsWithValues(userID, weekPeriod, totalRounds, avgFocusScore, totalWorkMin, totalBreakMin, sessionCount, daysActive, updatedAt), nil
}
