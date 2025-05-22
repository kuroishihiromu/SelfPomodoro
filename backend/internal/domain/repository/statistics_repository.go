package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
)

// StatisticsRepository は統計情報を取得するためのリポジトリインターフェース
type StatisticsRepository interface {
	// GetFocusTrend は指定期間内の日別集中度統計を取得する
	GetFocusTrend(ctx context.Context, userID uuid.UUID, period *model.StatisticsPeriod) ([]*model.FocusTrendItem, error)

	// GetFocusHeatmap は指定期間内の時間帯別集中度統計を取得する
	GetFocusHeatmap(ctx context.Context, userID uuid.UUID, period *model.StatisticsPeriod) ([]*model.FocusHeatmapItem, error)

	// GetAvgFocusScoreByDate は指定日の平均集中度を取得する
	GetAvgFocusScoreByDate(ctx context.Context, userID uuid.UUID, date time.Time) (float64, error)

	// GetAvgFocusScoreByHour は指定日時の平均集中度を取得する
	GetAvgFocusScoreByHour(ctx context.Context, userID uuid.UUID, date time.Time, hour int) (float64, error)
}
