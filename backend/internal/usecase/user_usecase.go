package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/auth"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// UserUseCase はユーザーに関するユースケースを定義するインターフェース
type UserUseCase interface {
	// GetOrCreateUserFromCognito はCognito認証後にユーザーを取得または作成する（最重要）
	GetOrCreateUserFromCognito(ctx context.Context, claims *auth.CognitoClaims) (*model.User, error)

	// GetUserProfile はユーザープロフィールを取得する
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*model.UserResponse, error)

	// UpdateUserProfile はユーザープロフィールを更新する
	UpdateUserProfile(ctx context.Context, userID uuid.UUID, req *model.UpdateUserRequest) (*model.UserResponse, error)

	// GetUserByEmail はメールアドレスでユーザーを取得する
	GetUserByEmail(ctx context.Context, email string) (*model.UserResponse, error)

	// DeleteUser はユーザーを削除する（GDPR対応）
	DeleteUser(ctx context.Context, userID uuid.UUID) error

	// ExistsUser はユーザーの存在確認を行う
	ExistsUser(ctx context.Context, userID uuid.UUID) (bool, error)
}

// userUseCase はUserUseCaseインターフェースの実装
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

// GetOrCreateUserFromCognito はCognito認証後にユーザーを取得または作成する（認証フロー用）
func (uc *userUseCase) GetOrCreateUserFromCognito(ctx context.Context, claims *auth.CognitoClaims) (*model.User, error) {
	if uc.userRepo == nil {
		return nil, fmt.Errorf("UserRepositoryが初期化されていません")
	}

	// ClaimsからユーザーIDを取得
	userID, err := claims.GetUserID()
	if err != nil {
		uc.logger.Errorf("ユーザーID取得エラー: %v", err)
		return nil, fmt.Errorf("ユーザーID取得エラー: %w", err)
	}

	// ユーザーの取得または作成
	user, err := uc.userRepo.GetOrCreateUser(ctx, userID, claims)
	if err != nil {
		uc.logger.Errorf("ユーザー取得・作成エラー: %v", err)
		return nil, fmt.Errorf("ユーザー取得・作成エラー: %w", err)
	}

	// プロバイダー別ログ出力
	if user.IsGoogleUser() {
		uc.logger.Infof("Googleユーザー認証成功: %s (%s)", user.Name, user.Email)
	} else {
		uc.logger.Infof("Cognitoユーザー認証成功: %s (%s)", user.Name, user.Email)
	}

	return user, nil
}

// GetUserProfile はユーザープロフィールを取得する
func (uc *userUseCase) GetUserProfile(ctx context.Context, userID uuid.UUID) (*model.UserResponse, error) {
	if uc.userRepo == nil {
		return nil, fmt.Errorf("UserRepositoryが初期化されていません")
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		uc.logger.Errorf("ユーザープロフィール取得エラー: %v", err)
		return nil, err
	}

	return user.ToResponse(), nil
}

// UpdateUserProfile はユーザープロフィールを更新する
func (uc *userUseCase) UpdateUserProfile(ctx context.Context, userID uuid.UUID, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	if uc.userRepo == nil {
		return nil, fmt.Errorf("UserRepositoryが初期化されていません")
	}

	// 現在のユーザー情報を取得
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		uc.logger.Errorf("ユーザー取得エラー: %v", err)
		return nil, err
	}

	// 更新内容の適用
	name := user.Name
	email := user.Email

	if req.Name != nil && *req.Name != "" {
		name = *req.Name
	}
	if req.Email != nil && *req.Email != "" {
		email = *req.Email
	}

	// プロフィール更新
	updatedUser, err := uc.userRepo.UpdateProfile(ctx, userID, name, email)
	if err != nil {
		uc.logger.Errorf("ユーザープロフィール更新エラー: %v", err)
		return nil, err
	}

	uc.logger.Infof("ユーザープロフィール更新成功: %s (%s)", updatedUser.Name, updatedUser.Email)

	return updatedUser.ToResponse(), nil
}

// GetUserByEmail はメールアドレスでユーザーを取得する
func (uc *userUseCase) GetUserByEmail(ctx context.Context, email string) (*model.UserResponse, error) {
	if uc.userRepo == nil {
		return nil, fmt.Errorf("UserRepositoryが初期化されていません")
	}

	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		uc.logger.Errorf("ユーザー（Email）取得エラー: %v", err)
		return nil, err
	}

	return user.ToResponse(), nil
}

// DeleteUser はユーザーを削除する（GDPR対応）
func (uc *userUseCase) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if uc.userRepo == nil {
		return fmt.Errorf("UserRepositoryが初期化されていません")
	}

	// ユーザーが存在するか確認
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		uc.logger.Errorf("削除対象ユーザー取得エラー: %v", err)
		return err
	}

	// ユーザー削除実行
	err = uc.userRepo.Delete(ctx, userID)
	if err != nil {
		uc.logger.Errorf("ユーザー削除エラー: %v", err)
		return err
	}

	uc.logger.Infof("ユーザー削除成功: %s", userID.String())
	return nil
}

// ExistsUser はユーザーの存在確認を行う
func (uc *userUseCase) ExistsUser(ctx context.Context, userID uuid.UUID) (bool, error) {
	if uc.userRepo == nil {
		return false, fmt.Errorf("UserRepositoryが初期化されていません")
	}

	exists, err := uc.userRepo.ExistsByID(ctx, userID)
	if err != nil {
		uc.logger.Errorf("ユーザー存在確認エラー: %v", err)
		return false, err
	}

	return exists, nil
}

// GetUserStatistics はユーザーの統計情報を取得する（将来の拡張用）
func (uc *userUseCase) GetUserStatistics(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	if uc.userRepo == nil {
		return nil, fmt.Errorf("UserRepositoryが初期化されていません")
	}

	// TODO: 統計情報の実装
	// - 総セッション数
	// - 総作業時間
	// - 平均集中度
	// - 最も活動的な時間帯

	stats := map[string]interface{}{
		"total_sessions":   0,
		"total_work_hours": 0,
		"average_focus":    0.0,
		"most_active_hour": nil,
		"provider":         "",
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return stats, err
	}

	stats["provider"] = user.Provider

	return stats, nil
}

// ValidateUserForOperation は操作前のユーザー検証を行う（共通処理）
func (uc *userUseCase) ValidateUserForOperation(ctx context.Context, userID uuid.UUID, operation string) error {
	exists, err := uc.ExistsUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("ユーザー存在確認エラー: %w", err)
	}

	if !exists {
		return fmt.Errorf("ユーザーが存在しません（操作: %s, ユーザーID: %s）", operation, userID.String()[:8]+"...")
	}

	return nil
}
