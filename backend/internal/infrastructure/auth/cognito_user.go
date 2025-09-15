package auth

import (
	"github.com/google/uuid"
)

// CognitoUser は認証されたユーザー情報を表すインフラモデル（API Gateway Authorizer対応版）
type CognitoUser struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
	Provider string    `json:"provider"` // "Cognito_UserPool", "Google"
}

// AuthClaims は認証クレーム情報（API Gateway Authorizer使用時は不要）
type AuthClaims map[string]interface{}

// CognitoUserCreationParams はCognito PostConfirmation時のユーザー作成パラメータ
type CognitoUserCreationParams struct {
	UserID     uuid.UUID
	Email      string
	Name       string
	GivenName  string
	FamilyName string
	// Google SSO関連（Cognito + Google統合時）
	Picture          string
	Locale           string
	IdentityProvider string
	// その他の属性
	EmailVerified bool
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
	User  *CognitoUser       `json:"user"`
	Token *AuthTokenResponse `json:"token,omitempty"`
}