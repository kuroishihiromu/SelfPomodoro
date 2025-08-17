package valueobject

import (
	"errors"
	"fmt"
)

// TotalWorkMinutes は総作業時間を表すValue Object
type TotalWorkMinutes struct {
	minutes int
}

// NewTotalWorkMinutes は新しいTotalWorkMinutesを作成する
func NewTotalWorkMinutes(minutes int) (TotalWorkMinutes, error) {
	if minutes < 0 {
		return TotalWorkMinutes{}, errors.New("総作業時間は0分以上である必要があります")
	}
	if minutes > 24*60 { // 24時間上限
		return TotalWorkMinutes{}, errors.New("総作業時間は24時間以下である必要があります")
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

// IsProductiveSession は生産的なセッションかを判定する（60分以上）
func (t TotalWorkMinutes) IsProductiveSession() bool {
	return t.minutes >= 60
}

// IsLongSession は長時間セッションかを判定する（2時間以上）
func (t TotalWorkMinutes) IsLongSession() bool {
	return t.minutes >= 120
}

// String は文字列表現を返す
func (t TotalWorkMinutes) String() string {
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

// Equals は別のTotalWorkMinutesと等価かを判定する
func (t TotalWorkMinutes) Equals(other TotalWorkMinutes) bool {
	return t.minutes == other.minutes
}