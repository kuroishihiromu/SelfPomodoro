package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// 最適化関連定数
const (
	MaxFocusScore = 100
	MinFocusScore = 0
)

// Round はポモドーロラウンドを表す構造体（完了時のみ作成版）
type Round struct {
	ID         uuid.UUID  `json:"id"`
	SessionID  uuid.UUID  `json:"session_id"`
	RoundOrder int        `json:"round_order"`
	StartTime  time.Time  `json:"start_time"`
	EndTime    *time.Time `json:"end_time,omitempty"`
	WorkTime   *int       `json:"work_time,omitempty"`
	BreakTime  *int       `json:"break_time,omitempty"`
	FocusScore *int       `json:"focus_score,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// NewRound は新しいラウンドを作成する（完了時作成用ファクトリーメソッド）
func NewRound(sessionID uuid.UUID, roundOrder int) *Round {
	now := time.Now()
	return &Round{
		ID:         uuid.New(),
		SessionID:  sessionID,
		RoundOrder: roundOrder,
		StartTime:  now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// ドメインルール：状態管理メソッド群

// IsCompleted はラウンドが完了しているかを判定する（完了時のみ作成されるため常にtrue）
func (r *Round) IsCompleted() bool {
	return r.EndTime != nil
}

// HasFocusScore は集中度スコアが設定されているかを判定する
func (r *Round) HasFocusScore() bool {
	return r.FocusScore != nil && *r.FocusScore >= MinFocusScore
}

// GetFocusScoreOrZero は集中度スコアまたは0を返す
func (r *Round) GetFocusScoreOrZero() int {
	if r.FocusScore == nil {
		return 0
	}
	return *r.FocusScore
}

// ドメインルール：完了処理

// CanBeCompleted は完了可能かを判定する（完了時のみ作成されるため基本的に不要だが安全性のため）
func (r *Round) CanBeCompleted() error {
	if r.IsCompleted() {
		return errors.New("ラウンドは既に完了しています")
	}
	return nil
}

// CompleteWith はラウンドを指定されたパラメータで完了する（ドメインルール適用）
func (r *Round) CompleteWith(focusScore *int, workTime, breakTime int) error {
	if err := r.CanBeCompleted(); err != nil {
		return err
	}

	// 集中度スコアのバリデーション
	if focusScore != nil {
		if *focusScore < MinFocusScore || *focusScore > MaxFocusScore {
			return errors.New("集中度スコアは0から100の間で設定してください")
		}
	}

	// 作業時間・休憩時間のバリデーション
	if workTime <= 0 || breakTime < 0 {
		return errors.New("作業時間は1分以上、休憩時間は0分以上で設定してください")
	}

	now := time.Now()
	r.EndTime = &now
	r.FocusScore = focusScore
	r.WorkTime = &workTime
	r.BreakTime = &breakTime
	r.UpdatedAt = now

	return nil
}


// ドメインルール：最適化メッセージ送信判定

// ShouldSendOptimizationMessage は最適化メッセージを送信すべきかを判定する
func (r *Round) ShouldSendOptimizationMessage() bool {
	return r.IsCompleted() &&
		r.HasFocusScore() &&
		*r.FocusScore >= MinFocusScore
}

// GetOptimizationMessageData は最適化メッセージ用のデータを返す
func (r *Round) GetOptimizationMessageData() (focusScore int, hasValidScore bool) {
	if !r.ShouldSendOptimizationMessage() {
		return 0, false
	}
	return *r.FocusScore, true
}

// ドメインルール：統計計算用メソッド

// ContributesToStatistics は統計計算に含めるべきかを判定する
func (r *Round) ContributesToStatistics() bool {
	return r.IsCompleted()
}

// GetWorkTimeForStats は統計用の作業時間を返す
func (r *Round) GetWorkTimeForStats() int {
	if r.WorkTime == nil {
		return 0
	}
	return *r.WorkTime
}

// GetBreakTimeForStats は統計用の休憩時間を返す
func (r *Round) GetBreakTimeForStats() int {
	if r.BreakTime == nil {
		return 0
	}
	return *r.BreakTime
}

// ドメインルール：ビジネス計算

// GetDuration はラウンドの実行時間を返す（分）
func (r *Round) GetDuration() time.Duration {
	if r.EndTime == nil {
		return time.Since(r.StartTime)
	}
	return r.EndTime.Sub(r.StartTime)
}

// GetDurationMinutes はラウンドの実行時間を分で返す
func (r *Round) GetDurationMinutes() int {
	return int(r.GetDuration().Minutes())
}

// IsOvertime は予定時間を超過しているかを判定する
func (r *Round) IsOvertime(expectedMinutes int) bool {
	return r.GetDurationMinutes() > expectedMinutes
}

// 既存のレスポンス変換メソッドは維持

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
	}
}

// RoundsResponse はラウンドのリストAPIレスポンスを表す
type RoundsResponse struct {
	Rounds []*RoundResponse `json:"rounds"`
}
