package handler

import (
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/usecase"
)

// Handlers はすべてのハンドラーをまとめた構造体
type Handlers struct {
	Auth *AuthHandler
	Task *TaskHandler
	// TODO: 他のハンドラーを追加する場合はここにフィールドを追加
}

// NewHandlers はすべてのハンドラーを初期化する
func NewHandlers(useCases *usecase.UseCases, logger logger.Logger) *Handlers {
	return &Handlers{
		Auth: NewAuthHandler(useCases.Auth, logger),
		Task: NewTaskHandler(useCases.Task, logger),
		// TODO: 他のハンドラーを初期化する場合はここに追加
	}
}
