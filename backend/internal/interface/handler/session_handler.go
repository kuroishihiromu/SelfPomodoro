package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/postgres"
	"github.com/tsunakit99/selfpomodoro/internal/usecase"
)

// SessionHandler はセッションに関するHTTPリクエストを処理するハンドラ
type SessionHandler struct {
	sessionUseCase usecase.SessionUseCase
	logger         logger.Logger
}

// NewSessionHandler は新しいSessionHandlerインスタンスを作成する
func NewSessionHandler(sessionUseCase usecase.SessionUseCase, logger logger.Logger) *SessionHandler {
	return &SessionHandler{
		sessionUseCase: sessionUseCase,
		logger:         logger,
	}
}

// getIUserIDFromContext はコンテキストからユーザIDを取得する
func (h *SessionHandler) getIUserIDFromContext(c echo.Context) (uuid.UUID, error) {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return uuid.Nil, errors.New("ユーザIDが取得できません")
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, errors.New("ユーザIDが不正です")
	}
	return id, nil
}

// StartSession は新しいセッション開始エンドポイントを処理する
// POST /sessions
func (h *SessionHandler) StartSession(c echo.Context) error {
	// コンテキストからユーザIDを取得
	userID, err := h.getIUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// ユースケースを呼び出してセッションを開始
	sessionResponse, err := h.sessionUseCase.StartSession(c.Request().Context(), userID)
	if err != nil {
		h.logger.Errorf("セッション開始エラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "セッションの開始に失敗しました"})
	}
	return c.JSON(http.StatusCreated, sessionResponse)
}

// GetAllSessions は全セッション取得エンドポイントを処理する
// GET /sessions
func (h *SessionHandler) GetAllSessions(c echo.Context) error {
	// コンテキストからユーザIDを取得
	userID, err := h.getIUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// ユースケースを呼び出して全セッションを取得
	sessionsResponse, err := h.sessionUseCase.GetAllSessions(c.Request().Context(), userID)
	if err != nil {
		h.logger.Errorf("セッション取得エラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "セッションの取得に失敗しました"})
	}
	return c.JSON(http.StatusOK, sessionsResponse)
}

// GetSession はセッション取得エンドポイントを処理する
// GET /sessions/:session_id
func (h *SessionHandler) GetSession(c echo.Context) error {
	// コンテキストからユーザIDを取得
	userID, err := h.getIUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// URLパラメータからセッションIDを取得
	sessionID, err := uuid.Parse(c.Param("session_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "セッションIDが不正です"})
	}

	// ユースケースを呼び出してセッションを取得
	sessionResponse, err := h.sessionUseCase.GetSession(c.Request().Context(), sessionID, userID)
	if err != nil {
		if errors.Is(err, postgres.ErrSessionNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "セッションが見つかりません"})
		}
		if errors.Is(err, postgres.ErrSessionAccessDenied) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "セッションへのアクセスが拒否されました"})
		}
		h.logger.Errorf("セッション取得エラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "セッションの取得に失敗しました"})
	}
	return c.JSON(http.StatusOK, sessionResponse)
}

// CompleteSession はセッション完了エンドポイントを処理する
// POST /sessions/:session_id/complete
func (h *SessionHandler) CompleteSession(c echo.Context) error {
	// コンテキストからユーザIDを取得
	userID, err := h.getIUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// URLパラメータからセッションIDを取得
	sessionID, err := uuid.Parse(c.Param("session_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "セッションIDが不正です"})
	}

	// ユースケースを呼び出してセッションを完了
	sessionResponse, err := h.sessionUseCase.CompleteSession(c.Request().Context(), sessionID, userID)
	if err != nil {
		if errors.Is(err, postgres.ErrSessionNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "セッションが見つかりません"})
		}
		if errors.Is(err, postgres.ErrSessionAccessDenied) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "セッションへのアクセスが拒否されました"})
		}
		h.logger.Errorf("セッション完了エラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "セッションの完了に失敗しました"})
	}
	return c.JSON(http.StatusOK, sessionResponse)
}

// DeleteSession はセッション削除エンドポイントを処理する
// DELETE /sessions/:session_id
func (h *SessionHandler) DeleteSession(c echo.Context) error {
	// コンテキストからユーザIDを取得
	userID, err := h.getIUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// URLパラメータからセッションIDを取得
	sessionID, err := uuid.Parse(c.Param("session_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "セッションIDが不正です"})
	}

	// ユースケースを呼び出してセッションを削除
	err = h.sessionUseCase.DeleteSession(c.Request().Context(), sessionID, userID)
	if err != nil {
		if errors.Is(err, postgres.ErrSessionNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "セッションが見つかりません"})
		}
		if errors.Is(err, postgres.ErrSessionAccessDenied) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "セッションへのアクセスが拒否されました"})
		}
		h.logger.Errorf("セッション削除エラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "セッションの削除に失敗しました"})
	}
	return c.NoContent(http.StatusNoContent)
}
