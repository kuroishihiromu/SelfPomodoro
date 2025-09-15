package statistics

import (
	"errors"
	"fmt"
	"time"
)

// WeekPeriod は週期間を表すValue Object（月曜日開始〜日曜日終了）
type WeekPeriod struct {
	weekStart Date
	weekEnd   Date
}

// NewWeekPeriod は新しいWeekPeriodを作成する
func NewWeekPeriod(weekStart, weekEnd Date) (WeekPeriod, error) {
	// 開始日が終了日より後でないことを確認
	if weekStart.After(weekEnd) {
		return WeekPeriod{}, errors.New("週の開始日は終了日より前である必要があります")
	}
	
	// 期間が7日を超えないことを確認
	startTime := weekStart.ToTime()
	endTime := weekEnd.ToTime()
	duration := endTime.Sub(startTime)
	if duration.Hours() > 7*24 {
		return WeekPeriod{}, errors.New("週期間は7日以内である必要があります")
	}
	
	// 開始日が月曜日であることを確認
	if weekStart.GetWeekday() != time.Monday {
		return WeekPeriod{}, errors.New("週の開始日は月曜日である必要があります")
	}
	
	// 終了日が日曜日であることを確認
	if weekEnd.GetWeekday() != time.Sunday {
		return WeekPeriod{}, errors.New("週の終了日は日曜日である必要があります")
	}
	
	return WeekPeriod{
		weekStart: weekStart,
		weekEnd:   weekEnd,
	}, nil
}

// NewWeekPeriodFromDate は指定した日付を含む週のWeekPeriodを作成する
func NewWeekPeriodFromDate(date Date) WeekPeriod {
	dateTime := date.ToTime()
	
	// 月曜日を見つける
	daysSinceMonday := int(dateTime.Weekday() - time.Monday)
	if daysSinceMonday < 0 {
		daysSinceMonday += 7 // 日曜日の場合
	}
	
	mondayTime := dateTime.AddDate(0, 0, -daysSinceMonday)
	sundayTime := mondayTime.AddDate(0, 0, 6)
	
	weekStart := NewDateFromTime(mondayTime)
	weekEnd := NewDateFromTime(sundayTime)
	
	period, _ := NewWeekPeriod(weekStart, weekEnd)
	return period
}

// NewCurrentWeek は現在の週のWeekPeriodを作成する
func NewCurrentWeek() WeekPeriod {
	return NewWeekPeriodFromDate(NewToday())
}

// WeekStart は週の開始日を返す
func (w WeekPeriod) WeekStart() Date {
	return w.weekStart
}

// WeekEnd は週の終了日を返す
func (w WeekPeriod) WeekEnd() Date {
	return w.weekEnd
}

// WeekStartString は週の開始日の文字列を返す
func (w WeekPeriod) WeekStartString() string {
	return w.weekStart.String()
}

// WeekEndString は週の終了日の文字列を返す
func (w WeekPeriod) WeekEndString() string {
	return w.weekEnd.String()
}

// Contains は指定した日付がこの週に含まれるかを判定する
func (w WeekPeriod) Contains(date Date) bool {
	return !date.Before(w.weekStart) && !date.After(w.weekEnd)
}

// IsCurrentWeek は現在の週かを判定する
func (w WeekPeriod) IsCurrentWeek() bool {
	currentWeek := NewCurrentWeek()
	return w.Equals(currentWeek)
}

// IsLastWeek は先週かを判定する
func (w WeekPeriod) IsLastWeek() bool {
	lastWeek := NewCurrentWeek().Previous()
	return w.Equals(lastWeek)
}

// IsFuture は未来の週かを判定する
func (w WeekPeriod) IsFuture() bool {
	return w.weekStart.IsFuture()
}

// IsPast は過去の週かを判定する
func (w WeekPeriod) IsPast() bool {
	return w.weekEnd.IsPast()
}

// Next は次の週のWeekPeriodを返す
func (w WeekPeriod) Next() WeekPeriod {
	nextWeekStart := w.weekStart.AddDays(7)
	nextWeekEnd := w.weekEnd.AddDays(7)
	period, _ := NewWeekPeriod(nextWeekStart, nextWeekEnd)
	return period
}

// Previous は前の週のWeekPeriodを返す
func (w WeekPeriod) Previous() WeekPeriod {
	prevWeekStart := w.weekStart.AddDays(-7)
	prevWeekEnd := w.weekEnd.AddDays(-7)
	period, _ := NewWeekPeriod(prevWeekStart, prevWeekEnd)
	return period
}

// GetDates は週に含まれる全ての日付を返す
func (w WeekPeriod) GetDates() []Date {
	dates := make([]Date, 7)
	for i := 0; i < 7; i++ {
		dates[i] = w.weekStart.AddDays(i)
	}
	return dates
}

// GetWeekdayDates は平日のみを返す
func (w WeekPeriod) GetWeekdayDates() []Date {
	var weekdays []Date
	for i := 0; i < 5; i++ { // 月〜金
		weekdays = append(weekdays, w.weekStart.AddDays(i))
	}
	return weekdays
}

// GetWeekendDates は週末のみを返す
func (w WeekPeriod) GetWeekendDates() []Date {
	return []Date{
		w.weekStart.AddDays(5), // 土曜日
		w.weekStart.AddDays(6), // 日曜日
	}
}

// GetWeekNumber は年の第何週かを返す
func (w WeekPeriod) GetWeekNumber() (int, int) {
	return w.weekStart.GetWeekOfYear()
}

// GetYear は年を返す
func (w WeekPeriod) GetYear() int {
	return w.weekStart.GetYear()
}

// GetMonth は月を返す（週の開始日の月）
func (w WeekPeriod) GetMonth() time.Month {
	return w.weekStart.GetMonth()
}

// ToDisplayString は表示用の文字列を返す
func (w WeekPeriod) ToDisplayString() string {
	return fmt.Sprintf("%s 〜 %s", w.weekStart.String(), w.weekEnd.String())
}

// ToJapaneseDisplayString は日本語表示用の文字列を返す
func (w WeekPeriod) ToJapaneseDisplayString() string {
	return fmt.Sprintf("%s 〜 %s", w.weekStart.ToJapaneseFormat(), w.weekEnd.ToJapaneseFormat())
}

// String は文字列表現を返す
func (w WeekPeriod) String() string {
	return w.ToDisplayString()
}

// Equals は別のWeekPeriodと等価かを判定する
func (w WeekPeriod) Equals(other WeekPeriod) bool {
	return w.weekStart.Equals(other.weekStart) && w.weekEnd.Equals(other.weekEnd)
}

// Before は他のWeekPeriodより前かを判定する
func (w WeekPeriod) Before(other WeekPeriod) bool {
	return w.weekEnd.Before(other.weekStart)
}

// After は他のWeekPeriodより後かを判定する
func (w WeekPeriod) After(other WeekPeriod) bool {
	return w.weekStart.After(other.weekEnd)
}

// Compare は他のWeekPeriodとの大小を比較する
// 戻り値: -1 (小さい), 0 (等しい), 1 (大きい)
func (w WeekPeriod) Compare(other WeekPeriod) int {
	if w.Before(other) {
		return -1
	}
	if w.After(other) {
		return 1
	}
	return 0
}