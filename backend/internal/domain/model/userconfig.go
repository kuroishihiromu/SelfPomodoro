package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// デフォルト値定数（ドメインルール）
const (
	DefaultWorkTime         = 25 // デフォルト作業時間（分）
	DefaultBreakTime        = 5  // デフォルト休憩時間（分）
	DefaultSessionRounds    = 3  // デフォルトラウンド数
	DefaultSessionBreakTime = 15 // デフォルト長休憩時間（分）

	// バリデーション制約
	MinWorkTime         = 1
	MaxWorkTime         = 120
	MinBreakTime        = 1
	MaxBreakTime        = 60
	MinSessionRounds    = 1
	MaxSessionRounds    = 10
	MinSessionBreakTime = 5
	MaxSessionBreakTime = 120
)

// UserConfig はユーザーのポモドーロ設定を表すドメインモデル（強化版）
type UserConfig struct {
	UserID           string    `dynamodb:"user_id" json:"user_id"`
	RoundWorkTime    int       `dynamodb:"round_work_time" json:"round_work_time"`       // 作業時間（分）
	RoundBreakTime   int       `dynamodb:"round_break_time" json:"round_break_time"`     // 休憩時間（分）
	SessionRounds    int       `dynamodb:"session_rounds" json:"session_rounds"`         // セッション内ラウンド数
	SessionBreakTime int       `dynamodb:"session_break_time" json:"session_break_time"` // セッション後長休憩（分）
	CreatedAt        time.Time `dynamodb:"created_at" json:"created_at"`
	UpdatedAt        time.Time `dynamodb:"updated_at" json:"updated_at"`
}

