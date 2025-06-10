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
	domainErrors "github.com/tsunakit99/selfpomodoro/internal/domain/errors"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/usecase"
)

// Global container (Lambda再利用最適化)
var globalContainer container.Container

// init はLambda init phaseで実行
func init() {
	globalContainer = container.NewLambdaContainer()
}

// TaskHandler はDI Container使用版のタスクハンドラー
type TaskHandler struct {
	useCases  *usecase.UseCases
	logger    logger.Logger
	validator *validator.Validate
}

// handler はLambdaのエントリーポイント（DI Container版）
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 1. Container初期化（遅延初期化でエラー安全）
	if err := globalContainer.Initialize(ctx); err != nil {
		// 初期化エラーは詳細ログ + 汎用エラーレスポンス
		return errorResponse(http.StatusInternalServerError, "サービス初期化エラー"), nil
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
		return errorResponse(http.StatusMethodNotAllowed, "メソッドが許可されていません"), nil
	}
}

// handleGetTasks はタスク一覧取得（UseCaseに委譲）
func (h *TaskHandler) handleGetTasks(ctx context.Context, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	tasksResponse, err := h.useCases.Task.GetAllTasks(ctx, userID)
	if err != nil {
		h.logger.Errorf("タスク一覧取得エラー: %v", err)
		return h.handleError(err), nil
	}
	return successResponse(http.StatusOK, tasksResponse), nil
}

// handleCreateTask はタスク作成（UseCaseに委譲）
func (h *TaskHandler) handleCreateTask(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	var req model.CreateTaskRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return errorResponse(http.StatusBadRequest, "無効なリクエスト形式"), nil
	}

	if err := h.validator.Struct(req); err != nil {
		return errorResponse(http.StatusBadRequest, "タスク詳細は必須です"), nil
	}

	taskResponse, err := h.useCases.Task.CreateTask(ctx, userID, &req)
	if err != nil {
		h.logger.Errorf("タスク作成エラー: %v", err)
		return h.handleError(err), nil
	}
	return successResponse(http.StatusCreated, taskResponse), nil
}

// handleUpdateOrToggleTask はタスク更新・切り替え（UseCaseに委譲）
func (h *TaskHandler) handleUpdateOrToggleTask(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	taskIDStr := request.PathParameters["task_id"]
	if taskIDStr == "" {
		return errorResponse(http.StatusBadRequest, "タスクIDが指定されていません"), nil
	}

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "無効なタスクID"), nil
	}

	if strings.Contains(request.Path, "/toggle") {
		// 完了状態切り替え
		taskResponse, err := h.useCases.Task.ToggleTaskCompletion(ctx, taskID, userID)
		if err != nil {
			h.logger.Errorf("タスク完了状態切り替えエラー: %v", err)
			return h.handleError(err), nil
		}
		return successResponse(http.StatusOK, taskResponse), nil
	} else {
		// タスク更新
		var req model.UpdateTaskRequest
		if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
			return errorResponse(http.StatusBadRequest, "無効なリクエスト形式"), nil
		}

		if err := h.validator.Struct(req); err != nil {
			return errorResponse(http.StatusBadRequest, "タスク詳細は必須です"), nil
		}

		taskResponse, err := h.useCases.Task.UpdateTask(ctx, taskID, userID, &req)
		if err != nil {
			h.logger.Errorf("タスク更新エラー: %v", err)
			return h.handleError(err), nil
		}
		return successResponse(http.StatusOK, taskResponse), nil
	}
}

// handleDeleteTask はタスク削除（UseCaseに委譲）
func (h *TaskHandler) handleDeleteTask(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	taskIDStr := request.PathParameters["task_id"]
	if taskIDStr == "" {
		return errorResponse(http.StatusBadRequest, "タスクIDが指定されていません"), nil
	}

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "無効なタスクID"), nil
	}

	err = h.useCases.Task.DeleteTask(ctx, taskID, userID)
	if err != nil {
		h.logger.Errorf("タスク削除エラー: %v", err)
		return h.handleError(err), nil
	}
	return successResponse(http.StatusOK, map[string]string{"message": "タスクが削除されました"}), nil
}

// handleError はドメインエラーを統一処理
func (h *TaskHandler) handleError(err error) events.APIGatewayProxyResponse {
	if appErr, ok := err.(*domainErrors.AppError); ok {
		return errorResponse(appErr.Status, appErr.Error())
	}

	h.logger.Errorf("予期しないエラータイプ: %T, %v", err, err)
	return errorResponse(http.StatusInternalServerError, "内部エラーが発生しました")
}

// Response helper functions (変更なし)
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
