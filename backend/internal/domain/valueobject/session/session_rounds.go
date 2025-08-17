package valueobject

import (
	"errors"
)

// SessionRounds はセッション内の計画ラウンド数を表すValue Object
// RoundCountとは異なり、セッション開始前に設定される目標値を管理する
type SessionRounds struct {
	count int
}

const (
	// MinSessionRounds は最小セッションラウンド数
	MinSessionRounds = 1
	// MaxSessionRounds は最大セッションラウンド数
	MaxSessionRounds = 10
	// DefaultSessionRounds はデフォルトセッションラウンド数
	DefaultSessionRounds = 3
	// OptimalMinSessionRounds は最適な最小ラウンド数
	OptimalMinSessionRounds = 3
	// OptimalMaxSessionRounds は最適な最大ラウンド数
	OptimalMaxSessionRounds = 5
)

// NewSessionRounds は新しいSessionRoundsを作成する
func NewSessionRounds(count int) (SessionRounds, error) {
	if count < MinSessionRounds {
		return SessionRounds{}, errors.New("セッションラウンド数は1回以上である必要があります")
	}
	
	if count > MaxSessionRounds {
		return SessionRounds{}, errors.New("セッションラウンド数は10回以下である必要があります")
	}
	
	return SessionRounds{count: count}, nil
}

// NewDefaultSessionRounds はデフォルトのSessionRoundsを作成する
func NewDefaultSessionRounds() SessionRounds {
	return SessionRounds{count: DefaultSessionRounds}
}

// MustNewSessionRounds はパニックを起こす可能性があるSessionRounds作成（テスト用）
func MustNewSessionRounds(count int) SessionRounds {
	sr, err := NewSessionRounds(count)
	if err != nil {
		panic(err)
	}
	return sr
}

// Count はラウンド数の値を返す
func (sr SessionRounds) Count() int {
	return sr.count
}

// String は文字列表現を返す
func (sr SessionRounds) String() string {
	return string(rune(sr.count + '0'))
}

// IsMinimum は最小ラウンド数かどうかを判定する
func (sr SessionRounds) IsMinimum() bool {
	return sr.count == MinSessionRounds
}

// IsMaximum は最大ラウンド数かどうかを判定する
func (sr SessionRounds) IsMaximum() bool {
	return sr.count == MaxSessionRounds
}

// IsDefault はデフォルトラウンド数かどうかを判定する
func (sr SessionRounds) IsDefault() bool {
	return sr.count == DefaultSessionRounds
}

// IsOptimal は最適なラウンド数範囲内かどうかを判定する
func (sr SessionRounds) IsOptimal() bool {
	return sr.count >= OptimalMinSessionRounds && sr.count <= OptimalMaxSessionRounds
}

// IsShort は短いセッション（2ラウンド以下）かどうかを判定する
func (sr SessionRounds) IsShort() bool {
	return sr.count <= 2
}

// IsLong は長いセッション（6ラウンド以上）かどうかを判定する
func (sr SessionRounds) IsLong() bool {
	return sr.count >= 6
}

// EstimatedDurationMinutes は推定所要時間を分で返す
// 作業時間25分 + 休憩時間5分 = 30分/ラウンドで計算
func (sr SessionRounds) EstimatedDurationMinutes() int {
	return sr.count * 30
}

// EstimatedDurationWithCustomTimes はカスタム時間での推定所要時間を分で返す
func (sr SessionRounds) EstimatedDurationWithCustomTimes(workMinutes, breakMinutes int) int {
	return sr.count * (workMinutes + breakMinutes)
}

// GetIntensityLevel は集中度レベルを返す
func (sr SessionRounds) GetIntensityLevel() string {
	switch {
	case sr.count <= 2:
		return "軽度"
	case sr.count <= 4:
		return "標準"
	case sr.count <= 6:
		return "集中"
	default:
		return "超集中"
	}
}

// GetRecommendedBreakMinutes は推奨休憩時間を返す
func (sr SessionRounds) GetRecommendedBreakMinutes() int {
	switch {
	case sr.count <= 2:
		return 10 // 短いセッション：短い休憩
	case sr.count <= 4:
		return 15 // 標準セッション：標準休憩
	case sr.count <= 6:
		return 20 // 長いセッション：長い休憩
	default:
		return 30 // 超長セッション：長い休憩
	}
}

// CanIncrement はインクリメント可能かを判定する
func (sr SessionRounds) CanIncrement() bool {
	return sr.count < MaxSessionRounds
}

// CanDecrement はデクリメント可能かを判定する
func (sr SessionRounds) CanDecrement() bool {
	return sr.count > MinSessionRounds
}

// Increment はラウンド数を1増やした新しいインスタンスを返す
func (sr SessionRounds) Increment() (SessionRounds, error) {
	if !sr.CanIncrement() {
		return sr, errors.New("これ以上ラウンド数を増やせません")
	}
	return NewSessionRounds(sr.count + 1)
}

// Decrement はラウンド数を1減らした新しいインスタンスを返す
func (sr SessionRounds) Decrement() (SessionRounds, error) {
	if !sr.CanDecrement() {
		return sr, errors.New("これ以上ラウンド数を減らせません")
	}
	return NewSessionRounds(sr.count - 1)
}

// Equals は他のSessionRoundsと等しいかを判定する
func (sr SessionRounds) Equals(other SessionRounds) bool {
	return sr.count == other.count
}

// IsLessThan は他のSessionRoundsより小さいかを判定する
func (sr SessionRounds) IsLessThan(other SessionRounds) bool {
	return sr.count < other.count
}

// IsGreaterThan は他のSessionRoundsより大きいかを判定する
func (sr SessionRounds) IsGreaterThan(other SessionRounds) bool {
	return sr.count > other.count
}

// ToOptimizationRange は最適化用の範囲内に調整した値を返す
func (sr SessionRounds) ToOptimizationRange() SessionRounds {
	if sr.count < OptimalMinSessionRounds {
		return SessionRounds{count: OptimalMinSessionRounds}
	}
	if sr.count > OptimalMaxSessionRounds {
		return SessionRounds{count: OptimalMaxSessionRounds}
	}
	return sr
}