package model

import "time"

// FocusTrendItem は日付ごとの集中度を表す
type FocusTrendItem struct {
	Date       string  `json:"date"` // YYYY-MM-DD
	FocusScore float64 `json:"focus_score"`
}

// FocusTrendResponse は集中度のトレンドのレスポンス形式
type FocusTrendResponse struct {
	Items []*FocusTrendItem `json:"items"`
}

// FocusHeatmapItem は時間帯ごとの集中度を表す
type FocusHeatmapItem struct {
	Date       string  `json:"date"` // YYYY-MM-DD
	Hour       int     `json:"hour"` // 0-23
	FocusScore float64 `json:"focus_score"`
}

// FocusHeatmapResponse は集中度のヒートマップのレスポンス形式
type FocusHeatmapResponse struct {
	Items []*FocusHeatmapItem `json:"items"`
}

// StatisticsPeriod は統計情報の期間を表す
type StatisticsPeriod struct {
	StartDate time.Time
	EndDate   time.Time
}

// NewLastWeekPeriod は過去1週間の期間を生成する
func NewLastWeekPeriod() *StatisticsPeriod {
	now := time.Now()
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	startDate := endDate.AddDate(0, 0, -6) // 7日前の23:59:59
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

	return &StatisticsPeriod{
		StartDate: startDate,
		EndDate:   endDate,
	}
}

// NewLastMonthPeriod は過去1ヶ月の期間を生成する
func NewLastMonthPeriod() *StatisticsPeriod {
	now := time.Now()
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	startDate := endDate.AddDate(0, -1, 0) // 1ヶ月前の23:59:59
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

	return &StatisticsPeriod{
		StartDate: startDate,
		EndDate:   endDate,
	}
}

// NewCustomPeriod は指定された期間を生成する
func NewCustomPeriod(startDate, endDate time.Time) *StatisticsPeriod {
	// 日付の時刻部分を正規化(開始日は00:00:00、終了日は23:59:59)
	normalizedStartDate := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	normalizedEndDate := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	return &StatisticsPeriod{
		StartDate: normalizedStartDate,
		EndDate:   normalizedEndDate,
	}
}
