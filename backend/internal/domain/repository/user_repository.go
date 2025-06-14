package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
)

// UserRepository はユーザー永続化のためのインターフェース
type UserRepository interface {
	// GetByID はIDによってユーザーを取得する
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)

	// GetByEmail はメールアドレスによってユーザーを取得する
	GetByEmail(ctx context.Context, email string) (*model.User, error)

	// Create は新しいユーザーを作成する
	Create(ctx context.Context, user *model.User) error

	// Update はユーザー情報を更新する
	Update(ctx context.Context, user *model.User) error

	// Delete はユーザーを削除する（GDPR対応等）
	Delete(ctx context.Context, id uuid.UUID) error

	// UpdateProfile はユーザープロフィール（名前・メール）を更新する
	UpdateProfile(ctx context.Context, id uuid.UUID, name, email string) (*model.User, error)

	// ExistsByID はユーザーの存在確認を行う（軽量版）
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)

	// GetUsersByProvider はプロバイダー別にユーザーを取得する（管理用）
	GetUsersByProvider(ctx context.Context, provider string, limit, offset int) ([]*model.User, error)
}
