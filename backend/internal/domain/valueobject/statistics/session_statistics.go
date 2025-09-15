package statistics

import (
	"fmt"
)

// SessionStatistics はセッション統計を表すValue Object
type SessionStatistics struct {
	averageFocus float64
	totalWorkMin int
	roundCount   int
	breakTime    int
}

// NewSessionStatistics は新しいセッション統計を作成する
func NewSessionStatistics(averageFocus float64, totalWorkMin int, roundCount int, breakTime int) SessionStatistics {
	return SessionStatistics{
		averageFocus: averageFocus,
		totalWorkMin: totalWorkMin,
		roundCount:   roundCount,
		breakTime:    breakTime,
	}
}

// AverageFocus は平均集中度を返す
func (ss SessionStatistics) AverageFocus() float64 {
	return ss.averageFocus
}

// TotalWorkMin は合計作業時間（分）を返す
func (ss SessionStatistics) TotalWorkMin() int {
	return ss.totalWorkMin
}

// RoundCount はラウンド数を返す
func (ss SessionStatistics) RoundCount() int {
	return ss.roundCount
}

// BreakTime は休憩時間を返す
func (ss SessionStatistics) BreakTime() int {
	return ss.breakTime
}

// Equals は別のSessionStatisticsと等価かを判定する
func (ss SessionStatistics) Equals(other SessionStatistics) bool {
	return ss.averageFocus == other.averageFocus &&
		ss.totalWorkMin == other.totalWorkMin &&
		ss.roundCount == other.roundCount &&
		ss.breakTime == other.breakTime
}

// IsHighPerformance は高パフォーマンスセッションかを判定する
func (ss SessionStatistics) IsHighPerformance() bool {
	return ss.averageFocus >= 80 && ss.totalWorkMin >= 60 // 1時間以上の高集中セッション
}

// IsValidSession は有効なセッションかを判定する
func (ss SessionStatistics) IsValidSession() bool {
	return ss.roundCount > 0 && ss.totalWorkMin > 0
}

// GetEfficiency は効率性を計算する（集中度 / 時間）
func (ss SessionStatistics) GetEfficiency() float64 {
	if ss.totalWorkMin == 0 {
		return 0
	}
	return ss.averageFocus / float64(ss.totalWorkMin) * 60 // 1時間あたりの集中度
}

// GetAverageWorkTimePerRound は1ラウンドあたりの平均作業時間を返す
func (ss SessionStatistics) GetAverageWorkTimePerRound() float64 {
	if ss.roundCount == 0 {
		return 0
	}
	return float64(ss.totalWorkMin) / float64(ss.roundCount)
}

// GetQualityLevel は品質レベルを返す
func (ss SessionStatistics) GetQualityLevel() string {
	if ss.averageFocus >= 80 {
		return "高品質"
	} else if ss.averageFocus >= 60 {
		return "中品質"
	} else {
		return "低品質"
	}
}

// String は文字列表現を返す
func (ss SessionStatistics) String() string {
	return fmt.Sprintf("SessionStats[AvgFocus=%.1f, TotalWork=%dmin, Rounds=%d, Break=%dmin]",
		ss.averageFocus, ss.totalWorkMin, ss.roundCount, ss.breakTime)
}