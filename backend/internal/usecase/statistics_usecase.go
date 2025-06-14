package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// StatisticsUsecase は統計情報を取得するためのユースケースインターフェース
type StatisticsUsecase interface {
	// GetFocusTrend は集中度推移を取得する
	GetFocusTrend(ctx context.Context, userID uuid.UUID, period string, startDate, endDate *time.Time) (*model.FocusTrendResponse, error)

	// GetFocusHeatmap は集中度ヒートマップを取得する
	GetFocusHeatmap(ctx context.Context, userID uuid.UUID, period string, startDate, endDate *time.Time) (*model.FocusHeatmapResponse, error)
}

// statisticsUsecase は統計情報を取得するためのユースケースの実装（新エラーハンドリング対応版）
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

// GetFocusTrend は集中度推移を取得する（新エラーハンドリング対応版）
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
			// カスタム期間が指定されているが日付が指定されていない場合は1週間に設定
			statsPeriod = model.NewLastWeekPeriod()
		} else {
			statsPeriod = model.NewCustomPeriod(*startDate, *endDate)
		}
	default:
		// デフォルトは1週間
		statsPeriod = model.NewLastWeekPeriod()
	}

	// リポジトリから集中度推移を取得
	trendItems, err := uc.statsRepo.GetFocusTrend(ctx, userID, statsPeriod)
	if err != nil {
		uc.logger.Errorf("集中度推移取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			// データが存在しない場合は空のレスポンスを返す（エラーとしない）
			return &model.FocusTrendResponse{Items: []*model.FocusTrendItem{}}, nil
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// 集中度推移をレスポンス形式に変換
	response := &model.FocusTrendResponse{
		Items: trendItems,
	}

	uc.logger.Infof("集中度推移取得成功: period=%s, items=%d", period, len(trendItems))
	return response, nil
}

// GetFocusHeatmap は集中度ヒートマップを取得する（新エラーハンドリング対応版）
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
		uc.logger.Errorf("集中度ヒートマップ取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			// データが存在しない場合は空のレスポンスを返す（エラーとしない）
			return &model.FocusHeatmapResponse{Items: []*model.FocusHeatmapItem{}}, nil
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// 集中度ヒートマップをレスポンス形式に変換
	response := &model.FocusHeatmapResponse{
		Items: heatmapItems,
	}

	uc.logger.Infof("集中度ヒートマップ取得成功: period=%s, items=%d", period, len(heatmapItems))
	return response, nil
}
