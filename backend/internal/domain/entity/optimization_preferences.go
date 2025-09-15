package entity

import (
	"fmt"
	"time"

	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
	roundVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/round"
	sessionVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/session"
)

// デフォルト値定数（ドメインルール）
const (
	DefaultOptWorkTime         = 25 // デフォルト作業時間（分）
	DefaultOptBreakTime        = 5  // デフォルト休憩時間（分）
	DefaultOptSessionBreakTime = 15 // デフォルト長休憩時間（分）

	// バリデーション制約
	MinOptWorkTime         = 1
	MaxOptWorkTime         = 120
	MinOptBreakTime        = 1
	MaxOptBreakTime        = 60
	MinOptSessionBreakTime = 5
	MaxOptSessionBreakTime = 120
)

// OptimizationPreferences はユーザーの最適化された設定を表すエンティティ
type OptimizationPreferences struct {
	UserID           userVO.UserID        `dynamodb:"user_id" json:"user_id"`
	RoundWorkTime    roundVO.WorkTime     `dynamodb:"round_work_time" json:"round_work_time"`       // 最適化された作業時間
	RoundBreakTime   roundVO.BreakTime    `dynamodb:"round_break_time" json:"round_break_time"`     // 最適化された休憩時間
	SessionRounds    sessionVO.SessionRounds `dynamodb:"session_rounds" json:"session_rounds"`         // 最適化されたセッション内ラウンド数
	SessionBreakTime sessionVO.BreakTime  `dynamodb:"session_break_time" json:"session_break_time"` // 最適化されたセッション後長休憩
	CreatedAt        time.Time            `dynamodb:"created_at" json:"created_at"`
	UpdatedAt        time.Time            `dynamodb:"updated_at" json:"updated_at"`
}

// NewOptimizationPreferences は新しい最適化設定を作成する（デフォルト値付き）
func NewOptimizationPreferences(userID userVO.UserID) (*OptimizationPreferences, error) {
	workTime, err := roundVO.NewWorkTime(DefaultOptWorkTime)
	if err != nil {
		return nil, fmt.Errorf("デフォルト作業時間の作成に失敗: %w", err)
	}
	
	breakTime, err := roundVO.NewBreakTime(DefaultOptBreakTime)
	if err != nil {
		return nil, fmt.Errorf("デフォルト休憩時間の作成に失敗: %w", err)
	}
	
	sessionBreakTime, err := sessionVO.NewBreakTime(DefaultOptSessionBreakTime)
	if err != nil {
		return nil, fmt.Errorf("デフォルトセッション休憩時間の作成に失敗: %w", err)
	}

	now := time.Now()
	return &OptimizationPreferences{
		UserID:           userID,
		RoundWorkTime:    workTime,
		RoundBreakTime:   breakTime,
		SessionRounds:    sessionVO.NewDefaultSessionRounds(),
		SessionBreakTime: sessionBreakTime,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// GetSessionRoundsOrDefault はセッションラウンド数またはデフォルト値を返す
func (op *OptimizationPreferences) GetSessionRoundsOrDefault() sessionVO.SessionRounds {
	if op == nil {
		return sessionVO.NewDefaultSessionRounds()
	}
	return op.SessionRounds
}

// GetWorkTimeOrDefault は作業時間またはデフォルト値を返す
func (op *OptimizationPreferences) GetWorkTimeOrDefault() roundVO.WorkTime {
	if op == nil {
		workTime, _ := roundVO.NewWorkTime(DefaultOptWorkTime)
		return workTime
	}
	return op.RoundWorkTime
}

// GetBreakTimeOrDefault は休憩時間またはデフォルト値を返す
func (op *OptimizationPreferences) GetBreakTimeOrDefault() roundVO.BreakTime {
	if op == nil {
		breakTime, _ := roundVO.NewBreakTime(DefaultOptBreakTime)
		return breakTime
	}
	return op.RoundBreakTime
}

// UpdateSettings は設定を更新する（ドメインルール適用）
func (op *OptimizationPreferences) UpdateSettings(workTime roundVO.WorkTime, breakTime roundVO.BreakTime, sessionRounds sessionVO.SessionRounds, sessionBreakTime sessionVO.BreakTime) error {
	op.RoundWorkTime = workTime
	op.RoundBreakTime = breakTime
	op.SessionRounds = sessionRounds
	op.SessionBreakTime = sessionBreakTime
	op.UpdatedAt = time.Now()
	
	return nil
}

// IsValid は設定値が有効かどうかチェックする（ドメインルール）
func (op *OptimizationPreferences) IsValid() bool {
	return op.RoundWorkTime.Minutes() >= MinOptWorkTime && op.RoundWorkTime.Minutes() <= MaxOptWorkTime &&
		op.RoundBreakTime.Minutes() >= MinOptBreakTime && op.RoundBreakTime.Minutes() <= MaxOptBreakTime &&
		op.SessionBreakTime.Minutes() >= MinOptSessionBreakTime && op.SessionBreakTime.Minutes() <= MaxOptSessionBreakTime
		// SessionRoundsのValue Object自体がバリデーションを保証するため、ここでのチェックは不要
}

