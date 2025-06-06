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
