package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
)

// SampleOptimizationDataRepository はオンボーディング時のサンプル最適化データ作成用Repository
type SampleOptimizationDataRepository interface {
	// CreateSampleOptimizationData は新規ユーザー向けのサンプル最適化データを作成する
	CreateSampleOptimizationData(ctx context.Context, userID uuid.UUID) error

	// CreateRoundOptimizationLogs はラウンド最適化ログを一括作成する
	CreateRoundOptimizationLogs(ctx context.Context, logs []*model.RoundOptimizationLog) error

	// CreateSessionOptimizationLogs はセッション最適化ログを一括作成する
	CreateSessionOptimizationLogs(ctx context.Context, logs []*model.SessionOptimizationLog) error
}
