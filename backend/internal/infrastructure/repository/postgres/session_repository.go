package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/database"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	sessionerrors "github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/postgres/errors"
)

// SessionRepositoryImpl はSessionRepositoryインターフェースの実装部分
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

// Create は新しいセッションを作成するメソッド
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
		r.logger.Errorf("セッション作成エラー: %v", err)
		return fmt.Errorf("%w: %v", sessionerrors.ErrSessionCreationFailed, err)
	}
	return nil
}

// GetByID はIDからセッションを取得するメソッド
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
			return nil, sessionerrors.ErrSessionNotFound
		}
		r.logger.Errorf("セッション取得エラー: %v", err)
		return nil, err
	}

	// ユーザーIDの確認
	if session.UserID != userID {
		return nil, sessionerrors.ErrSessionAccessDenied
	}
	return &session, nil
}

// GetAllByUserID はユーザIDからセッションを取得するメソッド
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
		r.logger.Errorf("セッション取得エラー: %v", err)
		return nil, err
	}
	return sessions, nil
}

// Update はセッションを更新するメソッド
func (r *SessionRepositoryImpl) Update(ctx context.Context, session *model.Session) error {
	// 最初にセッションの所有者を確認
	_, err := r.GetByID(ctx, session.ID, session.UserID)
	if err != nil {
		return err // GetByIDのエラーをそのまま返す
	}

	query := `
		UPDATE sessions
		SET end_time = $1, average_focus = $2, total_work_min = $3, round_count = $4, break_time = $5, updated_at = $6
		WHERE id = $7
	`

	_, err = r.db.DB.ExecContext(ctx, query,
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
		return fmt.Errorf("%w: %v", sessionerrors.ErrSessionUpdateFailed, err)
	}
	return nil
}

// Complete はセッションを完了するメソッド
func (r *SessionRepositoryImpl) Complete(ctx context.Context, id, userID uuid.UUID, averageFocus float64, totalWorkMin, roundCount, breakTime int) error {
	// 最初にセッションの所有者を確認
	session, err := r.GetByID(ctx, id, userID)
	if err != nil {
		return err // GetByIDのエラーをそのまま返す
	}

	endTime := time.Now()
	session.EndTime = &endTime

	// 計算結果をセット
	avgFocus := averageFocus
	session.AverageFocus = &avgFocus

	totalWork := totalWorkMin
	session.TotalWorkMin = &totalWork

	rounds := roundCount
	session.RoundCount = &rounds

	breakT := breakTime
	session.BreakTime = &breakT

	return r.Update(ctx, session)
}

// Delete はセッションを削除するメソッド
func (r *SessionRepositoryImpl) Delete(ctx context.Context, id, userID uuid.UUID) error {
	// 最初にセッションの所有者を確認
	_, err := r.GetByID(ctx, id, userID)
	if err != nil {
		return err // GetByIDのエラーをそのまま返す
	}

	query := `
		DELETE FROM sessions
		WHERE id = $1
	`

	_, err = r.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Errorf("セッション削除エラー: %v", err)
		return fmt.Errorf("%w: %v", sessionerrors.ErrSessionDeleteFailed, err)
	}
	return nil
}
