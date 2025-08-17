package valueobject

import (
	"errors"
	"fmt"
	"time"
)

// デフォルト休憩時間定数（セッション用）
const (
	DefaultSessionBreakTime = 15
)

// BreakTime は休憩時間を表すValue Object（セッション用）
type BreakTime struct {
	minutes int
}

// NewBreakTime は新しいBreakTimeを作成する
func NewBreakTime(minutes int) (BreakTime, error) {
	if minutes < 0 || minutes > 60 {
		return BreakTime{}, errors.New("休憩時間は0分から60分の間で設定してください")
	}
	return BreakTime{minutes: minutes}, nil
}

// NewDefaultBreakTime はデフォルトのBreakTimeを作成する（15分）
func NewDefaultBreakTime() BreakTime {
	return BreakTime{minutes: DefaultSessionBreakTime}
}

// MustNewBreakTime はパニックを起こす可能性があるBreakTime作成（テスト用）
func MustNewBreakTime(minutes int) BreakTime {
	bt, err := NewBreakTime(minutes)
	if err != nil {
		panic(err)
	}
	return bt
}

// Minutes は分単位の値を返す
func (b BreakTime) Minutes() int {
	return b.minutes
}

// Duration はtime.Durationを返す
func (b BreakTime) Duration() time.Duration {
	return time.Duration(b.minutes) * time.Minute
}

// Equals は別のBreakTimeと等価かを判定する
func (b BreakTime) Equals(other BreakTime) bool {
	return b.minutes == other.minutes
}

// IsStandard は標準的な休憩時間かを判定する（15分）
func (b BreakTime) IsStandard() bool {
	return b.minutes == DefaultSessionBreakTime
}

// IsLongBreak は長い休憩時間かを判定する（15分以上）
func (b BreakTime) IsLongBreak() bool {
	return b.minutes >= DefaultSessionBreakTime
}

// IsNoBreak は休憩なしかを判定する（0分）
func (b BreakTime) IsNoBreak() bool {
	return b.minutes == 0
}

// String は文字列表現を返す
func (b BreakTime) String() string {
	if b.IsNoBreak() {
		return "休憩なし"
	}
	return fmt.Sprintf("%d分", b.minutes)
}