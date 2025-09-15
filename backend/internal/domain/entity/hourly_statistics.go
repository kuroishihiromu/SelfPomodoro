package entity

import (
	"time"

	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
	statisticsVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/statistics"
)

// HourlyStatistics は時間別統計を表すエンティティ
type HourlyStatistics struct {
	userID        userVO.UserID
	date          statisticsVO.Date
	hour          statisticsVO.Hour
	totalRounds   statisticsVO.TotalRounds
	avgFocusScore statisticsVO.AverageFocusScore
	totalWorkMin  statisticsVO.TotalWorkMinutes
	totalBreakMin statisticsVO.TotalBreakMinutes
	updatedAt     time.Time
}

// NewHourlyStatistics は新しい時間別統計を作成する
func NewHourlyStatistics(userID userVO.UserID, date string, hour int) *HourlyStatistics {
	dateVO, _ := statisticsVO.NewDate(date)
	hourVO, _ := statisticsVO.NewHour(hour)
	return &HourlyStatistics{
		userID:        userID,
		date:          dateVO,
		hour:          hourVO,
		totalRounds:   statisticsVO.NewZeroTotalRounds(),
		avgFocusScore: statisticsVO.NewZeroAverageFocusScore(),
		totalWorkMin:  statisticsVO.NewZeroTotalWorkMinutes(),
		totalBreakMin: statisticsVO.NewZeroTotalBreakMinutes(),
		updatedAt:     time.Now(),
	}
}

// NewHourlyStatisticsWithValues は既存の値で時間別統計を作成する
func NewHourlyStatisticsWithValues(
	userID userVO.UserID,
	date statisticsVO.Date,
	hour statisticsVO.Hour,
	totalRounds statisticsVO.TotalRounds,
	avgFocusScore statisticsVO.AverageFocusScore,
	totalWorkMin statisticsVO.TotalWorkMinutes,
	totalBreakMin statisticsVO.TotalBreakMinutes,
	updatedAt time.Time,
) *HourlyStatistics {
	return &HourlyStatistics{
		userID:        userID,
		date:          date,
		hour:          hour,
		totalRounds:   totalRounds,
		avgFocusScore: avgFocusScore,
		totalWorkMin:  totalWorkMin,
		totalBreakMin: totalBreakMin,
		updatedAt:     updatedAt,
	}
}

// UserID はユーザーIDを返す
func (hs *HourlyStatistics) UserID() userVO.UserID {
	return hs.userID
}

// Date は日付を返す
func (hs *HourlyStatistics) Date() statisticsVO.Date {
	return hs.date
}

// DateString は日付文字列を返す
func (hs *HourlyStatistics) DateString() string {
	return hs.date.String()
}

// Hour は時間を返す
func (hs *HourlyStatistics) Hour() statisticsVO.Hour {
	return hs.hour
}

// HourValue は時間の数値を返す
func (hs *HourlyStatistics) HourValue() int {
	return hs.hour.Hour()
}

// TotalRounds は総ラウンド数を返す
func (hs *HourlyStatistics) TotalRounds() statisticsVO.TotalRounds {
	return hs.totalRounds
}

// TotalRoundsCount は総ラウンド数の数値を返す
func (hs *HourlyStatistics) TotalRoundsCount() int {
	return hs.totalRounds.Count()
}

// AvgFocusScore は平均集中度スコアを返す
func (hs *HourlyStatistics) AvgFocusScore() statisticsVO.AverageFocusScore {
	return hs.avgFocusScore
}

// TotalWorkMin は総作業時間（分）を返す
func (hs *HourlyStatistics) TotalWorkMin() statisticsVO.TotalWorkMinutes {
	return hs.totalWorkMin
}

// TotalBreakMin は総休憩時間（分）を返す
func (hs *HourlyStatistics) TotalBreakMin() statisticsVO.TotalBreakMinutes {
	return hs.totalBreakMin
}

// UpdatedAt は更新日時を返す
func (hs *HourlyStatistics) UpdatedAt() time.Time {
	return hs.updatedAt
}

// UpdateWithRound はラウンドデータで統計を更新する
func (hs *HourlyStatistics) UpdateWithRound(round *Round) {
	if round.FocusScore == nil {
		return // スコアなしは統計から除外
	}

	// 平均集中度の再計算
	if newAvgFocus, err := hs.avgFocusScore.UpdateAverage(float64(round.FocusScore.Value()), hs.totalRounds.Count()); err == nil {
		hs.avgFocusScore = newAvgFocus
	}
	if newTotalRounds, err := hs.totalRounds.Increment(); err == nil {
		hs.totalRounds = newTotalRounds
	}

	// 作業時間・休憩時間の更新
	if round.WorkTime != nil {
		if newWorkMin, err := hs.totalWorkMin.AddMinutes(round.WorkTime.Minutes()); err == nil {
			hs.totalWorkMin = newWorkMin
		}
	}
	if round.BreakTime != nil {
		if newBreakMin, err := hs.totalBreakMin.AddMinutes(round.BreakTime.Minutes()); err == nil {
			hs.totalBreakMin = newBreakMin
		}
	}

	hs.updatedAt = time.Now()
}

