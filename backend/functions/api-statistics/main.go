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
	httpError "github.com/tsunakit99/selfpomodoro/internal/handler"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/auth"
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
		return createErrorResponse(http.StatusInternalServerError, "INTERNAL_ERROR", "サービス初期化エラー"), nil
	}

	// 2. Dependencies取得（Infrastructure依存なし！）
	useCases := globalContainer.GetUseCases()
	logger := globalContainer.GetLogger()

	// 3. Handler初期化
	statsHandler := &StatisticsHandler{
		useCases: useCases,
		logger:   logger,
	}

	// 4. 認証済みユーザーID取得（API Gateway Authorizer経由）
	userID, err := statsHandler.getUserIDFromContext(request)
	if err != nil {
		return statsHandler.handleError(err), nil
	}

	// 5. ルーティング
	return statsHandler.routeOperation(ctx, request, userID)
}

// getUserIDFromContext はAPI Gateway Authorizerから認証済みユーザーIDを取得
func (h *StatisticsHandler) getUserIDFromContext(request events.APIGatewayProxyRequest) (uuid.UUID, error) {
	return auth.GetUserIDFromAPIGatewayContext(request, h.logger)
}

// routeOperation は操作ルーティング
func (h *StatisticsHandler) routeOperation(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	// GET以外は許可しない
	if request.HTTPMethod != "GET" {
		return createErrorResponse(http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "GETメソッドのみ許可されています"), nil
	}

	// focus-trendのパスパラメータ処理
	if strings.Contains(request.Path, "/focus-trend/") {
		return h.handleGetFocusTrend(ctx, request, userID)
	} else if strings.Contains(request.Path, "/focus-heatmap") {
		return h.handleGetFocusHeatmap(ctx, request, userID)
	}

	return createErrorResponse(http.StatusNotFound, "NOT_FOUND", "無効なパス"), nil
}

// handleGetFocusTrend は集中度トレンド取得を処理（パスパラメータ形式）
func (h *StatisticsHandler) handleGetFocusTrend(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	// パスから日付を抽出 (/focus-trend/{date} 形式)
	pathParts := strings.Split(request.Path, "/")
	if len(pathParts) < 2 {
		return createErrorResponse(http.StatusBadRequest, "MISSING_DATE", "日付が指定されていません"), nil
	}

	// 最後の部分が日付（開始日）
	startDate := pathParts[len(pathParts)-1]
	
	// 日付フォーマットの基本検証（8桁の数字かチェック）
	if len(startDate) != 8 {
		return createErrorResponse(http.StatusBadRequest, "INVALID_DATE_FORMAT", "日付はYYYY-MMDD形式で指定してください"), nil
	}

	// 数字のみかチェック
	for _, char := range startDate {
		if char < '0' || char > '9' {
			return createErrorResponse(http.StatusBadRequest, "INVALID_DATE_FORMAT", "日付はYYYY-MMDD形式（数字のみ）で指定してください"), nil
		}
	}

	h.logger.Infof("集中度トレンド取得リクエスト: userID=%s, startDate=%s", userID.String()[:8]+"...", startDate)

	// ユースケースを呼び出して集中度トレンドを取得
	response, err := h.useCases.Statistics.GetFocusTrend(ctx, userID, startDate)
	if err != nil {
		h.logger.Errorf("集中度トレンド取得エラー: %v", err)
		return h.handleError(err), nil
	}

	return createSuccessResponse(http.StatusOK, response), nil
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

	return createSuccessResponse(http.StatusOK, response), nil
}

// handleError はエラーを統一処理（error_mapper.go使用版）
func (h *StatisticsHandler) handleError(err error) events.APIGatewayProxyResponse {
	// error_mapper.goを使用してHTTPエラーにマッピング
	httpErr := httpError.MapErrorToHTTP(err)

	// ログ出力（サーバーエラーのみ詳細ログ）
	if httpError.IsServerError(err) {
		h.logger.Errorf("サーバーエラー: %v", err)
	} else if httpError.IsClientError(err) {
		h.logger.Warnf("クライアントエラー: %s - %s", httpErr.Code, httpErr.Message)
	}

	// 統一されたエラーレスポンス作成
	return createErrorResponse(httpErr.StatusCode, httpErr.Code, httpErr.Message)
}

// createSuccessResponse は成功レスポンスを作成
func createSuccessResponse(statusCode int, data interface{}) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(data)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    getCORSHeaders(),
		Body:       string(body),
	}
}

// createErrorResponse はエラーレスポンスを作成（統一フォーマット）
func createErrorResponse(statusCode int, code, message string) events.APIGatewayProxyResponse {
	// エラーレスポンスの統一フォーマット
	errorBody := map[string]interface{}{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	}

	body, _ := json.Marshal(errorBody)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    getCORSHeaders(),
		Body:       string(body),
	}
}

func getCORSHeaders() map[string]string {
	return map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type,Authorization",
		"Access-Control-Allow-Methods": "GET,POST,PATCH,DELETE,OPTIONS",
	}
}

func main() {
	lambda.Start(handler)
}
