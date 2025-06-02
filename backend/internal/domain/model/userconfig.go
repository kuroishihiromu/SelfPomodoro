package model

import (
	"time"

	"github.com/google/uuid"
)

// UserConfig はユーザーのポモドーロ設定を表すドメインモデル
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
		RoundWorkTime:    25, // デフォルト25分
		RoundBreakTime:   5,  // デフォルト5分
		SessionRounds:    3,  // デフォルト3ラウンド
		SessionBreakTime: 15, // デフォルト15分長休憩
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// UpdateSettings は設定を更新する
func (uc *UserConfig) UpdateSettings(workTime, breakTime, sessionRounds, sessionBreakTime int) {
	uc.RoundWorkTime = workTime
	uc.RoundBreakTime = breakTime
	uc.SessionRounds = sessionRounds
	uc.SessionBreakTime = sessionBreakTime
	uc.UpdatedAt = time.Now()
}

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

// IsValid は設定値が有効かどうかチェックする
func (uc *UserConfig) IsValid() bool {
	return uc.RoundWorkTime > 0 && uc.RoundWorkTime <= 120 &&
		uc.RoundBreakTime > 0 && uc.RoundBreakTime <= 60 &&
		uc.SessionRounds > 0 && uc.SessionRounds <= 10 &&
		uc.SessionBreakTime >= 5 && uc.SessionBreakTime <= 120
}
