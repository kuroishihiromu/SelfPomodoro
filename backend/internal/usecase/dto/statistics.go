package dto

import "time"

// DailyStatisticsResponse は日別統計のAPIレスポンスを表すDTO
type DailyStatisticsResponse struct {
	UserID        string    `json:"user_id"`
	Date          string    `json:"date"`
	TotalRounds   int       `json:"total_rounds"`
	AvgFocusScore float64   `json:"avg_focus_score"`
	TotalWorkMin  int       `json:"total_work_min"`
	TotalBreakMin int       `json:"total_break_min"`
	SessionCount  int       `json:"session_count"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// WeeklyStatisticsResponse は週別統計のAPIレスポンスを表すDTO
type WeeklyStatisticsResponse struct {
	UserID        string    `json:"user_id"`
	WeekStart     string    `json:"week_start"`
	WeekEnd       string    `json:"week_end"`
	TotalRounds   int       `json:"total_rounds"`
	AvgFocusScore float64   `json:"avg_focus_score"`
	TotalWorkMin  int       `json:"total_work_min"`
	TotalBreakMin int       `json:"total_break_min"`
	SessionCount  int       `json:"session_count"`
	DaysActive    int       `json:"days_active"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// FocusTrendResponse はAPIレスポンス用のDTO
type FocusTrendResponse struct {
	Date       string  `json:"date"`
	FocusScore float64 `json:"focus_score"`
}

// FocusHeatmapResponse はAPIレスポンス用のDTO
type FocusHeatmapResponse struct {
	Date       string  `json:"date"`
	Hour       int     `json:"hour"`
	FocusScore float64 `json:"focus_score"`
}