package entity

import (
	"time"

	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
	statisticsVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/statistics"
)

// WeeklyStatistics は週別統計を表すエンティティ
type WeeklyStatistics struct {
	userID        userVO.UserID
	weekPeriod    statisticsVO.WeekPeriod
	totalRounds   statisticsVO.TotalRounds
	avgFocusScore statisticsVO.AverageFocusScore
	totalWorkMin  statisticsVO.TotalWorkMinutes
	totalBreakMin statisticsVO.TotalBreakMinutes
	sessionCount  statisticsVO.SessionCount
	daysActive    statisticsVO.DaysActive
	updatedAt     time.Time
}

// NewWeeklyStatistics は新しい週別統計を作成する
func NewWeeklyStatistics(userID userVO.UserID, weekStart, weekEnd string) *WeeklyStatistics {
	weekStartDate, _ := statisticsVO.NewDate(weekStart)
	weekEndDate, _ := statisticsVO.NewDate(weekEnd)
	weekPeriodVO, _ := statisticsVO.NewWeekPeriod(weekStartDate, weekEndDate)
	return &WeeklyStatistics{
		userID:        userID,
		weekPeriod:    weekPeriodVO,
		totalRounds:   statisticsVO.NewZeroTotalRounds(),
		avgFocusScore: statisticsVO.NewZeroAverageFocusScore(),
		totalWorkMin:  statisticsVO.NewZeroTotalWorkMinutes(),
		totalBreakMin: statisticsVO.NewZeroTotalBreakMinutes(),
		sessionCount:  statisticsVO.NewZeroSessionCount(),
		daysActive:    statisticsVO.NewZeroDaysActive(),
		updatedAt:     time.Now(),
	}
}

// NewWeeklyStatisticsWithValues は既存の値で週別統計を作成する
func NewWeeklyStatisticsWithValues(
	userID userVO.UserID,
	weekPeriod statisticsVO.WeekPeriod,
	totalRounds statisticsVO.TotalRounds,
	avgFocusScore statisticsVO.AverageFocusScore,
	totalWorkMin statisticsVO.TotalWorkMinutes,
	totalBreakMin statisticsVO.TotalBreakMinutes,
	sessionCount statisticsVO.SessionCount,
	daysActive statisticsVO.DaysActive,
	updatedAt time.Time,
) *WeeklyStatistics {
	return &WeeklyStatistics{
		userID:        userID,
		weekPeriod:    weekPeriod,
		totalRounds:   totalRounds,
		avgFocusScore: avgFocusScore,
		totalWorkMin:  totalWorkMin,
		totalBreakMin: totalBreakMin,
		sessionCount:  sessionCount,
		daysActive:    daysActive,
		updatedAt:     updatedAt,
	}
}

// UserID はユーザーIDを返す
func (ws *WeeklyStatistics) UserID() userVO.UserID {
	return ws.userID
}

// WeekStart は週の開始日を返す
func (ws *WeeklyStatistics) WeekStart() string {
	return ws.weekPeriod.WeekStartString()
}

// WeekEnd は週の終了日を返す
func (ws *WeeklyStatistics) WeekEnd() string {
	return ws.weekPeriod.WeekEndString()
}

// TotalRounds は総ラウンド数を返す
func (ws *WeeklyStatistics) TotalRounds() statisticsVO.TotalRounds {
	return ws.totalRounds
}

// TotalRoundsCount は総ラウンド数の数値を返す
func (ws *WeeklyStatistics) TotalRoundsCount() int {
	return ws.totalRounds.Count()
}

// AvgFocusScore は平均集中度スコアを返す
func (ws *WeeklyStatistics) AvgFocusScore() statisticsVO.AverageFocusScore {
	return ws.avgFocusScore
}

// TotalWorkMin は総作業時間（分）を返す
func (ws *WeeklyStatistics) TotalWorkMin() statisticsVO.TotalWorkMinutes {
	return ws.totalWorkMin
}

