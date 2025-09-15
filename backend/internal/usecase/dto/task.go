package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateTaskRequest はタスク作成リクエストを表すDTO
type CreateTaskRequest struct {
	Detail string `json:"detail" validate:"required"`
}

// UpdateTaskRequest はタスク更新リクエストを表すDTO
type UpdateTaskRequest struct {
	Detail string `json:"detail" validate:"required"`
}

// TaskResponse はタスクAPIレスポンスを表すDTO
type TaskResponse struct {
	ID          uuid.UUID `json:"id"`
	Detail      string    `json:"detail"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TasksResponse はタスクのリストAPIレスポンスを表すDTO
type TasksResponse struct {
	Tasks []*TaskResponse `json:"tasks"`
}