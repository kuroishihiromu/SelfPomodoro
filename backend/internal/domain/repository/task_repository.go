package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
)

// TaskRepository はTask集約の永続化を担当するリポジトリ
// DDD原則: Task集約ルートに対するリポジトリ
type TaskRepository interface {
	// Create はタスクを作成する
	Create(ctx context.Context, task *entity.Task) error

	// GetByID はIDによってタスクを取得する
	GetByID(ctx context.Context, id uuid.UUID, userID userVO.UserID) (*entity.Task, error)

	// GetAllByUserID はユーザーIDに紐づくすべてのタスクを取得する
	GetAllByUserID(ctx context.Context, userID userVO.UserID) ([]*entity.Task, error)

	// Update はタスクの詳細を更新する
	Update(ctx context.Context, task *entity.Task) error

	// ToggleCompletion はタスクの完了状態を切り替える
	ToggleCompletion(ctx context.Context, id uuid.UUID, userID userVO.UserID) error

	// Delete はタスクを削除する
	Delete(ctx context.Context, id uuid.UUID, userID userVO.UserID) error
}
