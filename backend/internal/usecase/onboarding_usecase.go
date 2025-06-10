package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// PostConfirmationParams はPostConfirmation時のパラメータ
type PostConfirmationParams struct {
	UserID     uuid.UUID
	Email      string
	Name       string
	GivenName  string
	FamilyName string
	Provider   string // "Cognito_UserPool", "Google"
}

// OnboardingUseCase はユーザーオンボーディングに関するユースケース（Clean Architecture準拠）
type OnboardingUseCase interface {
	// CompletePostConfirmationSetup はCognito PostConfirmation後の初期セットアップを完了する
	CompletePostConfirmationSetup(ctx context.Context, params PostConfirmationParams) error
}

// onboardingUseCase はOnboardingUseCaseの実装（Repository直接依存）
type onboardingUseCase struct {
	userRepo       repository.UserRepository                   // User作成用
	userConfigRepo repository.UserConfigRepository             // UserConfig作成用
	sampleDataRepo repository.SampleOptimizationDataRepository // サンプルデータ作成用（直接使用）
	logger         logger.Logger
}

// NewOnboardingUseCase は新しいOnboardingUseCaseを作成する（Repository直接依存版）
func NewOnboardingUseCase(
	userRepo repository.UserRepository,
	userConfigRepo repository.UserConfigRepository,
	sampleDataRepo repository.SampleOptimizationDataRepository,
	logger logger.Logger,
) OnboardingUseCase {
	return &onboardingUseCase{
		userRepo:       userRepo,
		userConfigRepo: userConfigRepo,
		sampleDataRepo: sampleDataRepo,
		logger:         logger,
	}
}

// CompletePostConfirmationSetup はCognito PostConfirmation後の完全な初期セットアップを実行する
func (uc *onboardingUseCase) CompletePostConfirmationSetup(ctx context.Context, params PostConfirmationParams) error {
	uc.logger.Infof("ユーザーオンボーディング開始: UserID=%s, Email=%s",
		params.UserID.String()[:8]+"...", params.Email)

	// 1. User作成（Repository直接使用）
	err := uc.createUser(ctx, params)
	if err != nil {
		return fmt.Errorf("User作成エラー: %w", err)
	}

	// 2. UserConfig作成（Repository直接使用）
	err = uc.createUserConfig(ctx, params.UserID)
	if err != nil {
		// UserConfigは重要だが、失敗しても処理継続（デフォルト値フォールバック可能）
		uc.logger.Warnf("UserConfig作成エラー（処理継続）: %v", err)
	}

	// 3. サンプル最適化データ作成（Repository直接使用）
	err = uc.createSampleOptimizationData(ctx, params.UserID)
	if err != nil {
		// サンプルデータは重要ではないので、失敗しても処理継続
		uc.logger.Warnf("サンプル最適化データ作成エラー（処理継続）: %v", err)
	}

	uc.logger.Infof("ユーザーオンボーディング完了: UserID=%s", params.UserID.String()[:8]+"...")
	return nil
}

// createUser はUserを作成する（PostConfirmation専用・新規作成のみ）
func (uc *onboardingUseCase) createUser(ctx context.Context, params PostConfirmationParams) error {
	if uc.userRepo == nil {
		return fmt.Errorf("UserRepositoryが初期化されていません")
	}

	// PostConfirmation時は必ず新規ユーザーなので、Create専用処理
	user := uc.createUserFromParams(params)

	// UserRepositoryに直接Create（曖昧なGetOrCreateは使用しない）
	err := uc.userRepo.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("User作成失敗: %w", err)
	}

	uc.logger.Infof("User作成成功: %s (%s)", params.Name, params.Email)
	return nil
}

// createUserFromParams はPostConfirmationParamsからUserドメインモデルを作成する
func (uc *onboardingUseCase) createUserFromParams(params PostConfirmationParams) *model.User {
	// 表示名の決定（優先順位: name > given_name + family_name > email）
	displayName := params.Name
	if displayName == "" && params.GivenName != "" {
		displayName = params.GivenName
		if params.FamilyName != "" {
			displayName += " " + params.FamilyName
		}
	}
	if displayName == "" {
		displayName = params.Email
	}

	return model.NewUser(model.UserCreationParams{
		UserID:       params.UserID,
		Name:         displayName,
		Email:        params.Email,
		Provider:     params.Provider,
		ProviderID:   nil,
		IsGoogleUser: false,
	})
}

// createUserConfig はデフォルトのUserConfigを作成する（Repository直接使用）
func (uc *onboardingUseCase) createUserConfig(ctx context.Context, userID uuid.UUID) error {
	if uc.userConfigRepo == nil {
		uc.logger.Warn("UserConfigRepository が利用できません")
		return nil
	}

	// UserConfigが既に存在するかチェック（冪等性の確保）
	existingConfig, err := uc.userConfigRepo.GetUserConfig(ctx, userID)
	if err == nil && existingConfig != nil {
		uc.logger.Infof("UserConfigは既に存在します: %s", userID.String()[:8]+"...")
		return nil
	}

	// デフォルトのUserConfigを作成
	userConfig := model.NewUserConfig(userID)
	err = uc.userConfigRepo.CreateUserConfig(ctx, userConfig)
	if err != nil {
		return fmt.Errorf("UserConfig作成失敗: %w", err)
	}

	uc.logger.Infof("UserConfig作成成功: work=%d分, break=%d分, rounds=%d",
		userConfig.RoundWorkTime, userConfig.RoundBreakTime, userConfig.SessionRounds)
	return nil
}

// createSampleOptimizationData はサンプル最適化データを作成する（Repository直接使用）
func (uc *onboardingUseCase) createSampleOptimizationData(ctx context.Context, userID uuid.UUID) error {
	if uc.sampleDataRepo == nil {
		uc.logger.Warn("SampleOptimizationDataRepository が利用できません")
		return nil // 失敗してもオンボーディングは継続
	}

	uc.logger.Infof("サンプル最適化データ作成開始: UserID=%s", userID.String()[:8]+"...")

	// Repository直接呼び出し（UseCase層を介さない）
	err := uc.sampleDataRepo.CreateSampleOptimizationData(ctx, userID)
	if err != nil {
		return fmt.Errorf("サンプル最適化データ作成失敗: %w", err)
	}

	uc.logger.Infof("サンプル最適化データ作成完了: UserID=%s", userID.String()[:8]+"...")
	return nil
}
