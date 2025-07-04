package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/container"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	httpError "github.com/tsunakit99/selfpomodoro/internal/handler"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/usecase"
)

// Global container (Lambda再利用最適化)
var globalContainer container.Container

// init はLambda init phaseで実行
func init() {
	globalContainer = container.NewLambdaContainer()
}

// TaskHandler はDI Container使用版のタスクハンドラー（新エラーハンドリング対応版）
type TaskHandler struct {
	useCases  *usecase.UseCases
	logger    logger.Logger
	validator *validator.Validate
}

// handler はLambdaのエントリーポイント（DI Container版・新エラーハンドリング対応版）
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 1. Container初期化（遅延初期化でエラー安全）
	if err := globalContainer.Initialize(ctx); err != nil {
		// 初期化エラーは詳細ログ + 汎用エラーレスポンス
		return createErrorResponse(http.StatusInternalServerError, "INTERNAL_ERROR", "サービス初期化エラー"), nil
	}

	// 2. Dependencies取得（Infrastructure依存なし！）
	useCases := globalContainer.GetUseCases()
	logger := globalContainer.GetLogger()

	// 3. Handler初期化（軽量）
	taskHandler := &TaskHandler{
		useCases:  useCases,
		logger:    logger,
		validator: validator.New(),
	}

	// 4. 認証・User存在確認（統一処理）
	userID, err := taskHandler.authenticateAndValidateUser(ctx, request)
	if err != nil {
		return taskHandler.handleError(err), nil
	}

	// 5. 操作ルーティング（現在のロジック維持）
	return taskHandler.routeOperation(ctx, request, userID)
}

// authenticateAndValidateUser は認証・User存在確認の統一処理
func (h *TaskHandler) authenticateAndValidateUser(ctx context.Context, request events.APIGatewayProxyRequest) (uuid.UUID, error) {
	// UseCaseに完全委譲（Infrastructure詳細なし）
	return h.useCases.Auth.AuthenticateAndValidateUser(ctx, request)
}

// routeOperation は操作ルーティング（現在のロジック維持）
func (h *TaskHandler) routeOperation(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case "GET":
		return h.handleGetTasks(ctx, userID)
	case "POST":
		return h.handleCreateTask(ctx, request, userID)
	case "PATCH":
		return h.handleUpdateOrToggleTask(ctx, request, userID)
	case "DELETE":
		return h.handleDeleteTask(ctx, request, userID)
	default:
		return createErrorResponse(http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "メソッドが許可されていません"), nil
	}
}

// handleGetTasks はタスク一覧取得（UseCaseに委譲）
func (h *TaskHandler) handleGetTasks(ctx context.Context, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	tasksResponse, err := h.useCases.Task.GetAllTasks(ctx, userID)
	if err != nil {
		h.logger.Errorf("タスク一覧取得エラー: %v", err)
		return h.handleError(err), nil
	}
	return createSuccessResponse(http.StatusOK, tasksResponse), nil
}

// handleCreateTask はタスク作成（UseCaseに委譲）
func (h *TaskHandler) handleCreateTask(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	var req model.CreateTaskRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return createErrorResponse(http.StatusBadRequest, "INVALID_REQUEST_FORMAT", "無効なリクエスト形式"), nil
	}

	if err := h.validator.Struct(req); err != nil {
		return createErrorResponse(http.StatusBadRequest, "VALIDATION_ERROR", "タスク詳細は必須です"), nil
	}

	taskResponse, err := h.useCases.Task.CreateTask(ctx, userID, &req)
	if err != nil {
		h.logger.Errorf("タスク作成エラー: %v", err)
		return h.handleError(err), nil
	}
	return createSuccessResponse(http.StatusCreated, taskResponse), nil
}

// handleUpdateOrToggleTask はタスク更新・切り替え（UseCaseに委譲）
func (h *TaskHandler) handleUpdateOrToggleTask(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	taskIDStr := request.PathParameters["task_id"]
	if taskIDStr == "" {
		return createErrorResponse(http.StatusBadRequest, "MISSING_TASK_ID", "タスクIDが指定されていません"), nil
	}

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return createErrorResponse(http.StatusBadRequest, "INVALID_TASK_ID", "無効なタスクID"), nil
	}

	if strings.Contains(request.Path, "/toggle") {
		// 完了状態切り替え
		taskResponse, err := h.useCases.Task.ToggleTaskCompletion(ctx, taskID, userID)
		if err != nil {
			h.logger.Errorf("タスク完了状態切り替えエラー: %v", err)
			return h.handleError(err), nil
		}
		return createSuccessResponse(http.StatusOK, taskResponse), nil
	} else {
		// タスク更新
		var req model.UpdateTaskRequest
		if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
			return createErrorResponse(http.StatusBadRequest, "INVALID_REQUEST_FORMAT", "無効なリクエスト形式"), nil
		}

		if err := h.validator.Struct(req); err != nil {
			return createErrorResponse(http.StatusBadRequest, "VALIDATION_ERROR", "タスク詳細は必須です"), nil
		}

		taskResponse, err := h.useCases.Task.UpdateTask(ctx, taskID, userID, &req)
		if err != nil {
			h.logger.Errorf("タスク更新エラー: %v", err)
			return h.handleError(err), nil
		}
		return createSuccessResponse(http.StatusOK, taskResponse), nil
	}
}

// handleDeleteTask はタスク削除（UseCaseに委譲）
func (h *TaskHandler) handleDeleteTask(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	taskIDStr := request.PathParameters["task_id"]
	if taskIDStr == "" {
		return createErrorResponse(http.StatusBadRequest, "MISSING_TASK_ID", "タスクIDが指定されていません"), nil
	}

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return createErrorResponse(http.StatusBadRequest, "INVALID_TASK_ID", "無効なタスクID"), nil
	}

	err = h.useCases.Task.DeleteTask(ctx, taskID, userID)
	if err != nil {
		h.logger.Errorf("タスク削除エラー: %v", err)
		return h.handleError(err), nil
	}
	return createSuccessResponse(http.StatusOK, map[string]string{"message": "タスクが削除されました"}), nil
}

// handleError はエラーを統一処理（error_mapper.go使用版）
func (h *TaskHandler) handleError(err error) events.APIGatewayProxyResponse {
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
		"Access-Control-Allow-Methods": "GET,POST,PATCH,DELETE,OPTIONS",
	}
}

func main() {
	lambda.Start(handler)
}
