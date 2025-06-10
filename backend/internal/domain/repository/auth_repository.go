// internal/domain/repository/auth_repository.go
package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
)

// AuthRepository は認証に関するリポジトリインターフェース
type AuthRepository interface {
	// ValidateToken はトークンを検証してユーザーIDを返す
	ValidateToken(ctx context.Context, token string) (uuid.UUID, error)

	// ValidateTokenAndGetClaims はトークンを検証してクレーム情報を返す
	ValidateTokenAndGetClaims(ctx context.Context, token string) (*model.AuthClaims, error)

	// // RefreshToken はリフレッシュトークンを使って新しいアクセストークンを取得する
	// RefreshToken(ctx context.Context, refreshToken string) (*model.AuthTokenResponse, error)

	// // RevokeToken はトークンを無効化する
	// RevokeToken(ctx context.Context, token string) error

	// HealthCheck は認証サービスの接続確認を行う
	HealthCheck(ctx context.Context) error
}
