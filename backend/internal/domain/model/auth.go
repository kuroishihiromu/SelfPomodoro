// internal/domain/model/auth.go
package model

import (
	"github.com/google/uuid"
)

// AuthUser は認証されたユーザー情報を表すドメインモデル
type AuthUser struct {
	UserID   uuid.UUID   `json:"user_id"`
	Email    string      `json:"email"`
	Name     string      `json:"name"`
	Provider string      `json:"provider"` // "Cognito_UserPool", "Google"
	Claims   *AuthClaims `json:"claims,omitempty"`
}

// AuthClaims は認証トークンのクレーム情報を表すドメインモデル
type AuthClaims struct {
	UserID          uuid.UUID `json:"user_id"`
	Email           string    `json:"email"`
	Name            string    `json:"name"`
	Provider        string    `json:"provider"`
	CognitoUsername string    `json:"cognito_username,omitempty"`
	EmailVerified   bool      `json:"email_verified"`
	TokenUse        string    `json:"token_use"` // "access" or "id"
	ExpiresAt       int64     `json:"expires_at"`
	IssuedAt        int64     `json:"issued_at"`
	Subject         string    `json:"subject"`
	Issuer          string    `json:"issuer"`
	Audience        string    `json:"audience"`
}

// IsValid はクレームの有効性を確認する
func (c *AuthClaims) IsValid() bool {
	return c.UserID != uuid.Nil && c.Email != "" && c.Subject != ""
}

// IsExpired はトークンが期限切れかどうかを確認する
func (c *AuthClaims) IsExpired() bool {
	// TODO: 現在時刻と比較してチェック
	return false
}

// AuthTokenRequest は認証トークンリクエストを表す
type AuthTokenRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// AuthTokenResponse は認証トークンレスポンスを表す
type AuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// AuthResponse は認証成功レスポンスを表す
type AuthResponse struct {
	User  *AuthUser          `json:"user"`
	Token *AuthTokenResponse `json:"token,omitempty"`
}
