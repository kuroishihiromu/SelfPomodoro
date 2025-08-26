package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/domain/service/optimization"
	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
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
	userRepo                    repository.UserRepository                       // User作成用
	optimizationPreferencesRepo repository.OptimizationPreferencesRepository // OptimizationPreferences作成用
	sampleDataService           *optimization.SampleOptimizationDataGenerationDomainService
	logger                      logger.Logger
}

// NewOnboardingUseCase は新しいOnboardingUseCaseを作成する（新エラーハンドリング対応版）
func NewOnboardingUseCase(
	userRepo repository.UserRepository,
	optimizationPreferencesRepo repository.OptimizationPreferencesRepository,
	sessionRepo repository.SessionRepository,
	logger logger.Logger,
) OnboardingUseCase {
	sampleDataService := optimization.NewSampleOptimizationDataGenerationDomainService(sessionRepo, logger)
	return &onboardingUseCase{
		userRepo:                    userRepo,
		optimizationPreferencesRepo: optimizationPreferencesRepo,
		sampleDataService:           sampleDataService,
		logger:                      logger,
	}
}

// CompletePostConfirmationSetup はCognito PostConfirmation後の完全な初期セットアップを実行する（新エラーハンドリング対応版）
func (uc *onboardingUseCase) CompletePostConfirmationSetup(ctx context.Context, params PostConfirmationParams) error {
	// UUIDをValue Objectに変換
	userIDVO, err := userVO.NewUserIDFromString(params.UserID.String())
	if err != nil {
		return appErrors.NewInternalError(err)
	}

	uc.logger.Infof("ユーザーオンボーディング開始: UserID=%s, Email=%s",
		params.UserID.String()[:8]+"...", params.Email)

	// 1. ✅ ドメインロジック活用：User作成
	err = uc.createUserWithDomainLogic(ctx, params)
	if err != nil {
		uc.logger.Errorf("User作成エラー: %v", err)
		return appErrors.NewUserCreationFailedError() // Domain Error
	}

	// 2. ✅ ドメインロジック活用：UserConfig作成
	err = uc.createUserConfigWithDomainLogic(ctx, userIDVO)
	if err != nil {
		// UserConfigは重要だが、失敗しても処理継続（デフォルト値フォールバック可能）
		uc.logger.Warnf("UserConfig作成エラー（処理継続）: %v", err)
	}

	// 3. サンプル最適化データ作成（Repository直接使用）
	err = uc.createSampleOptimizationData(ctx, userIDVO)
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

	// ✅ PostConfirmationParamsをValue Objectsに変換
	userIDVO, err := userVO.NewUserIDFromString(params.UserID.String())
	if err != nil {
		return appErrors.NewInternalError(err)
	}

	nameVO, err := userVO.NewUserName(params.Name)
	if err != nil {
		return appErrors.NewInvalidUserDataError()
	}

	emailVO, err := userVO.NewEmailAddress(params.Email)
	if err != nil {
		return appErrors.NewInvalidUserDataError()
	}

	providerVO, err := userVO.NewProvider(params.Provider)
	if err != nil {
		return appErrors.NewInvalidUserDataError()
	}

	// ✅ ドメインファクトリー使用：User作成
	user := entity.NewUser(entity.UserCreationParams{
		UserID:     userIDVO,
		Name:       nameVO,
		Email:      emailVO,
		Provider:   providerVO,
		ProviderID: nil, // PostConfirmation時はnilでOK
	})

	// ✅ ドメインロジック活用：作成前バリデーション
	if !user.IsValidForCreation() {
		uc.logger.Error("User作成バリデーションエラー: 必須項目が不足しています")
		return appErrors.NewInvalidUserDataError()
	}

	// UserRepositoryに直接Create（曖昧なGetOrCreateは使用しない）
	err = uc.userRepo.Create(ctx, user)
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
		user.Name.Value(), user.Email.Value(), providerDisplay)

	return nil
}

// ✅ ドメインロジック活用：UserConfig作成（デフォルト値使用・新エラーハンドリング対応版）
func (uc *onboardingUseCase) createUserConfigWithDomainLogic(ctx context.Context, userID userVO.UserID) error {
	if uc.optimizationPreferencesRepo == nil {
		uc.logger.Warn("OptimizationPreferencesRepository が利用できません")
		return nil
	}

	// OptimizationPreferencesが既に存在するかチェック（冪等性の確保）
	existingConfig, err := uc.optimizationPreferencesRepo.Get(ctx, userID)
	if err == nil && existingConfig != nil {
		uc.logger.Infof("UserConfigは既に存在します: %s", userID.String()[:8]+"...")
		return nil
	}

	// ✅ ドメインファクトリー使用：デフォルトのOptimizationPreferencesを作成
	optimizationPreferences, err := entity.NewOptimizationPreferences(userID)
	if err != nil {
		uc.logger.Errorf("OptimizationPreferences作成エラー: %v", err)
		return appErrors.NewInvalidUserConfigError()
	}

	// ✅ ドメインロジック活用：作成前バリデーション
	if !optimizationPreferences.IsValid() {
		uc.logger.Error("UserConfig作成バリデーションエラー: 無効な設定値")
		return appErrors.NewInvalidUserConfigError()
	}

	err = uc.optimizationPreferencesRepo.Create(ctx, optimizationPreferences)
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
		optimizationPreferences.GetWorkTimeOrDefault().Minutes(),
		optimizationPreferences.GetBreakTimeOrDefault().Minutes(),
		optimizationPreferences.GetSessionRoundsOrDefault().Count(),
		optimizationPreferences.SessionBreakTime.Minutes())

	return nil
}

// createSampleOptimizationData はサンプル最適化データを作成する（ドメインサービス使用版）
func (uc *onboardingUseCase) createSampleOptimizationData(ctx context.Context, userID userVO.UserID) error {
	uc.logger.Infof("サンプル最適化データ作成開始: UserID=%s", userID.String()[:8]+"...")

	if uc.sampleDataService == nil {
		uc.logger.Warn("SampleDataService が利用できません")
		return nil
	}

	err := uc.sampleDataService.GenerateAndStoreSampleData(ctx, userID)
	if err != nil {
		uc.logger.Errorf("サンプルデータ生成エラー: %v", err)
		return err
	}

	uc.logger.Infof("サンプル最適化データ作成完了: UserID=%s", userID.String()[:8]+"...")
	return nil
}
