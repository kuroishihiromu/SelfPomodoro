package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/database"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository"
	"github.com/tsunakit99/selfpomodoro/internal/usecase"
)

// StatisticsHandler はLambda用の統計ハンドラー
type StatisticsHandler struct {
	statsUseCase usecase.StatisticsUsecase
	logger       logger.Logger
}

// handler はLambdaのエントリーポイント
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 依存関係の初期化
	cfg, err := config.Load()
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "設定読み込みエラー"), nil
	}

	appLogger, err := logger.NewLogger(cfg.LogLevel, cfg.Environment)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "ロガー初期化エラー"), nil
	}

	// PostgreSQL接続
	postgresDB, err := database.NewPostgresDB(cfg, appLogger)
	if err != nil {
		appLogger.Errorf("PostgreSQL接続エラー: %v", err)
		return errorResponse(http.StatusInternalServerError, "データベース接続エラー"), nil
	}
	defer postgresDB.Close()

	// リポジトリとユースケースの初期化
	repositoryFactory := repository.NewRepositoryFactory(postgresDB, nil, appLogger)
	statsUseCase := usecase.NewStatisticsUsecase(repositoryFactory.Statistics, appLogger)

	// ハンドラーの初期化
	statsHandler := &StatisticsHandler{
		statsUseCase: statsUseCase,
		logger:       appLogger,
	}

	// ユーザーID取得
	userID, err := getUserIDFromRequest(request)
	if err != nil {
		appLogger.Errorf("ユーザーID取得エラー: %v", err)
		return errorResponse(http.StatusUnauthorized, "認証エラー"), nil
	}

	// GET以外は許可しない
	if request.HTTPMethod != "GET" {
		return errorResponse(http.StatusMethodNotAllowed, "GETメソッドのみ許可されています"), nil
	}

	// パスによるルーティング
	if strings.Contains(request.Path, "/focus-trend") {
		return statsHandler.handleGetFocusTrend(ctx, request, userID)
	} else if strings.Contains(request.Path, "/focus-heatmap") {
		return statsHandler.handleGetFocusHeatmap(ctx, request, userID)
	}

	return errorResponse(http.StatusNotFound, "無効なパス"), nil
}

// getUserIDFromRequest はリクエストからユーザーIDを取得する
func getUserIDFromRequest(request events.APIGatewayProxyRequest) (uuid.UUID, error) {
	// 開発環境用の簡易認証（dev-tokenの場合）
	authHeader := request.Headers["Authorization"]
	if authHeader == "" {
		authHeader = request.Headers["authorization"] // 小文字の場合もチェック
	}

	if authHeader == "Bearer dev-token" {
		// 開発用固定ユーザーID
		return uuid.Parse("00000000-0000-0000-0000-000000000001")
	}

	return uuid.Nil, fmt.Errorf("認証情報が見つかりません")
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
	response, err := h.statsUseCase.GetFocusTrend(ctx, userID, period, startDate, endDate)
	if err != nil {
		h.logger.Errorf("集中度トレンド取得エラー: %v", err)
		return errorResponse(http.StatusInternalServerError, "集中度トレンド取得に失敗しました"), nil
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
	response, err := h.statsUseCase.GetFocusHeatmap(ctx, userID, period, startDate, endDate)
	if err != nil {
		h.logger.Errorf("集中度ヒートマップ取得エラー: %v", err)
		return errorResponse(http.StatusInternalServerError, "集中度ヒートマップ取得に失敗しました"), nil
	}

	return successResponse(http.StatusOK, response), nil
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
