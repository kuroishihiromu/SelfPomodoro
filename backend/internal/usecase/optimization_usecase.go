package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// OptimizationLogEntry は最適化ログの基本構造（一時的定義）
type OptimizationLogEntry struct {
	UserID    userVO.UserID `json:"user_id"`
	Timestamp time.Time     `json:"timestamp"`
	LogType   string        `json:"log_type"` // "round" or "session"
}

// OptimizationUseCase は最適化に関するビジネスロジックを担当
type OptimizationUseCase interface {
	// UserConfigへの最適化結果適用
	ApplyOptimizationToUserConfig(ctx context.Context, userID uuid.UUID) error
	
	// 最適化履歴の取得（簡易版）
	GetOptimizationHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*OptimizationLogEntry, error)
}

// optimizationUseCase はOptimizationUseCaseの実装
type optimizationUseCase struct {
	optimizationPreferencesRepo repository.OptimizationPreferencesRepository
	logger                      logger.Logger
}

// NewOptimizationUseCase は新しい最適化UseCaseを作成する
func NewOptimizationUseCase(
	optimizationPreferencesRepo repository.OptimizationPreferencesRepository,
	logger logger.Logger,
) OptimizationUseCase {
	return &optimizationUseCase{
		optimizationPreferencesRepo: optimizationPreferencesRepo,
		logger:                      logger,
	}
}

// GetOptimizationHistory は最適化履歴を取得する（簡易版）
func (uc *optimizationUseCase) GetOptimizationHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*OptimizationLogEntry, error) {
	// UUIDをValue Objectに変換
	userIDVO, err := userVO.NewUserIDFromString(userID.String())
	if err != nil {
		return nil, appErrors.NewInternalError(err)
	}

	uc.logger.Infof("最適化履歴取得開始: UserID=%s, Limit=%d", userID.String()[:8]+"...", limit)

	if limit <= 0 || limit > 100 {
		limit = 50 // デフォルト制限
	}

	// 現在は空の履歴を返す（将来の実装用）
	history := make([]*OptimizationLogEntry, 0)

	uc.logger.Infof("最適化履歴取得完了: UserID=%s, 件数=%d", userIDVO.String()[:8]+"...", len(history))
	return history, nil
}

// ApplyOptimizationToUserConfig は最適化結果をUserConfigに適用する
func (uc *optimizationUseCase) ApplyOptimizationToUserConfig(ctx context.Context, userID uuid.UUID) error {
	// UUIDをValue Objectに変換
	userIDVO, err := userVO.NewUserIDFromString(userID.String())
	if err != nil {
		return appErrors.NewInternalError(err)
	}

	uc.logger.Infof("最適化結果のUserConfig適用開始: UserID=%s", userIDVO.String()[:8]+"...")

	// 現在のUserConfigを取得
	currentConfig, err := uc.optimizationPreferencesRepo.Get(ctx, userIDVO)
	if err != nil {
		uc.logger.Errorf("UserConfig取得失敗: %v", err)
		return err
	}

	// 現在は何も変更せずに完了（将来の最適化ロジック実装用）
	uc.logger.Infof("最適化結果適用（現在は変更なし）: UserID=%s", userIDVO.String()[:8]+"...")
	
	// 変更があった場合の更新処理のテンプレート
	configUpdated := false
	if configUpdated {
		err = uc.optimizationPreferencesRepo.Update(ctx, currentConfig)
		if err != nil {
			uc.logger.Errorf("UserConfig更新失敗: %v", err)
			return err
		}
		uc.logger.Infof("最適化結果のUserConfig適用完了: UserID=%s", userIDVO.String()[:8]+"...")
	} else {
		uc.logger.Infof("最適化結果の変更なし: UserID=%s", userIDVO.String()[:8]+"...")
	}

	return nil
}