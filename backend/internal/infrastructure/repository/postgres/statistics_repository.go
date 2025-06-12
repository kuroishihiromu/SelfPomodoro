package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/database"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// StatisticsRepositoryImpl はStatisticsRepositoryインターフェースの実装（新エラーハンドリング対応版）
type StatisticsRepositoryImpl struct {
	db     *database.PostgresDB
	logger logger.Logger
}

// NewStatisticsRepository はStatisticsRepositoryImplの新しいインスタンスを作成する
func NewStatisticsRepository(db *database.PostgresDB, logger logger.Logger) repository.StatisticsRepository {
	return &StatisticsRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

// GetFocusTrend は指定期間内の日別集中度統計を取得する（新エラーハンドリング対応版）
func (r *StatisticsRepositoryImpl) GetFocusTrend(ctx context.Context, userID uuid.UUID, period *model.StatisticsPeriod) ([]*model.FocusTrendItem, error) {
	query := `
		SELECT
			TO_CHAR(r.start_time, 'YYYY-MM-DD') AS date,
			COALESCE(AVG(r.focus_score), 0) AS avg_focus_score
		FROM
			rounds r
		JOIN
			sessions s ON r.session_id = s.id
		WHERE
			s.user_id = $1
			AND r.start_time BETWEEN $2 AND $3
			AND r.focus_score IS NOT NULL
			AND r.is_aborted = FALSE
		GROUP BY
			TO_CHAR(r.start_time, 'YYYY-MM-DD')
		ORDER BY
			date ASC
	`

	type dbResult struct {
		Date          string  `db:"date"`
		AvgFocusScore float64 `db:"avg_focus_score"`
	}

	var results []dbResult
	err := r.db.DB.SelectContext(ctx, &results, query, userID, period.StartDate, period.EndDate)
	if err != nil {
		r.logger.Errorf("集中度トレンド取得エラー: %v", err)
		return nil, appErrors.NewDatabaseQueryError(err) // Infrastructure Error
	}

	// 結果をモデルに変換
	trendItems := make([]*model.FocusTrendItem, len(results))
	for i, result := range results {
		trendItems[i] = &model.FocusTrendItem{
			Date:       result.Date,
			FocusScore: result.AvgFocusScore,
		}
	}

	// データがない日付を補完
	trendItems = r.fillMissingDates(trendItems, period)

	r.logger.Infof("集中度トレンド取得成功: period=%s～%s, items=%d",
		period.StartDate.Format("2006-01-02"), period.EndDate.Format("2006-01-02"), len(trendItems))

	return trendItems, nil
}

// GetFocusHeatmap は指定期間内の時間帯別集中度統計を取得する（新エラーハンドリング対応版）
func (r *StatisticsRepositoryImpl) GetFocusHeatmap(ctx context.Context, userID uuid.UUID, period *model.StatisticsPeriod) ([]*model.FocusHeatmapItem, error) {
	query := `
		SELECT
			TO_CHAR(r.start_time, 'YYYY-MM-DD') AS date,
			EXTRACT(HOUR FROM r.start_time) AS hour,
			COALESCE(AVG(r.focus_score), 0) AS avg_focus_score
		FROM
			rounds r
		JOIN
			sessions s ON r.session_id = s.id
		WHERE
			s.user_id = $1
			AND r.start_time BETWEEN $2 AND $3
			AND r.focus_score IS NOT NULL
			AND r.is_aborted = FALSE
		GROUP BY
			date, hour
		ORDER BY
			date ASC, hour ASC
	`

	type dbResult struct {
		Date          string  `db:"date"`
		Hour          int     `db:"hour"`
		AvgFocusScore float64 `db:"avg_focus_score"`
	}

	var results []dbResult
	err := r.db.DB.SelectContext(ctx, &results, query, userID, period.StartDate, period.EndDate)
	if err != nil {
		r.logger.Errorf("集中度ヒートマップ取得エラー: %v", err)
		return nil, appErrors.NewDatabaseQueryError(err) // Infrastructure Error
	}

	// 結果をモデルに変換
	heatmapItems := make([]*model.FocusHeatmapItem, len(results))
	for i, result := range results {
		heatmapItems[i] = &model.FocusHeatmapItem{
			Date:       result.Date,
			Hour:       result.Hour,
			FocusScore: result.AvgFocusScore,
		}
	}

	r.logger.Infof("集中度ヒートマップ取得成功: period=%s～%s, items=%d",
		period.StartDate.Format("2006-01-02"), period.EndDate.Format("2006-01-02"), len(heatmapItems))

	return heatmapItems, nil
}

// GetAvgFocusScoreByDate は指定日の平均集中度を取得する（新エラーハンドリング対応版）
func (r *StatisticsRepositoryImpl) GetAvgFocusScoreByDate(ctx context.Context, userID uuid.UUID, date time.Time) (float64, error) {
	// 日付の範囲を設定
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())

	query := `
		SELECT
			COALESCE(AVG(r.focus_score), 0) AS avg_focus_score
		FROM
			rounds r
		JOIN
			sessions s ON r.session_id = s.id
		WHERE
			s.user_id = $1
			AND r.start_time BETWEEN $2 AND $3
			AND r.focus_score IS NOT NULL
			AND r.is_aborted = FALSE
	`

	var avgFocusScore float64
	err := r.db.DB.GetContext(ctx, &avgFocusScore, query, userID, startOfDay, endOfDay)
	if err != nil {
		r.logger.Errorf("指定日の平均集中度取得エラー: %v", err)
		return 0, appErrors.NewDatabaseQueryError(err) // Infrastructure Error
	}

	r.logger.Debugf("指定日の平均集中度取得成功: date=%s, avgFocus=%.1f",
		date.Format("2006-01-02"), avgFocusScore)

	return avgFocusScore, nil
}

