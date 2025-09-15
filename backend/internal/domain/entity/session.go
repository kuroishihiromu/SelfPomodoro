package entity

import (
	"errors"
	"time"

	sessionVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/session"
	statisticsVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/statistics"
	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
)

// Session はユーザのセッションを表すAggregateRoot（DDD強化版）
type Session struct {
	ID           sessionVO.SessionID          `json:"id"`
	UserID       userVO.UserID               `json:"user_id"`
	StartTime    time.Time                   `json:"start_time"`
	EndTime      *time.Time                  `json:"end_time"`
	AverageFocus *sessionVO.AverageFocus     `json:"average_focus"`
	TotalWorkMin *sessionVO.TotalWorkMinutes `json:"total_work_min"`
	RoundCount   *sessionVO.RoundCount       `json:"round_count"`
	BreakTime    *sessionVO.BreakTime        `json:"break_time"`
	CreatedAt    time.Time                   `json:"created_at"`
	UpdatedAt    time.Time                   `json:"updated_at"`

	// Aggregate内のエンティティ（メモリ上でのみ保持）
	rounds []*Round `json:"-"`
}


// NewSession は新しいセッションを作成する（ファクトリーメソッド）
func NewSession(userID userVO.UserID) *Session {
	now := time.Now()
	return &Session{
		ID:        sessionVO.NewSessionID(),
		UserID:    userID,
		StartTime: now,
		CreatedAt: now,
		UpdatedAt: now,
		rounds:    make([]*Round, 0), // 空のラウンドスライス初期化
	}
}

