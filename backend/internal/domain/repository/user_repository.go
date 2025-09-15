package repository

import (
	"context"

	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
)

// UserRepository はUser集約の永続化を担当するリポジトリ
// DDD原則: User集約ルートに対するリポジトリ
type UserRepository interface {
	// GetByID はIDによってユーザーを取得する
	GetByID(ctx context.Context, id userVO.UserID) (*entity.User, error)

	// GetByEmail はメールアドレスによってユーザーを取得する
	GetByEmail(ctx context.Context, email userVO.EmailAddress) (*entity.User, error)

	// Create は新しいユーザーを作成する
	Create(ctx context.Context, user *entity.User) error

	// Update はユーザー情報を更新する
	Update(ctx context.Context, user *entity.User) error

	// Delete はユーザーを削除する（GDPR対応等）
	Delete(ctx context.Context, id userVO.UserID) error

	// UpdateProfile はユーザープロフィール（名前・メール）を更新する
	UpdateProfile(ctx context.Context, id userVO.UserID, name, email string) (*entity.User, error)

	// ExistsByID はユーザーの存在確認を行う（軽量版）
	ExistsByID(ctx context.Context, id userVO.UserID) (bool, error)

	// GetUsersByProvider はプロバイダー別にユーザーを取得する（管理用）
	GetUsersByProvider(ctx context.Context, provider string, limit, offset int) ([]*entity.User, error)

	// DeleteAllUserData はユーザーに関連するすべてのデータを削除する（統合テーブル対応）
	DeleteAllUserData(ctx context.Context, userID userVO.UserID) error

	// CountUserItems はユーザーに関連するアイテム数を取得する（テスト用）
	CountUserItems(ctx context.Context, userID userVO.UserID) (int, error)
}
