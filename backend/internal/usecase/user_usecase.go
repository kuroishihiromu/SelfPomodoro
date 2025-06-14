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

// UserUseCase はユーザーに関するユースケースを定義するインターフェース（ドメイン強化版）
type UserUseCase interface {
	// GetUserProfile はユーザープロフィールを取得する
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*model.UserResponse, error)

	// UpdateUserProfile はユーザープロフィール（名前・メール）を更新する
	UpdateUserProfile(ctx context.Context, userID uuid.UUID, req *model.UpdateUserRequest) (*model.UserResponse, error)

	// GetUserByEmail はメールアドレスでユーザーを検索する（管理・デバッグ用）
	GetUserByEmail(ctx context.Context, email string) (*model.UserResponse, error)

	// CheckUserExists はユーザーの存在確認を行う（軽量版・認証用）
	CheckUserExists(ctx context.Context, userID uuid.UUID) (bool, error)

	// DeleteUser はユーザーを削除する（GDPR対応等）
	DeleteUser(ctx context.Context, userID uuid.UUID) error

	// GetUsersByProvider はプロバイダー別ユーザー一覧を取得する（管理用）
	GetUsersByProvider(ctx context.Context, provider string, limit, offset int) ([]*model.UserResponse, error)
}

// userUseCase はUserUseCaseインターフェースの実装（ドメイン強化版）
type userUseCase struct {
	userRepo repository.UserRepository
	logger   logger.Logger
}

// NewUserUseCase は新しいUserUseCaseインスタンスを作成する
func NewUserUseCase(userRepo repository.UserRepository, logger logger.Logger) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
		logger:   logger,
	}
}

// GetUserProfile はユーザープロフィールを取得する（ドメイン強化版）
func (uc *userUseCase) GetUserProfile(ctx context.Context, userID uuid.UUID) (*model.UserResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		uc.logger.Errorf("ユーザープロフィール取得エラー: %v", err)
		if errors.Is(err, appErrors.ErrUserNotFound) {
			return nil, appErrors.NewUserNotFoundError()
		}
		return nil, appErrors.NewInternalError(err)
	}

	// ✅ ドメインロジック活用：プロバイダー情報ログ出力
	providerDisplay := user.GetProviderDisplayName()
	uc.logger.Infof("ユーザープロフィール取得成功: %s (%s) - プロバイダー: %s",
		user.Name, user.Email, providerDisplay)

	return user.ToResponse(), nil
}

// UpdateUserProfile はユーザープロフィールを更新する（ドメイン強化版）
func (uc *userUseCase) UpdateUserProfile(ctx context.Context, userID uuid.UUID, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	// 既存ユーザーの取得
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		uc.logger.Errorf("ユーザー取得エラー: %v", err)
		if errors.Is(err, appErrors.ErrUserNotFound) {
			return nil, appErrors.NewUserNotFoundError()
		}
		return nil, appErrors.NewInternalError(err)
	}

	// ✅ ドメインロジック活用：メール変更権限チェック
	if req.Email != nil && !user.CanChangeEmail() {
		uc.logger.Errorf("Googleユーザーはメールアドレスを変更できません: UserID=%s", userID.String()[:8]+"...")
		return nil, appErrors.NewForbiddenError("Googleユーザーはメールアドレスを変更できません")
	}

	// 更新前の状態をログ出力
	oldName := user.Name
	oldEmail := user.Email

	// ✅ ドメインロジック活用：プロフィール更新処理
	var newName, newEmail string
	if req.Name != nil {
		newName = *req.Name
	} else {
		newName = user.Name
	}

	if req.Email != nil {
		newEmail = *req.Email
	} else {
		newEmail = user.Email
	}

	user.UpdateProfile(newName, newEmail)

	// ✅ ドメインロジック活用：更新後バリデーション
	if !user.ValidateName() {
		return nil, appErrors.NewValidationError("有効な名前を入力してください")
	}

	if !user.ValidateEmail() {
		return nil, appErrors.NewValidationError("有効なメールアドレスを入力してください")
	}

	// データベース更新
	err = uc.userRepo.Update(ctx, user)
	if err != nil {
		uc.logger.Errorf("ユーザープロフィール更新エラー: %v", err)
		if errors.Is(err, appErrors.ErrUserEmailDuplicate) {
			return nil, appErrors.NewEmailAlreadyExistsError()
		}
		if errors.Is(err, appErrors.ErrUserNotFound) {
			return nil, appErrors.NewUserNotFoundError()
		}
		return nil, appErrors.NewInternalError(err)
	}

	// 更新内容をログ出力
	updateCount := 0
	if oldName != newName {
		uc.logger.Infof("名前更新: %s → %s", oldName, newName)
		updateCount++
	}
	if oldEmail != newEmail {
		uc.logger.Infof("メールアドレス更新: %s → %s", oldEmail, newEmail)
		updateCount++
	}

	uc.logger.Infof("ユーザープロフィール更新成功: UserID=%s, 更新項目数=%d",
		userID.String()[:8]+"...", updateCount)

	return user.ToResponse(), nil
}

