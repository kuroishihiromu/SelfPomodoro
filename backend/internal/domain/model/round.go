package model

import (
	"time"

	"github.com/google/uuid"
)

// Round はポモドーロラウンドを表す構造体
type Round struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	SessionID  uuid.UUID  `db:"session_id" json:"session_id"`
	RoundOrder int        `db:"round_order" json:"round_order"`
	StartTime  time.Time  `db:"start_time" json:"start_time"`
	EndTime    *time.Time `db:"end_time" json:"end_time,omitempty"`
	WorkTime   *int       `db:"work_time" json:"work_time,omitempty"`
	BreakTime  *int       `db:"break_time" json:"break_time,omitempty"`
	FocusScore *int       `db:"focus_score" json:"focus_score,omitempty"`
	IsAborted  bool       `db:"is_aborted" json:"is_aborted"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
}

// NewRound は新しいラウンドを作成する
func NewRound(sessionID uuid.UUID, roundOrder int) *Round {
	now := time.Now()
	return &Round{
		ID:         uuid.New(),
		SessionID:  sessionID,
		RoundOrder: roundOrder,
		StartTime:  now,
		IsAborted:  false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// RoundCreateRequest はラウンド作成リクエストを表す
type RoundCreateRequest struct {
}

// RoundCompleteRequest はラウンド完了リクエストを表す
type RoundCompleteRequest struct {
	FocusScore *int `json:"focus_score" validate:"omitempty,min=0,max=100"`
}

// RoundResponse はラウンドのレスポンスを表す
type RoundResponse struct {
	ID         uuid.UUID  `json:"id"`
	SessionID  uuid.UUID  `json:"session_id"`
	RoundOrder int        `json:"round_order"`
	StartTime  time.Time  `json:"start_time"`
	EndTime    *time.Time `json:"end_time,omitempty"`
	WorkTime   *int       `json:"work_time,omitempty"`
	BreakTime  *int       `json:"break_time,omitempty"`
	FocusScore *int       `json:"focus_score,omitempty"`
	IsAborted  bool       `json:"is_aborted"`
}

// ToResponse はラウンドのドメインモデルをレスポンス形式に変換する
func (r *Round) ToResponse() *RoundResponse {
	return &RoundResponse{
		ID:         r.ID,
		SessionID:  r.SessionID,
		RoundOrder: r.RoundOrder,
		StartTime:  r.StartTime,
		EndTime:    r.EndTime,
		WorkTime:   r.WorkTime,
		BreakTime:  r.BreakTime,
		FocusScore: r.FocusScore,
		IsAborted:  r.IsAborted,
	}
}

// RoundsResponse はラウンドのリストAPIレスポンスを表す
type RoundsResponse struct {
	Rounds []*RoundResponse `json:"rounds"`
}
