package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// UserConfigUseCase はユーザー設定に関するユースケースを定義するインターフェース
type UserConfigUseCase interface {
	// GetUserConfig はユーザー設定を取得する（デフォルト値フォールバック付き）
	GetUserConfig(ctx context.Context, userID uuid.UUID) (*model.UserConfigResponse, error)

	// CreateUserConfig は新しいユーザー設定を作成する（PostConfirmation専用）
	CreateUserConfig(ctx context.Context, userID uuid.UUID, req *model.CreateUserConfigRequest) (*model.UserConfigResponse, error)

	// UpdateUserConfig はユーザー設定を更新する
	UpdateUserConfig(ctx context.Context, userID uuid.UUID, req *model.UpdateUserConfigRequest) (*model.UserConfigResponse, error)

	// DeleteUserConfig はユーザー設定を削除する
	DeleteUserConfig(ctx context.Context, userID uuid.UUID) error

	// GetUserConfigForOptimization は最適化処理用にユーザー設定を取得する（内部使用・デフォルト値フォールバック）
	GetUserConfigForOptimization(ctx context.Context, userID uuid.UUID) (*model.UserConfig, error)
}

// userConfigUseCase はUserConfigUseCaseインターフェースの実装（ドメイン強化版）
type userConfigUseCase struct {
	userConfigRepo repository.UserConfigRepository
	logger         logger.Logger
}

// NewUserConfigUseCase は新しいUserConfigUseCaseインスタンスを作成する
func NewUserConfigUseCase(userConfigRepo repository.UserConfigRepository, logger logger.Logger) UserConfigUseCase {
	return &userConfigUseCase{
		userConfigRepo: userConfigRepo,
		logger:         logger,
	}
}

// GetUserConfig はユーザー設定を取得する（ドメイン強化版・デフォルト値フォールバック付き）
func (uc *userConfigUseCase) GetUserConfig(ctx context.Context, userID uuid.UUID) (*model.UserConfigResponse, error) {
	if uc.userConfigRepo == nil {
		uc.logger.Warn("UserConfigRepository が初期化されていません。デフォルト設定を返します")
		// ✅ ドメインファクトリー使用：DynamoDBが利用できない場合のフォールバック
		defaultConfig := model.NewDefaultUserConfig(userID)
		return defaultConfig.ToResponse(), nil
	}

	config, err := uc.userConfigRepo.GetUserConfig(ctx, userID)
	if err != nil {
		uc.logger.Errorf("ユーザー設定取得エラー、デフォルト設定を返します: %v", err)
		// ✅ ドメインファクトリー使用：設定が存在しない場合のデフォルト値フォールバック
		defaultConfig := model.NewDefaultUserConfig(userID)

		// ✅ ドメインロジック活用：デフォルト値ログ出力
		uc.logger.Infof("デフォルト設定使用: work=%d分, break=%d分, rounds=%d, sessionBreak=%d分",
			defaultConfig.GetWorkTimeOrDefault(),
			defaultConfig.GetBreakTimeOrDefault(),
			defaultConfig.GetSessionRoundsOrDefault(),
			defaultConfig.GetSessionBreakTimeOrDefault())

		return defaultConfig.ToResponse(), nil
	}

	// ✅ ドメインロジック活用：取得した設定値の安全性確認
	uc.logger.Infof("ユーザー設定取得成功: work=%d分, break=%d分, rounds=%d, sessionBreak=%d分",
		config.GetWorkTimeOrDefault(),
		config.GetBreakTimeOrDefault(),
		config.GetSessionRoundsOrDefault(),
		config.GetSessionBreakTimeOrDefault())

	return config.ToResponse(), nil
}

// CreateUserConfig は新しいユーザー設定を作成する（ドメイン強化版）
func (uc *userConfigUseCase) CreateUserConfig(ctx context.Context, userID uuid.UUID, req *model.CreateUserConfigRequest) (*model.UserConfigResponse, error) {
	if uc.userConfigRepo == nil {
		return nil, appErrors.NewInternalError(fmt.Errorf("ユーザー設定機能は現在利用できません"))
	}

	// ✅ ドメインファクトリー使用：リクエストから新しい設定を作成
	config := model.NewDefaultUserConfig(userID)
	config.UpdateSettings(req.RoundWorkTime, req.RoundBreakTime, req.SessionRounds, req.SessionBreakTime)

	// ✅ ドメインロジック活用：設定の有効性をチェック
	if err := config.ValidateSettings(); err != nil {
		uc.logger.Errorf("ユーザー設定バリデーションエラー: %v", err)
		return nil, appErrors.NewValidationError(err.Error())
	}

	// リポジトリに保存
	if err := uc.userConfigRepo.CreateUserConfig(ctx, config); err != nil {
		uc.logger.Errorf("ユーザー設定作成エラー: %v", err)
		if errors.Is(err, appErrors.ErrUserConfigCreateFailed) {
			return nil, appErrors.NewUserConfigCreateFailedError()
		}
		return nil, appErrors.NewInternalError(err)
	}

	uc.logger.Infof("ユーザー設定作成成功: UserID=%s, work=%d分, break=%d分, rounds=%d",
		userID.String()[:8]+"...", config.RoundWorkTime, config.RoundBreakTime, config.SessionRounds)

	return config.ToResponse(), nil
}

