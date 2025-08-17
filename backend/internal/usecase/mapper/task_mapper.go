package mapper

import (
	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	"github.com/tsunakit99/selfpomodoro/internal/usecase/dto"
)

// TaskMapper はタスク関連のマッピングを担当する
type TaskMapper struct{}

// NewTaskMapper は新しいTaskMapperを作成する
func NewTaskMapper() *TaskMapper {
	return &TaskMapper{}
}

// ToTaskResponse はドメインモデルからAPIレスポンス形式に変換する
func (m *TaskMapper) ToTaskResponse(task *entity.Task) *dto.TaskResponse {
	return &dto.TaskResponse{
		ID:          task.ID.Value(),
		Detail:      task.Detail.Value(),
		IsCompleted: task.Status.IsCompleted(),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

// ToTasksResponse はタスクのリストをレスポンス形式に変換する
func (m *TaskMapper) ToTasksResponse(tasks []*entity.Task) *dto.TasksResponse {
	responses := make([]*dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = m.ToTaskResponse(task)
	}
	
	return &dto.TasksResponse{
		Tasks: responses,
	}
}

// ToTaskResponseList はタスクのリストをレスポンス形式に変換する（複数パターン対応）
func (m *TaskMapper) ToTaskResponseList(tasks []*entity.Task) []*dto.TaskResponse {
	responses := make([]*dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = m.ToTaskResponse(task)
	}
	return responses
}