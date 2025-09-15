package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateOptimizationPreferencesRequest は最適化設定作成リクエストを表すDTO
type CreateOptimizationPreferencesRequest struct {
	RoundWorkTime     *int `json:"round_work_time" validate:"omitempty,min=1,max=180"`
	RoundBreakTime    *int `json:"round_break_time" validate:"omitempty,min=1,max=60"`
	SessionRounds     *int `json:"session_rounds" validate:"omitempty,min=1,max=20"`
	SessionBreakTime  *int `json:"session_break_time" validate:"omitempty,min=1,max=180"`
}

// UpdateOptimizationPreferencesRequest は最適化設定更新リクエストを表すDTO
type UpdateOptimizationPreferencesRequest struct {
	RoundWorkTime     *int `json:"round_work_time" validate:"omitempty,min=1,max=180"`
	RoundBreakTime    *int `json:"round_break_time" validate:"omitempty,min=1,max=60"`
	SessionRounds     *int `json:"session_rounds" validate:"omitempty,min=1,max=20"`
	SessionBreakTime  *int `json:"session_break_time" validate:"omitempty,min=1,max=180"`
}

// OptimizationPreferencesResponse は最適化設定のレスポンスを表すDTO
type OptimizationPreferencesResponse struct {
	UserID            uuid.UUID `json:"user_id"`
	RoundWorkTime     int       `json:"round_work_time"`
	RoundBreakTime    int       `json:"round_break_time"`
	SessionRounds     int       `json:"session_rounds"`
	SessionBreakTime  int       `json:"session_break_time"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}