package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/errors"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// StatsticsUsecase は統計情報を取得するためのユースケースインターフェース
type StatisticsUsecase interface {
	// GetFocusTrend は集中度推移を取得する
	GetFocusTrend(ctx context.Context, userID uuid.UUID, period string, startDate, endDate *time.Time) (*model.FocusTrendResponse, error)

	// GetFocusHeatmap は集中度ヒートマップを取得する
	GetFocusHeatmap(ctx context.Context, userID uuid.UUID, period string, startDate, endDate *time.Time) (*model.FocusHeatmapResponse, error)
}

// statisticsUsecase は統計情報を取得するためのユースケースの実装
type statisticsUsecase struct {
	statsRepo repository.StatisticsRepository
	logger    logger.Logger
}

// NewStatisticsUsecase は統計情報を取得するためのユースケースを生成する
func NewStatisticsUsecase(statsRepo repository.StatisticsRepository, logger logger.Logger) StatisticsUsecase {
	return &statisticsUsecase{
		statsRepo: statsRepo,
		logger:    logger,
	}
}

// GetFocusTrend は集中度推移を取得する
func (uc *statisticsUsecase) GetFocusTrend(ctx context.Context, userID uuid.UUID, period string, startDate, endDate *time.Time) (*model.FocusTrendResponse, error) {
	// 期間を決定
	var statsPeriod *model.StatisticsPeriod

	switch period {
	case "week":
		statsPeriod = model.NewLastWeekPeriod()
	case "month":
		statsPeriod = model.NewLastMonthPeriod()
	case "custom":
		if startDate == nil || endDate == nil {
			//　カスタム期間が指定されているが日付が指定されていない場合は1週間に設定
			statsPeriod = model.NewLastWeekPeriod()
		} else {
			statsPeriod = model.NewCustomPeriod(*startDate, *endDate)
		}
	default:
		// デフォルトは1週間
		statsPeriod = model.NewLastWeekPeriod()
	}

	// リポジトリから集中度数推移を取得
	trendItems, err := uc.statsRepo.GetFocusTrend(ctx, userID, statsPeriod)
	if err != nil {
		uc.logger.Error("集中度推移取得エラー", err)
		return nil, errors.NewInternalError(err)
	}

	// 集中度推移をレスポンス形式に変換
	response := &model.FocusTrendResponse{
		Items: trendItems,
	}

	return response, nil
}

// GetFocusHeatmap は集中度ヒートマップを取得する
func (uc *statisticsUsecase) GetFocusHeatmap(ctx context.Context, userID uuid.UUID, period string, startDate, endDate *time.Time) (*model.FocusHeatmapResponse, error) {
	// 期間を決定
	var statsPeriod *model.StatisticsPeriod

	switch period {
	case "week":
		statsPeriod = model.NewLastWeekPeriod()
	case "month":
		statsPeriod = model.NewLastMonthPeriod()
	case "custom":
		if startDate == nil || endDate == nil {
			// カスタム期間が指定されているが日付が指定されていない場合は1週間に設定
			statsPeriod = model.NewLastWeekPeriod()
		} else {
			statsPeriod = model.NewCustomPeriod(*startDate, *endDate)
		}
	default:
		// デフォルトは1週間
		statsPeriod = model.NewLastWeekPeriod()
	}

	// リポジトリから集中度ヒートマップを取得
	heatmapItems, err := uc.statsRepo.GetFocusHeatmap(ctx, userID, statsPeriod)
	if err != nil {
		uc.logger.Error("集中度ヒートマップ取得エラー", err)
		return nil, errors.NewInternalError(err)
	}

	// 集中度ヒートマップをレスポンス形式に変換
	response := &model.FocusHeatmapResponse{
		Items: heatmapItems,
	}

	return response, nil
}