// UpdateUserConfig はユーザー設定を更新する（ドメイン強化版）
func (uc *userConfigUseCase) UpdateUserConfig(ctx context.Context, userID uuid.UUID, req *model.UpdateUserConfigRequest) (*model.UserConfigResponse, error) {
	if uc.userConfigRepo == nil {
		return nil, appErrors.NewInternalError(fmt.Errorf("ユーザー設定機能は現在利用できません"))
	}

	// 既存の設定を取得（存在前提）
	config, err := uc.userConfigRepo.GetUserConfig(ctx, userID)
	if err != nil {
		uc.logger.Errorf("ユーザー設定取得エラー: %v", err)
		if errors.Is(err, appErrors.ErrUserConfigNotFound) {
			return nil, appErrors.NewUserConfigNotFoundError()
		}
		return nil, appErrors.NewInternalError(err)
	}

	// ✅ ドメインロジック活用：設定を部分更新（部分更新対応）
	uc.applyPartialUpdate(config, req)

	// ✅ ドメインロジック活用：更新後の設定の有効性をチェック
	if err := config.ValidateSettings(); err != nil {
		uc.logger.Errorf("ユーザー設定更新バリデーションエラー: %v", err)
		return nil, appErrors.NewValidationError(err.Error())
	}

	// リポジトリに保存
	if err := uc.userConfigRepo.UpdateUserConfig(ctx, config); err != nil {
		uc.logger.Errorf("ユーザー設定更新エラー: %v", err)
		if errors.Is(err, appErrors.ErrUserConfigUpdateFailed) {
			return nil, appErrors.NewUserConfigUpdateFailedError()
		}
		return nil, appErrors.NewInternalError(err)
	}

	uc.logger.Infof("ユーザー設定更新成功: UserID=%s, work=%d分, break=%d分, rounds=%d",
		userID.String()[:8]+"...", config.RoundWorkTime, config.RoundBreakTime, config.SessionRounds)

	return config.ToResponse(), nil
}

// DeleteUserConfig はユーザー設定を削除する
func (uc *userConfigUseCase) DeleteUserConfig(ctx context.Context, userID uuid.UUID) error {
	if uc.userConfigRepo == nil {
		return appErrors.NewInternalError(fmt.Errorf("ユーザー設定機能は現在利用できません"))
	}

	if err := uc.userConfigRepo.DeleteUserConfig(ctx, userID); err != nil {
		uc.logger.Errorf("ユーザー設定削除エラー: %v", err)
		if errors.Is(err, appErrors.ErrUserConfigNotFound) {
			return appErrors.NewUserConfigNotFoundError()
		}
		if errors.Is(err, appErrors.ErrUserConfigDeleteFailed) {
			return appErrors.NewUserConfigDeleteFailedError()
		}
		return appErrors.NewInternalError(err)
	}

	uc.logger.Infof("ユーザー設定削除成功: UserID=%s", userID.String()[:8]+"...")
	return nil
}

// GetUserConfigForOptimization は最適化処理用にユーザー設定を取得する（ドメイン強化版・デフォルト値フォールバック）
func (uc *userConfigUseCase) GetUserConfigForOptimization(ctx context.Context, userID uuid.UUID) (*model.UserConfig, error) {
	if uc.userConfigRepo == nil {
		uc.logger.Warn("UserConfigRepository が初期化されていません。デフォルト設定を返します")
		// ✅ ドメインファクトリー使用
		return model.NewDefaultUserConfig(userID), nil
	}

	config, err := uc.userConfigRepo.GetUserConfig(ctx, userID)
	if err != nil {
		uc.logger.Warnf("最適化用ユーザー設定取得エラー、デフォルト設定にフォールバックします: %v", err)
		// ✅ ドメインファクトリー使用：エラーの場合もデフォルト設定でフォールバック（最適化処理は継続）
		defaultConfig := model.NewDefaultUserConfig(userID)

		// ✅ ドメインロジック活用：最適化用のベース値ログ出力
		work, breakTime, rounds, sessionBreak := defaultConfig.GetOptimizationBaseValues()
		uc.logger.Infof("最適化用デフォルト設定: work=%d分, break=%d分, rounds=%d, sessionBreak=%d分",
			work, breakTime, rounds, sessionBreak)

		return defaultConfig, nil
	}

	// ✅ ドメインロジック活用：最適化用設定確認ログ
	work, breakTime, rounds, sessionBreak := config.GetOptimizationBaseValues()
	uc.logger.Infof("最適化用設定取得成功: work=%d分, break=%d分, rounds=%d, sessionBreak=%d分",
		work, breakTime, rounds, sessionBreak)

	return config, nil
}

// ✅ ドメインロジック活用：部分更新適用
func (uc *userConfigUseCase) applyPartialUpdate(config *model.UserConfig, req *model.UpdateUserConfigRequest) {
	updateCount := 0

	if req.RoundWorkTime != nil {
		oldValue := config.RoundWorkTime
		config.RoundWorkTime = *req.RoundWorkTime
		uc.logger.Debugf("作業時間更新: %d分 → %d分", oldValue, *req.RoundWorkTime)
		updateCount++
	}

	if req.RoundBreakTime != nil {
		oldValue := config.RoundBreakTime
		config.RoundBreakTime = *req.RoundBreakTime
		uc.logger.Debugf("休憩時間更新: %d分 → %d分", oldValue, *req.RoundBreakTime)
		updateCount++
	}

	if req.SessionRounds != nil {
		oldValue := config.SessionRounds
		config.SessionRounds = *req.SessionRounds
		uc.logger.Debugf("セッションラウンド数更新: %d回 → %d回", oldValue, *req.SessionRounds)
		updateCount++
	}

	if req.SessionBreakTime != nil {
		oldValue := config.SessionBreakTime
		config.SessionBreakTime = *req.SessionBreakTime
		uc.logger.Debugf("セッション長休憩時間更新: %d分 → %d分", oldValue, *req.SessionBreakTime)
		updateCount++
	}

	// ✅ ドメインロジック活用：UpdatedAt自動更新
	config.UpdatedAt = time.Now()

	uc.logger.Infof("設定項目更新数: %d項目", updateCount)
}
