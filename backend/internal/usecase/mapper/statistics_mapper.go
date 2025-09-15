package mapper

import (
	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	"github.com/tsunakit99/selfpomodoro/internal/usecase/dto"
)

// StatisticsMapper は統計関連のマッピングを担当する
type StatisticsMapper struct{}

// NewStatisticsMapper は新しいStatisticsMapperを作成する
func NewStatisticsMapper() *StatisticsMapper {
	return &StatisticsMapper{}
}

// ToDailyStatisticsResponse はドメインモデルからAPIレスポンス形式に変換する
func (m *StatisticsMapper) ToDailyStatisticsResponse(daily *entity.DailyStatistics) *dto.DailyStatisticsResponse {
	return &dto.DailyStatisticsResponse{
		UserID:        daily.UserID().String(),
		Date:          daily.DateString(),
		TotalRounds:   daily.TotalRoundsCount(),
		AvgFocusScore: daily.AvgFocusScore().Score(),
		TotalWorkMin:  daily.TotalWorkMin().Minutes(),
		TotalBreakMin: daily.TotalBreakMin().Minutes(),
		SessionCount:  daily.SessionCountValue(),
		UpdatedAt:     daily.UpdatedAt(),
	}
}

// ToWeeklyStatisticsResponse はドメインモデルからAPIレスポンス形式に変換する
func (m *StatisticsMapper) ToWeeklyStatisticsResponse(weekly *entity.WeeklyStatistics) *dto.WeeklyStatisticsResponse {
	return &dto.WeeklyStatisticsResponse{
		UserID:        weekly.UserID().String(),
		WeekStart:     weekly.WeekStart(),
		WeekEnd:       weekly.WeekEnd(),
		TotalRounds:   weekly.TotalRoundsCount(),
		AvgFocusScore: weekly.AvgFocusScore().Score(),
		TotalWorkMin:  weekly.TotalWorkMin().Minutes(),
		TotalBreakMin: weekly.TotalBreakMin().Minutes(),
		SessionCount:  weekly.SessionCountValue(),
		DaysActive:    weekly.DaysActiveValue(),
		UpdatedAt:     weekly.UpdatedAt(),
	}
}

// ToFocusTrendResponse は統計データをFocusTrendResponse形式に変換する（汎用メソッド）
func (m *StatisticsMapper) ToFocusTrendResponse(daily *entity.DailyStatistics) *dto.FocusTrendResponse {
	return &dto.FocusTrendResponse{
		Date:       daily.DateString(),
		FocusScore: daily.AvgFocusScore().Score(),
	}
}

// ToFocusTrendResponseFromDaily は統計データをFocusTrendResponse形式に変換する
func (m *StatisticsMapper) ToFocusTrendResponseFromDaily(daily *entity.DailyStatistics) *dto.FocusTrendResponse {
	return &dto.FocusTrendResponse{
		Date:       daily.DateString(),
		FocusScore: daily.AvgFocusScore().Score(),
	}
}

// ToFocusTrendResponseFromWeekly は統計データをFocusTrendResponse形式に変換する
func (m *StatisticsMapper) ToFocusTrendResponseFromWeekly(weekly *entity.WeeklyStatistics) *dto.FocusTrendResponse {
	return &dto.FocusTrendResponse{
		Date:       weekly.WeekStart(), // 週の開始日（月曜日）を表示
		FocusScore: weekly.AvgFocusScore().Score(),
	}
}

// ToFocusHeatmapResponse は統計データをFocusHeatmapResponse形式に変換する
func (m *StatisticsMapper) ToFocusHeatmapResponse(hourly *entity.HourlyStatistics) *dto.FocusHeatmapResponse {
	return &dto.FocusHeatmapResponse{
		Date:       hourly.DateString(),
		Hour:       hourly.HourValue(),
		FocusScore: hourly.AvgFocusScore().Score(),
	}
}