package entity

import (
	"time"

	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
	taskVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/task"
)

// Task はユーザのタスクを表すドメインエンティティ
type Task struct {
	ID        taskVO.TaskID
	UserID    userVO.UserID
	Detail    taskVO.TaskDetail
	Status    taskVO.TaskStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewTask は新しいタスクを作成する
func NewTask(userID userVO.UserID, detail string) *Task {
	now := time.Now()
	taskDetail, _ := taskVO.NewTaskDetail(detail)
	return &Task{
		ID:        taskVO.NewTaskID(),
		UserID:    userID,
		Detail:    taskDetail,
		Status:    taskVO.NewIncompleteTaskStatus(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ToggleCompletion はタスクの完了状態を切り替える
func (t *Task) ToggleCompletion() {
	t.Status = t.Status.Toggle()
	t.UpdatedAt = time.Now()
}

// UpdateDetail はタスクの詳細を更新する
func (t *Task) UpdateDetail(detail string) error {
	newDetail, err := taskVO.NewTaskDetail(detail)
	if err != nil {
		return err
	}
	t.Detail = newDetail
	t.UpdatedAt = time.Now()
	return nil
}

// IsCompleted はタスクが完了しているかを判定する
func (t *Task) IsCompleted() bool {
	return t.Status.IsCompleted()
}

// IsIncomplete はタスクが未完了かを判定する
func (t *Task) IsIncomplete() bool {
	return t.Status.IsIncomplete()
}

// Complete はタスクを完了状態にする
func (t *Task) Complete() {
	t.Status = t.Status.Complete()
	t.UpdatedAt = time.Now()
}

// Incomplete はタスクを未完了状態にする
func (t *Task) Incomplete() {
	t.Status = t.Status.Incomplete()
	t.UpdatedAt = time.Now()
}

