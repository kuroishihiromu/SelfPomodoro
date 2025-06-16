package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
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

// OnboardingUseCase はユーザーオンボーディングに関するユースケース（新エラーハンドリング対応版）
type OnboardingUseCase interface {
	// CompletePostConfirmationSetup はCognito PostConfirmation後の初期セットアップを完了する
	CompletePostConfirmationSetup(ctx context.Context, params PostConfirmationParams) error
}

// onboardingUseCase はOnboardingUseCaseの実装（新エラーハンドリング対応版）
type onboardingUseCase struct {
	userRepo       repository.UserRepository                   // User作成用
	userConfigRepo repository.UserConfigRepository             // UserConfig作成用
	sampleDataRepo repository.SampleOptimizationDataRepository // サンプルデータ作成用（直接使用）
	logger         logger.Logger
}

// NewOnboardingUseCase は新しいOnboardingUseCaseを作成する（新エラーハンドリング対応版）
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

// CompletePostConfirmationSetup はCognito PostConfirmation後の完全な初期セットアップを実行する（新エラーハンドリング対応版）
func (uc *onboardingUseCase) CompletePostConfirmationSetup(ctx context.Context, params PostConfirmationParams) error {
	uc.logger.Infof("ユーザーオンボーディング開始: UserID=%s, Email=%s",
		params.UserID.String()[:8]+"...", params.Email)

	// 1. ✅ ドメインロジック活用：User作成
	err := uc.createUserWithDomainLogic(ctx, params)
	if err != nil {
		uc.logger.Errorf("User作成エラー: %v", err)
		return appErrors.NewUserCreationFailedError() // Domain Error
	}

	// 2. ✅ ドメインロジック活用：UserConfig作成
	err = uc.createUserConfigWithDomainLogic(ctx, params.UserID)
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

// ✅ ドメインロジック活用：User作成（PostConfirmation専用・新規作成のみ・新エラーハンドリング対応版）
func (uc *onboardingUseCase) createUserWithDomainLogic(ctx context.Context, params PostConfirmationParams) error {
	if uc.userRepo == nil {
		return appErrors.NewInternalError(errors.New("UserRepositoryが初期化されていません"))
	}

	// ✅ PostConfirmationParamsをCognitoUserParamsに変換
	cognitoParams := model.CognitoUserParams{
		UserID:     params.UserID,
		Email:      params.Email,
		Name:       params.Name,
		GivenName:  params.GivenName,
		FamilyName: params.FamilyName,
		// 現在のPostConfirmationParamsには含まれていないが、将来拡張可能
		Picture:          "",
		Locale:           "",
		IdentityProvider: "",
		EmailVerified:    true, // PostConfirmation済みなので確認済み
	}

	// ✅ ドメインファクトリー使用：Cognito属性からUser作成
	user := model.NewUserFromCognitoAttributes(cognitoParams)

	// ✅ ドメインロジック活用：作成前バリデーション
	if !user.IsValidForCreation() {
		uc.logger.Error("User作成バリデーションエラー: 必須項目が不足しています")
		return appErrors.NewInvalidUserDataError()
	}

	// UserRepositoryに直接Create（曖昧なGetOrCreateは使用しない）
	err := uc.userRepo.Create(ctx, user)
	if err != nil {
		uc.logger.Errorf("User作成失敗: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrUniqueConstraint) {
			// 一意制約違反（既に存在する）は成功とみなす（冪等性）
			uc.logger.Infof("User既存（冪等性確保）: %s", user.Email)
			return nil
		}
		if appErrors.IsDatabaseError(err) {
			return appErrors.NewInternalError(err)
		}

		return appErrors.NewInternalError(err)
	}

	// ✅ ドメインロジック活用：プロバイダー別ログ出力
	providerDisplay := user.GetProviderDisplayName()
	uc.logger.Infof("User作成成功: %s (%s) - プロバイダー: %s",
		user.Name, user.Email, providerDisplay)

	return nil
}

// ✅ ドメインロジック活用：UserConfig作成（デフォルト値使用・新エラーハンドリング対応版）
func (uc *onboardingUseCase) createUserConfigWithDomainLogic(ctx context.Context, userID uuid.UUID) error {
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

	// ✅ ドメインファクトリー使用：デフォルトのUserConfigを作成
	userConfig := model.NewDefaultUserConfig(userID)

	// ✅ ドメインロジック活用：作成前バリデーション
	if !userConfig.IsValid() {
		uc.logger.Error("UserConfig作成バリデーションエラー: 無効な設定値")
		return appErrors.NewInvalidUserConfigError()
	}

	err = uc.userConfigRepo.CreateUserConfig(ctx, userConfig)
	if err != nil {
		uc.logger.Errorf("UserConfig作成失敗: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrDynamoDBCondition) {
			// 条件チェック失敗（既に存在する）は成功とみなす（冪等性）
			uc.logger.Infof("UserConfig既存（冪等性確保）: %s", userID.String()[:8]+"...")
			return nil
		}
		if appErrors.IsDynamoDBError(err) {
			return appErrors.NewUserConfigCreateFailedError()
		}

		return appErrors.NewInternalError(err)
	}

	// ✅ ドメインロジック活用：設定値ログ出力
	uc.logger.Infof("UserConfig作成成功: work=%d分, break=%d分, rounds=%d, sessionBreak=%d分",
		userConfig.GetWorkTimeOrDefault(),
		userConfig.GetBreakTimeOrDefault(),
		userConfig.GetSessionRoundsOrDefault(),
		userConfig.GetSessionBreakTimeOrDefault())

	return nil
}

// createSampleOptimizationData はサンプル最適化データを作成する（Repository直接使用・新エラーハンドリング対応版）
func (uc *onboardingUseCase) createSampleOptimizationData(ctx context.Context, userID uuid.UUID) error {
	if uc.sampleDataRepo == nil {
		uc.logger.Warn("SampleOptimizationDataRepository が利用できません")
		return nil // 失敗してもオンボーディングは継続
	}

	uc.logger.Infof("サンプル最適化データ作成開始: UserID=%s", userID.String()[:8]+"...")

	// Repository直接呼び出し（UseCase層を介さない）
	err := uc.sampleDataRepo.CreateSampleOptimizationData(ctx, userID)
	if err != nil {
		uc.logger.Errorf("サンプル最適化データ作成失敗: %v", err)

		// Infrastructure Error → Domain Error 変換
		if appErrors.IsDynamoDBError(err) {
			return appErrors.NewInternalError(err)
		}
		if appErrors.IsInfrastructureError(err) {
			return appErrors.NewInternalError(err)
		}

		return appErrors.NewInternalError(err)
	}

	uc.logger.Infof("サンプル最適化データ作成完了: UserID=%s", userID.String()[:8]+"...")
	return nil
}
