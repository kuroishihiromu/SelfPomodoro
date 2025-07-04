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

// NewLastWeekPeriod は過去1週間の期間を生成する（日別統計用）
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

// NewCurrentWeekPeriod は現在の週の期間を生成する（月曜基準・週別統計用）
func NewCurrentWeekPeriod() *StatisticsPeriod {
	now := time.Now()
	weekStart, weekEnd := GetWeekBoundaries(now)
	
	// 文字列から時刻に変換
	startDate, _ := time.Parse("2006-01-02", weekStart)
	endDate, _ := time.Parse("2006-01-02", weekEnd)
	
	// 時刻部分を適切に設定
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

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

// NewMonthlyHeatmapPeriod は指定された年月の月間ヒートマップ期間を生成する
func NewMonthlyHeatmapPeriod(year int, month int) *StatisticsPeriod {
	// 月の最初の日
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	
	// 月の最後の日
	endDate := startDate.AddDate(0, 1, -1) // 翌月の1日から1日引く
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	return &StatisticsPeriod{
		StartDate: startDate,
		EndDate:   endDate,
	}
}

// AggregatedDailyStats は日別事前集約統計を表すドメインモデル
type AggregatedDailyStats struct {
	UserID        string    `json:"user_id"`
	Date          string    `json:"date"`         // YYYY-MM-DD
	TotalRounds   int       `json:"total_rounds"`
	AvgFocusScore float64   `json:"avg_focus_score"`
	TotalWorkMin  int       `json:"total_work_min"`
	TotalBreakMin int       `json:"total_break_min"`
	SessionCount  int       `json:"session_count"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// AggregatedHourlyStats は時間別事前集約統計を表すドメインモデル
type AggregatedHourlyStats struct {
	UserID        string    `json:"user_id"`
	Date          string    `json:"date"`         // YYYY-MM-DD
	Hour          int       `json:"hour"`         // 0-23
	TotalRounds   int       `json:"total_rounds"`
	AvgFocusScore float64   `json:"avg_focus_score"`
	TotalWorkMin  int       `json:"total_work_min"`
	TotalBreakMin int       `json:"total_break_min"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// AggregatedWeeklyStats は週別事前集約統計を表すドメインモデル
type AggregatedWeeklyStats struct {
	UserID        string    `json:"user_id"`
	WeekStart     string    `json:"week_start"`   // YYYY-MM-DD (Monday of the week)
	WeekEnd       string    `json:"week_end"`     // YYYY-MM-DD (Sunday of the week)
	TotalRounds   int       `json:"total_rounds"`
	AvgFocusScore float64   `json:"avg_focus_score"`
	TotalWorkMin  int       `json:"total_work_min"`
	TotalBreakMin int       `json:"total_break_min"`
	SessionCount  int       `json:"session_count"`
	DaysActive    int       `json:"days_active"`  // Number of days with activity in the week
	UpdatedAt     time.Time `json:"updated_at"`
}

// Session統計計算のためのヘルパー構造体
type SessionStats struct {
	AverageFocus float64
	TotalWorkMin int
	RoundCount   int
	BreakTime    int
}

// NewAggregatedDailyStats は新しい日別統計を作成する
func NewAggregatedDailyStats(userID, date string) *AggregatedDailyStats {
	return &AggregatedDailyStats{
		UserID:        userID,
		Date:          date,
		TotalRounds:   0,
		AvgFocusScore: 0.0,
		TotalWorkMin:  0,
		TotalBreakMin: 0,
		SessionCount:  0,
		UpdatedAt:     time.Now(),
	}
}

// NewAggregatedHourlyStats は新しい時間別統計を作成する
func NewAggregatedHourlyStats(userID, date string, hour int) *AggregatedHourlyStats {
	return &AggregatedHourlyStats{
		UserID:        userID,
		Date:          date,
		Hour:          hour,
		TotalRounds:   0,
		AvgFocusScore: 0.0,
		TotalWorkMin:  0,
		TotalBreakMin: 0,
		UpdatedAt:     time.Now(),
	}
}

// NewAggregatedWeeklyStats は新しい週別統計を作成する
func NewAggregatedWeeklyStats(userID string, weekStart, weekEnd string) *AggregatedWeeklyStats {
	return &AggregatedWeeklyStats{
		UserID:        userID,
		WeekStart:     weekStart,
		WeekEnd:       weekEnd,
		TotalRounds:   0,
		AvgFocusScore: 0.0,
		TotalWorkMin:  0,
		TotalBreakMin: 0,
		SessionCount:  0,
		DaysActive:    0,
		UpdatedAt:     time.Now(),
	}
}

// ToFocusTrendItem は日別統計をFocusTrendItemに変換する
func (ds *AggregatedDailyStats) ToFocusTrendItem() *FocusTrendItem {
	return &FocusTrendItem{
		Date:       ds.Date,
		FocusScore: ds.AvgFocusScore,
	}
}

// ToFocusHeatmapItem は時間別統計をFocusHeatmapItemに変換する
func (hs *AggregatedHourlyStats) ToFocusHeatmapItem() *FocusHeatmapItem {
	return &FocusHeatmapItem{
		Date:       hs.Date,
		Hour:       hs.Hour,
		FocusScore: hs.AvgFocusScore,
	}
}

// ToFocusTrendItem は週別統計をFocusTrendItemに変換する
func (ws *AggregatedWeeklyStats) ToFocusTrendItem() *FocusTrendItem {
	return &FocusTrendItem{
		Date:       ws.WeekStart, // 週の開始日（月曜日）を表示
		FocusScore: ws.AvgFocusScore,
	}
}

// UpdateWithRound はラウンドデータで日別統計を更新する（ドメインRound用）
func (ds *AggregatedDailyStats) UpdateWithRound(round *Round) {
	if round.FocusScore == nil {
		return // スコアなしは統計から除外
	}

	// 平均集中度の再計算
	totalScore := ds.AvgFocusScore * float64(ds.TotalRounds)
	totalScore += float64(*round.FocusScore)
	ds.TotalRounds++
	ds.AvgFocusScore = totalScore / float64(ds.TotalRounds)

	// 作業時間・休憩時間の更新
	if round.WorkTime != nil {
		ds.TotalWorkMin += *round.WorkTime
	}
	if round.BreakTime != nil {
		ds.TotalBreakMin += *round.BreakTime
	}

	ds.UpdatedAt = time.Now()
}

// CalculateSessionStats はセッションの統計を計算する
// PostgreSQLのCalculateSessionStatsロジックを統合
func CalculateSessionStats(rounds []*Round) *SessionStats {
	if len(rounds) == 0 {
		return &SessionStats{}
	}

	var totalFocusScore float64
	var focusRoundCount int
	var totalWorkMin int
	var totalBreakMin int
	var completedRoundCount int

	for _, round := range rounds {
		if round.EndTime == nil {
			continue // 未完了ラウンドは除外
		}

		completedRoundCount++

		// 集中度の計算（フォーカススコアがある場合のみ）
		if round.FocusScore != nil {
			totalFocusScore += float64(*round.FocusScore)
			focusRoundCount++
		}

		// 作業時間・休憩時間の集計
		if round.WorkTime != nil {
			totalWorkMin += *round.WorkTime
		}
		if round.BreakTime != nil {
			totalBreakMin += *round.BreakTime
		}
	}

	// 平均集中度の計算
	var averageFocus float64
	if focusRoundCount > 0 {
		averageFocus = totalFocusScore / float64(focusRoundCount)
	}

	return &SessionStats{
		AverageFocus: averageFocus,
		TotalWorkMin: totalWorkMin,
		RoundCount:   completedRoundCount,
		BreakTime:    totalBreakMin,
	}
}

// UpdateWithRound はラウンドデータで時間別統計を更新する
func (hs *AggregatedHourlyStats) UpdateWithRound(round *Round) {
	if round.FocusScore == nil {
		return // スコアなしは統計から除外
	}

	// 平均集中度の再計算
	totalScore := hs.AvgFocusScore * float64(hs.TotalRounds)
	totalScore += float64(*round.FocusScore)
	hs.TotalRounds++
	hs.AvgFocusScore = totalScore / float64(hs.TotalRounds)

	// 作業時間・休憩時間の更新
	if round.WorkTime != nil {
		hs.TotalWorkMin += *round.WorkTime
	}
	if round.BreakTime != nil {
		hs.TotalBreakMin += *round.BreakTime
	}

	hs.UpdatedAt = time.Now()
}

// IncrementSessionCount はセッション数をインクリメントする
func (ds *AggregatedDailyStats) IncrementSessionCount() {
	ds.SessionCount++
	ds.UpdatedAt = time.Now()
}

// UpdateWithRound はラウンドデータで週別統計を更新する
func (ws *AggregatedWeeklyStats) UpdateWithRound(round *Round) {
	if round.FocusScore == nil {
		return // スコアなしは統計から除外
	}

	// 平均集中度の再計算
	totalScore := ws.AvgFocusScore * float64(ws.TotalRounds)
	totalScore += float64(*round.FocusScore)
	ws.TotalRounds++
	ws.AvgFocusScore = totalScore / float64(ws.TotalRounds)

	// 作業時間・休憩時間の更新
	if round.WorkTime != nil {
		ws.TotalWorkMin += *round.WorkTime
	}
	if round.BreakTime != nil {
		ws.TotalBreakMin += *round.BreakTime
	}

	ws.UpdatedAt = time.Now()
}

// IncrementSessionCount はセッション数をインクリメントする（週別）
func (ws *AggregatedWeeklyStats) IncrementSessionCount() {
	ws.SessionCount++
	ws.UpdatedAt = time.Now()
}

// IncrementDaysActive はアクティブ日数をインクリメントする
func (ws *AggregatedWeeklyStats) IncrementDaysActive() {
	ws.DaysActive++
	ws.UpdatedAt = time.Now()
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
