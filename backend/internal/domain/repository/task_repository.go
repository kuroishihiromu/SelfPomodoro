package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
)

// TaskRepository はタスク永続化のためのインターフェース
type TaskRepository interface {
	// Create はタスクを作成する
	Create(ctx context.Context, task *model.Task) error

	// GetByID はIDによってタスクを取得する
	GetByID(ctx context.Context, id, userID uuid.UUID) (*model.Task, error)

	// GetAllByUserID はユーザーIDに紐づくすべてのタスクを取得する
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Task, error)

	// Update はタスクの詳細を更新する
	Update(ctx context.Context, task *model.Task) error

	// ToggleCompletion はタスクの完了状態を切り替える
	ToggleCompletion(ctx context.Context, id, userID uuid.UUID) error

	// Delete はタスクを削除する
	Delete(ctx context.Context, id, userID uuid.UUID) error
}
