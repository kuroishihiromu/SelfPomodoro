package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/usecase"
)

// StatisticsHandler は統計情報に関するリクエストを処理するハンドラー
type StatisticsHandler struct {
	statsUseCase usecase.StatisticsUsecase
	logger       logger.Logger
}

// NewStatisticsHandler は新しい StatisticsHandler を生成する
func NewStatisticsHandler(statsUseCase usecase.StatisticsUsecase, logger logger.Logger) *StatisticsHandler {
	return &StatisticsHandler{
		statsUseCase: statsUseCase,
		logger:       logger,
	}
}

// getUserIDFromContext はコンテキストからユーザーIDを取得する
func (h *StatisticsHandler) getUserIDFromContext(c echo.Context) (uuid.UUID, error) {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return uuid.Nil, errors.New("ユーザーIDが見つかりません")
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, errors.New("無効なユーザーID")
	}
	return id, nil
}

// GetFocusTrend は集中度トレンドを取得するエンドポイントを処理
// GET /statistics/focus_trend
func (h *StatisticsHandler) GetFocusTrend(c echo.Context) error {
	// ユーザーIDを取得
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// クエリパラメータから期間を取得
	period := c.QueryParam("period") // "week", "month", "custom"

	var startDate, endDate *time.Time

	// カスタム期間が指定されている場合は、開始日と終了日を取得
	if period == "custom" {
		startDateStr := c.QueryParam("start_date")
		endDateStr := c.QueryParam("end_date")

		if startDateStr != "" {
			parsedStartDate, err := time.Parse("2006-01-02", startDateStr)
			if err == nil {
				startDate = &parsedStartDate
			} else {
				h.logger.Warnf("無効な開始日: %s", startDateStr)
			}
		}

		if endDateStr != "" {
			parsedEndDate, err := time.Parse("2006-01-02", endDateStr)
			if err == nil {
				endDate = &parsedEndDate
			} else {
				h.logger.Warnf("無効な終了日: %s", endDateStr)
			}
		}
	}

	// ユースケースを呼び出して集中度トレンドを取得
	response, err := h.statsUseCase.GetFocusTrend(c.Request().Context(), userID, period, startDate, endDate)
	if err != nil {
		h.logger.Errorf("集中度トレンド取得エラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "集中度トレンド取得に失敗しました"})
	}

	return c.JSON(http.StatusOK, response)
}

// GetFocusHeatmap は集中度ヒートマップを取得するエンドポイントを処理
// GET /statistics/focus_heatmap
func (h *StatisticsHandler) GetFocusHeatmap(c echo.Context) error {
	// ユーザーIDを取得
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// クエリパラメータから期間を取得
	period := c.QueryParam("period") // "week", "month", "custom"

	var startDate, endDate *time.Time

	// カスタム期間が指定されている場合は、開始日と終了日を取得
	if period == "custom" {
		startDateStr := c.QueryParam("start_date")
		endDateStr := c.QueryParam("end_date")

		if startDateStr != "" {
			parsedStartDate, err := time.Parse("2006-01-02", startDateStr)
			if err == nil {
				startDate = &parsedStartDate
			} else {
				h.logger.Warnf("無効な開始日: %s", startDateStr)
			}
		}

		if endDateStr != "" {
			parsedEndDate, err := time.Parse("2006-01-02", endDateStr)
			if err == nil {
				endDate = &parsedEndDate
			} else {
				h.logger.Warnf("無効な終了日: %s", endDateStr)
			}
		}
	}

	// ユースケースを呼び出して集中度ヒートマップを取得
	response, err := h.statsUseCase.GetFocusHeatmap(c.Request().Context(), userID, period, startDate, endDate)
	if err != nil {
		h.logger.Errorf("集中度ヒートマップ取得エラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "集中度ヒートマップ取得に失敗しました"})
	}

	return c.JSON(http.StatusOK, response)
}
