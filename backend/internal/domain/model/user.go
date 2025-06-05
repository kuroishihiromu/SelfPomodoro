package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/auth"
)

// User はユーザーを表すドメインモデル
type User struct {
	ID         uuid.UUID `db:"id" json:"id"`                   // Cognito sub
	Name       string    `db:"name" json:"name"`               // 表示名
	Email      string    `db:"email" json:"email"`             // メールアドレス
	Provider   string    `db:"provider" json:"provider"`       // "Cognito_UserPool", "Google"
	ProviderID *string   `db:"provider_id" json:"provider_id"` // Google sub等
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// NewUserFromCognito はCognito Claimsから新しいユーザーを作成する
func NewUserFromCognito(userID uuid.UUID, claims *auth.CognitoClaims) *User {
	now := time.Now()

	// プロバイダーの判定
	provider := "Cognito_UserPool" // デフォルト
	var providerID *string

	// Google SSOの場合の判定（identities claimから判定）
	if claims.IdentityProvider != "" {
		provider = "Google"
		if claims.Subject != "" {
			providerID = &claims.Subject
		}
	}

	// 名前の取得（優先順位: name > given_name + family_name > cognito:username > email）
	name := claims.Name
	if name == "" && claims.GivenName != "" {
		name = claims.GivenName
		if claims.FamilyName != "" {
			name += " " + claims.FamilyName
		}
	}
	if name == "" {
		name = claims.CognitoUsername
	}
	if name == "" {
		name = claims.Email
	}

	return &User{
		ID:         userID,
		Name:       name,
		Email:      claims.Email,
		Provider:   provider,
		ProviderID: providerID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// CreateUserRequest はユーザー作成リクエストを表す
type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// UpdateUserRequest はユーザー更新リクエストを表す
type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty" validate:"omitempty"`
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
}

// UserResponse はユーザーのAPIレスポンスを表す
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse はドメインモデルからAPIレスポンス形式に変換する
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Provider:  u.Provider,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// UpdateProfile はユーザープロフィールを更新する
func (u *User) UpdateProfile(name, email string) {
	if name != "" {
		u.Name = name
	}
	if email != "" {
		u.Email = email
	}
	u.UpdatedAt = time.Now()
}

// IsGoogleUser はGoogleユーザーかどうかを返す
func (u *User) IsGoogleUser() bool {
	return u.Provider == "Google"
}

// IsCognitoUser はCognitoユーザーかどうかを返す
func (u *User) IsCognitoUser() bool {
	return u.Provider == "Cognito_UserPool"
}
