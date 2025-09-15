package valueobject

import (
	"errors"
	"fmt"
	"time"
)

// デフォルト作業時間定数
const (
	DefaultWorkTime = 25
)

// WorkTime は作業時間を表すValue Object
type WorkTime struct {
	minutes int
}

// NewWorkTime は新しいWorkTimeを作成する
func NewWorkTime(minutes int) (WorkTime, error) {
	if minutes < 1 || minutes > 120 {
		return WorkTime{}, errors.New("作業時間は1分から120分の間で設定してください")
	}
	return WorkTime{minutes: minutes}, nil
}

// MustNewWorkTime はパニックを起こす可能性があるWorkTime作成（テスト用）
func MustNewWorkTime(minutes int) WorkTime {
	wt, err := NewWorkTime(minutes)
	if err != nil {
		panic(err)
	}
	return wt
}

// Minutes は分単位の値を返す
func (w WorkTime) Minutes() int {
	return w.minutes
}

// Duration はtime.Durationを返す
func (w WorkTime) Duration() time.Duration {
	return time.Duration(w.minutes) * time.Minute
}

// Equals は別のWorkTimeと等価かを判定する
func (w WorkTime) Equals(other WorkTime) bool {
	return w.minutes == other.minutes
}

// IsStandard は標準的な作業時間かを判定する（25分）
func (w WorkTime) IsStandard() bool {
	return w.minutes == DefaultWorkTime
}

// IsShort は短い作業時間かを判定する（15分以下）
func (w WorkTime) IsShort() bool {
	return w.minutes <= 15
}

// IsLong は長い作業時間かを判定する（45分以上）
func (w WorkTime) IsLong() bool {
	return w.minutes >= 45
}

// String は文字列表現を返す
func (w WorkTime) String() string {
	return fmt.Sprintf("%d分", w.minutes)
}