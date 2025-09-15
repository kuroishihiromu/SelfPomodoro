package valueobject

import (
	"errors"
	"fmt"
)

// RoundCount はセッション内のラウンド数を表すValue Object
type RoundCount struct {
	count int
}

// NewRoundCount は新しいRoundCountを作成する
func NewRoundCount(count int) (RoundCount, error) {
	if count < 0 {
		return RoundCount{}, errors.New("ラウンド数は0以上である必要があります")
	}
	if count > 50 { // 現実的な上限
		return RoundCount{}, errors.New("ラウンド数は50以下である必要があります")
	}
	return RoundCount{count: count}, nil
}

// NewZeroRoundCount は0のRoundCountを作成する
func NewZeroRoundCount() RoundCount {
	return RoundCount{count: 0}
}

// Count は数値を返す
func (r RoundCount) Count() int {
	return r.count
}

// Increment はラウンド数を1増やす
func (r RoundCount) Increment() (RoundCount, error) {
	return NewRoundCount(r.count + 1)
}

// Decrement はラウンド数を1減らす
func (r RoundCount) Decrement() (RoundCount, error) {
	if r.count <= 0 {
		return RoundCount{}, errors.New("ラウンド数は0以下にできません")
	}
	return NewRoundCount(r.count - 1)
}

// IsEmpty は空（0ラウンド）かを判定する
func (r RoundCount) IsEmpty() bool {
	return r.count == 0
}

// IsMinimal は最小ラウンド数（1ラウンド）かを判定する
func (r RoundCount) IsMinimal() bool {
	return r.count == 1
}

// IsStandard は標準的なラウンド数（3-4ラウンド）かを判定する
func (r RoundCount) IsStandard() bool {
	return r.count >= 3 && r.count <= 4
}

// IsProductiveSession は生産的なセッション（3ラウンド以上）かを判定する
func (r RoundCount) IsProductiveSession() bool {
	return r.count >= 3
}

// IsIntensiveSession は集中セッション（6ラウンド以上）かを判定する
func (r RoundCount) IsIntensiveSession() bool {
	return r.count >= 6
}

// CanAccommodateRoundOrder は指定したラウンド順序を受け入れ可能かを判定する
func (r RoundCount) CanAccommodateRoundOrder(roundOrder int) bool {
	return roundOrder >= 1 && roundOrder <= r.count
}

// HasReachedMaximum は最大ラウンド数に達しているかを判定する
func (r RoundCount) HasReachedMaximum(maxRounds int) bool {
	return r.count >= maxRounds
}

// RemainingRounds は最大ラウンド数までの残りラウンド数を返す
func (r RoundCount) RemainingRounds(maxRounds int) int {
	remaining := maxRounds - r.count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// ToSessionQuality はセッション品質を評価する
func (r RoundCount) ToSessionQuality() string {
	if r.count == 0 {
		return "未開始"
	} else if r.count == 1 {
		return "短時間"
	} else if r.count >= 2 && r.count <= 4 {
		return "標準"
	} else if r.count >= 5 && r.count <= 8 {
		return "集中"
	} else {
		return "超集中"
	}
}

// String は文字列表現を返す
func (r RoundCount) String() string {
	if r.count == 0 {
		return "ラウンドなし"
	}
	return fmt.Sprintf("%dラウンド", r.count)
}

// Equals は別のRoundCountと等価かを判定する
func (r RoundCount) Equals(other RoundCount) bool {
	return r.count == other.count
}

// Compare は他のRoundCountとの大小を比較する
// 戻り値: -1 (小さい), 0 (等しい), 1 (大きい)
func (r RoundCount) Compare(other RoundCount) int {
	if r.count < other.count {
		return -1
	}
	if r.count > other.count {
		return 1
	}
	return 0
}