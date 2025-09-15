package statistics

import (
	"errors"
	"fmt"
)

// DaysActive は統計のアクティブ日数を表すValue Object
type DaysActive struct {
	days int
}

// NewDaysActive は新しいDaysActiveを作成する
func NewDaysActive(days int) (DaysActive, error) {
	if days < 0 {
		return DaysActive{}, errors.New("アクティブ日数は0以上である必要があります")
	}
	if days > 7 {
		return DaysActive{}, errors.New("アクティブ日数は7日以下である必要があります")
	}
	return DaysActive{days: days}, nil
}

// NewZeroDaysActive は0日のDaysActiveを作成する
func NewZeroDaysActive() DaysActive {
	return DaysActive{days: 0}
}

// Days は日数を返す
func (d DaysActive) Days() int {
	return d.days
}

// Increment はアクティブ日数を1増やす
func (d DaysActive) Increment() (DaysActive, error) {
	return NewDaysActive(d.days + 1)
}

// Add は指定した日数だけアクティブ日数を増やす
func (d DaysActive) Add(days int) (DaysActive, error) {
	return NewDaysActive(d.days + days)
}

// IsZero は0日かを判定する
func (d DaysActive) IsZero() bool {
	return d.days == 0
}

// IsWeekend は週末のみ活動（1-2日）かを判定する
func (d DaysActive) IsWeekend() bool {
	return d.days >= 1 && d.days <= 2
}

// IsRegular は定期的な活動（3-4日）かを判定する
func (d DaysActive) IsRegular() bool {
	return d.days >= 3 && d.days <= 4
}

// IsFrequent は頻繁な活動（5-6日）かを判定する
func (d DaysActive) IsFrequent() bool {
	return d.days >= 5 && d.days <= 6
}

// IsDaily は毎日活動（7日）かを判定する
func (d DaysActive) IsDaily() bool {
	return d.days == 7
}

// IsConsistent は一貫した活動（5日以上）かを判定する
func (d DaysActive) IsConsistent() bool {
	return d.days >= 5
}

// GetActivityRate は活動率を返す（0.0-1.0）
func (d DaysActive) GetActivityRate() float64 {
	return float64(d.days) / 7.0
}

// GetActivityPercentage は活動率をパーセンテージで返す（0-100）
func (d DaysActive) GetActivityPercentage() int {
	return int(d.GetActivityRate() * 100)
}

// ToActivityLevel はアクティブ日数からアクティビティレベルを返す
func (d DaysActive) ToActivityLevel() string {
	if d.days == 0 {
		return "非アクティブ"
	} else if d.days <= 2 {
		return "低活動"
	} else if d.days <= 4 {
		return "中活動"
	} else if d.days <= 6 {
		return "高活動"
	} else {
		return "毎日活動"
	}
}

// String は文字列表現を返す
func (d DaysActive) String() string {
	if d.days == 0 {
		return "アクティブ日なし"
	}
	return fmt.Sprintf("%d日アクティブ", d.days)
}

// Equals は別のDaysActiveと等価かを判定する
func (d DaysActive) Equals(other DaysActive) bool {
	return d.days == other.days
}

// Compare は他のDaysActiveとの大小を比較する
// 戻り値: -1 (小さい), 0 (等しい), 1 (大きい)
func (d DaysActive) Compare(other DaysActive) int {
	if d.days < other.days {
		return -1
	}
	if d.days > other.days {
		return 1
	}
	return 0
}