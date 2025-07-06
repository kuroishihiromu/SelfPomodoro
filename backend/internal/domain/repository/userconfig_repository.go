package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
)

// UserConfigRepository はユーザー設定永続化のためのインターフェース
type UserConfigRepository interface {
	// GetUserConfig はユーザーIDからユーザー設定を取得する
	GetUserConfig(ctx context.Context, userID uuid.UUID) (*entity.UserConfig, error)

	// CreateUserConfig は新しいユーザー設定を作成する
	CreateUserConfig(ctx context.Context, config *entity.UserConfig) error

	// UpdateUserConfig はユーザー設定を更新する
	UpdateUserConfig(ctx context.Context, config *entity.UserConfig) error

	// DeleteUserConfig はユーザー設定を削除する（GDPR対応等）
	DeleteUserConfig(ctx context.Context, userID uuid.UUID) error
}
