package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/database"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	rounderrors "github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/postgres/errors"
)

// RoundRepositoryImpl はRoundRepositoryインターフェースの実装部分
type RoundRepositoryImpl struct {
	db     *database.PostgresDB
	logger logger.Logger
}

// NewRoundRepository は新しいRoundRepositoryImplを作成する
func NewRoundRepository(db *database.PostgresDB, logger logger.Logger) repository.RoundRepository {
	return &RoundRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

// Create は新しいラウンドを作成する
func (r *RoundRepositoryImpl) Create(ctx context.Context, round *model.Round) error {
	query := `
		INSERT INTO rounds (id, session_id, round_order, start_time, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		`

	_, err := r.db.DB.ExecContext(ctx, query,
		round.ID,
		round.SessionID,
		round.RoundOrder,
		round.StartTime,
		round.CreatedAt,
		round.UpdatedAt,
	)

	if err != nil {
		r.logger.Errorf("ラウンド作成エラー: %v", err)
		return err
	}
	return nil
}

// GetByID は指定されたIDのラウンドを取得する
// GetByID はIDによってラウンドを取得するメソッド
func (r *RoundRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.Round, error) {
	query := `
		SELECT id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at
		FROM rounds
		WHERE id = $1
	`

	var round model.Round
	err := r.db.DB.GetContext(ctx, &round, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, rounderrors.ErrRoundNotFound
		}
		r.logger.Errorf("ラウンド取得エラー: %v", err)
		return nil, err
	}

	return &round, nil
}

// GetAllBySessionID はセッションIDに紐づくすべてのラウンドを取得するメソッド
func (r *RoundRepositoryImpl) GetAllBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*model.Round, error) {
	query := `
		SELECT id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at
		FROM rounds
		WHERE session_id = $1
		ORDER BY round_order ASC
	`

	var rounds []*model.Round
	err := r.db.DB.SelectContext(ctx, &rounds, query, sessionID)
	if err != nil {
		r.logger.Errorf("ラウンド一覧取得エラー: %v", err)
		return nil, err
	}
	return rounds, nil
}

// GetLastRoundBySessionID はセッションの最後のラウンドを取得するメソッド
func (r *RoundRepositoryImpl) GetLastRoundBySessionID(ctx context.Context, sessionID uuid.UUID) (*model.Round, error) {
	query := `
		SELECT id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at
		FROM rounds
		WHERE session_id = $1
		ORDER BY round_order DESC
		LIMIT 1
	`

	var round model.Round
	err := r.db.DB.GetContext(ctx, &round, query, sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, rounderrors.ErrNoRoundsInSession
		}
		r.logger.Errorf("最終ラウンド取得エラー: %v", err)
		return nil, err
	}

	return &round, nil
}

// Complete はラウンドを完了するメソッド
func (r *RoundRepositoryImpl) Complete(ctx context.Context, id uuid.UUID, focusScore *int, workTime, breakTime int) error {
	// ラウンドが存在するか確認
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err // GetByIDのエラーをそのまま返す
	}

	query := `
		UPDATE rounds
		SET end_time = $1, focus_score = $2, work_time = $3, break_time = $4, updated_at = $5
		WHERE id = $6
	`

	_, err = r.db.DB.ExecContext(ctx, query,
		time.Now(),
		focusScore,
		workTime,
		breakTime,
		time.Now(),
		id,
	)
	if err != nil {
		r.logger.Errorf("ラウンド完了エラー: %v", err)
		return err
	}
	return nil
}

// AbortRound はラウンドを中止としてマークするメソッド
func (r *RoundRepositoryImpl) AbortRound(ctx context.Context, id uuid.UUID) error {
	// ラウンドが存在するか確認
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err // GetByIDのエラーをそのまま返す
	}

	query := `
		UPDATE rounds
		SET end_time = $1, is_aborted = TRUE, updated_at = $2
		WHERE id = $3
	`

	_, err = r.db.DB.ExecContext(ctx, query,
		time.Now(),
		time.Now(),
		id,
	)
	if err != nil {
		r.logger.Errorf("ラウンドスキップエラー: %v", err)
		return err
	}
	return nil
}

// CalculateSessionStats はセッションIDに基づいてラウンドの統計情報を計算するメソッド
func (r *RoundRepositoryImpl) CalculateSessionStats(ctx context.Context, sessionID uuid.UUID) (float64, int, int, int, error) {
	// セッションに属するすべての完了済み（スキップされていない）ラウンドを取得
	query := `
		SELECT id, session_id, round_order, start_time, end_time, work_time, break_time, focus_score, is_aborted, created_at, updated_at
		FROM rounds
		WHERE session_id = $1 AND end_time IS NOT NULL AND is_aborted = FALSE
		ORDER BY round_order ASC
	`

	var rounds []*model.Round
	err := r.db.DB.SelectContext(ctx, &rounds, query, sessionID)
	if err != nil {
		r.logger.Errorf("ラウンド統計情報取得エラー: %v", err)
		return 0, 0, 0, 0, err
	}

	// ラウンドが無い場合
	if len(rounds) == 0 {
		return 0, 0, 0, 0, nil
	}

	// 統計情報を計算
	var totalFocus int
	var validFocusCount int
	var totalWorkMin int
	var totalBreakTime int

	for _, round := range rounds {
		// 集中度の合計（スコアが設定されている場合のみ）
		if round.FocusScore != nil {
			totalFocus += *round.FocusScore
			validFocusCount++
		}

		// 作業時間の合計
		if round.WorkTime != nil {
			totalWorkMin += *round.WorkTime
		}

		// 休憩時間の合計
		if round.BreakTime != nil {
			totalBreakTime += *round.BreakTime
		}
	}

	// 平均集中度の計算（スコアが1つ以上ある場合）
	var averageFocus float64
	if validFocusCount > 0 {
		averageFocus = float64(totalFocus) / float64(validFocusCount)
	}

	return averageFocus, totalWorkMin, len(rounds), totalBreakTime, nil
}
