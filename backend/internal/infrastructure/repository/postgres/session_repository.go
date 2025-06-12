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

// SessionRepositoryImpl はSessionRepositoryインターフェースの実装（新エラーハンドリング対応版）
type SessionRepositoryImpl struct {
	db     *database.PostgresDB
	logger logger.Logger
}

// NewSessionRepository は新しいSessionRepositoryImplインスタンスを作成する
func NewSessionRepository(db *database.PostgresDB, logger logger.Logger) repository.SessionRepository {
	return &SessionRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

// Create は新しいセッションを作成する（新エラーハンドリング対応版）
func (r *SessionRepositoryImpl) Create(ctx context.Context, session *model.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, start_time, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.DB.ExecContext(ctx, query,
		session.ID,
		session.UserID,
		session.StartTime,
		session.CreatedAt,
		session.UpdatedAt,
	)

	if err != nil {
		// PostgreSQL固有のエラーハンドリング
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				r.logger.Errorf("セッション作成一意制約違反: %v", err)
				return appErrors.NewUniqueConstraintError(err)
			case "23503": // foreign_key_violation
				r.logger.Errorf("セッション作成外部キー制約違反: %v", err)
				return appErrors.NewDatabaseError("create_session_fk", err)
			}
		}
		r.logger.Errorf("セッション作成エラー: %v", err)
		return appErrors.NewDatabaseError("create_session", err)
	}
	return nil
}

// GetByID はIDからセッションを取得する（新エラーハンドリング対応版）
func (r *SessionRepositoryImpl) GetByID(ctx context.Context, id, userID uuid.UUID) (*model.Session, error) {
	query := `
		SELECT id, user_id, start_time, end_time, average_focus, total_work_min, round_count, break_time, created_at, updated_at
		FROM sessions
		WHERE id = $1
	`

	var session model.Session
	err := r.db.DB.GetContext(ctx, &session, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Debugf("セッションが見つかりません: %s", id.String())
			return nil, appErrors.ErrRecordNotFound // Infrastructure Error
		}
		r.logger.Errorf("セッション取得エラー: %v", err)
		return nil, appErrors.NewDatabaseQueryError(err)
	}

	// ユーザーIDの確認（アクセス権限チェック）
	if session.UserID != userID {
		r.logger.Warnf("セッションアクセス権限エラー: SessionID=%s, RequestUserID=%s, OwnerUserID=%s",
			id.String(), userID.String(), session.UserID.String())
		return nil, appErrors.ErrRecordNotFound // アクセス権限エラーも404として扱う
	}

	return &session, nil
}

// GetAllByUserID はユーザIDからセッション一覧を取得する（新エラーハンドリング対応版）
func (r *SessionRepositoryImpl) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Session, error) {
	query := `
		SELECT id, user_id, start_time, end_time, average_focus, total_work_min, round_count, break_time, created_at, updated_at
		FROM sessions
		WHERE user_id = $1
		ORDER BY start_time DESC
	`

	var sessions []*model.Session
	err := r.db.DB.SelectContext(ctx, &sessions, query, userID)
	if err != nil {
		r.logger.Errorf("セッション一覧取得エラー: %v", err)
		return nil, appErrors.NewDatabaseQueryError(err)
	}
	return sessions, nil
}

// Update はセッションを更新する（新エラーハンドリング対応版）
func (r *SessionRepositoryImpl) Update(ctx context.Context, session *model.Session) error {
	// 最初にセッションの所有者を確認
	_, err := r.GetByID(ctx, session.ID, session.UserID)
	if err != nil {
		return err // GetByIDのInfrastructure Errorをそのまま返す
	}

	query := `
		UPDATE sessions
		SET end_time = $1, average_focus = $2, total_work_min = $3, round_count = $4, break_time = $5, updated_at = $6
		WHERE id = $7
	`

	result, err := r.db.DB.ExecContext(ctx, query,
		session.EndTime,
		session.AverageFocus,
		session.TotalWorkMin,
		session.RoundCount,
		session.BreakTime,
		time.Now(),
		session.ID,
	)

	if err != nil {
		r.logger.Errorf("セッション更新エラー: %v", err)
		return appErrors.NewDatabaseError("update_session", err)
	}

	// 更新された行数の確認
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Errorf("セッション更新結果確認エラー: %v", err)
		return appErrors.NewDatabaseError("update_session_check", err)
	}

	if rowsAffected == 0 {
		r.logger.Warnf("セッション更新対象なし: %s", session.ID.String())
		return appErrors.ErrRecordNotFound
	}

	return nil
}

// Complete はセッションを完了する（新エラーハンドリング対応版）
func (r *SessionRepositoryImpl) Complete(ctx context.Context, id, userID uuid.UUID, averageFocus float64, totalWorkMin, roundCount, breakTime int) error {
	// 最初にセッションの所有者を確認
	session, err := r.GetByID(ctx, id, userID)
	if err != nil {
		return err // GetByIDのInfrastructure Errorをそのまま返す
	}

	// セッションが既に完了しているかチェック
	if session.IsCompleted() {
		r.logger.Warnf("セッションは既に完了しています: %s", id.String())
		return appErrors.ErrRecordNotFound // ビジネス的には409だが、Infrastructure層では404として扱う
	}

	endTime := time.Now()
	query := `
		UPDATE sessions
		SET end_time = $1, average_focus = $2, total_work_min = $3, round_count = $4, break_time = $5, updated_at = $6
		WHERE id = $7
	`

	result, err := r.db.DB.ExecContext(ctx, query,
		endTime,
		averageFocus,
		totalWorkMin,
		roundCount,
		breakTime,
		time.Now(),
		id,
	)

	if err != nil {
		r.logger.Errorf("セッション完了エラー: %v", err)
		return appErrors.NewDatabaseError("complete_session", err)
	}

	// 更新された行数の確認
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Errorf("セッション完了結果確認エラー: %v", err)
		return appErrors.NewDatabaseError("complete_session_check", err)
	}

	if rowsAffected == 0 {
		r.logger.Warnf("セッション完了対象なし: %s", id.String())
		return appErrors.ErrRecordNotFound
	}

	return nil
}

// Delete はセッションを削除する（新エラーハンドリング対応版）
func (r *SessionRepositoryImpl) Delete(ctx context.Context, id, userID uuid.UUID) error {
	// 最初にセッションの所有者を確認
	_, err := r.GetByID(ctx, id, userID)
	if err != nil {
		return err // GetByIDのInfrastructure Errorをそのまま返す
	}

	query := `
		DELETE FROM sessions
		WHERE id = $1
	`

	result, err := r.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Errorf("セッション削除エラー: %v", err)
		return appErrors.NewDatabaseError("delete_session", err)
	}

	// 削除された行数の確認
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Errorf("セッション削除結果確認エラー: %v", err)
		return appErrors.NewDatabaseError("delete_session_check", err)
	}

	if rowsAffected == 0 {
		r.logger.Warnf("セッション削除対象なし: %s", id.String())
		return appErrors.ErrRecordNotFound
	}

	return nil
}
