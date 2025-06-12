package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/database"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// TaskRepositoryImpl はTaskRepositoryインターフェースの実際の実装部分
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
		// PostgreSQL固有のエラーハンドリング
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				r.logger.Errorf("タスク作成一意制約違反: %v", err)
				return appErrors.NewUniqueConstraintError(err)
			case "23503": // foreign_key_violation
				r.logger.Errorf("タスク作成外部キー制約違反: %v", err)
				return appErrors.NewDatabaseError("create_task_fk", err)
			}
		}
		r.logger.Errorf("タスク作成エラー: %v", err)
		return appErrors.NewDatabaseError("create_task", err)
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
			r.logger.Debugf("タスクが見つかりません: %s", id.String())
			return nil, appErrors.ErrRecordNotFound // Infrastructure Error
		}
		r.logger.Errorf("タスク取得エラー: %v", err)
		return nil, appErrors.NewDatabaseQueryError(err)
	}

	// ユーザーIDの確認（アクセス権限チェック）
	if task.UserID != userID {
		r.logger.Warnf("タスクアクセス権限エラー: TaskID=%s, RequestUserID=%s, OwnerUserID=%s",
			id.String(), userID.String(), task.UserID.String())
		return nil, appErrors.ErrRecordNotFound // アクセス権限エラーも404として扱う
	}

	return &task, nil
}

// GetAllByUserID はユーザーIDに紐づくすべてのタスクを取得するメソッド
func (r *TaskRepositoryImpl) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Task, error) {
	query := `
		SELECT id, user_id, detail, is_completed, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	var tasks []*model.Task
	err := r.db.DB.SelectContext(ctx, &tasks, query, userID)
	if err != nil {
		r.logger.Errorf("タスク一覧取得エラー: %v", err)
		return nil, appErrors.NewDatabaseQueryError(err)
	}
	return tasks, nil
}

// Update はタスクを更新するメソッド
func (r *TaskRepositoryImpl) Update(ctx context.Context, task *model.Task) error {
	// 最初にタスクの所有者を確認
	_, err := r.GetByID(ctx, task.ID, task.UserID)
	if err != nil {
		return err // GetByIDのInfrastructure Errorをそのまま返す
	}

	query := `
		UPDATE tasks
		SET detail = $1, updated_at = $2
		WHERE id = $3
	`

	result, err := r.db.DB.ExecContext(ctx, query,
		task.Detail,
		time.Now(),
		task.ID,
	)
	if err != nil {
		r.logger.Errorf("タスク更新エラー: %v", err)
		return appErrors.NewDatabaseError("update_task", err)
	}

	// 更新された行数の確認
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Errorf("タスク更新結果確認エラー: %v", err)
		return appErrors.NewDatabaseError("update_task_check", err)
	}

	if rowsAffected == 0 {
		r.logger.Warnf("タスク更新対象なし: %s", task.ID.String())
		return appErrors.ErrRecordNotFound
	}

	return nil
}

// ToggleCompletion はタスクの完了状態をトグルするメソッド
func (r *TaskRepositoryImpl) ToggleCompletion(ctx context.Context, id, userID uuid.UUID) error {
	// 最初にタスクの所有者を確認
	_, err := r.GetByID(ctx, id, userID)
	if err != nil {
		return err // GetByIDのInfrastructure Errorをそのまま返す
	}

	query := `
		UPDATE tasks
		SET is_completed = NOT is_completed, updated_at = $1
		WHERE id = $2
	`

	result, err := r.db.DB.ExecContext(ctx, query,
		time.Now(),
		id,
	)
	if err != nil {
		r.logger.Errorf("タスク完了状態更新エラー: %v", err)
		return appErrors.NewDatabaseError("toggle_task", err)
	}

	// 更新された行数の確認
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Errorf("タスク完了状態更新結果確認エラー: %v", err)
		return appErrors.NewDatabaseError("toggle_task_check", err)
	}

	if rowsAffected == 0 {
		r.logger.Warnf("タスク完了状態更新対象なし: %s", id.String())
		return appErrors.ErrRecordNotFound
	}

	return nil
}

// Delete はタスクを削除するメソッド
func (r *TaskRepositoryImpl) Delete(ctx context.Context, id, userID uuid.UUID) error {
	// 最初にタスクの所有者を確認
	_, err := r.GetByID(ctx, id, userID)
	if err != nil {
		return err // GetByIDのInfrastructure Errorをそのまま返す
	}

	query := `
		DELETE FROM tasks
		WHERE id = $1
	`

	result, err := r.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Errorf("タスク削除エラー: %v", err)
		return appErrors.NewDatabaseError("delete_task", err)
	}

	// 削除された行数の確認
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Errorf("タスク削除結果確認エラー: %v", err)
		return appErrors.NewDatabaseError("delete_task_check", err)
	}

	if rowsAffected == 0 {
		r.logger.Warnf("タスク削除対象なし: %s", id.String())
		return appErrors.ErrRecordNotFound
	}

	return nil
}
