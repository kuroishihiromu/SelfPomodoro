package handler

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/postgres"
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

// GetAllTasks はタスク一覧取得エンドポイントを処理する
// GET /tasks
func (h *TaskHandler) GetAllTasks(c echo.Context) error {
	// ユーザーIDを取得
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// ユースケースを呼び出してタスク一覧を取得
	tasksResponse, err := h.taskUseCase.GetAllTasks(c.Request().Context(), userID)
	if err != nil {
		h.logger.Errorf("タスク一覧取得エラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "タスク一覧の取得に失敗しました"})
	}
	return c.JSON(http.StatusOK, tasksResponse)
}

// UpdateTask はタスク更新エンドポイントを処理する
// PATCH /tasks/:task_id/edit
func (h *TaskHandler) UpdateTask(c echo.Context) error {
	// ユーザーIDを取得
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// タスクIDのパース
	taskID, err := uuid.Parse(c.Param("task_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "無効なタスクID"})
	}

	// リクエストボディのバインド
	var req model.UpdateTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "無効なリクエスト形式"})
	}

	// バリデーション
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "タスク詳細は必須です"})
	}

	// ユースケースを呼び出してタスクを更新
	taskResponse, err := h.taskUseCase.UpdateTask(c.Request().Context(), taskID, userID, &req)
	if err != nil {
		if errors.Is(err, postgres.ErrTaskNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "タスクが見つかりません"})
		}
		if errors.Is(err, postgres.ErrTaskAccessDenied) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "このタスクへのアクセス権限がありません"})
		}

		h.logger.Errorf("タスク更新エラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "タスク更新に失敗しました"})
	}

	return c.JSON(http.StatusOK, taskResponse)
}

// ToggleTask はタスク完了状態切り替えエンドポイントを処理する
// PATCH /tasks/:task_id/toggle
func (h *TaskHandler) ToggleTask(c echo.Context) error {
	// ユーザーIDを取得
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// タスクIDのパース
	taskID, err := uuid.Parse(c.Param("task_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "無効なタスクID"})
	}

	// ユースケースを呼び出してタスクの完了状態を切り替え
	taskResponse, err := h.taskUseCase.ToggleTaskCompletion(c.Request().Context(), taskID, userID)
	if err != nil {
		if errors.Is(err, postgres.ErrTaskNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "タスクが見つかりません"})
		}
		if errors.Is(err, postgres.ErrTaskAccessDenied) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "このタスクへのアクセス権限がありません"})
		}

		h.logger.Errorf("タスク完了状態切り替えエラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "タスク完了状態の切り替えに失敗しました"})
	}

	return c.JSON(http.StatusOK, taskResponse)
}

// DeleteTask はタスク削除エンドポイントを処理する
// DELETE /tasks/:task_id
func (h *TaskHandler) DeleteTask(c echo.Context) error {
	// ユーザーIDを取得
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// タスクIDのパース
	taskID, err := uuid.Parse(c.Param("task_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "無効なタスクID"})
	}

	// ユースケースを呼び出してタスクを削除
	err = h.taskUseCase.DeleteTask(c.Request().Context(), taskID, userID)
	if err != nil {
		if errors.Is(err, postgres.ErrTaskNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "タスクが見つかりません"})
		}
		if errors.Is(err, postgres.ErrTaskAccessDenied) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "このタスクへのアクセス権限がありません"})
		}

		h.logger.Errorf("タスク削除エラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "タスク削除に失敗しました"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "タスクが削除されました"})
}
