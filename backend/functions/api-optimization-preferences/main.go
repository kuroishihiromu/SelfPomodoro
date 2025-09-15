package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/container"
	"github.com/tsunakit99/selfpomodoro/internal/usecase/dto"
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

// OptimizationPreferencesHandler はDI Container使用版のOptimizationPreferencesハンドラー
type OptimizationPreferencesHandler struct {
	useCases  *usecase.UseCases
	logger    logger.Logger
	validator *validator.Validate
}

// handler はLambdaのエントリーポイント（DI Container版）
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 1. Container初期化
	if err := globalContainer.Initialize(ctx); err != nil {
		return createErrorResponse(http.StatusInternalServerError, "INTERNAL_ERROR", "サービス初期化エラー"), nil
	}

	// 2. Dependencies取得
	useCases := globalContainer.GetUseCases()
	logger := globalContainer.GetLogger()

	// 3. Handler初期化
	optimizationPreferencesHandler := &OptimizationPreferencesHandler{
		useCases:  useCases,
		logger:    logger,
		validator: validator.New(),
	}

	// 4. 認証済みユーザーID取得（API Gateway Authorizer経由）
	userID, err := optimizationPreferencesHandler.getUserIDFromContext(request)
	if err != nil {
		return optimizationPreferencesHandler.handleError(err), nil
	}

	// 5. 操作ルーティング
	return optimizationPreferencesHandler.routeOperation(ctx, request, userID)
}

// getUserIDFromContext はAPI Gateway Authorizerから認証済みユーザーIDを取得
func (h *OptimizationPreferencesHandler) getUserIDFromContext(request events.APIGatewayProxyRequest) (uuid.UUID, error) {
	return auth.GetUserIDFromAPIGatewayContext(request, h.logger)
}

// routeOperation は操作ルーティング（OptimizationPreferences専用）
func (h *OptimizationPreferencesHandler) routeOperation(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case "GET":
		return h.handleGetOptimizationPreferences(ctx, userID)
	case "PUT":
		return h.handleUpdateOptimizationPreferences(ctx, request, userID)
	default:
		return createErrorResponse(http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "GETまたはPUTメソッドのみ許可されています"), nil
	}
}

// handleGetOptimizationPreferences はOptimizationPreferences取得を処理
func (h *OptimizationPreferencesHandler) handleGetOptimizationPreferences(ctx context.Context, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	// ユースケースを呼び出してOptimizationPreferencesを取得
	response, err := h.useCases.OptimizationPreferences.GetOptimizationPreferences(ctx, userID)
	if err != nil {
		h.logger.Errorf("OptimizationPreferences取得エラー: %v", err)
		return h.handleError(err), nil
	}

	return createSuccessResponse(http.StatusOK, response), nil
}

// handleUpdateOptimizationPreferences はOptimizationPreferences更新を処理
func (h *OptimizationPreferencesHandler) handleUpdateOptimizationPreferences(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	// リクエストボディをパース
	var updateRequest dto.UpdateOptimizationPreferencesRequest
	if err := json.Unmarshal([]byte(request.Body), &updateRequest); err != nil {
		h.logger.Warnf("OptimizationPreferences更新リクエストのパースエラー: %v", err)
		return createErrorResponse(http.StatusBadRequest, "INVALID_REQUEST", "無効なリクエストボディ"), nil
	}

	// バリデーション
	if err := h.validator.Struct(&updateRequest); err != nil {
		h.logger.Warnf("OptimizationPreferences更新バリデーションエラー: %v", err)
		return createErrorResponse(http.StatusBadRequest, "VALIDATION_ERROR", "リクエストデータが無効です"), nil
	}

	// ユースケースを呼び出してOptimizationPreferencesを更新
	response, err := h.useCases.OptimizationPreferences.UpdateOptimizationPreferences(ctx, userID, &updateRequest)
	if err != nil {
		h.logger.Errorf("OptimizationPreferences更新エラー: %v", err)
		return h.handleError(err), nil
	}

	return createSuccessResponse(http.StatusOK, response), nil
}

// handleError はエラーを統一処理（error_mapper.go使用版）
func (h *OptimizationPreferencesHandler) handleError(err error) events.APIGatewayProxyResponse {
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

// getCORSHeaders はCORSヘッダーを取得
func getCORSHeaders() map[string]string {
	return map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type,Authorization",
		"Access-Control-Allow-Methods": "GET,PUT,OPTIONS",
	}
}

func main() {
	lambda.Start(handler)
}
