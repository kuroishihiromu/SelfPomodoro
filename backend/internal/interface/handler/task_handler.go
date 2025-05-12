package handler

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/usecase"
)

// TaskHandler はタスク関連のエンドポイントを処理するハンドラー
type TaskHandler struct {
	taskUseCase usecase.TaskUseCase
	logger      logger.Logger
	validator   *validator.Validate
}

// NewTaskHandler は新しいTaskHandlerインスタンスを作成する
func NewTaskHandler(taskUseCase usecase.TaskUseCase, logger logger.Logger) *TaskHandler {
	return &TaskHandler{
		taskUseCase: taskUseCase,
		logger:      logger,
		validator:   validator.New(),
	}
}

// getUserIDFromContext はコンテキストからユーザーIDを取得する
func (h *TaskHandler) getUserIDFromContext(c echo.Context) (uuid.UUID, error) {
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

// CreateTask はタスク作成エンドポイントを処理する
// POST /tasks
func (h *TaskHandler) CreateTask(c echo.Context) error {
	// ユーザーIDを取得
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// リクエストボディのバインド
	var req model.CreateTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "無効なリクエスト形式"})
	}

	// バリデーション
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "タスク詳細は必須です"})
	}

	// ユースケースを呼び出してタスクを作成
	taskResponse, err := h.taskUseCase.CreateTask(c.Request().Context(), userID, &req)
	if err != nil {
		h.logger.Errorf("タスク作成エラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "タスク作成に失敗しました"})
	}
	return c.JSON(http.StatusCreated, taskResponse)
}
