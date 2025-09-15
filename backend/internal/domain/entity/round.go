package entity

import (
	"errors"
	"time"

	roundVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/round"
	sessionVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/session"
)

// 最適化関連定数
const (
	MaxFocusScore = 100
	MinFocusScore = 0
)

// Round はポモドーロラウンドを表すエンティティ（Session Aggregate内）
type Round struct {
	ID         roundVO.RoundID     `json:"id"`
	SessionID  sessionVO.SessionID `json:"session_id"`  // Aggregateルートへの参照
	RoundOrder roundVO.RoundOrder  `json:"round_order"` // Aggregate内での順序
	StartTime  time.Time           `json:"start_time"`
	EndTime    *time.Time          `json:"end_time,omitempty"`
	WorkTime   *roundVO.WorkTime   `json:"work_time,omitempty"`
	BreakTime  *roundVO.BreakTime  `json:"break_time,omitempty"`
	FocusScore *roundVO.FocusScore `json:"focus_score,omitempty"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

// NewRound は新しいラウンドを作成する（完了時作成用ファクトリーメソッド）
func NewRound(sessionID sessionVO.SessionID, roundOrder int) *Round {
	now := time.Now()
	roundOrderVO, _ := roundVO.NewRoundOrder(roundOrder) // Assuming valid input
	return &Round{
		ID:         roundVO.NewRoundID(),
		SessionID:  sessionID,
		RoundOrder: roundOrderVO,
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
	return r.FocusScore != nil
}

// GetFocusScoreOrZero は集中度スコアまたは0を返す
func (r *Round) GetFocusScoreOrZero() int {
	if r.FocusScore == nil {
		return 0
	}
	return r.FocusScore.Value()
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

	// 集中度スコアのバリデーション・変換
	var focusScoreVO *roundVO.FocusScore
	if focusScore != nil {
		fs, err := roundVO.NewFocusScore(*focusScore)
		if err != nil {
			return err
		}
		focusScoreVO = &fs
	}

	// 作業時間のバリデーション・変換
	workTimeVO, err := roundVO.NewWorkTime(workTime)
	if err != nil {
		return err
	}

	// 休憩時間のバリデーション・変換
	breakTimeVO, err := roundVO.NewBreakTime(breakTime)
	if err != nil {
		return err
	}

	now := time.Now()
	r.EndTime = &now
	r.FocusScore = focusScoreVO
	r.WorkTime = &workTimeVO
	r.BreakTime = &breakTimeVO
	r.UpdatedAt = now

	return nil
}

// ドメインルール：最適化メッセージ送信判定

// ShouldSendOptimizationMessage は最適化メッセージを送信すべきかを判定する
func (r *Round) ShouldSendOptimizationMessage() bool {
	return r.IsCompleted() && r.HasFocusScore()
}

// GetOptimizationMessageData は最適化メッセージ用のデータを返す
func (r *Round) GetOptimizationMessageData() (focusScore int, hasValidScore bool) {
	if !r.ShouldSendOptimizationMessage() {
		return 0, false
	}
	return r.FocusScore.Value(), true
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
	return r.WorkTime.Minutes()
}

// GetBreakTimeForStats は統計用の休憩時間を返す
func (r *Round) GetBreakTimeForStats() int {
	if r.BreakTime == nil {
		return 0
	}
	return r.BreakTime.Minutes()
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

