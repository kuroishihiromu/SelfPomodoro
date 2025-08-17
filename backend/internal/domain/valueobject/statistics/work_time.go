package statistics

import (
	"fmt"
	"time"

	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
)

// WorkTime は作業時間を表すValue Object
type WorkTime struct {
	minutes int
}

// NewWorkTime は新しい作業時間を作成する（分単位）
func NewWorkTime(minutes int) (WorkTime, error) {
	if minutes < 0 {
		return WorkTime{}, appErrors.NewValidationError(fmt.Sprintf("作業時間は0以上である必要があります: %d", minutes))
	}
	if minutes > 24*60 { // 24時間を超える場合はエラー
		return WorkTime{}, appErrors.NewValidationError(fmt.Sprintf("作業時間は24時間以下である必要があります: %d分", minutes))
	}
	return WorkTime{minutes: minutes}, nil
}

// NewWorkTimeFromDuration は期間から新しい作業時間を作成する
func NewWorkTimeFromDuration(duration time.Duration) (WorkTime, error) {
	minutes := int(duration.Minutes())
	return NewWorkTime(minutes)
}

// Minutes は作業時間を分で返す
func (wt WorkTime) Minutes() int {
	return wt.minutes
}

// Hours は作業時間を時間で返す
func (wt WorkTime) Hours() float64 {
	return float64(wt.minutes) / 60.0
}

// Duration は作業時間をtime.Durationで返す
func (wt WorkTime) Duration() time.Duration {
	return time.Duration(wt.minutes) * time.Minute
}

// IsShort は短時間作業かを判定する（30分未満）
func (wt WorkTime) IsShort() bool {
	return wt.minutes < 30
}

// IsMedium は中程度の作業時間かを判定する（30分以上2時間未満）
func (wt WorkTime) IsMedium() bool {
	return wt.minutes >= 30 && wt.minutes < 120
}

// IsLong は長時間作業かを判定する（2時間以上）
func (wt WorkTime) IsLong() bool {
	return wt.minutes >= 120
}

// Add は作業時間を加算する
func (wt WorkTime) Add(other WorkTime) WorkTime {
	total := wt.minutes + other.minutes
	// エラーは発生しないように既存の値を返す（ドメインロジックで制御）
	result, _ := NewWorkTime(total)
	return result
}

// Equals は別のWorkTimeと等価かを判定する
func (wt WorkTime) Equals(other WorkTime) bool {
	return wt.minutes == other.minutes
}

// String は文字列表現を返す
func (wt WorkTime) String() string {
	if wt.minutes >= 60 {
		hours := wt.minutes / 60
		mins := wt.minutes % 60
		if mins > 0 {
			return fmt.Sprintf("%d時間%d分", hours, mins)
		}
		return fmt.Sprintf("%d時間", hours)
	}
	return fmt.Sprintf("%d分", wt.minutes)
}