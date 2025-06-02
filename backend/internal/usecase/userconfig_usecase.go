package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// UserConfigUseCase はユーザー設定に関するユースケースを定義するインターフェース
type UserConfigUseCase interface {
	// GetUserConfig はユーザー設定を取得する（存在しない場合はデフォルト値で作成）
	GetUserConfig(ctx context.Context, userID uuid.UUID) (*model.UserConfigResponse, error)

	// CreateUserConfig は新しいユーザー設定を作成する
	CreateUserConfig(ctx context.Context, userID uuid.UUID, req *model.CreateUserConfigRequest) (*model.UserConfigResponse, error)

	// UpdateUserConfig はユーザー設定を更新する
	UpdateUserConfig(ctx context.Context, userID uuid.UUID, req *model.UpdateUserConfigRequest) (*model.UserConfigResponse, error)

	// DeleteUserConfig はユーザー設定を削除する
	DeleteUserConfig(ctx context.Context, userID uuid.UUID) error

	// GetUserConfigForOptimization は最適化処理用にユーザー設定を取得する（内部使用）
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

// GetUserConfig はユーザー設定を取得する
func (uc *userConfigUseCase) GetUserConfig(ctx context.Context, userID uuid.UUID) (*model.UserConfigResponse, error) {
	if uc.userConfigRepo == nil {
		uc.logger.Warn("UserConfigRepository が初期化されていません。デフォルト設定を返します")
		// DynamoDBが利用できない場合のフォールバック
		defaultConfig := model.NewUserConfig(userID)
		return defaultConfig.ToResponse(), nil
	}

	config, err := uc.userConfigRepo.GetOrCreateUserConfig(ctx, userID)
	if err != nil {
		uc.logger.Errorf("ユーザー設定取得エラー: %v", err)
		return nil, err
	}

	return config.ToResponse(), nil
}

// CreateUserConfig は新しいユーザー設定を作成する
func (uc *userConfigUseCase) CreateUserConfig(ctx context.Context, userID uuid.UUID, req *model.CreateUserConfigRequest) (*model.UserConfigResponse, error) {
	if uc.userConfigRepo == nil {
		return nil, fmt.Errorf("ユーザー設定機能は現在利用できません")
	}

	// 新しい設定を作成
	config := model.NewUserConfig(userID)
	config.UpdateSettings(req.RoundWorkTime, req.RoundBreakTime, req.SessionRounds, req.SessionBreakTime)

	// 設定の有効性をチェック
	if !config.IsValid() {
		return nil, fmt.Errorf("設定値が無効です")
	}

	// リポジトリに保存
	if err := uc.userConfigRepo.CreateUserConfig(ctx, config); err != nil {
		uc.logger.Errorf("ユーザー設定作成エラー: %v", err)
		return nil, err
	}

	uc.logger.Infof("ユーザー設定作成成功: %s", userID.String())
	return config.ToResponse(), nil
}

// UpdateUserConfig はユーザー設定を更新する
func (uc *userConfigUseCase) UpdateUserConfig(ctx context.Context, userID uuid.UUID, req *model.UpdateUserConfigRequest) (*model.UserConfigResponse, error) {
	if uc.userConfigRepo == nil {
		return nil, fmt.Errorf("ユーザー設定機能は現在利用できません")
	}

	// 既存の設定を取得
	config, err := uc.userConfigRepo.GetUserConfig(ctx, userID)
	if err != nil {
		uc.logger.Errorf("ユーザー設定取得エラー: %v", err)
		return nil, err
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
		return nil, fmt.Errorf("更新後の設定値が無効です")
	}

	// リポジトリに保存
	if err := uc.userConfigRepo.UpdateUserConfig(ctx, config); err != nil {
		uc.logger.Errorf("ユーザー設定更新エラー: %v", err)
		return nil, err
	}

	uc.logger.Infof("ユーザー設定更新成功: %s", userID.String())
	return config.ToResponse(), nil
}

// DeleteUserConfig はユーザー設定を削除する
func (uc *userConfigUseCase) DeleteUserConfig(ctx context.Context, userID uuid.UUID) error {
	if uc.userConfigRepo == nil {
		return fmt.Errorf("ユーザー設定機能は現在利用できません")
	}

	if err := uc.userConfigRepo.DeleteUserConfig(ctx, userID); err != nil {
		uc.logger.Errorf("ユーザー設定削除エラー: %v", err)
		return err
	}

	uc.logger.Infof("ユーザー設定削除成功: %s", userID.String())
	return nil
}

// GetUserConfigForOptimization は最適化処理用にユーザー設定を取得する
func (uc *userConfigUseCase) GetUserConfigForOptimization(ctx context.Context, userID uuid.UUID) (*model.UserConfig, error) {
	if uc.userConfigRepo == nil {
		uc.logger.Warn("UserConfigRepository が初期化されていません。デフォルト設定を返します")
		return model.NewUserConfig(userID), nil
	}

	config, err := uc.userConfigRepo.GetOrCreateUserConfig(ctx, userID)
	if err != nil {
		uc.logger.Errorf("最適化用ユーザー設定取得エラー: %v", err)
		// エラーの場合もデフォルト設定でフォールバック
		uc.logger.Warn("デフォルト設定にフォールバックします")
		return model.NewUserConfig(userID), nil
	}

	return config, nil
}