// NewUserConfig は新しいユーザー設定を作成する（デフォルト値付き）
func NewUserConfig(userID uuid.UUID) *UserConfig {
	now := time.Now()
	return &UserConfig{
		UserID:           userID.String(),
		RoundWorkTime:    DefaultWorkTime,
		RoundBreakTime:   DefaultBreakTime,
		SessionRounds:    DefaultSessionRounds,
		SessionBreakTime: DefaultSessionBreakTime,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// NewDefaultUserConfig はデフォルト設定のUserConfigを作成（明示的なファクトリー）
func NewDefaultUserConfig(userID uuid.UUID) *UserConfig {
	return NewUserConfig(userID)
}

// ドメインルール：安全なデフォルト値取得メソッド群

// GetWorkTimeOrDefault は作業時間またはデフォルト値を返す
func (uc *UserConfig) GetWorkTimeOrDefault() int {
	if uc == nil || uc.RoundWorkTime <= 0 {
		return DefaultWorkTime
	}
	return uc.RoundWorkTime
}

// GetBreakTimeOrDefault は休憩時間またはデフォルト値を返す
func (uc *UserConfig) GetBreakTimeOrDefault() int {
	if uc == nil || uc.RoundBreakTime <= 0 {
		return DefaultBreakTime
	}
	return uc.RoundBreakTime
}

// GetSessionRoundsOrDefault はセッションラウンド数またはデフォルト値を返す
func (uc *UserConfig) GetSessionRoundsOrDefault() int {
	if uc == nil || uc.SessionRounds <= 0 {
		return DefaultSessionRounds
	}
	return uc.SessionRounds
}

// GetSessionBreakTimeOrDefault はセッション長休憩時間またはデフォルト値を返す
func (uc *UserConfig) GetSessionBreakTimeOrDefault() int {
	if uc == nil || uc.SessionBreakTime <= 0 {
		return DefaultSessionBreakTime
	}
	return uc.SessionBreakTime
}

// ドメインルール：最適化データのベース値提供

// GetOptimizationBaseValues は最適化用のベース値を返す
func (uc *UserConfig) GetOptimizationBaseValues() (workTime int, breakTime int, sessionRounds int, sessionBreakTime int) {
	return uc.GetWorkTimeOrDefault(),
		uc.GetBreakTimeOrDefault(),
		uc.GetSessionRoundsOrDefault(),
		uc.GetSessionBreakTimeOrDefault()
}

// UpdateSettings は設定を更新する（ドメインルール適用）
func (uc *UserConfig) UpdateSettings(workTime, breakTime, sessionRounds, sessionBreakTime int) {
	uc.RoundWorkTime = workTime
	uc.RoundBreakTime = breakTime
	uc.SessionRounds = sessionRounds
	uc.SessionBreakTime = sessionBreakTime
	uc.UpdatedAt = time.Now()
}

// IsValid は設定値が有効かどうかチェックする（ドメインルール）
func (uc *UserConfig) IsValid() bool {
	return uc.RoundWorkTime >= MinWorkTime && uc.RoundWorkTime <= MaxWorkTime &&
		uc.RoundBreakTime >= MinBreakTime && uc.RoundBreakTime <= MaxBreakTime &&
		uc.SessionRounds >= MinSessionRounds && uc.SessionRounds <= MaxSessionRounds &&
		uc.SessionBreakTime >= MinSessionBreakTime && uc.SessionBreakTime <= MaxSessionBreakTime
}

// ValidateSettings は設定値のバリデーションを行う
func (uc *UserConfig) ValidateSettings() error {
	if uc.RoundWorkTime < MinWorkTime || uc.RoundWorkTime > MaxWorkTime {
		return fmt.Errorf("作業時間は%d分から%d分の間で設定してください", MinWorkTime, MaxWorkTime)
	}
	if uc.RoundBreakTime < MinBreakTime || uc.RoundBreakTime > MaxBreakTime {
		return fmt.Errorf("休憩時間は%d分から%d分の間で設定してください", MinBreakTime, MaxBreakTime)
	}
	if uc.SessionRounds < MinSessionRounds || uc.SessionRounds > MaxSessionRounds {
		return fmt.Errorf("セッションラウンド数は%d回から%d回の間で設定してください", MinSessionRounds, MaxSessionRounds)
	}
	if uc.SessionBreakTime < MinSessionBreakTime || uc.SessionBreakTime > MaxSessionBreakTime {
		return fmt.Errorf("セッション長休憩時間は%d分から%d分の間で設定してください", MinSessionBreakTime, MaxSessionBreakTime)
	}
	return nil
}

// 既存のレスポンス変換メソッドは維持

// CreateUserConfigRequest はユーザー設定作成リクエストを表す
type CreateUserConfigRequest struct {
	RoundWorkTime    int `json:"round_work_time" validate:"required,min=1,max=120"`    // 1-120分
	RoundBreakTime   int `json:"round_break_time" validate:"required,min=1,max=60"`    // 1-60分
	SessionRounds    int `json:"session_rounds" validate:"required,min=1,max=10"`      // 1-10ラウンド
	SessionBreakTime int `json:"session_break_time" validate:"required,min=5,max=120"` // 5-120分
}

// UpdateUserConfigRequest はユーザー設定更新リクエストを表す
type UpdateUserConfigRequest struct {
	RoundWorkTime    *int `json:"round_work_time,omitempty" validate:"omitempty,min=1,max=120"`
	RoundBreakTime   *int `json:"round_break_time,omitempty" validate:"omitempty,min=1,max=60"`
	SessionRounds    *int `json:"session_rounds,omitempty" validate:"omitempty,min=1,max=10"`
	SessionBreakTime *int `json:"session_break_time,omitempty" validate:"omitempty,min=5,max=120"`
}

// UserConfigResponse はユーザー設定のAPIレスポンスを表す
type UserConfigResponse struct {
	UserID           string    `json:"user_id"`
	RoundWorkTime    int       `json:"round_work_time"`
	RoundBreakTime   int       `json:"round_break_time"`
	SessionRounds    int       `json:"session_rounds"`
	SessionBreakTime int       `json:"session_break_time"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// ToResponse はドメインモデルからAPIレスポンス形式に変換する
func (uc *UserConfig) ToResponse() *UserConfigResponse {
	return &UserConfigResponse{
		UserID:           uc.UserID,
		RoundWorkTime:    uc.RoundWorkTime,
		RoundBreakTime:   uc.RoundBreakTime,
		SessionRounds:    uc.SessionRounds,
		SessionBreakTime: uc.SessionBreakTime,
		CreatedAt:        uc.CreatedAt,
		UpdatedAt:        uc.UpdatedAt,
	}
}
