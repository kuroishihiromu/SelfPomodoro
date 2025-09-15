package optimization

import (
	"errors"
	"fmt"
	"time"

	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
	sessionVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/session"
)

// SessionOptimizationRecord はセッション最適化記録を表すValue Object
type SessionOptimizationRecord struct {
	userID        userVO.UserID
	timestamp     time.Time
	roundCount    int
	breakTime     sessionVO.BreakTime
	avgFocusScore float64
	totalWorkTime int // 分単位
}

// NewSessionOptimizationRecord は新しいセッション最適化記録を作成する
func NewSessionOptimizationRecord(userID userVO.UserID, roundCount int, breakTime sessionVO.BreakTime, avgFocusScore float64, totalWorkTime int) (SessionOptimizationRecord, error) {
	// バリデーション
	if roundCount <= 0 || roundCount > 10 {
		return SessionOptimizationRecord{}, errors.New("ラウンド数は1-10の間で設定してください")
	}
	
	if avgFocusScore < 0 || avgFocusScore > 100 {
		return SessionOptimizationRecord{}, errors.New("平均集中度スコアは0-100の間で設定してください")
	}
	
	if totalWorkTime <= 0 {
		return SessionOptimizationRecord{}, errors.New("合計作業時間は1分以上である必要があります")
	}

	return SessionOptimizationRecord{
		userID:        userID,
		timestamp:     time.Now(),
		roundCount:    roundCount,
		breakTime:     breakTime,
		avgFocusScore: avgFocusScore,
		totalWorkTime: totalWorkTime,
	}, nil
}

// NewSessionOptimizationRecordWithTime は指定時刻でセッション最適化記録を作成する（テスト・サンプルデータ用）
func NewSessionOptimizationRecordWithTime(userID userVO.UserID, timestamp time.Time, roundCount int, breakTime sessionVO.BreakTime, avgFocusScore float64, totalWorkTime int) (SessionOptimizationRecord, error) {
	if roundCount <= 0 || roundCount > 10 {
		return SessionOptimizationRecord{}, errors.New("ラウンド数は1-10の間で設定してください")
	}
	
	if avgFocusScore < 0 || avgFocusScore > 100 {
		return SessionOptimizationRecord{}, errors.New("平均集中度スコアは0-100の間で設定してください")
	}
	
	if totalWorkTime <= 0 {
		return SessionOptimizationRecord{}, errors.New("合計作業時間は1分以上である必要があります")
	}

	return SessionOptimizationRecord{
		userID:        userID,
		timestamp:     timestamp,
		roundCount:    roundCount,
		breakTime:     breakTime,
		avgFocusScore: avgFocusScore,
		totalWorkTime: totalWorkTime,
	}, nil
}

// UserID はユーザーIDを返す
func (s SessionOptimizationRecord) UserID() userVO.UserID {
	return s.userID
}

// Timestamp はタイムスタンプを返す
func (s SessionOptimizationRecord) Timestamp() time.Time {
	return s.timestamp
}

// RoundCount はラウンド数を返す
func (s SessionOptimizationRecord) RoundCount() int {
	return s.roundCount
}

// BreakTime は休憩時間を返す
func (s SessionOptimizationRecord) BreakTime() sessionVO.BreakTime {
	return s.breakTime
}

// AvgFocusScore は平均集中度スコアを返す
func (s SessionOptimizationRecord) AvgFocusScore() float64 {
	return s.avgFocusScore
}

// TotalWorkTime は合計作業時間（分）を返す
func (s SessionOptimizationRecord) TotalWorkTime() int {
	return s.totalWorkTime
}

// Equals は別のSessionOptimizationRecordと等価かを判定する
func (s SessionOptimizationRecord) Equals(other SessionOptimizationRecord) bool {
	return s.userID.Equals(other.userID) &&
		s.timestamp.Equal(other.timestamp) &&
		s.roundCount == other.roundCount &&
		s.breakTime.Equals(other.breakTime) &&
		s.avgFocusScore == other.avgFocusScore &&
		s.totalWorkTime == other.totalWorkTime
}

// IsHighPerformance は高いパフォーマンスのセッション最適化記録かを判定する
func (s SessionOptimizationRecord) IsHighPerformance() bool {
	return s.avgFocusScore >= 80
}

// IsEffectiveSession は効果的なセッションかを判定する
func (s SessionOptimizationRecord) IsEffectiveSession() bool {
	// 平均集中度が60以上で、適切な作業時間
	return s.avgFocusScore >= 60 && 
		   s.totalWorkTime >= 30 && // 最低30分は作業
		   s.roundCount >= 2 // 複数ラウンド実施
}

// GetSessionEfficiency はセッション効率を計算する（分あたりの集中度）
func (s SessionOptimizationRecord) GetSessionEfficiency() float64 {
	if s.totalWorkTime == 0 {
		return 0
	}
	return s.avgFocusScore / float64(s.totalWorkTime)
}

// GetOptimizationQuality は最適化の品質レベルを返す
func (s SessionOptimizationRecord) GetOptimizationQuality() string {
	if s.avgFocusScore >= 80 {
		return "高品質"
	} else if s.avgFocusScore >= 60 {
		return "中品質"
	} else {
		return "低品質"
	}
}

// IsLongBreak は長休憩時間かを判定する
func (s SessionOptimizationRecord) IsLongBreak() bool {
	return s.breakTime.IsLongBreak()
}

// ToTimestamp は永続化用のタイムスタンプ文字列を返す
func (s SessionOptimizationRecord) ToTimestamp() string {
	return s.timestamp.Format(time.RFC3339)
}

// String は文字列表現を返す
func (s SessionOptimizationRecord) String() string {
	return fmt.Sprintf("SessionOptimization[User=%s, Time=%s, Rounds=%d, AvgFocus=%.1f, TotalWork=%dmin]",
		s.userID.String()[:8],
		s.timestamp.Format("2006-01-02 15:04"),
		s.roundCount,
		s.avgFocusScore,
		s.totalWorkTime)
}