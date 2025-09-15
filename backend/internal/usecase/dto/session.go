package dto

import (
	"time"

	"github.com/google/uuid"
)

// SessionResponse はセッションのレスポンス形式を表すDTO
type SessionResponse struct {
	ID           uuid.UUID  `json:"id"`
	StartTime    time.Time  `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	AverageFocus *float64   `json:"average_focus"`
	TotalWorkMin *int       `json:"total_work_min"`
	RoundCount   *int       `json:"round_count"`
	BreakTime    *int       `json:"break_time"`
}

// SessionsResponse はセッションのリストAPIレスポンス形式を表すDTO
type SessionsResponse struct {
	Sessions []*SessionResponse `json:"sessions"`
}