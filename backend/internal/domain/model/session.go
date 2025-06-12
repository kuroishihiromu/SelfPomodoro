package model

import (
	"time"

	"github.com/google/uuid"
)

// Session はユーザのセッションを表すドメインモデル（強化版）
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

// SessionStatistics はセッション統計情報を表す
type SessionStatistics struct {
	AverageFocus   float64
	TotalWorkMin   int
	RoundCount     int
	TotalBreakTime int
}

// NewSession は新しいセッションを作成する（ファクトリーメソッド）
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

// ドメインルール：状態管理メソッド群

// IsCompleted はセッションが完了しているかを判定する
func (s *Session) IsCompleted() bool {
	return s.EndTime != nil
}

// IsInProgress はセッションが進行中かを判定する
func (s *Session) IsInProgress() bool {
	return s.EndTime == nil
}

// HasRounds はセッションにラウンドが存在するかを判定する
func (s *Session) HasRounds() bool {
	return s.RoundCount != nil && *s.RoundCount > 0
}

// GetRoundCountOrZero はラウンド数または0を返す
func (s *Session) GetRoundCountOrZero() int {
	if s.RoundCount == nil {
		return 0
	}
	return *s.RoundCount
}

// GetAverageFocusOrZero は平均集中度または0を返す
func (s *Session) GetAverageFocusOrZero() float64 {
	if s.AverageFocus == nil {
		return 0.0
	}
	return *s.AverageFocus
}

// GetTotalWorkMinOrZero は総作業時間または0を返す
func (s *Session) GetTotalWorkMinOrZero() int {
	if s.TotalWorkMin == nil {
		return 0
	}
	return *s.TotalWorkMin
}

// ドメインルール：統計計算（核心ビジネスロジック）

// CalculateStatistics はラウンドリストからセッション統計を計算する
func (s *Session) CalculateStatistics(rounds []*Round) *SessionStatistics {
	if len(rounds) == 0 {
		return &SessionStatistics{
			AverageFocus:   0.0,
			TotalWorkMin:   0,
			RoundCount:     0,
			TotalBreakTime: 0,
		}
	}

	var totalFocus int
	var validFocusCount int
	var totalWorkMin int
	var totalBreakTime int
	var completedRounds int

	for _, round := range rounds {
		// 完了したラウンドのみを統計対象とする
		if !round.ContributesToStatistics() {
			continue
		}

		completedRounds++

		// 集中度の集計（スコアが設定されている場合のみ）
		if round.HasFocusScore() {
			totalFocus += round.GetFocusScoreOrZero()
			validFocusCount++
		}

		// 作業時間の集計
		totalWorkMin += round.GetWorkTimeForStats()

		// 休憩時間の集計
		totalBreakTime += round.GetBreakTimeForStats()
	}

	// 平均集中度の計算
	averageFocus := 0.0
	if validFocusCount > 0 {
		averageFocus = float64(totalFocus) / float64(validFocusCount)
	}

	return &SessionStatistics{
		AverageFocus:   averageFocus,
		TotalWorkMin:   totalWorkMin,
		RoundCount:     completedRounds,
		TotalBreakTime: totalBreakTime,
	}
}

// CompleteWithStatistics は統計情報を使ってセッションを完了する
func (s *Session) CompleteWithStatistics(stats *SessionStatistics) {
	now := time.Now()
	s.EndTime = &now
	s.AverageFocus = &stats.AverageFocus
	s.TotalWorkMin = &stats.TotalWorkMin
	s.RoundCount = &stats.RoundCount
	s.BreakTime = &stats.TotalBreakTime
	s.UpdatedAt = now
}

// CompleteWithRounds はラウンドから統計を計算してセッションを完了する
func (s *Session) CompleteWithRounds(rounds []*Round) {
	stats := s.CalculateStatistics(rounds)
	s.CompleteWithStatistics(stats)
}

// ドメインルール：最適化メッセージ送信判定

// ShouldSendOptimizationMessage は最適化メッセージを送信すべきかを判定する
func (s *Session) ShouldSendOptimizationMessage() bool {
	return s.IsCompleted() &&
		s.HasRounds() &&
		s.AverageFocus != nil
}

// GetOptimizationMessageData は最適化メッセージ用のデータを返す
func (s *Session) GetOptimizationMessageData() (avgFocus float64, totalWork int, hasValidData bool) {
	if !s.ShouldSendOptimizationMessage() {
		return 0.0, 0, false
	}
	return s.GetAverageFocusOrZero(), s.GetTotalWorkMinOrZero(), true
}

// ドメインルール：ビジネス計算

// GetDuration はセッションの実行時間を返す
func (s *Session) GetDuration() time.Duration {
	if s.EndTime == nil {
		return time.Since(s.StartTime)
	}
	return s.EndTime.Sub(s.StartTime)
}

// GetDurationMinutes はセッションの実行時間を分で返す
func (s *Session) GetDurationMinutes() int {
	return int(s.GetDuration().Minutes())
}

// GetDurationHours はセッションの実行時間を時間で返す
func (s *Session) GetDurationHours() float64 {
	return s.GetDuration().Hours()
}

// GetEfficiency は作業効率を計算する（作業時間/総時間）
func (s *Session) GetEfficiency() float64 {
	if s.IsInProgress() {
		return 0.0
	}

	totalMinutes := s.GetDurationMinutes()
	if totalMinutes == 0 {
		return 0.0
	}

	workMinutes := s.GetTotalWorkMinOrZero()
	return float64(workMinutes) / float64(totalMinutes) * 100.0
}

// IsProductiveSession は生産的なセッションかを判定する（ビジネスルール）
func (s *Session) IsProductiveSession(minRounds int, minAverageFocus float64) bool {
	return s.IsCompleted() &&
		s.GetRoundCountOrZero() >= minRounds &&
		s.GetAverageFocusOrZero() >= minAverageFocus
}

// GetSessionQuality はセッション品質を評価する
func (s *Session) GetSessionQuality() string {
	if !s.IsCompleted() {
		return "進行中"
	}

	avgFocus := s.GetAverageFocusOrZero()
	rounds := s.GetRoundCountOrZero()

	if rounds == 0 {
		return "未完了"
	}

	if avgFocus >= 80 && rounds >= 3 {
		return "優秀"
	} else if avgFocus >= 60 && rounds >= 2 {
		return "良好"
	} else if avgFocus >= 40 || rounds >= 1 {
		return "普通"
	} else {
		return "要改善"
	}
}

// 既存のレスポンス変換メソッドは維持

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
