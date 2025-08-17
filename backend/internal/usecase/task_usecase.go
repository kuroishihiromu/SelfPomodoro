package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/usecase/dto"
	"github.com/tsunakit99/selfpomodoro/internal/usecase/mapper"
	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// TaskUseCase はタスクに関するユースケースを定義するインターフェース
type TaskUseCase interface {
	// CreateTask は新しいタスクを作成する
	CreateTask(ctx context.Context, userID uuid.UUID, req *dto.CreateTaskRequest) (*dto.TaskResponse, error)

	// GetTask は指定されたIDのタスクを取得する
	GetTask(ctx context.Context, id, userID uuid.UUID) (*dto.TaskResponse, error)

	// GetAllTasks はユーザーIDに紐づくすべてのタスクを取得する
	GetAllTasks(ctx context.Context, userID uuid.UUID) (*dto.TasksResponse, error)

	// UpdateTask はタスクの詳細を更新する
	UpdateTask(ctx context.Context, id, userID uuid.UUID, req *dto.UpdateTaskRequest) (*dto.TaskResponse, error)

	// ToggleTaskCompletion はタスクの完了状態を切り替える
	ToggleTaskCompletion(ctx context.Context, id, userID uuid.UUID) (*dto.TaskResponse, error)

	// DeleteTask はタスクを削除する
	DeleteTask(ctx context.Context, id, userID uuid.UUID) error
}

// taskUseCase はTaskUseCaseインターフェースの実装（新エラーハンドリング対応版）
type taskUseCase struct {
	taskRepo   repository.TaskRepository
	taskMapper *mapper.TaskMapper
	logger     logger.Logger
}

// NewTaskUseCase は新しいTaskUseCaseImplインスタンスを作成する
func NewTaskUseCase(taskRepo repository.TaskRepository, logger logger.Logger) TaskUseCase {
	return &taskUseCase{
		taskRepo:   taskRepo,
		taskMapper: mapper.NewTaskMapper(),
		logger:     logger,
	}
}

// CreateTask は新しいタスクを作成する（新エラーハンドリング対応版）
func (uc *taskUseCase) CreateTask(ctx context.Context, userID uuid.UUID, req *dto.CreateTaskRequest) (*dto.TaskResponse, error) {
	// UUIDをValue Objectに変換
	userIDVO, err := userVO.NewUserIDFromString(userID.String())
	if err != nil {
		return nil, appErrors.NewInternalError(err)
	}

	task := entity.NewTask(userIDVO, req.Detail)

	if err := uc.taskRepo.Create(ctx, task); err != nil {
		uc.logger.Errorf("タスク作成エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrUniqueConstraint) {
			return nil, appErrors.NewTaskAlreadyDoneError() // ビジネス観点でのエラー
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}
	return uc.taskMapper.ToTaskResponse(task), nil
}

// GetTask は指定されたIDのタスクを取得する（新エラーハンドリング対応版）
func (uc *taskUseCase) GetTask(ctx context.Context, id, userID uuid.UUID) (*dto.TaskResponse, error) {
	// UUIDをValue Objectに変換
	userIDVO, err := userVO.NewUserIDFromString(userID.String())
	if err != nil {
		return nil, appErrors.NewInternalError(err)
	}

	task, err := uc.taskRepo.GetByID(ctx, id, userIDVO)
	if err != nil {
		uc.logger.Errorf("タスク取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return nil, appErrors.NewTaskNotFoundError() // Domain Error
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}
	return uc.taskMapper.ToTaskResponse(task), nil
}

// GetAllTasks はユーザーIDに紐づくすべてのタスクを取得する（新エラーハンドリング対応版）
func (uc *taskUseCase) GetAllTasks(ctx context.Context, userID uuid.UUID) (*dto.TasksResponse, error) {
	// UUIDをValue Objectに変換
	userIDVO, err := userVO.NewUserIDFromString(userID.String())
	if err != nil {
		return nil, appErrors.NewInternalError(err)
	}

	tasks, err := uc.taskRepo.GetAllByUserID(ctx, userIDVO)
	if err != nil {
		uc.logger.Errorf("タスク一覧取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// Mapperを使用してタスク一覧をレスポンス形式に変換（削除された行）
	return uc.taskMapper.ToTasksResponse(tasks), nil
}

// UpdateTask はタスクの詳細を更新する（新エラーハンドリング対応版）
func (uc *taskUseCase) UpdateTask(ctx context.Context, id, userID uuid.UUID, req *dto.UpdateTaskRequest) (*dto.TaskResponse, error) {
	// UUIDをValue Objectに変換
	userIDVO, err := userVO.NewUserIDFromString(userID.String())
	if err != nil {
		return nil, appErrors.NewInternalError(err)
	}

	task, err := uc.taskRepo.GetByID(ctx, id, userIDVO)
	if err != nil {
		uc.logger.Errorf("タスク取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return nil, appErrors.NewTaskNotFoundError() // Domain Error
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// タスク詳細を更新（Value Object使用）
	if err := task.UpdateDetail(req.Detail); err != nil {
		return nil, appErrors.NewInternalError(err)
	}
	if err := uc.taskRepo.Update(ctx, task); err != nil {
		uc.logger.Errorf("タスク更新エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return nil, appErrors.NewTaskNotFoundError() // Domain Error
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// 更新後のタスクを取得
	updatedTask, err := uc.taskRepo.GetByID(ctx, id, userIDVO)
	if err != nil {
		uc.logger.Errorf("更新後のタスク取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return nil, appErrors.NewTaskNotFoundError() // Domain Error
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}
	return uc.taskMapper.ToTaskResponse(updatedTask), nil
}

// ToggleTaskCompletion はタスクの完了状態を切り替える（新エラーハンドリング対応版）
func (uc *taskUseCase) ToggleTaskCompletion(ctx context.Context, id, userID uuid.UUID) (*dto.TaskResponse, error) {
	// UUIDをValue Objectに変換
	userIDVO, err := userVO.NewUserIDFromString(userID.String())
	if err != nil {
		return nil, appErrors.NewInternalError(err)
	}

	if err := uc.taskRepo.ToggleCompletion(ctx, id, userIDVO); err != nil {
		uc.logger.Errorf("タスク完了状態切り替えエラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return nil, appErrors.NewTaskNotFoundError() // Domain Error
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// 切り替え後のタスクを取得
	task, err := uc.taskRepo.GetByID(ctx, id, userIDVO)
	if err != nil {
		uc.logger.Errorf("完了状態切り替え後のタスク取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return nil, appErrors.NewTaskNotFoundError() // Domain Error
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}
	return uc.taskMapper.ToTaskResponse(task), nil
}

// DeleteTask はタスクを削除する（新エラーハンドリング対応版）
func (uc *taskUseCase) DeleteTask(ctx context.Context, id, userID uuid.UUID) error {
	// UUIDをValue Objectに変換
	userIDVO, err := userVO.NewUserIDFromString(userID.String())
	if err != nil {
		return appErrors.NewInternalError(err)
	}

	if err := uc.taskRepo.Delete(ctx, id, userIDVO); err != nil {
		uc.logger.Errorf("タスク削除エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return appErrors.NewTaskNotFoundError() // Domain Error
		}
		if appErrors.IsDatabaseError(err) {
			return appErrors.NewInternalError(err)
		}

		return appErrors.NewInternalError(err)
	}
	return nil
}
