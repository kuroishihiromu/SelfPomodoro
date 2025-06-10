package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/container"
	domainErrors "github.com/tsunakit99/selfpomodoro/internal/domain/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/usecase"
)

// Global container (Lambda再利用最適化)
var globalContainer container.Container

// init はLambda init phaseで実行
func init() {
	globalContainer = container.NewLambdaContainer()
}

// StatisticsHandler はDI Container使用版の統計ハンドラー
type StatisticsHandler struct {
	useCases *usecase.UseCases
	logger   logger.Logger
}

// handler はLambdaのエントリーポイント（DI Container版）
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 1. Container初期化
	if err := globalContainer.Initialize(ctx); err != nil {
		return errorResponse(http.StatusInternalServerError, "サービス初期化エラー"), nil
	}

	// 2. Dependencies取得（Infrastructure依存なし！）
	useCases := globalContainer.GetUseCases()
	logger := globalContainer.GetLogger()

	// 3. Handler初期化
	statsHandler := &StatisticsHandler{
		useCases: useCases,
		logger:   logger,
	}

	// 4. 認証・User存在確認（統一処理）
	userID, err := statsHandler.authenticateAndValidateUser(ctx, request)
	if err != nil {
		return statsHandler.handleError(err), nil
	}

	// 5. GET以外は許可しない
	if request.HTTPMethod != "GET" {
		return errorResponse(http.StatusMethodNotAllowed, "GETメソッドのみ許可されています"), nil
	}

	// 6. パスによるルーティング
	if strings.Contains(request.Path, "/focus-trend") {
		return statsHandler.handleGetFocusTrend(ctx, request, userID)
	} else if strings.Contains(request.Path, "/focus-heatmap") {
		return statsHandler.handleGetFocusHeatmap(ctx, request, userID)
	}

	return errorResponse(http.StatusNotFound, "無効なパス"), nil
}

// authenticateAndValidateUser は認証・User存在確認の統一処理
func (h *StatisticsHandler) authenticateAndValidateUser(ctx context.Context, request events.APIGatewayProxyRequest) (uuid.UUID, error) {
	return h.useCases.Auth.AuthenticateAndValidateUser(ctx, request)
}

// handleGetFocusTrend は集中度トレンド取得を処理
func (h *StatisticsHandler) handleGetFocusTrend(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	// クエリパラメータから期間を取得
	period := request.QueryStringParameters["period"] // "week", "month", "custom"

	var startDate, endDate *time.Time

	// カスタム期間が指定されている場合は、開始日と終了日を取得
	if period == "custom" {
		startDateStr := request.QueryStringParameters["start_date"]
		endDateStr := request.QueryStringParameters["end_date"]

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
	response, err := h.useCases.Statistics.GetFocusTrend(ctx, userID, period, startDate, endDate)
	if err != nil {
		h.logger.Errorf("集中度トレンド取得エラー: %v", err)
		return h.handleError(err), nil
	}

	return successResponse(http.StatusOK, response), nil
}

// handleGetFocusHeatmap は集中度ヒートマップ取得を処理
func (h *StatisticsHandler) handleGetFocusHeatmap(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	// クエリパラメータから期間を取得
	period := request.QueryStringParameters["period"] // "week", "month", "custom"

	var startDate, endDate *time.Time

	// カスタム期間が指定されている場合は、開始日と終了日を取得
	if period == "custom" {
		startDateStr := request.QueryStringParameters["start_date"]
		endDateStr := request.QueryStringParameters["end_date"]

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
	response, err := h.useCases.Statistics.GetFocusHeatmap(ctx, userID, period, startDate, endDate)
	if err != nil {
		h.logger.Errorf("集中度ヒートマップ取得エラー: %v", err)
		return h.handleError(err), nil
	}

	return successResponse(http.StatusOK, response), nil
}

// handleError はドメインエラーを統一処理
func (h *StatisticsHandler) handleError(err error) events.APIGatewayProxyResponse {
	if appErr, ok := err.(*domainErrors.AppError); ok {
		return errorResponse(appErr.Status, appErr.Error())
	}

	h.logger.Errorf("予期しないエラータイプ: %T, %v", err, err)
	return errorResponse(http.StatusInternalServerError, "内部エラーが発生しました")
}

func successResponse(statusCode int, data interface{}) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(data)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type,Authorization",
			"Access-Control-Allow-Methods": "GET,POST,PATCH,DELETE,OPTIONS",
		},
		Body: string(body),
	}
}

func errorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(map[string]string{"error": message})
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type,Authorization",
			"Access-Control-Allow-Methods": "GET,POST,PATCH,DELETE,OPTIONS",
		},
		Body: string(body),
	}
}

func main() {
	lambda.Start(handler)
}
