package statistics

import (
	"errors"
	"fmt"
)

// TotalRounds は統計の総ラウンド数を表すValue Object
type TotalRounds struct {
	count int
}

// NewTotalRounds は新しいTotalRoundsを作成する
func NewTotalRounds(count int) (TotalRounds, error) {
	if count < 0 {
		return TotalRounds{}, errors.New("総ラウンド数は0以上である必要があります")
	}
	if count > 1000 { // 現実的な上限
		return TotalRounds{}, errors.New("総ラウンド数は1000以下である必要があります")
	}
	return TotalRounds{count: count}, nil
}

// NewZeroTotalRounds は0のTotalRoundsを作成する
func NewZeroTotalRounds() TotalRounds {
	return TotalRounds{count: 0}
}

// Count は数値を返す
func (t TotalRounds) Count() int {
	return t.count
}

// Increment はラウンド数を1増やす
func (t TotalRounds) Increment() (TotalRounds, error) {
	return NewTotalRounds(t.count + 1)
}

// Add は指定した数だけラウンド数を増やす
func (t TotalRounds) Add(rounds int) (TotalRounds, error) {
	return NewTotalRounds(t.count + rounds)
}

// IsZero は0ラウンドかを判定する
func (t TotalRounds) IsZero() bool {
	return t.count == 0
}

// IsMinimal は最小ラウンド数（1ラウンド）かを判定する
func (t TotalRounds) IsMinimal() bool {
	return t.count == 1
}

// IsProductiveDay は生産的な日（10ラウンド以上）かを判定する
func (t TotalRounds) IsProductiveDay() bool {
	return t.count >= 10
}

// IsProductiveWeek は生産的な週（50ラウンド以上）かを判定する
func (t TotalRounds) IsProductiveWeek() bool {
	return t.count >= 50
}

// ToQualityLevel はラウンド数から品質レベルを返す
func (t TotalRounds) ToQualityLevel() string {
	if t.count == 0 {
		return "未開始"
	} else if t.count <= 3 {
		return "軽量"
	} else if t.count <= 10 {
		return "標準"
	} else if t.count <= 20 {
		return "集中"
	} else {
		return "超集中"
	}
}

// String は文字列表現を返す
func (t TotalRounds) String() string {
	if t.count == 0 {
		return "ラウンドなし"
	}
	return fmt.Sprintf("%dラウンド", t.count)
}

// Equals は別のTotalRoundsと等価かを判定する
func (t TotalRounds) Equals(other TotalRounds) bool {
	return t.count == other.count
}

// Compare は他のTotalRoundsとの大小を比較する
// 戻り値: -1 (小さい), 0 (等しい), 1 (大きい)
func (t TotalRounds) Compare(other TotalRounds) int {
	if t.count < other.count {
		return -1
	}
	if t.count > other.count {
		return 1
	}
	return 0
}