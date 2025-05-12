package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/database"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// タスクリポジトリに関するエラー
var (
	ErrTaskNotFound       = errors.New("タスクが見つかりません")
	ErrTaskAccessDenied   = errors.New("このタスクへのアクセス権限がありません")
	ErrTaskCreationFailed = errors.New("タスクの作成に失敗しました")
	ErrTaskUpdateFailed   = errors.New("タスクの更新に失敗しました")
	ErrTaskDeleteFailed   = errors.New("タスクの削除に失敗しました")
)

// TaskRepositoryImpl はTaskRepositoryインターフェースのの実際の実装部分
type TaskRepositoryImpl struct {
	db     *database.PostgresDB
	logger logger.Logger
}

// NewTaskRepository は新しいTaskRepositoryImplインスタンスを作成する
func NewTaskRepository(db *database.PostgresDB, logger logger.Logger) repository.TaskRepository {
	return &TaskRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

// Create はタスクを作成するメソッド
func (r *TaskRepositoryImpl) Create(ctx context.Context, task *model.Task) error {
	query := `
		INSERT INTO tasks (id, user_id, detail, is_completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.DB.ExecContext(ctx, query,
		task.ID,
		task.UserID,
		task.Detail,
		task.IsCompleted,
		task.CreatedAt,
		task.UpdatedAt,
	)

	if err != nil {
		r.logger.Errorf(("タスク作成エラー: %v"), err)
		return fmt.Errorf("%w: %v", ErrTaskCreationFailed, err)
	}
	return nil
}

// GetByID はIDによってタスクを取得するメソッド
func (r *TaskRepositoryImpl) GetByID(ctx context.Context, id, userID uuid.UUID) (*model.Task, error) {
	query := `
		SELECT id, user_id, detail, is_completed, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	var task model.Task
	err := r.db.DB.GetContext(ctx, &task, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTaskNotFound
		}
		r.logger.Errorf("タスク取得エラー: %v", err)
		return nil, err
	}

	// ユーザーIDの確認
	if task.UserID != userID {
		return nil, ErrTaskAccessDenied
	}
	return &task, nil
}

// GetAllByUserID はユーザーIDに紐づくすべてのタスクを取得するメソッド
func (r *TaskRepositoryImpl) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Task, error) {
	query := `
		SELECT id, user_id, detail, is_completed, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
	`

	var tasks []*model.Task
	err := r.db.DB.SelectContext(ctx, &tasks, query, userID)
	if err != nil {
		r.logger.Errorf("タスク一覧取得エラー: %v", err)
		return nil, err
	}
	return tasks, nil
}

// Update はタスクを更新するメソッド
func (r *TaskRepositoryImpl) Update(ctx context.Context, task *model.Task) error {
	// 最初にタスクの所有者を確認
	_, err := r.GetByID(ctx, task.ID, task.UserID)
	if err != nil {
		return err // GetByIDのエラーをそのまま返す
	}

	query := `
		UPDATE tasks
		SET detail = $1, updated_at = $2
		WHERE id = $3
	`

	_, err = r.db.DB.ExecContext(ctx, query,
		task.Detail,
		time.Now(),
		task.ID,
	)
	if err != nil {
		r.logger.Errorf("タスク更新エラー: %v", err)
		return fmt.Errorf("%w: %v", ErrTaskUpdateFailed, err)
	}
	return nil
}

// ToggleCompletion はタスクの完了状態をトグルするメソッド
func (r *TaskRepositoryImpl) ToggleCompletion(ctx context.Context, id, userID uuid.UUID) error {
	// 最初にタスクの所有者を確認
	_, err := r.GetByID(ctx, id, userID)
	if err != nil {
		return err // GetByIDのエラーをそのまま返す
	}

	query := `
		UPDATE tasks
		SET is_completed = NOT is_completed, updated_at = $1
		WHERE id = $2
	`

	_, err = r.db.DB.ExecContext(ctx, query,
		time.Now(),
		id,
	)
	if err != nil {
		r.logger.Errorf("タスク完了状態更新エラー: %v", err)
		return fmt.Errorf("%w: %v", ErrTaskUpdateFailed, err)
	}
	return nil
}

// Delete はタスクを削除するメソッド
func (r *TaskRepositoryImpl) Delete(ctx context.Context, id, userID uuid.UUID) error {
	// 最初にタスクの所有者を確認
	_, err := r.GetByID(ctx, id, userID)
	if err != nil {
		return err // GetByIDのエラーをそのまま返す
	}

	query := `
		DELETE FROM tasks
		WHERE id = $1
	`

	_, err = r.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Errorf("タスク削除エラー: %v", err)
		return fmt.Errorf("%w: %v", ErrTaskDeleteFailed, err)
	}
	return nil
}
