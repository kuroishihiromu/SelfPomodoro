package statistics

import (
	"time"
)

// StatisticsPeriod は統計情報の期間を表すValue Object
type StatisticsPeriod struct {
	startDate time.Time
	endDate   time.Time
}

// NewStatisticsPeriod は新しい統計期間を作成する
func NewStatisticsPeriod(startDate, endDate time.Time) StatisticsPeriod {
	// 日付の時刻部分を正規化(開始日は00:00:00、終了日は23:59:59)
	normalizedStartDate := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	normalizedEndDate := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	return StatisticsPeriod{
		startDate: normalizedStartDate,
		endDate:   normalizedEndDate,
	}
}

// NewLastWeekPeriod は過去1週間の期間を生成する（日別統計用）
func NewLastWeekPeriod() StatisticsPeriod {
	now := time.Now()
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	startDate := endDate.AddDate(0, 0, -6) // 7日前の23:59:59
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

	return StatisticsPeriod{
		startDate: startDate,
		endDate:   endDate,
	}
}

// NewCurrentWeekPeriod は現在の週の期間を生成する（月曜基準・週別統計用）
func NewCurrentWeekPeriod() StatisticsPeriod {
	now := time.Now()
	weekStart, weekEnd := GetWeekBoundaries(now)

	// 文字列から時刻に変換
	startDate, _ := time.Parse("2006-01-02", weekStart)
	endDate, _ := time.Parse("2006-01-02", weekEnd)

	// 時刻部分を適切に設定
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	return StatisticsPeriod{
		startDate: startDate,
		endDate:   endDate,
	}
}

// NewLastMonthPeriod は過去1ヶ月の期間を生成する
func NewLastMonthPeriod() StatisticsPeriod {
	now := time.Now()
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	startDate := endDate.AddDate(0, -1, 0) // 1ヶ月前の23:59:59
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

	return StatisticsPeriod{
		startDate: startDate,
		endDate:   endDate,
	}
}

// NewMonthlyHeatmapPeriod は指定された年月の月間ヒートマップ期間を生成する
func NewMonthlyHeatmapPeriod(year int, month int) StatisticsPeriod {
	// 月の最初の日
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)

	// 月の最後の日
	endDate := startDate.AddDate(0, 1, -1) // 翌月の1日から1日引く
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	return StatisticsPeriod{
		startDate: startDate,
		endDate:   endDate,
	}
}

// StartDate は開始日を返す
func (sp StatisticsPeriod) StartDate() time.Time {
	return sp.startDate
}

// EndDate は終了日を返す
func (sp StatisticsPeriod) EndDate() time.Time {
	return sp.endDate
}

// Duration は期間の長さを返す
func (sp StatisticsPeriod) Duration() time.Duration {
	return sp.endDate.Sub(sp.startDate)
}

// DurationDays は期間の日数を返す
func (sp StatisticsPeriod) DurationDays() int {
	return int(sp.Duration().Hours() / 24)
}

// Contains は指定した日時が期間内かを判定する
func (sp StatisticsPeriod) Contains(t time.Time) bool {
	return !t.Before(sp.startDate) && !t.After(sp.endDate)
}

// Equals は別のStatisticsPeriodと等価かを判定する
func (sp StatisticsPeriod) Equals(other StatisticsPeriod) bool {
	return sp.startDate.Equal(other.startDate) && sp.endDate.Equal(other.endDate)
}

// IsValid は有効な期間かを判定する（開始日 <= 終了日）
func (sp StatisticsPeriod) IsValid() bool {
	return !sp.startDate.After(sp.endDate)
}

// GetWeekBoundaries は指定された日付の週の境界を返す（月曜日開始）
func GetWeekBoundaries(date time.Time) (startDate, endDate string) {
	// 月曜日を週の開始とする
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7 // 日曜日を7として扱う
	}

	// 月曜日までの日数を計算
	daysToMonday := weekday - 1

	// 週の開始（月曜日）
	monday := date.AddDate(0, 0, -daysToMonday)
	monday = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())

	// 週の終了（日曜日）
	sunday := monday.AddDate(0, 0, 6)
	sunday = time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 23, 59, 59, 999999999, sunday.Location())

	return monday.Format("2006-01-02"), sunday.Format("2006-01-02")
}

// String は文字列表現を返す
func (sp StatisticsPeriod) String() string {
	return sp.startDate.Format("2006-01-02") + " to " + sp.endDate.Format("2006-01-02")
}