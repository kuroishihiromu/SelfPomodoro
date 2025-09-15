package entity

import (
	"time"

	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
	statisticsVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/statistics"
)

// DailyStatistics は日別統計を表すエンティティ
type DailyStatistics struct {
	userID        userVO.UserID
	date          statisticsVO.Date
	totalRounds   statisticsVO.TotalRounds
	avgFocusScore statisticsVO.AverageFocusScore
	totalWorkMin  statisticsVO.TotalWorkMinutes
	totalBreakMin statisticsVO.TotalBreakMinutes
	sessionCount  statisticsVO.SessionCount
	updatedAt     time.Time
}

// NewDailyStatistics は新しい日別統計を作成する
func NewDailyStatistics(userID userVO.UserID, date string) *DailyStatistics {
	dateVO, _ := statisticsVO.NewDate(date)
	return &DailyStatistics{
		userID:        userID,
		date:          dateVO,
		totalRounds:   statisticsVO.NewZeroTotalRounds(),
		avgFocusScore: statisticsVO.NewZeroAverageFocusScore(),
		totalWorkMin:  statisticsVO.NewZeroTotalWorkMinutes(),
		totalBreakMin: statisticsVO.NewZeroTotalBreakMinutes(),
		sessionCount:  statisticsVO.NewZeroSessionCount(),
		updatedAt:     time.Now(),
	}
}

// NewDailyStatisticsWithValues は既存の値で日別統計を作成する
func NewDailyStatisticsWithValues(
	userID userVO.UserID,
	date statisticsVO.Date,
	totalRounds statisticsVO.TotalRounds,
	avgFocusScore statisticsVO.AverageFocusScore,
	totalWorkMin statisticsVO.TotalWorkMinutes,
	totalBreakMin statisticsVO.TotalBreakMinutes,
	sessionCount statisticsVO.SessionCount,
	updatedAt time.Time,
) *DailyStatistics {
	return &DailyStatistics{
		userID:        userID,
		date:          date,
		totalRounds:   totalRounds,
		avgFocusScore: avgFocusScore,
		totalWorkMin:  totalWorkMin,
		totalBreakMin: totalBreakMin,
		sessionCount:  sessionCount,
		updatedAt:     updatedAt,
	}
}

// UserID はユーザーIDを返す
func (ds *DailyStatistics) UserID() userVO.UserID {
	return ds.userID
}

// Date は日付を返す
func (ds *DailyStatistics) Date() statisticsVO.Date {
	return ds.date
}

// DateString は日付文字列を返す
func (ds *DailyStatistics) DateString() string {
	return ds.date.String()
}

// TotalRounds は総ラウンド数を返す
func (ds *DailyStatistics) TotalRounds() statisticsVO.TotalRounds {
	return ds.totalRounds
}

// TotalRoundsCount は総ラウンド数の数値を返す
func (ds *DailyStatistics) TotalRoundsCount() int {
	return ds.totalRounds.Count()
}

// AvgFocusScore は平均集中度スコアを返す
func (ds *DailyStatistics) AvgFocusScore() statisticsVO.AverageFocusScore {
	return ds.avgFocusScore
}

// TotalWorkMin は総作業時間（分）を返す
func (ds *DailyStatistics) TotalWorkMin() statisticsVO.TotalWorkMinutes {
	return ds.totalWorkMin
}

// TotalBreakMin は総休憩時間（分）を返す
func (ds *DailyStatistics) TotalBreakMin() statisticsVO.TotalBreakMinutes {
	return ds.totalBreakMin
}

// SessionCount はセッション数を返す
func (ds *DailyStatistics) SessionCount() statisticsVO.SessionCount {
	return ds.sessionCount
}

// SessionCountValue はセッション数の数値を返す
func (ds *DailyStatistics) SessionCountValue() int {
	return ds.sessionCount.Count()
}

// UpdatedAt は更新日時を返す
func (ds *DailyStatistics) UpdatedAt() time.Time {
	return ds.updatedAt
}

// UpdateWithRound はラウンドデータで統計を更新する
func (ds *DailyStatistics) UpdateWithRound(round *Round) {
	if round.FocusScore == nil {
		return // スコアなしは統計から除外
	}

	// 平均集中度の再計算
	if newAvgFocus, err := ds.avgFocusScore.UpdateAverage(float64(round.FocusScore.Value()), ds.totalRounds.Count()); err == nil {
		ds.avgFocusScore = newAvgFocus
	}
	if newTotalRounds, err := ds.totalRounds.Increment(); err == nil {
		ds.totalRounds = newTotalRounds
	}

	// 作業時間・休憩時間の更新
	if round.WorkTime != nil {
		if newWorkMin, err := ds.totalWorkMin.AddMinutes(round.WorkTime.Minutes()); err == nil {
			ds.totalWorkMin = newWorkMin
		}
	}
	if round.BreakTime != nil {
		if newBreakMin, err := ds.totalBreakMin.AddMinutes(round.BreakTime.Minutes()); err == nil {
			ds.totalBreakMin = newBreakMin
		}
	}

	ds.updatedAt = time.Now()
}

// IncrementSessionCount はセッション数をインクリメントする
func (ds *DailyStatistics) IncrementSessionCount() {
	if newSessionCount, err := ds.sessionCount.Increment(); err == nil {
		ds.sessionCount = newSessionCount
	}
	ds.updatedAt = time.Now()
}

// HasData はデータが存在するかを判定する
func (ds *DailyStatistics) HasData() bool {
	return !ds.totalRounds.IsZero() || !ds.sessionCount.IsZero()
}

// IsHighProductivity は高い生産性の日かを判定する
func (ds *DailyStatistics) IsHighProductivity() bool {
	return ds.avgFocusScore.IsHigh() && ds.totalWorkMin.IsProductiveDay()
}

// GetEfficiency は効率性を計算する（集中度 / 時間）
func (ds *DailyStatistics) GetEfficiency() float64 {
	if ds.totalWorkMin.IsZero() {
		return 0
	}
	return ds.avgFocusScore.Score() / float64(ds.totalWorkMin.Minutes()) * 60 // 1時間あたりの集中度
}