// GetAvgFocusScoreByHour は指定日時の平均集中度を取得する（新エラーハンドリング対応版）
func (r *StatisticsRepositoryImpl) GetAvgFocusScoreByHour(ctx context.Context, userID uuid.UUID, date time.Time, hour int) (float64, error) {
	// 指定時間の範囲を設定
	startOfHour := time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, date.Location())
	endOfHour := time.Date(date.Year(), date.Month(), date.Day(), hour, 59, 59, 999999999, date.Location())

	query := `
		SELECT
			COALESCE(AVG(r.focus_score), 0) AS avg_focus_score
		FROM
			rounds r
		JOIN
			sessions s ON r.session_id = s.id
		WHERE
			s.user_id = $1
			AND r.start_time BETWEEN $2 AND $3
			AND r.focus_score IS NOT NULL
			AND r.is_aborted = FALSE
	`

	var avgFocusScore float64
	err := r.db.DB.GetContext(ctx, &avgFocusScore, query, userID, startOfHour, endOfHour)
	if err != nil {
		r.logger.Errorf("指定日時の平均集中度取得エラー: %v", err)
		return 0, appErrors.NewDatabaseQueryError(err) // Infrastructure Error
	}

	r.logger.Debugf("指定日時の平均集中度取得成功: date=%s %02d:00, avgFocus=%.1f",
		date.Format("2006-01-02"), hour, avgFocusScore)

	return avgFocusScore, nil
}

// fillMissingDates は指定された期間内の日付を補完する(0値)（ヘルパーメソッド）
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

// sortByDate は日付でソート（ヘルパーメソッド）
func (r *StatisticsRepositoryImpl) sortByDate(items []*model.FocusTrendItem) []*model.FocusTrendItem {
	// バブルソート（シンプルな実装）
	n := len(items)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if items[j].Date > items[j+1].Date {
				items[j], items[j+1] = items[j+1], items[j]
			}
		}
	}
	return items
}
