package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
)

// SessionRepository はセッション永続化のためのインターフェース
type SessionRepository interface {
	// Create は新しいセッションを作成する
	Create(ctx context.Context, session *model.Session) error

	// GetByID はIDからセッションを取得する
	GetByID(ctx context.Context, id, userID uuid.UUID) (*model.Session, error)

	// GetAllByUserID はユーザIDからセッションを取得する
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Session, error)

	// Update はセッションを更新する
	Update(ctx context.Context, session *model.Session) error

	// Complete はセッションを完了する(終了時刻、平均集中度、総作業時間を設定)
	Complete(ctx context.Context, id, userID uuid.UUID, averageFocus float64, totalWorkMin, roundCount, breakTime int) error

	// Delete はセッションを削除する
	Delete(ctx context.Context, id, userID uuid.UUID) error
}