// GetUserByEmail はメールアドレスでユーザーを検索する（ドメイン強化版）
func (uc *userUseCase) GetUserByEmail(ctx context.Context, email string) (*model.UserResponse, error) {
	// 基本的なメールアドレス形式チェック
	if email == "" || !uc.isValidEmailFormat(email) {
		return nil, appErrors.NewValidationError("有効なメールアドレスを入力してください")
	}

	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		uc.logger.Errorf("メールアドレスによるユーザー検索エラー: %v", err)
		if errors.Is(err, appErrors.ErrUserNotFound) {
			return nil, appErrors.NewUserNotFoundError()
		}
		return nil, appErrors.NewInternalError(err)
	}

	// ✅ ドメインロジック活用：プロバイダー情報ログ出力
	providerDisplay := user.GetProviderDisplayName()
	uc.logger.Infof("メールアドレス検索成功: %s (%s) - プロバイダー: %s",
		user.Name, user.Email, providerDisplay)

	return user.ToResponse(), nil
}

// CheckUserExists はユーザーの存在確認を行う（ドメイン強化版・軽量版）
func (uc *userUseCase) CheckUserExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	exists, err := uc.userRepo.ExistsByID(ctx, userID)
	if err != nil {
		uc.logger.Errorf("ユーザー存在確認エラー: %v", err)
		return false, appErrors.NewInternalError(err)
	}

	if exists {
		uc.logger.Debugf("ユーザー存在確認: UserID=%s - 存在", userID.String()[:8]+"...")
	} else {
		uc.logger.Debugf("ユーザー存在確認: UserID=%s - 存在しない", userID.String()[:8]+"...")
	}

	return exists, nil
}

// DeleteUser はユーザーを削除する（ドメイン強化版・GDPR対応）
func (uc *userUseCase) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	// 削除前にユーザー情報を取得（ログ出力用）
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		uc.logger.Errorf("削除対象ユーザー取得エラー: %v", err)
		if errors.Is(err, appErrors.ErrUserNotFound) {
			return appErrors.NewUserNotFoundError()
		}
		return appErrors.NewInternalError(err)
	}

	// ✅ ドメインロジック活用：削除前情報ログ
	providerDisplay := user.GetProviderDisplayName()
	uc.logger.Infof("ユーザー削除開始: %s (%s) - プロバイダー: %s",
		user.Name, user.Email, providerDisplay)

	// ユーザー削除実行
	err = uc.userRepo.Delete(ctx, userID)
	if err != nil {
		uc.logger.Errorf("ユーザー削除エラー: %v", err)
		if errors.Is(err, appErrors.ErrUserNotFound) {
			return appErrors.NewUserNotFoundError()
		}
		return appErrors.NewInternalError(err)
	}

	uc.logger.Infof("ユーザー削除成功: UserID=%s (%s)",
		userID.String()[:8]+"...", providerDisplay)

	return nil
}

// GetUsersByProvider はプロバイダー別ユーザー一覧を取得する（ドメイン強化版・管理用）
func (uc *userUseCase) GetUsersByProvider(ctx context.Context, provider string, limit, offset int) ([]*model.UserResponse, error) {
	// プロバイダー名の検証
	if provider == "" {
		return nil, appErrors.NewValidationError("プロバイダー名を指定してください")
	}

	// 有効なプロバイダー名かチェック
	if !uc.isValidProviderName(provider) {
		return nil, appErrors.NewValidationError("無効なプロバイダー名です")
	}

	// ページネーションパラメータの検証
	if limit <= 0 || limit > 100 {
		limit = 20 // デフォルト値
	}
	if offset < 0 {
		offset = 0
	}

	users, err := uc.userRepo.GetUsersByProvider(ctx, provider, limit, offset)
	if err != nil {
		uc.logger.Errorf("プロバイダー別ユーザー一覧取得エラー: %v", err)
		return nil, appErrors.NewInternalError(err)
	}

	// レスポンス形式に変換
	userResponses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	uc.logger.Infof("プロバイダー別ユーザー一覧取得成功: provider=%s, count=%d, limit=%d, offset=%d",
		provider, len(userResponses), limit, offset)

	return userResponses, nil
}

// ✅ ドメインロジック活用：ヘルパーメソッド群

// isValidEmailFormat はメールアドレス形式の基本チェックを行う
func (uc *userUseCase) isValidEmailFormat(email string) bool {
	// 簡易チェック（より詳細なチェックが必要な場合はvalidatorライブラリを使用）
	return len(email) > 3 &&
		email != "" &&
		email[0] != '@' &&
		email[len(email)-1] != '@' &&
		uc.containsAtSymbol(email) &&
		!uc.containsMultipleAtSymbols(email)
}

// containsAtSymbol は@マークが含まれているかチェック
func (uc *userUseCase) containsAtSymbol(email string) bool {
	for _, char := range email {
		if char == '@' {
			return true
		}
	}
	return false
}

// containsMultipleAtSymbols は複数の@マークが含まれているかチェック
func (uc *userUseCase) containsMultipleAtSymbols(email string) bool {
	count := 0
	for _, char := range email {
		if char == '@' {
			count++
		}
	}
	return count > 1
}

// isValidProviderName は有効なプロバイダー名かチェック
func (uc *userUseCase) isValidProviderName(provider string) bool {
	validProviders := []string{
		"Cognito_UserPool",
		"Google",
	}

	for _, validProvider := range validProviders {
		if provider == validProvider {
			return true
		}
	}
	return false
}
