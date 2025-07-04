package dynamodb

import (
	"fmt"
	"time"
)

// DynamoDBキー生成ユーティリティ
// Single Table Designのパーティションキー・ソートキー生成を担当

// UserPartitionKey はユーザー関連エンティティのパーティションキーを生成する
func UserPartitionKey(userID string) string {
	return "USER#" + userID
}

// UserConfigSortKey はユーザー設定のソートキーを生成する
func UserConfigSortKey() string {
	return "CONFIG"
}

// TaskSortKey はタスクのソートキーを生成する
func TaskSortKey(taskID string) string {
	return "TASK#" + taskID
}

// SessionSortKey はセッションのソートキーを生成する
func SessionSortKey(date, sessionID string) string {
	return "SESSION#" + date + "#" + sessionID
}

// RoundSortKey はラウンドのソートキーを生成する
func RoundSortKey(date, sessionID string, roundOrder int) string {
	return "ROUND#" + date + "#" + sessionID + "#" + fmt.Sprintf("%03d", roundOrder)
}

// DailyStatsSortKey は日別統計のソートキーを生成する
func DailyStatsSortKey(date string) string {
	return "STATS_DAILY#" + date
}

// HourlyStatsSortKey は時間別統計のソートキーを生成する
func HourlyStatsSortKey(date string, hour int) string {
	return "STATS_HOURLY#" + date + "#" + fmt.Sprintf("%02d", hour)
}

// WeeklyStatsSortKey は週別統計のソートキーを生成する
func WeeklyStatsSortKey(weekStart string) string {
	return "STATS_WEEKLY#" + weekStart
}

// OptimizationRoundSortKey はラウンド最適化ログのソートキーを生成する（統一形式）
func OptimizationRoundSortKey(timestamp string) string {
	return "OPTIMIZATION_LOG#ROUND#" + timestamp
}

// OptimizationSessionSortKey はセッション最適化ログのソートキーを生成する（統一形式）
func OptimizationSessionSortKey(timestamp string) string {
	return "OPTIMIZATION_LOG#SESSION#" + timestamp
}


// DateFromTime は時刻から日付文字列(YYYY-MM-DD)を生成する
func DateFromTime(t time.Time) string {
	return t.Format("2006-01-02")
}

// QueryPrefixes はクエリ用のプレフィックスを定義

// TaskQueryPrefix はタスク検索用のプレフィックスを返す
func TaskQueryPrefix() string {
	return "TASK#"
}

// SessionQueryPrefix はセッション検索用のプレフィックスを返す
func SessionQueryPrefix(date string) string {
	if date == "" {
		return "SESSION#"
	}
	return "SESSION#" + date + "#"
}

// RoundQueryPrefix はラウンド検索用のプレフィックスを返す
func RoundQueryPrefix(date string) string {
	if date == "" {
		return "ROUND#"
	}
	return "ROUND#" + date + "#"
}

// DailyStatsQueryPrefix は日別統計検索用のプレフィックスを返す
func DailyStatsQueryPrefix() string {
	return "STATS_DAILY#"
}

// HourlyStatsQueryPrefix は時間別統計検索用のプレフィックスを返す
func HourlyStatsQueryPrefix() string {
	return "STATS_HOURLY#"
}

// WeeklyStatsQueryPrefix は週別統計検索用のプレフィックスを返す
func WeeklyStatsQueryPrefix() string {
	return "STATS_WEEKLY#"
}

// OptimizationLogQueryPrefix は最適化ログ検索用のプレフィックスを返す
func OptimizationLogQueryPrefix(logType string) string {
	if logType == "" {
		return "OPTIMIZATION_LOG#"
	}
	return "OPTIMIZATION_LOG#" + logType + "#"
}

// Range queries support

// SessionQueryRange はセッション期間検索のキー範囲を返す
func SessionQueryRange(startDate, endDate string) (string, string) {
	startSK := "SESSION#" + startDate
	endSK := "SESSION#" + endDate + "#zzz"
	return startSK, endSK
}

// RoundQueryRange はラウンド期間検索のキー範囲を返す
func RoundQueryRange(startDate, endDate string) (string, string) {
	startSK := "ROUND#" + startDate
	endSK := "ROUND#" + endDate + "#zzz"
	return startSK, endSK
}

// DailyStatsQueryRange は日別統計期間検索のキー範囲を返す
func DailyStatsQueryRange(startDate, endDate string) (string, string) {
	startSK := "STATS_DAILY#" + startDate
	endSK := "STATS_DAILY#" + endDate
	return startSK, endSK
}

// HourlyStatsQueryRange は時間別統計期間検索のキー範囲を返す
func HourlyStatsQueryRange(startDate, endDate string) (string, string) {
	startSK := "STATS_HOURLY#" + startDate + "#00"
	endSK := "STATS_HOURLY#" + endDate + "#23"
	return startSK, endSK
}

// WeeklyStatsQueryRange は週別統計期間検索のキー範囲を返す
func WeeklyStatsQueryRange(startWeek, endWeek string) (string, string) {
	startSK := "STATS_WEEKLY#" + startWeek
	endSK := "STATS_WEEKLY#" + endWeek
	return startSK, endSK
}