// TotalBreakMin は総休憩時間（分）を返す
func (ws *WeeklyStatistics) TotalBreakMin() statisticsVO.TotalBreakMinutes {
	return ws.totalBreakMin
}

// SessionCount はセッション数を返す
func (ws *WeeklyStatistics) SessionCount() statisticsVO.SessionCount {
	return ws.sessionCount
}

// SessionCountValue はセッション数の数値を返す
func (ws *WeeklyStatistics) SessionCountValue() int {
	return ws.sessionCount.Count()
}

// DaysActive はアクティブ日数を返す
func (ws *WeeklyStatistics) DaysActive() statisticsVO.DaysActive {
	return ws.daysActive
}

// DaysActiveValue はアクティブ日数の数値を返す
func (ws *WeeklyStatistics) DaysActiveValue() int {
	return ws.daysActive.Days()
}

// UpdatedAt は更新日時を返す
func (ws *WeeklyStatistics) UpdatedAt() time.Time {
	return ws.updatedAt
}

// UpdateWithRound はラウンドデータで統計を更新する
func (ws *WeeklyStatistics) UpdateWithRound(round *Round) {
	if round.FocusScore == nil {
		return // スコアなしは統計から除外
	}

	// 平均集中度の再計算
	if newAvgFocus, err := ws.avgFocusScore.UpdateAverage(float64(round.FocusScore.Value()), ws.totalRounds.Count()); err == nil {
		ws.avgFocusScore = newAvgFocus
	}
	if newTotalRounds, err := ws.totalRounds.Increment(); err == nil {
		ws.totalRounds = newTotalRounds
	}

	// 作業時間・休憩時間の更新
	if round.WorkTime != nil {
		if newWorkMin, err := ws.totalWorkMin.AddMinutes(round.WorkTime.Minutes()); err == nil {
			ws.totalWorkMin = newWorkMin
		}
	}
	if round.BreakTime != nil {
		if newBreakMin, err := ws.totalBreakMin.AddMinutes(round.BreakTime.Minutes()); err == nil {
			ws.totalBreakMin = newBreakMin
		}
	}

	ws.updatedAt = time.Now()
}

// IncrementSessionCount はセッション数をインクリメントする
func (ws *WeeklyStatistics) IncrementSessionCount() {
	if newSessionCount, err := ws.sessionCount.Increment(); err == nil {
		ws.sessionCount = newSessionCount
	}
	ws.updatedAt = time.Now()
}

// IncrementDaysActive はアクティブ日数をインクリメントする
func (ws *WeeklyStatistics) IncrementDaysActive() {
	if newDaysActive, err := ws.daysActive.Increment(); err == nil {
		ws.daysActive = newDaysActive
	}
	ws.updatedAt = time.Now()
}

// HasData はデータが存在するかを判定する
func (ws *WeeklyStatistics) HasData() bool {
	return !ws.totalRounds.IsZero() || !ws.sessionCount.IsZero()
}

// IsHighProductivityWeek は高い生産性の週かを判定する
func (ws *WeeklyStatistics) IsHighProductivityWeek() bool {
	return ws.avgFocusScore.IsProductiveWeek() && ws.daysActive.IsConsistent() && ws.totalWorkMin.IsProductiveWeek()
}

// GetWeeklyEfficiency は週の効率性を計算する
func (ws *WeeklyStatistics) GetWeeklyEfficiency() float64 {
	if ws.totalWorkMin.IsZero() {
		return 0
	}
	return ws.avgFocusScore.Score() / float64(ws.totalWorkMin.Minutes()) * 60 // 1時間あたりの集中度
}

// GetAverageSessionLength は平均セッション長を計算する
func (ws *WeeklyStatistics) GetAverageSessionLength() float64 {
	if ws.sessionCount.IsZero() {
		return 0
	}
	return float64(ws.totalWorkMin.Minutes()) / float64(ws.sessionCount.Count())
}

