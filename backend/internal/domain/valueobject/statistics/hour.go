package statistics

import (
	"errors"
	"fmt"
)

// Hour は時間別統計の時間を表すValue Object（0-23）
type Hour struct {
	hour int
}

// NewHour は新しいHourを作成する
func NewHour(hour int) (Hour, error) {
	if hour < 0 || hour > 23 {
		return Hour{}, errors.New("時間は0-23の間で設定してください")
	}
	return Hour{hour: hour}, nil
}

// Hour は時間値を返す
func (h Hour) Hour() int {
	return h.hour
}

// Next は次の時間のHourを返す（23の次は0）
func (h Hour) Next() Hour {
	nextHour := (h.hour + 1) % 24
	hour, _ := NewHour(nextHour)
	return hour
}

// Previous は前の時間のHourを返す（0の前は23）
func (h Hour) Previous() Hour {
	prevHour := (h.hour - 1 + 24) % 24
	hour, _ := NewHour(prevHour)
	return hour
}

// IsMorning は午前中（6-11時）かを判定する
func (h Hour) IsMorning() bool {
	return h.hour >= 6 && h.hour <= 11
}

// IsAfternoon は午後（12-17時）かを判定する
func (h Hour) IsAfternoon() bool {
	return h.hour >= 12 && h.hour <= 17
}

// IsEvening は夕方（18-21時）かを判定する
func (h Hour) IsEvening() bool {
	return h.hour >= 18 && h.hour <= 21
}

// IsNight は夜間（22-5時）かを判定する
func (h Hour) IsNight() bool {
	return h.hour >= 22 || h.hour <= 5
}

// IsBusinessHours は営業時間（9-17時）かを判定する
func (h Hour) IsBusinessHours() bool {
	return h.hour >= 9 && h.hour <= 17
}

// IsEarlyMorning は早朝（5-8時）かを判定する
func (h Hour) IsEarlyMorning() bool {
	return h.hour >= 5 && h.hour <= 8
}

// IsLateNight は深夜（22-4時）かを判定する
func (h Hour) IsLateNight() bool {
	return h.hour >= 22 || h.hour <= 4
}

// GetTimeOfDay は時間帯の文字列を返す
func (h Hour) GetTimeOfDay() string {
	if h.IsEarlyMorning() {
		return "早朝"
	} else if h.IsMorning() {
		return "午前"
	} else if h.IsAfternoon() {
		return "午後"
	} else if h.IsEvening() {
		return "夕方"
	} else {
		return "夜間"
	}
}

// To12HourFormat は12時間形式の文字列を返す
func (h Hour) To12HourFormat() string {
	if h.hour == 0 {
		return "12 AM"
	} else if h.hour < 12 {
		return fmt.Sprintf("%d AM", h.hour)
	} else if h.hour == 12 {
		return "12 PM"
	} else {
		return fmt.Sprintf("%d PM", h.hour-12)
	}
}

// To24HourFormat は24時間形式の文字列を返す
func (h Hour) To24HourFormat() string {
	return fmt.Sprintf("%02d:00", h.hour)
}

// String は文字列表現を返す
func (h Hour) String() string {
	return h.To24HourFormat()
}

// Equals は別のHourと等価かを判定する
func (h Hour) Equals(other Hour) bool {
	return h.hour == other.hour
}

// Compare は他のHourとの大小を比較する
// 戻り値: -1 (小さい), 0 (等しい), 1 (大きい)
func (h Hour) Compare(other Hour) int {
	if h.hour < other.hour {
		return -1
	}
	if h.hour > other.hour {
		return 1
	}
	return 0
}