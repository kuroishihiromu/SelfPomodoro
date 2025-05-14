package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tsunakit99/selfpomodoro/internal/interface/handler"
	customMiddleware "github.com/tsunakit99/selfpomodoro/internal/middleware"
)

func (s *Server) SetupRouter(handlers *handler.Handlers) {
	e := s.echo

	// ミドルウェアの設定
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(customMiddleware.LogggerMiddleware(s.logger))

	// ヘルスチェックエンドポイント
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	// APIグループ（認証なし）
	api := e.Group("/api/v1")

	// 認証ルート(開発用ダミー)
	auth := api.Group("/auth")
	auth.POST("/dev-login", handlers.Auth.DevLogin)

	//　認証が必要なルート
	secured := api.Group("")
	secured.Use(customMiddleware.AuthMiddleware())

	// タスク関連のルート
	tasks := secured.Group("/tasks")
	tasks.POST("", handlers.Task.CreateTask)
	tasks.GET("", handlers.Task.GetAllTasks)
	tasks.PATCH("/:task_id/edit", handlers.Task.UpdateTask)
	tasks.PATCH("/:task_id/toggle", handlers.Task.ToggleTask)
	tasks.DELETE("/:task_id", handlers.Task.DeleteTask)
}
