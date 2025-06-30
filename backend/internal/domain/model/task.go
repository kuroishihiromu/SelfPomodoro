package model

import (
	"time"

	"github.com/google/uuid"
)

// Task　はユーザのタスクを表すドメインモデル
type Task struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Detail      string    `json:"detail"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewTask は新しいタスクを作成する
func NewTask(userID uuid.UUID, detail string) *Task {
	now := time.Now()
	return &Task{
		ID:          uuid.New(),
		UserID:      userID,
		Detail:      detail,
		IsCompleted: false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// ToggleCompletion はタスクの完了状態を切り替える
func (t *Task) ToggleCompletion() {
	t.IsCompleted = !t.IsCompleted
	t.UpdatedAt = time.Now()
}

// UpdateDetail はタスクの詳細を更新する
func (t *Task) UpdateDetail(detail string) {
	t.Detail = detail
	t.UpdatedAt = time.Now()
}

// CreateTaskRequest はタスク作成リクエストを表す
type CreateTaskRequest struct {
	Detail string `json:"detail" validate:"required"`
}

// UpdateTaskRequest はタスク更新リクエストを表す
type UpdateTaskRequest struct {
	Detail string `json:"detail" validate:"required"`
}

// TaskResponse はタスクAPIレスポンスを表す
type TaskResponse struct {
	ID          uuid.UUID `json:"id"`
	Detail      string    `json:"detail"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse はドメインモデルからAPIレスポンス形式に変換する
func (t *Task) ToResponse() *TaskResponse {
	return &TaskResponse{
		ID:          t.ID,
		Detail:      t.Detail,
		IsCompleted: t.IsCompleted,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// TasksResponse はタスクのリストAPIレスポンスを表す
type TasksResponse struct {
	Tasks []*TaskResponse `json:"tasks"`
}
