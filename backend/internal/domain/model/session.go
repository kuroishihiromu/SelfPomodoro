package model

import (
	"time"

	"github.com/google/uuid"
)

// Session はユーザのセッションを表すドメインモデル
type Session struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	UserID       uuid.UUID  `db:"user_id" json:"user_id"`
	StartTime    time.Time  `db:"start_time" json:"start_time"`
	EndTime      *time.Time `db:"end_time" json:"end_time"`
	AverageFocus *float64   `db:"average_focus" json:"average_focus"`
	TotalWorkMin *int       `db:"total_work_min" json:"total_work_min"`
	RoundCount   *int       `db:"round_count" json:"round_count"`
	BreakTime    *int       `db:"break_time" json:"break_time"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}

// NewSession は新しいセッションを作成する
func NewSession(userID uuid.UUID) *Session {
	now := time.Now()
	return &Session{
		ID:        uuid.New(),
		UserID:    userID,
		StartTime: now,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// SessionResponse はセッションのレスポンス形式を表す構造体
type SessionResponse struct {
	ID           uuid.UUID  `json:"id"`
	StartTime    time.Time  `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	AverageFocus *float64   `json:"average_focus"`
	TotalWorkMin *int       `json:"total_work_min"`
	RoundCount   *int       `json:"round_count"`
	BreakTime    *int       `json:"break_time"`
}

// ToResponse はドメインモデルからAPIレスポンス形式に変換する
func (s *Session) ToResponse() *SessionResponse {
	return &SessionResponse{
		ID:           s.ID,
		StartTime:    s.StartTime,
		EndTime:      s.EndTime,
		AverageFocus: s.AverageFocus,
		TotalWorkMin: s.TotalWorkMin,
		RoundCount:   s.RoundCount,
		BreakTime:    s.BreakTime,
	}
}

// SessionsResponse はセッションのリストAPIレスポンス形式を表す構造体
type SessionsResponse struct {
	Sessions []*SessionResponse `json:"sessions"`
}
