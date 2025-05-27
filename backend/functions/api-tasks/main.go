package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/database"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/postgres"
	"github.com/tsunakit99/selfpomodoro/internal/usecase"
)

// TaskHandler はLambda用のタスクハンドラー
type TaskHandler struct {
	taskUseCase usecase.TaskUseCase
	logger      logger.Logger
	validator   *validator.Validate
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
	taskUseCase := usecase.NewTaskUseCase(repositoryFactory.Task, appLogger)

	// ハンドラーの初期化
	taskHandler := &TaskHandler{
		taskUseCase: taskUseCase,
		logger:      appLogger,
		validator:   validator.New(),
	}

	// ユーザーID取得
	userID, err := getUserIDFromRequest(request)
	if err != nil {
		appLogger.Errorf("ユーザーID取得エラー: %v", err)
		return errorResponse(http.StatusUnauthorized, "認証エラー"), nil
	}

	// ルーティング
	switch request.HTTPMethod {
	case "GET":
		return taskHandler.handleGetTasks(ctx, userID)
	case "POST":
		return taskHandler.handleCreateTask(ctx, request, userID)
	case "PATCH":
		return taskHandler.handleUpdateOrToggleTask(ctx, request, userID)
	case "DELETE":
		return taskHandler.handleDeleteTask(ctx, request, userID)
	default:
		return errorResponse(http.StatusMethodNotAllowed, "メソッドが許可されていません"), nil
	}
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

	// 本番環境用Cognito認証（将来実装）
	// if claims, exists := request.RequestContext.Authorizer["claims"]; exists {
	//     claimsMap := claims.(map[string]interface{})
	//     if sub, ok := claimsMap["sub"].(string); ok {
	//         return uuid.Parse(sub)
	//     }
	// }

	return uuid.Nil, fmt.Errorf("認証情報が見つかりません")
}

// handleGetTasks はタスク一覧取得を処理
func (h *TaskHandler) handleGetTasks(ctx context.Context, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	tasksResponse, err := h.taskUseCase.GetAllTasks(ctx, userID)
	if err != nil {
		h.logger.Errorf("タスク一覧取得エラー: %v", err)
		return errorResponse(http.StatusInternalServerError, "タスク一覧の取得に失敗しました"), nil
	}

	return successResponse(http.StatusOK, tasksResponse), nil
}

// handleCreateTask はタスク作成を処理
func (h *TaskHandler) handleCreateTask(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	var req model.CreateTaskRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return errorResponse(http.StatusBadRequest, "無効なリクエスト形式"), nil
	}

	if err := h.validator.Struct(req); err != nil {
		return errorResponse(http.StatusBadRequest, "タスク詳細は必須です"), nil
	}

	taskResponse, err := h.taskUseCase.CreateTask(ctx, userID, &req)
	if err != nil {
		h.logger.Errorf("タスク作成エラー: %v", err)
		return errorResponse(http.StatusInternalServerError, "タスク作成に失敗しました"), nil
	}

	return successResponse(http.StatusCreated, taskResponse), nil
}

// handleUpdateOrToggleTask はタスク更新または完了状態切り替えを処理
func (h *TaskHandler) handleUpdateOrToggleTask(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	// パスパラメータからタスクIDを取得
	taskIDStr := request.PathParameters["task_id"]
	if taskIDStr == "" {
		return errorResponse(http.StatusBadRequest, "タスクIDが指定されていません"), nil
	}

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "無効なタスクID"), nil
	}

	// パスによって処理を分岐
	if strings.Contains(request.Path, "/toggle") {
		// 完了状態切り替え
		taskResponse, err := h.taskUseCase.ToggleTaskCompletion(ctx, taskID, userID)
		if err != nil {
			if isNotFoundError(err) {
				return errorResponse(http.StatusNotFound, "タスクが見つかりません"), nil
			}
			if isAccessDeniedError(err) {
				return errorResponse(http.StatusForbidden, "このタスクへのアクセス権限がありません"), nil
			}
			h.logger.Errorf("タスク完了状態切り替えエラー: %v", err)
			return errorResponse(http.StatusInternalServerError, "タスク完了状態の切り替えに失敗しました"), nil
		}
		return successResponse(http.StatusOK, taskResponse), nil
	} else if strings.Contains(request.Path, "/edit") {
		// タスク更新
		var req model.UpdateTaskRequest
		if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
			return errorResponse(http.StatusBadRequest, "無効なリクエスト形式"), nil
		}

		if err := h.validator.Struct(req); err != nil {
			return errorResponse(http.StatusBadRequest, "タスク詳細は必須です"), nil
		}

		taskResponse, err := h.taskUseCase.UpdateTask(ctx, taskID, userID, &req)
		if err != nil {
			if isNotFoundError(err) {
				return errorResponse(http.StatusNotFound, "タスクが見つかりません"), nil
			}
			if isAccessDeniedError(err) {
				return errorResponse(http.StatusForbidden, "このタスクへのアクセス権限がありません"), nil
			}
			h.logger.Errorf("タスク更新エラー: %v", err)
			return errorResponse(http.StatusInternalServerError, "タスク更新に失敗しました"), nil
		}
		return successResponse(http.StatusOK, taskResponse), nil
	}

	return errorResponse(http.StatusNotFound, "無効なパス"), nil
}

// handleDeleteTask はタスク削除を処理
func (h *TaskHandler) handleDeleteTask(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	taskIDStr := request.PathParameters["task_id"]
	if taskIDStr == "" {
		return errorResponse(http.StatusBadRequest, "タスクIDが指定されていません"), nil
	}

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "無効なタスクID"), nil
	}

	err = h.taskUseCase.DeleteTask(ctx, taskID, userID)
	if err != nil {
		if isNotFoundError(err) {
			return errorResponse(http.StatusNotFound, "タスクが見つかりません"), nil
		}
		if isAccessDeniedError(err) {
			return errorResponse(http.StatusForbidden, "このタスクへのアクセス権限がありません"), nil
		}
		h.logger.Errorf("タスク削除エラー: %v", err)
		return errorResponse(http.StatusInternalServerError, "タスク削除に失敗しました"), nil
	}

	return successResponse(http.StatusOK, map[string]string{"message": "タスクが削除されました"}), nil
}

// ヘルパー関数
func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), postgres.ErrTaskNotFound.Error())
}

func isAccessDeniedError(err error) bool {
	return strings.Contains(err.Error(), postgres.ErrTaskAccessDenied.Error())
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
