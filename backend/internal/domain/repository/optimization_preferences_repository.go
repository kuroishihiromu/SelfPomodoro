package repository

import (
	"context"

	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
)

// OptimizationPreferencesRepository は最適化設定永続化のためのリポジトリ
// DDD原則: OptimizationPreferences集約ルートに対するリポジトリ
type OptimizationPreferencesRepository interface {
	// Get は最適化設定を取得する
	Get(ctx context.Context, userID userVO.UserID) (*entity.OptimizationPreferences, error)

	// Create は新しい最適化設定を作成する
	Create(ctx context.Context, preferences *entity.OptimizationPreferences) error

	// Update は最適化設定を更新する
	Update(ctx context.Context, preferences *entity.OptimizationPreferences) error

	// Delete は最適化設定を削除する（GDPR対応等）
	Delete(ctx context.Context, userID userVO.UserID) error

	// GetOrCreateDefault は設定を取得し、存在しない場合はデフォルト設定を作成して返す
	GetOrCreateDefault(ctx context.Context, userID userVO.UserID) (*entity.OptimizationPreferences, error)
}

// Legacy: 下位互換性のためのエイリアス（段階的移行用）
type UserConfigRepository = OptimizationPreferencesRepository
