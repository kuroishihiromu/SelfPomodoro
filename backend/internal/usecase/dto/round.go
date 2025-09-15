package dto

import (
	"time"

	"github.com/google/uuid"
)

// RoundCreateRequest はラウンド作成リクエストを表すDTO
type RoundCreateRequest struct {
	// 現在は空（将来の拡張用）
}

// RoundCompleteRequest はラウンド完了リクエストを表すDTO
type RoundCompleteRequest struct {
	FocusScore *int `json:"focus_score" validate:"omitempty,min=0,max=100"`
	WorkTime   *int `json:"work_time" validate:"omitempty,min=1"`
	BreakTime  *int `json:"break_time" validate:"omitempty,min=1"`
}

// RoundResponse はラウンドのレスポンスを表すDTO
type RoundResponse struct {
	ID         uuid.UUID  `json:"id"`
	SessionID  uuid.UUID  `json:"session_id"`
	RoundOrder int        `json:"round_order"`
	StartTime  time.Time  `json:"start_time"`
	EndTime    *time.Time `json:"end_time,omitempty"`
	WorkTime   *int       `json:"work_time,omitempty"`
	BreakTime  *int       `json:"break_time,omitempty"`
	FocusScore *int       `json:"focus_score,omitempty"`
}

// RoundsResponse はラウンドのリストAPIレスポンス形式を表すDTO
type RoundsResponse struct {
	Rounds []*RoundResponse `json:"rounds"`
}