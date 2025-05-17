package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
)

// RoundRepository はラウンドに関するデータベース操作を定義するインターフェース
type RoundRepository interface {
	// Create は新しいラウンドを作成する
	Create(ctx context.Context, round *model.Round) error

	// GetByID は指定されたIDのラウンドを取得する
	GetByID(ctx context.Context, id uuid.UUID) (*model.Round, error)

	// GetAllBySessionID は指定されたセッションIDのラウンドを取得する
	GetAllBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*model.Round, error)

	// GetLastBySessionID は指定されたセッションIDの最新のラウンドを取得する
	GetLastRoundBySessionID(ctx context.Context, sessionID uuid.UUID) (*model.Round, error)

	// Complete はラウンドを完了する
	Complete(ctx context.Context, id uuid.UUID, focusScore *int, worktime, breaktime int) error

	// AbortRound はラウンドを中止する
	AbortRound(ctx context.Context, id uuid.UUID) error

	// CalculateSessionState はセッションIDに基づいてラウンドの統計情報を計算する
	CalculateSessionStats(ctx context.Context, sessionID uuid.UUID) (averageFocus float64, totalWorkMin, roundCount, breakTime int, err error)
}
