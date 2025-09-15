package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateUserRequest はユーザー作成リクエストを表すDTO
type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// UpdateUserRequest はユーザー更新リクエストを表すDTO
type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty" validate:"omitempty"`
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
}

// UserResponse はユーザーのAPIレスポンスを表すDTO
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}