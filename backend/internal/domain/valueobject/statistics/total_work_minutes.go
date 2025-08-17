package statistics

import (
	"errors"
	"fmt"
)

// TotalWorkMinutes は統計の総作業時間を表すValue Object
type TotalWorkMinutes struct {
	minutes int
}

// NewTotalWorkMinutes は新しいTotalWorkMinutesを作成する
func NewTotalWorkMinutes(minutes int) (TotalWorkMinutes, error) {
	if minutes < 0 {
		return TotalWorkMinutes{}, errors.New("総作業時間は0分以上である必要があります")
	}
	return TotalWorkMinutes{minutes: minutes}, nil
}

// NewZeroTotalWorkMinutes は0分のTotalWorkMinutesを作成する
func NewZeroTotalWorkMinutes() TotalWorkMinutes {
	return TotalWorkMinutes{minutes: 0}
}

// Minutes は分単位の値を返す
func (t TotalWorkMinutes) Minutes() int {
	return t.minutes
}

// Hours は時間単位の値を返す
func (t TotalWorkMinutes) Hours() float64 {
	return float64(t.minutes) / 60.0
}

// Add は他のTotalWorkMinutesを加算する
func (t TotalWorkMinutes) Add(other TotalWorkMinutes) (TotalWorkMinutes, error) {
	return NewTotalWorkMinutes(t.minutes + other.minutes)
}

// AddMinutes は指定分数を加算する
func (t TotalWorkMinutes) AddMinutes(minutes int) (TotalWorkMinutes, error) {
	return NewTotalWorkMinutes(t.minutes + minutes)
}

// IsZero は0分かを判定する
func (t TotalWorkMinutes) IsZero() bool {
	return t.minutes == 0
}

// IsProductiveDay は生産的な日かを判定する（2時間以上）
func (t TotalWorkMinutes) IsProductiveDay() bool {
	return t.minutes >= 120
}

// IsProductiveWeek は生産的な週かを判定する（10時間以上）
func (t TotalWorkMinutes) IsProductiveWeek() bool {
	return t.minutes >= 600
}

// IsIntensiveSession は集中的なセッションかを判定する（4時間以上）
func (t TotalWorkMinutes) IsIntensiveSession() bool {
	return t.minutes >= 240
}


// ToHourMinuteString は時間:分の文字列表現を返す
func (t TotalWorkMinutes) ToHourMinuteString() string {
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
func (t TotalWorkMinutes) String() string {
	return t.ToHourMinuteString()
}

// Equals は別のTotalWorkMinutesと等価かを判定する
func (t TotalWorkMinutes) Equals(other TotalWorkMinutes) bool {
	return t.minutes == other.minutes
}

// Compare は他のTotalWorkMinutesとの大小を比較する
// 戻り値: -1 (小さい), 0 (等しい), 1 (大きい)
func (t TotalWorkMinutes) Compare(other TotalWorkMinutes) int {
	if t.minutes < other.minutes {
		return -1
	}
	if t.minutes > other.minutes {
		return 1
	}
	return 0
}