package valueobject

import (
	"errors"
	"fmt"
	"time"
)

// デフォルト休憩時間定数（ラウンド用）
const (
	DefaultRoundBreakTime = 5
)

// BreakTime は休憩時間を表すValue Object（ラウンド用）
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

// NewDefaultBreakTime はデフォルトのBreakTimeを作成する（5分）
func NewDefaultBreakTime() BreakTime {
	return BreakTime{minutes: DefaultRoundBreakTime}
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

// IsStandard は標準的な休憩時間かを判定する（5分）
func (b BreakTime) IsStandard() bool {
	return b.minutes == DefaultRoundBreakTime
}

// IsNoBreak は休憩なしかを判定する（0分）
func (b BreakTime) IsNoBreak() bool {
	return b.minutes == 0
}

// IsShort は短い休憩時間かを判定する（3分以下）
func (b BreakTime) IsShort() bool {
	return b.minutes <= 3
}

// IsLong は長い休憩時間かを判定する（10分以上）
func (b BreakTime) IsLong() bool {
	return b.minutes >= 10
}

// String は文字列表現を返す
func (b BreakTime) String() string {
	if b.IsNoBreak() {
		return "休憩なし"
	}
	return fmt.Sprintf("%d分", b.minutes)
}