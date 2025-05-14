package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// TaskUseCase はタスクに関するユースケースを定義するインターフェース
type TaskUseCase interface {
	// CreateTask は新しいタスクを作成する
	CreateTask(ctx context.Context, userID uuid.UUID, req *model.CreateTaskRequest) (*model.TaskResponse, error)

	// GetTask　は時指定されたIDのタスクを取得する
	GetTask(ctx context.Context, id, userID uuid.UUID) (*model.TaskResponse, error)

	// GetAllTasks はユーザーIDに紐づくすべてのタスクを取得する
	GetAllTasks(ctx context.Context, userID uuid.UUID) (*model.TasksResponse, error)

	// UpdateTask はタスクの詳細を更新する
	UpdateTask(ctx context.Context, id, userID uuid.UUID, req *model.UpdateTaskRequest) (*model.TaskResponse, error)

	// ToggleTaskCompletion はタスクの完了状態を切り替える
	ToggleTaskCompletion(ctx context.Context, id, userID uuid.UUID) (*model.TaskResponse, error)

	// DeleteTask はタスクを削除する
	DeleteTask(ctx context.Context, id, userID uuid.UUID) error
}

// taskUseCase はTaskUseCaseインターフェースの実装
type taskUseCase struct {
	taskRepo repository.TaskRepository
	logger   logger.Logger
}

// NewTaskUseCase は新しいTaskUseCaseImplインスタンスを作成する
func NewTaskUseCase(taskRepo repository.TaskRepository, logger logger.Logger) TaskUseCase {
	return &taskUseCase{
		taskRepo: taskRepo,
		logger:   logger,
	}
}

// CreateTask は新しいタスクを作成する
func (uc *taskUseCase) CreateTask(ctx context.Context, userID uuid.UUID, req *model.CreateTaskRequest) (*model.TaskResponse, error) {
	task := model.NewTask(userID, req.Detail)

	// DBにタスクを保存
	if err := uc.taskRepo.Create(ctx, task); err != nil {
		uc.logger.Errorf("タスク作成エラー: %v", err)
		return nil, err
	}

	// タスクをAPIレスポンス形式に変換
	return task.ToResponse(), nil
}

// GetTask は指定されたIDのタスクを取得する
func (uc *taskUseCase) GetTask(ctx context.Context, id, userID uuid.UUID) (*model.TaskResponse, error) {
	task, err := uc.taskRepo.GetByID(ctx, id, userID)
	if err != nil {
		uc.logger.Errorf("タスク取得エラー: %v", err)
		return nil, err
	}

	return task.ToResponse(), nil
}

// GetAllTasks はユーザーIDに紐づくすべてのタスクを取得する
func (uc *taskUseCase) GetAllTasks(ctx context.Context, userID uuid.UUID) (*model.TasksResponse, error) {
	tasks, err := uc.taskRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		uc.logger.Errorf("タスク取得エラー: %v", err)
		return nil, err
	}

	// タスクをAPIレスポンス形式に変換
	taskResponses := make([]*model.TaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = task.ToResponse()
	}

	return &model.TasksResponse{Tasks: taskResponses}, nil
}

// UpdateTask はタスクの詳細を更新する
func (uc *taskUseCase) UpdateTask(ctx context.Context, id, userID uuid.UUID, req *model.UpdateTaskRequest) (*model.TaskResponse, error) {
	// 現在のタスクを取得
	task, err := uc.taskRepo.GetByID(ctx, id, userID)
	if err != nil {
		uc.logger.Errorf("タスク取得エラー: %v", err)
		return nil, err
	}

	// タスクの詳細を更新
	task.Detail = req.Detail

	// DBにタスクを保存
	if err := uc.taskRepo.Update(ctx, task); err != nil {
		uc.logger.Errorf("タスク更新エラー: %v", err)
		return nil, err
	}

	// 更新後のタスクを取得
	updatedTask, err := uc.taskRepo.GetByID(ctx, id, userID)
	if err != nil {
		uc.logger.Errorf("更新後のタスク取得エラー: %v", err)
		return nil, err
	}

	// タスクをAPIレスポンス形式に変換
	return updatedTask.ToResponse(), nil
}

// ToggleTaskCompletion はタスクの完了状態を切り替える
func (uc *taskUseCase) ToggleTaskCompletion(ctx context.Context, id, userID uuid.UUID) (*model.TaskResponse, error) {
	// タスクの完了状態を切り替える
	if err := uc.taskRepo.ToggleCompletion(ctx, id, userID); err != nil {
		uc.logger.Errorf("タスク完了状態切り替えエラー: %v", err)
		return nil, err
	}

	// 更新後のタスクを取得
	task, err := uc.taskRepo.GetByID(ctx, id, userID)
	if err != nil {
		uc.logger.Errorf("完了状態切り替え後のタスク取得エラー: %v", err)
		return nil, err
	}

	// タスクをAPIレスポンス形式に変換
	return task.ToResponse(), nil
}

// DeleteTask はタスクを削除する
func (uc *taskUseCase) DeleteTask(ctx context.Context, id, userID uuid.UUID) error {
	if err := uc.taskRepo.Delete(ctx, id, userID); err != nil {
		uc.logger.Errorf("タスク削除エラー: %v", err)
		return err
	}
	return nil
}
