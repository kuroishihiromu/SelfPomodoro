package statistics

import (
	"errors"
	"fmt"
)

// TotalBreakMinutes は統計の総休憩時間を表すValue Object
type TotalBreakMinutes struct {
	minutes int
}

// NewTotalBreakMinutes は新しいTotalBreakMinutesを作成する
func NewTotalBreakMinutes(minutes int) (TotalBreakMinutes, error) {
	if minutes < 0 {
		return TotalBreakMinutes{}, errors.New("総休憩時間は0分以上である必要があります")
	}
	return TotalBreakMinutes{minutes: minutes}, nil
}

// NewZeroTotalBreakMinutes は0分のTotalBreakMinutesを作成する
func NewZeroTotalBreakMinutes() TotalBreakMinutes {
	return TotalBreakMinutes{minutes: 0}
}

// Minutes は分単位の値を返す
func (t TotalBreakMinutes) Minutes() int {
	return t.minutes
}

// Hours は時間単位の値を返す
func (t TotalBreakMinutes) Hours() float64 {
	return float64(t.minutes) / 60.0
}

// Add は他のTotalBreakMinutesを加算する
func (t TotalBreakMinutes) Add(other TotalBreakMinutes) (TotalBreakMinutes, error) {
	return NewTotalBreakMinutes(t.minutes + other.minutes)
}

// AddMinutes は指定分数を加算する
func (t TotalBreakMinutes) AddMinutes(minutes int) (TotalBreakMinutes, error) {
	return NewTotalBreakMinutes(t.minutes + minutes)
}

// IsZero は0分かを判定する
func (t TotalBreakMinutes) IsZero() bool {
	return t.minutes == 0
}



// ToHourMinuteString は時間:分の文字列表現を返す
func (t TotalBreakMinutes) ToHourMinuteString() string {
	if t.minutes == 0 {
		return "休憩なし"
	}
	if t.minutes < 60 {
		return fmt.Sprintf("%d分", t.minutes)
	}
	hours := t.minutes / 60
	minutes := t.minutes % 60
	if minutes == 0 {
		return fmt.Sprintf("%d時間", hours)
	}
	return fmt.Sprintf("%d時間%d分", hours, minutes)
}

// String は文字列表現を返す
func (t TotalBreakMinutes) String() string {
	return t.ToHourMinuteString()
}

// Equals は別のTotalBreakMinutesと等価かを判定する
func (t TotalBreakMinutes) Equals(other TotalBreakMinutes) bool {
	return t.minutes == other.minutes
}

// Compare は他のTotalBreakMinutesとの大小を比較する
// 戻り値: -1 (小さい), 0 (等しい), 1 (大きい)
func (t TotalBreakMinutes) Compare(other TotalBreakMinutes) int {
	if t.minutes < other.minutes {
		return -1
	}
	if t.minutes > other.minutes {
		return 1
	}
	return 0
}