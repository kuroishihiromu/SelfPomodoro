package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/usecase"
)

// AuthHandler は認証に関するハンドラーを定義する構造体
type AuthHandler struct {
	authUseCase usecase.AuthUseCase
	logger      logger.Logger
}

// NewAuthHandler はAuthHandlerの新しいインスタンスを作成する
func NewAuthHandler(authUseCase usecase.AuthUseCase, logger logger.Logger) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		logger:      logger,
	}
}

// DevLoginでは、開発環境用のログインを行う
// TODO: 本番環境では、Cognito認証を行う
// POST /auth/dev-login
func (h *AuthHandler) DevLogin(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"token":   "dev-token",
		"user_id": "00000000-0000-0000-0000-000000000001",
	})
}
