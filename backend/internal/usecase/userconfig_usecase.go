package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	domainErrors "github.com/tsunakit99/selfpomodoro/internal/domain/errors"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	dynamodberrors "github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/dynamodb/errors"
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

// userConfigUseCase はUserConfigUseCaseインターフェースの実装
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

// GetUserConfig はユーザー設定を取得する（デフォルト値フォールバック付き）
func (uc *userConfigUseCase) GetUserConfig(ctx context.Context, userID uuid.UUID) (*model.UserConfigResponse, error) {
	if uc.userConfigRepo == nil {
		uc.logger.Warn("UserConfigRepository が初期化されていません。デフォルト設定を返します")
		// DynamoDBが利用できない場合のフォールバック
		defaultConfig := model.NewUserConfig(userID)
		return defaultConfig.ToResponse(), nil
	}

	config, err := uc.userConfigRepo.GetUserConfig(ctx, userID)
	if err != nil {
		uc.logger.Errorf("ユーザー設定取得エラー、デフォルト設定を返します: %v", err)
		// 設定が存在しない場合のデフォルト値フォールバック
		defaultConfig := model.NewUserConfig(userID)
		return defaultConfig.ToResponse(), nil
	}

	return config.ToResponse(), nil
}

// CreateUserConfig は新しいユーザー設定を作成する（PostConfirmation専用）
func (uc *userConfigUseCase) CreateUserConfig(ctx context.Context, userID uuid.UUID, req *model.CreateUserConfigRequest) (*model.UserConfigResponse, error) {
	if uc.userConfigRepo == nil {
		return nil, domainErrors.NewInternalError(fmt.Errorf("ユーザー設定機能は現在利用できません"))
	}

	// 新しい設定を作成
	config := model.NewUserConfig(userID)
	config.UpdateSettings(req.RoundWorkTime, req.RoundBreakTime, req.SessionRounds, req.SessionBreakTime)

	// 設定の有効性をチェック
	if !config.IsValid() {
		return nil, domainErrors.NewInvalidUserConfigError()
	}

	// リポジトリに保存
	if err := uc.userConfigRepo.CreateUserConfig(ctx, config); err != nil {
		uc.logger.Errorf("ユーザー設定作成エラー: %v", err)
		if domainErrors.Is(err, dynamodberrors.ErrUserConfigCreateFailed) {
			return nil, domainErrors.NewUserConfigCreateFailedError()
		}
		return nil, domainErrors.NewInternalError(err)
	}

	uc.logger.Infof("ユーザー設定作成成功: %s", userID.String())
	return config.ToResponse(), nil
}

// UpdateUserConfig はユーザー設定を更新する（存在前提）
func (uc *userConfigUseCase) UpdateUserConfig(ctx context.Context, userID uuid.UUID, req *model.UpdateUserConfigRequest) (*model.UserConfigResponse, error) {
	if uc.userConfigRepo == nil {
		return nil, domainErrors.NewInternalError(fmt.Errorf("ユーザー設定機能は現在利用できません"))
	}

	// 既存の設定を取得（存在前提）
	config, err := uc.userConfigRepo.GetUserConfig(ctx, userID)
	if err != nil {
		uc.logger.Errorf("ユーザー設定取得エラー: %v", err)
		if domainErrors.Is(err, dynamodberrors.ErrUserConfigNotFound) {
			return nil, domainErrors.NewUserConfigNotFoundError()
		}
		return nil, domainErrors.NewInternalError(err)
	}

	// 設定を更新（部分更新対応）
	if req.RoundWorkTime != nil {
		config.RoundWorkTime = *req.RoundWorkTime
	}
	if req.RoundBreakTime != nil {
		config.RoundBreakTime = *req.RoundBreakTime
	}
	if req.SessionRounds != nil {
		config.SessionRounds = *req.SessionRounds
	}
	if req.SessionBreakTime != nil {
		config.SessionBreakTime = *req.SessionBreakTime
	}

	// 設定の有効性をチェック
	if !config.IsValid() {
		return nil, domainErrors.NewInvalidUserConfigError()
	}

	// リポジトリに保存
	if err := uc.userConfigRepo.UpdateUserConfig(ctx, config); err != nil {
		uc.logger.Errorf("ユーザー設定更新エラー: %v", err)
		if domainErrors.Is(err, dynamodberrors.ErrUserConfigUpdateFailed) {
			return nil, domainErrors.NewUserConfigUpdateFailedError()
		}
		return nil, domainErrors.NewInternalError(err)
	}

	uc.logger.Infof("ユーザー設定更新成功: %s", userID.String())
	return config.ToResponse(), nil
}

// DeleteUserConfig はユーザー設定を削除する
func (uc *userConfigUseCase) DeleteUserConfig(ctx context.Context, userID uuid.UUID) error {
	if uc.userConfigRepo == nil {
		return domainErrors.NewInternalError(fmt.Errorf("ユーザー設定機能は現在利用できません"))
	}

	if err := uc.userConfigRepo.DeleteUserConfig(ctx, userID); err != nil {
		uc.logger.Errorf("ユーザー設定削除エラー: %v", err)
		if domainErrors.Is(err, dynamodberrors.ErrUserConfigNotFound) {
			return domainErrors.NewUserConfigNotFoundError()
		}
		if domainErrors.Is(err, dynamodberrors.ErrUserConfigDeleteFailed) {
			return domainErrors.NewUserConfigDeleteFailedError()
		}
		return domainErrors.NewInternalError(err)
	}

	uc.logger.Infof("ユーザー設定削除成功: %s", userID.String())
	return nil
}

// GetUserConfigForOptimization は最適化処理用にユーザー設定を取得する（デフォルト値フォールバック）
func (uc *userConfigUseCase) GetUserConfigForOptimization(ctx context.Context, userID uuid.UUID) (*model.UserConfig, error) {
	if uc.userConfigRepo == nil {
		uc.logger.Warn("UserConfigRepository が初期化されていません。デフォルト設定を返します")
		return model.NewUserConfig(userID), nil
	}

	config, err := uc.userConfigRepo.GetUserConfig(ctx, userID)
	if err != nil {
		uc.logger.Warnf("最適化用ユーザー設定取得エラー、デフォルト設定にフォールバックします: %v", err)
		// エラーの場合もデフォルト設定でフォールバック（最適化処理は継続）
		return model.NewUserConfig(userID), nil
	}

	return config, nil
}
