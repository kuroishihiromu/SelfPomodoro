package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/usecase/dto"
	"github.com/tsunakit99/selfpomodoro/internal/usecase/mapper"
	statisticsVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/statistics"
	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// StatisticsUsecase は統計情報を取得するためのユースケースインターフェース
type StatisticsUsecase interface {
	// GetFocusTrend は指定日を開始日として1週間分の集中度推移を取得する
	GetFocusTrend(ctx context.Context, userID uuid.UUID, startDate string) ([]*dto.FocusTrendResponse, error)

	// GetFocusHeatmap は集中度ヒートマップを取得する
	GetFocusHeatmap(ctx context.Context, userID uuid.UUID, period string, startDate, endDate *time.Time) ([]*dto.FocusHeatmapResponse, error)
}

// statisticsUsecase は統計情報を取得するためのユースケースの実装（新エラーハンドリング対応版）
type statisticsUsecase struct {
	statsRepo   repository.StatisticsRepository
	statsMapper *mapper.StatisticsMapper
	logger      logger.Logger
}

// NewStatisticsUsecase は統計情報を取得するためのユースケースを生成する
func NewStatisticsUsecase(statsRepo repository.StatisticsRepository, logger logger.Logger) StatisticsUsecase {
	return &statisticsUsecase{
		statsRepo:   statsRepo,
		statsMapper: mapper.NewStatisticsMapper(),
		logger:      logger,
	}
}

// GetFocusTrend は指定日を開始日として1週間分の集中度推移を取得する
func (uc *statisticsUsecase) GetFocusTrend(ctx context.Context, userID uuid.UUID, startDate string) ([]*dto.FocusTrendResponse, error) {
	// ユーザーIDをValue Objectに変換
	userIDVO, err := userVO.NewUserIDFromString(userID.String())
	if err != nil {
		return nil, appErrors.NewInternalError(err)
	}

	// 開始日をパース（YYYY-MMDD形式 -> YYYY-MM-DD形式に変換）
	var startDateTime time.Time
	if len(startDate) == 8 { // YYYY-MMDD形式
		// YYYY-MMDD -> YYYY-MM-DD に変換
		formattedDate := startDate[:4] + "-" + startDate[4:6] + "-" + startDate[6:8]
		startDateTime, err = time.Parse("2006-01-02", formattedDate)
		if err != nil {
			uc.logger.Errorf("無効な開始日形式: %s", startDate)
			return nil, appErrors.NewValidationError("無効な日付形式です。YYYY-MMDD形式で指定してください")
		}
	} else {
		return nil, appErrors.NewValidationError("日付はYYYY-MMDD形式で指定してください")
	}

	// 1週間後の終了日を計算
	endDateTime := startDateTime.AddDate(0, 0, 6) // 7日間（開始日含む）

	// 期間を作成
	statsPeriod := statisticsVO.NewStatisticsPeriod(startDateTime, endDateTime)

	// リポジトリから日別統計を取得
	dailyStats, err := uc.statsRepo.GetDailyStatisticsByPeriod(ctx, userIDVO, statsPeriod)
	uc.logger.Infof("1週間分日別統計取得: startDate=%s, endDate=%s", startDateTime.Format("2006-01-02"), endDateTime.Format("2006-01-02"))

	if err != nil {
		uc.logger.Errorf("集中度推移取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			// データが存在しない場合は空のスライスを返す（エラーとしない）
			return []*dto.FocusTrendResponse{}, nil
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// 日別統計からFocusTrendResponseに変換
	if len(dailyStats) == 0 {
		return []*dto.FocusTrendResponse{}, nil
	}

	// マッパーを使用して変換
	trendResponses := make([]*dto.FocusTrendResponse, len(dailyStats))
	for i, dailyStat := range dailyStats {
		trendResponses[i] = uc.statsMapper.ToFocusTrendResponse(dailyStat)
	}

	uc.logger.Infof("集中度推移取得成功: startDate=%s, items=%d", startDate, len(trendResponses))
	return trendResponses, nil
}

// GetFocusHeatmap は集中度ヒートマップを取得する（新エラーハンドリング対応版）
func (uc *statisticsUsecase) GetFocusHeatmap(ctx context.Context, userID uuid.UUID, period string, startDate, endDate *time.Time) ([]*dto.FocusHeatmapResponse, error) {
	// ユーザーIDをValue Objectに変掻
	userIDVO, err := userVO.NewUserIDFromString(userID.String())
	if err != nil {
		return nil, appErrors.NewInternalError(err)
	}

	// 期間を決定
	var statsPeriod statisticsVO.StatisticsPeriod

	switch period {
	case "week":
		statsPeriod = statisticsVO.NewLastWeekPeriod()
	case "month":
		statsPeriod = statisticsVO.NewLastMonthPeriod()
	case "custom":
		if startDate == nil || endDate == nil {
			// カスタム期間が指定されているが日付が指定されていない場合は1週間に設定
			statsPeriod = statisticsVO.NewLastWeekPeriod()
		} else {
			statsPeriod = statisticsVO.NewStatisticsPeriod(*startDate, *endDate)
		}
	default:
		// デフォルトは1週間
		statsPeriod = statisticsVO.NewLastWeekPeriod()
	}

	// リポジトリから時間別統計を取得（ヒートマップ用）
	hourlyStats, err := uc.statsRepo.GetHourlyStatisticsByPeriod(ctx, userIDVO, statsPeriod)
	if err != nil {
		uc.logger.Errorf("集中度ヒートマップ取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			// データが存在しない場合は空のスライスを返す（エラーとしない）
			return []*dto.FocusHeatmapResponse{}, nil
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// 時間別統計からFocusHeatmapResponseに変換
	if len(hourlyStats) == 0 {
		return []*dto.FocusHeatmapResponse{}, nil
	}

	// マッパーを使用して変換
	heatmapResponses := make([]*dto.FocusHeatmapResponse, len(hourlyStats))
	for i, hourlyStat := range hourlyStats {
		heatmapResponses[i] = uc.statsMapper.ToFocusHeatmapResponse(hourlyStat)
	}

	uc.logger.Infof("集中度ヒートマップ取得成功: period=%s, items=%d", period, len(heatmapResponses))
	return heatmapResponses, nil
}
