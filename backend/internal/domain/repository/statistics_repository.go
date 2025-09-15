package repository

import (
	"context"
	"time"

	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	statisticsVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/statistics"
	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
)

// StatisticsRepository は統計情報集約の永続化を担当するリポジトリ
// DDD原則: Statistics集約（Daily/Hourly/Weekly）に対するリポジトリ
type StatisticsRepository interface {
	// 日別統計
	GetDailyStatistics(ctx context.Context, userID userVO.UserID, date string) (*entity.DailyStatistics, error)
	SaveDailyStatistics(ctx context.Context, stats *entity.DailyStatistics) error
	GetDailyStatisticsByPeriod(ctx context.Context, userID userVO.UserID, period statisticsVO.StatisticsPeriod) ([]*entity.DailyStatistics, error)

	// 時間別統計
	GetHourlyStatistics(ctx context.Context, userID userVO.UserID, date string, hour int) (*entity.HourlyStatistics, error)
	SaveHourlyStatistics(ctx context.Context, stats *entity.HourlyStatistics) error
	GetHourlyStatisticsByPeriod(ctx context.Context, userID userVO.UserID, period statisticsVO.StatisticsPeriod) ([]*entity.HourlyStatistics, error)

	// 週別統計
	GetWeeklyStatistics(ctx context.Context, userID userVO.UserID, weekStart, weekEnd string) (*entity.WeeklyStatistics, error)
	SaveWeeklyStatistics(ctx context.Context, stats *entity.WeeklyStatistics) error
	GetWeeklyStatisticsByPeriod(ctx context.Context, userID userVO.UserID, period statisticsVO.StatisticsPeriod) ([]*entity.WeeklyStatistics, error)

	// 集約用クエリは削除済み - UseCase層でMapper使用に変更

	// 統計更新処理（Round完了時の集約更新）
	UpdateStatisticsWithRound(ctx context.Context, userID userVO.UserID, round *entity.Round, timestamp time.Time) error
}
