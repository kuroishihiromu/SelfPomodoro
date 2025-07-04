package model

import (
	"time"

	"github.com/google/uuid"
)

// RoundOptimizationLog はDynamoDBのround_optimization_logsテーブル用のドメインモデル
type RoundOptimizationLog struct {
	UserID     string `dynamodb:"user_id" json:"user_id"`         // PK
	Timestamp  string `dynamodb:"timestamp" json:"timestamp"`     // SK (ISO8601形式)
	WorkTime   int    `dynamodb:"work_time" json:"work_time"`     // 次ラウンドの作業時間（分）
	BreakTime  int    `dynamodb:"break_time" json:"break_time"`   // 次ラウンドの休憩時間（分）
	FocusScore int    `dynamodb:"focus_score" json:"focus_score"` // 入力スコア (0-100)
	CreatedAt  string `dynamodb:"created_at" json:"created_at"`   // 作成日時
}

// NewRoundOptimizationLog は新しいラウンド最適化ログを作成する
func NewRoundOptimizationLog(userID uuid.UUID, workTime, breakTime, focusScore int) *RoundOptimizationLog {
	now := time.Now()
	return &RoundOptimizationLog{
		UserID:     userID.String(),
		Timestamp:  now.Format(time.RFC3339),
		WorkTime:   workTime,
		BreakTime:  breakTime,
		FocusScore: focusScore,
		CreatedAt:  now.Format(time.RFC3339),
	}
}

// NewRoundOptimizationLogWithTime は指定時刻でラウンド最適化ログを作成する（サンプルデータ用）
func NewRoundOptimizationLogWithTime(userID uuid.UUID, timestamp time.Time, workTime, breakTime, focusScore int) *RoundOptimizationLog {
	return &RoundOptimizationLog{
		UserID:     userID.String(),
		Timestamp:  timestamp.Format(time.RFC3339),
		WorkTime:   workTime,
		BreakTime:  breakTime,
		FocusScore: focusScore,
		CreatedAt:  timestamp.Format(time.RFC3339),
	}
}

// SessionOptimizationLog はDynamoDBのsession_optimization_logsテーブル用のドメインモデル
type SessionOptimizationLog struct {
	UserID        string  `dynamodb:"user_id" json:"user_id"`                 // PK
	Timestamp     string  `dynamodb:"timestamp" json:"timestamp"`             // SK (ISO8601形式)
	RoundCount    int     `dynamodb:"round_count" json:"round_count"`         // 次セッションのラウンド数
	BreakTime     int     `dynamodb:"break_time" json:"break_time"`           // 次セッションの長休憩時間（分）
	AvgFocusScore float64 `dynamodb:"avg_focus_score" json:"avg_focus_score"` // 平均集中度スコア
	TotalWorkTime int     `dynamodb:"total_work_time" json:"total_work_time"` // 合計作業時間（分）
	CreatedAt     string  `dynamodb:"created_at" json:"created_at"`           // 作成日時
}

// NewSessionOptimizationLog は新しいセッション最適化ログを作成する
func NewSessionOptimizationLog(userID uuid.UUID, roundCount, breakTime int, avgFocusScore float64, totalWorkTime int) *SessionOptimizationLog {
	now := time.Now()
	return &SessionOptimizationLog{
		UserID:        userID.String(),
		Timestamp:     now.Format(time.RFC3339),
		RoundCount:    roundCount,
		BreakTime:     breakTime,
		AvgFocusScore: avgFocusScore,
		TotalWorkTime: totalWorkTime,
		CreatedAt:     now.Format(time.RFC3339),
	}
}

// NewSessionOptimizationLogWithTime は指定時刻でセッション最適化ログを作成する（サンプルデータ用）
func NewSessionOptimizationLogWithTime(userID uuid.UUID, timestamp time.Time, roundCount, breakTime int, avgFocusScore float64, totalWorkTime int) *SessionOptimizationLog {
	return &SessionOptimizationLog{
		UserID:        userID.String(),
		Timestamp:     timestamp.Format(time.RFC3339),
		RoundCount:    roundCount,
		BreakTime:     breakTime,
		AvgFocusScore: avgFocusScore,
		TotalWorkTime: totalWorkTime,
		CreatedAt:     timestamp.Format(time.RFC3339),
	}
}

// OptimizationEffectiveness は最適化の効果性を表すドメインモデル
type OptimizationEffectiveness struct {
	UserID                uuid.UUID `json:"user_id"`
	PeriodStart          time.Time `json:"period_start"`
	PeriodEnd            time.Time `json:"period_end"`
	RoundOptimizations   int       `json:"round_optimizations"`
	SessionOptimizations int       `json:"session_optimizations"`
	AvgFocusImprovement  float64   `json:"avg_focus_improvement"`
	OptimizationTrend    string    `json:"optimization_trend"` // "improving", "stable", "declining"
	LastOptimizedAt      time.Time `json:"last_optimized_at"`
}

// NewOptimizationEffectiveness は新しい最適化効果性を作成する
func NewOptimizationEffectiveness(userID uuid.UUID, periodStart, periodEnd time.Time) *OptimizationEffectiveness {
	return &OptimizationEffectiveness{
		UserID:      userID,
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		OptimizationTrend: "stable",
	}
}

// CalculateTrend は最適化トレンドを計算する
func (oe *OptimizationEffectiveness) CalculateTrend(recentImprovement, historicalImprovement float64) {
	diff := recentImprovement - historicalImprovement
	if diff > 2.0 {
		oe.OptimizationTrend = "improving"
	} else if diff < -2.0 {
		oe.OptimizationTrend = "declining"
	} else {
		oe.OptimizationTrend = "stable"
	}
}

// OptimizationSummary は最適化サマリーを表すドメインモデル
type OptimizationSummary struct {
	UserID                      uuid.UUID                `json:"user_id"`
	TotalRoundOptimizations     int                      `json:"total_round_optimizations"`
	TotalSessionOptimizations   int                      `json:"total_session_optimizations"`
	LatestRoundResult          *RoundOptimizationLog    `json:"latest_round_result,omitempty"`
	LatestSessionResult        *SessionOptimizationLog  `json:"latest_session_result,omitempty"`
	LastOptimizedAt            time.Time                `json:"last_optimized_at"`
	IsOptimizationActive       bool                     `json:"is_optimization_active"`
	CreatedAt                  time.Time                `json:"created_at"`
}

// NewOptimizationSummary は新しい最適化サマリーを作成する
func NewOptimizationSummary(userID uuid.UUID) *OptimizationSummary {
	return &OptimizationSummary{
		UserID:               userID,
		IsOptimizationActive: false,
		CreatedAt:            time.Now(),
	}
}

// HasRecentOptimization は最近の最適化があるかチェックする
func (os *OptimizationSummary) HasRecentOptimization(within time.Duration) bool {
	if os.LastOptimizedAt.IsZero() {
		return false
	}
	return time.Since(os.LastOptimizedAt) <= within
}
