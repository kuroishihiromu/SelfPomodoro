package statistics

import (
	"errors"
	"fmt"
	"time"
)

// Date は統計の日付を表すValue Object（YYYY-MM-DD形式）
type Date struct {
	date string
}

// NewDate は新しいDateを作成する
func NewDate(date string) (Date, error) {
	// YYYY-MM-DD形式の検証
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return Date{}, errors.New("日付はYYYY-MM-DD形式で入力してください")
	}
	return Date{date: date}, nil
}

// NewDateFromTime はtime.Timeから新しいDateを作成する
func NewDateFromTime(t time.Time) Date {
	return Date{date: t.Format("2006-01-02")}
}

// NewToday は今日の日付でDateを作成する
func NewToday() Date {
	return NewDateFromTime(time.Now())
}

// String は日付文字列を返す
func (d Date) String() string {
	return d.date
}

// ToTime はtime.Timeに変換する
func (d Date) ToTime() time.Time {
	t, _ := time.Parse("2006-01-02", d.date)
	return t
}

// IsToday は今日の日付かを判定する
func (d Date) IsToday() bool {
	today := NewToday()
	return d.Equals(today)
}

// IsYesterday は昨日の日付かを判定する
func (d Date) IsYesterday() bool {
	yesterday := NewDateFromTime(time.Now().AddDate(0, 0, -1))
	return d.Equals(yesterday)
}

// IsFuture は未来の日付かを判定する
func (d Date) IsFuture() bool {
	return d.ToTime().After(time.Now())
}

// IsPast は過去の日付かを判定する
func (d Date) IsPast() bool {
	return d.ToTime().Before(time.Now().Truncate(24 * time.Hour))
}

// DaysFromToday は今日からの日数を返す（負数は過去、正数は未来）
func (d Date) DaysFromToday() int {
	today := time.Now().Truncate(24 * time.Hour)
	targetDate := d.ToTime()
	diff := targetDate.Sub(today)
	return int(diff.Hours() / 24)
}

// GetWeekday は曜日を返す
func (d Date) GetWeekday() time.Weekday {
	return d.ToTime().Weekday()
}

// IsWeekend は週末（土日）かを判定する
func (d Date) IsWeekend() bool {
	weekday := d.GetWeekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// IsWeekday は平日かを判定する
func (d Date) IsWeekday() bool {
	return !d.IsWeekend()
}

// GetYear は年を返す
func (d Date) GetYear() int {
	return d.ToTime().Year()
}

// GetMonth は月を返す
func (d Date) GetMonth() time.Month {
	return d.ToTime().Month()
}

// GetDay は日を返す
func (d Date) GetDay() int {
	return d.ToTime().Day()
}

// GetWeekOfYear は年の第何週かを返す
func (d Date) GetWeekOfYear() (int, int) {
	return d.ToTime().ISOWeek()
}

// GetMonthStart は月の最初の日を返す
func (d Date) GetMonthStart() Date {
	t := d.ToTime()
	monthStart := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	return NewDateFromTime(monthStart)
}

// GetMonthEnd は月の最後の日を返す
func (d Date) GetMonthEnd() Date {
	t := d.ToTime()
	monthEnd := time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location())
	return NewDateFromTime(monthEnd)
}

// AddDays は指定した日数を加算する
func (d Date) AddDays(days int) Date {
	newTime := d.ToTime().AddDate(0, 0, days)
	return NewDateFromTime(newTime)
}

// AddMonths は指定した月数を加算する
func (d Date) AddMonths(months int) Date {
	newTime := d.ToTime().AddDate(0, months, 0)
	return NewDateFromTime(newTime)
}

// AddYears は指定した年数を加算する
func (d Date) AddYears(years int) Date {
	newTime := d.ToTime().AddDate(years, 0, 0)
	return NewDateFromTime(newTime)
}

// ToJapaneseFormat は日本語形式の文字列を返す
func (d Date) ToJapaneseFormat() string {
	t := d.ToTime()
	return fmt.Sprintf("%d年%d月%d日", t.Year(), t.Month(), t.Day())
}

// Equals は別のDateと等価かを判定する
func (d Date) Equals(other Date) bool {
	return d.date == other.date
}

// Before は他のDateより前かを判定する
func (d Date) Before(other Date) bool {
	return d.ToTime().Before(other.ToTime())
}

// After は他のDateより後かを判定する
func (d Date) After(other Date) bool {
	return d.ToTime().After(other.ToTime())
}

// Compare は他のDateとの大小を比較する
// 戻り値: -1 (小さい), 0 (等しい), 1 (大きい)
func (d Date) Compare(other Date) int {
	if d.Before(other) {
		return -1
	}
	if d.After(other) {
		return 1
	}
	return 0
}