// LoadSessionWithRounds は既存のセッションにラウンドを読み込む（Repository用）
func (s *Session) LoadRounds(rounds []*Round) {
	s.rounds = rounds
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

// ドメインルール：Aggregateとしてのラウンド管理

// CanAddRound は新しいラウンドを追加できるかを判定する
func (s *Session) CanAddRound(maxRounds int) error {
	if !s.IsInProgress() {
		return errors.New("完了したセッションにはラウンドを追加できません")
	}

	if len(s.rounds) >= maxRounds {
		return errors.New("セッションの最大ラウンド数に達しています")
	}

	// 最後のラウンドが未完了の場合は追加不可
	if len(s.rounds) > 0 {
		lastRound := s.rounds[len(s.rounds)-1]
		if !lastRound.IsCompleted() {
			return errors.New("前のラウンドが完了していません")
		}
	}

	return nil
}

// AddRound は新しいラウンドをセッションに追加する（Aggregateの一貫性制御）
func (s *Session) AddRound(maxRounds int) (*Round, error) {
	if err := s.CanAddRound(maxRounds); err != nil {
		return nil, err
	}

	nextOrder := len(s.rounds) + 1
	
	// ビジネス不変条件: RoundOrder <= RoundCount の制約チェック
	if err := s.validateRoundOrderConstraint(nextOrder); err != nil {
		return nil, err
	}
	
	round := NewRound(s.ID, nextOrder)
	s.rounds = append(s.rounds, round)
	s.UpdatedAt = time.Now()

	return round, nil
}

// validateRoundOrderConstraint はRoundOrderとRoundCountの整合性を検証する
func (s *Session) validateRoundOrderConstraint(newRoundOrder int) error {
	currentRoundCount := len(s.rounds)
	expectedRoundCount := currentRoundCount + 1 // 新しいラウンド追加後
	
	if newRoundOrder > expectedRoundCount {
		return errors.New("ラウンド順序がラウンド数を超えることはできません")
	}
	
	if newRoundOrder != expectedRoundCount {
		return errors.New("ラウンド順序は連続している必要があります")
	}
	
	return nil
}

// GetRounds はセッション内のラウンド一覧を返す
func (s *Session) GetRounds() []*Round {
	return s.rounds
}

// GetCurrentRound は現在進行中のラウンドを返す
func (s *Session) GetCurrentRound() *Round {
	if len(s.rounds) == 0 {
		return nil
	}

	lastRound := s.rounds[len(s.rounds)-1]
	if !lastRound.IsCompleted() {
		return lastRound
	}

	return nil
}

// CompleteCurrentRound は現在のラウンドを完了する
func (s *Session) CompleteCurrentRound(focusScore *int, workTime, breakTime int) error {
	currentRound := s.GetCurrentRound()
	if currentRound == nil {
		return errors.New("完了対象のラウンドがありません")
	}

	if err := currentRound.CompleteWith(focusScore, workTime, breakTime); err != nil {
		return err
	}

	s.UpdatedAt = time.Now()
	return nil
}

// HasRounds はセッションにラウンドが存在するかを判定する
func (s *Session) HasRounds() bool {
	return s.RoundCount != nil && !s.RoundCount.IsEmpty()
}

// GetRoundCountOrZero はラウンド数または0を返す
func (s *Session) GetRoundCountOrZero() int {
	if s.RoundCount == nil {
		return 0
	}
	return s.RoundCount.Count()
}

// GetAverageFocusOrZero は平均集中度または0を返す
func (s *Session) GetAverageFocusOrZero() float64 {
	if s.AverageFocus == nil {
		return 0.0
	}
	return s.AverageFocus.Score()
}

// GetTotalWorkMinOrZero は総作業時間または0を返す
func (s *Session) GetTotalWorkMinOrZero() int {
	if s.TotalWorkMin == nil {
		return 0
	}
	return s.TotalWorkMin.Minutes()
}

// ドメインルール：統計計算（核心ビジネスロジック）

// CalculateStatistics はAggregateのラウンドからセッション統計を計算する（DDD強化）
func (s *Session) CalculateStatistics() statisticsVO.SessionStatistics {
	if len(s.rounds) == 0 {
		return statisticsVO.NewSessionStatistics(0.0, 0, 0, 0)
	}

	var totalFocus int
	var validFocusCount int
	var totalWorkMin int
	var totalBreakTime int
	var completedRounds int

	for _, round := range s.rounds {
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

	return statisticsVO.NewSessionStatistics(averageFocus, totalWorkMin, completedRounds, totalBreakTime)
}

// CompleteWithStatistics は統計情報を使ってセッションを完了する
func (s *Session) CompleteWithStatistics(stats statisticsVO.SessionStatistics) {
	now := time.Now()
	avgFocus := stats.AverageFocus()
	totalWork := stats.TotalWorkMin()
	roundCount := stats.RoundCount()
	breakTimeMinutes := stats.BreakTime()
	
	// Value Objectsを作成
	var avgFocusVO *sessionVO.AverageFocus
	if avgFocusScore, err := sessionVO.NewAverageFocus(avgFocus); err == nil {
		avgFocusVO = &avgFocusScore
	}
	
	var totalWorkVO *sessionVO.TotalWorkMinutes
	if totalWorkMin, err := sessionVO.NewTotalWorkMinutes(totalWork); err == nil {
		totalWorkVO = &totalWorkMin
	}
	
	var breakTimeVO *sessionVO.BreakTime
	if breakTimeMinutes > 0 {
		bt, err := sessionVO.NewBreakTime(breakTimeMinutes)
		if err == nil {
			breakTimeVO = &bt
		}
	}
	
	var roundCountVO *sessionVO.RoundCount
	if roundCountValue, err := sessionVO.NewRoundCount(roundCount); err == nil {
		roundCountVO = &roundCountValue
	}
	
	s.EndTime = &now
	s.AverageFocus = avgFocusVO
	s.TotalWorkMin = totalWorkVO
	s.RoundCount = roundCountVO
	s.BreakTime = breakTimeVO
	s.UpdatedAt = now
}

// CompleteWithRounds はAggregateのラウンドから統計を計算してセッションを完了する（DDD強化）
func (s *Session) CompleteWithRounds() {
	stats := s.CalculateStatistics()
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

