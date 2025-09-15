package valueobject

import (
	"errors"
	"fmt"
)

// RoundOrder はセッション内でのラウンド順序を表すValue Object
type RoundOrder struct {
	order int
}

// NewRoundOrder は新しいRoundOrderを作成する
func NewRoundOrder(order int) (RoundOrder, error) {
	if order < 1 {
		return RoundOrder{}, errors.New("ラウンド順序は1以上である必要があります")
	}
	if order > 50 { // 現実的な上限
		return RoundOrder{}, errors.New("ラウンド順序は50以下である必要があります")
	}
	return RoundOrder{order: order}, nil
}

// NewFirstRoundOrder は最初のラウンド順序（1）を作成する
func NewFirstRoundOrder() RoundOrder {
	return RoundOrder{order: 1}
}

// Order は順序値を返す
func (r RoundOrder) Order() int {
	return r.order
}

// Next は次の順序のRoundOrderを返す
func (r RoundOrder) Next() (RoundOrder, error) {
	return NewRoundOrder(r.order + 1)
}

// Previous は前の順序のRoundOrderを返す
func (r RoundOrder) Previous() (RoundOrder, error) {
	if r.order <= 1 {
		return RoundOrder{}, errors.New("最初のラウンドには前の順序が存在しません")
	}
	return NewRoundOrder(r.order - 1)
}

// IsFirst は最初のラウンドかを判定する
func (r RoundOrder) IsFirst() bool {
	return r.order == 1
}

// IsWithinLimit は指定した上限内かを判定する（セッション制約チェック用）
func (r RoundOrder) IsWithinLimit(maxRounds int) bool {
	return r.order <= maxRounds
}

// CanBeNextOf は指定したRoundOrderの次として適切かを判定する
func (r RoundOrder) CanBeNextOf(previous RoundOrder) bool {
	return r.order == previous.order+1
}

// DistanceFrom は他のRoundOrderとの距離を返す
func (r RoundOrder) DistanceFrom(other RoundOrder) int {
	if r.order > other.order {
		return r.order - other.order
	}
	return other.order - r.order
}

// IsConsecutiveWith は連続する順序かを判定する
func (r RoundOrder) IsConsecutiveWith(other RoundOrder) bool {
	return r.DistanceFrom(other) == 1
}

// String は文字列表現を返す
func (r RoundOrder) String() string {
	return fmt.Sprintf("%d番目", r.order)
}

// Equals は別のRoundOrderと等価かを判定する
func (r RoundOrder) Equals(other RoundOrder) bool {
	return r.order == other.order
}

// Compare は他のRoundOrderとの大小を比較する
// 戻り値: -1 (小さい), 0 (等しい), 1 (大きい)
func (r RoundOrder) Compare(other RoundOrder) int {
	if r.order < other.order {
		return -1
	}
	if r.order > other.order {
		return 1
	}
	return 0
}