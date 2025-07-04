package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
)

// RoundRepository はラウンドに関するデータベース操作を定義するインターフェース
type RoundRepository interface {
	// Create は新しいラウンドを作成する
	Create(ctx context.Context, round *model.Round, userID uuid.UUID) error

	// GetByID は指定されたIDのラウンドを取得する
	GetByID(ctx context.Context, id uuid.UUID) (*model.Round, error)

	// GetBySessionIDWithUserID は指定されたセッションIDとユーザーIDのラウンドを効率的に取得する（上限チェック用）
	GetBySessionIDWithUserID(ctx context.Context, sessionID, userID uuid.UUID) ([]*model.Round, error)

	// Complete はラウンドを完了する
	Complete(ctx context.Context, id uuid.UUID, focusScore *int, worktime, breaktime int) error
}
