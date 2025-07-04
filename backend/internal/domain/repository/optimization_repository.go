package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
)

// OptimizationRepository は最適化データの永続化を担当するリポジトリインターフェース
type OptimizationRepository interface {
	// ラウンド最適化関連
	SaveRoundOptimizationLog(ctx context.Context, log *model.RoundOptimizationLog) error
	GetRoundOptimizationHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*model.RoundOptimizationLog, error)
	GetLatestRoundOptimizationResult(ctx context.Context, userID uuid.UUID) (*model.RoundOptimizationLog, error)

	// セッション最適化関連
	SaveSessionOptimizationLog(ctx context.Context, log *model.SessionOptimizationLog) error
	GetSessionOptimizationHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*model.SessionOptimizationLog, error)
	GetLatestSessionOptimizationResult(ctx context.Context, userID uuid.UUID) (*model.SessionOptimizationLog, error)

	// 最適化データ統計
	GetOptimizationEffectiveness(ctx context.Context, userID uuid.UUID, since time.Time) (*model.OptimizationEffectiveness, error)
	CountOptimizationLogs(ctx context.Context, userID uuid.UUID) (roundCount, sessionCount int, err error)
}