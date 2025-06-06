// internal/domain/model/user.go
package model

import (
	"time"

	"github.com/google/uuid"
)

// User はユーザーを表すドメインモデル（Pure Domain - Infrastructure依存なし）
type User struct {
	ID         uuid.UUID `db:"id" json:"id"`                   // Cognito sub
	Name       string    `db:"name" json:"name"`               // 表示名
	Email      string    `db:"email" json:"email"`             // メールアドレス
	Provider   string    `db:"provider" json:"provider"`       // "Cognito_UserPool", "Google"
	ProviderID *string   `db:"provider_id" json:"provider_id"` // Google sub等
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// UserCreationParams はユーザー作成時のパラメータ（外部システム非依存）
type UserCreationParams struct {
	UserID       uuid.UUID
	Name         string
	Email        string
	Provider     string
	ProviderID   *string
	IsGoogleUser bool
}

// NewUser は新しいユーザーを作成する（Pure Domain Logic）
func NewUser(params UserCreationParams) *User {
	now := time.Now()
	return &User{
		ID:         params.UserID,
		Name:       params.Name,
		Email:      params.Email,
		Provider:   params.Provider,
		ProviderID: params.ProviderID,
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

// UpdateProfile はユーザープロフィールを更新する（Pure Domain Logic）
func (u *User) UpdateProfile(name, email string) {
	if name != "" {
		u.Name = name
	}
	if email != "" {
		u.Email = email
	}
	u.UpdatedAt = time.Now()
}

// IsGoogleUser はGoogleユーザーかどうかを返す（Pure Domain Logic）
func (u *User) IsGoogleUser() bool {
	return u.Provider == "Google"
}

// IsCognitoUser はCognitoユーザーかどうかを返す（Pure Domain Logic）
func (u *User) IsCognitoUser() bool {
	return u.Provider == "Cognito_UserPool"
}

// ValidateEmail はメールアドレスの基本的な検証を行う（Pure Domain Logic）
func (u *User) ValidateEmail() bool {
	return u.Email != "" && len(u.Email) > 0
}

// ValidateName は名前の基本的な検証を行う（Pure Domain Logic）
func (u *User) ValidateName() bool {
	return u.Name != "" && len(u.Name) > 0
}
