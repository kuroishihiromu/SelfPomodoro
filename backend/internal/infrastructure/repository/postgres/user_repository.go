package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/auth"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/database"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// ユーザーリポジトリに関するエラー
var (
	ErrUserNotFound       = errors.New("ユーザーが見つかりません")
	ErrUserCreationFailed = errors.New("ユーザーの作成に失敗しました")
	ErrUserUpdateFailed   = errors.New("ユーザーの更新に失敗しました")
	ErrUserDeleteFailed   = errors.New("ユーザーの削除に失敗しました")
	ErrEmailAlreadyExists = errors.New("メールアドレスが既に使用されています")
)

// UserRepositoryImpl はUserRepositoryインターフェースの実装
type UserRepositoryImpl struct {
	db     *database.PostgresDB
	logger logger.Logger
}

// NewUserRepository は新しいUserRepositoryImplインスタンスを作成する
func NewUserRepository(db *database.PostgresDB, logger logger.Logger) repository.UserRepository {
	return &UserRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

// GetByID はIDによってユーザーを取得する
func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `
		SELECT id, name, email, provider, provider_id, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user model.User
	err := r.db.DB.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		r.logger.Errorf("ユーザー取得エラー: %v", err)
		return nil, err
	}

	return &user, nil
}

// GetByEmail はメールアドレスによってユーザーを取得する
func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, name, email, provider, provider_id, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user model.User
	err := r.db.DB.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		r.logger.Errorf("ユーザー（Email）取得エラー: %v", err)
		return nil, err
	}

	return &user, nil
}

// Create は新しいユーザーを作成する
func (r *UserRepositoryImpl) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, name, email, provider, provider_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.DB.ExecContext(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.Provider,
		user.ProviderID,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		// PostgreSQL重複エラーチェック
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if pqErr.Constraint == "users_pkey" {
					return fmt.Errorf("ユーザーID重複: %v", err)
				}
				if pqErr.Constraint == "users_email_key" {
					return ErrEmailAlreadyExists
				}
			}
		}
		r.logger.Errorf("ユーザー作成エラー: %v", err)
		return fmt.Errorf("%w: %v", ErrUserCreationFailed, err)
	}

	return nil
}

// Update はユーザー情報を更新する
func (r *UserRepositoryImpl) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, provider = $3, provider_id = $4, updated_at = $5
		WHERE id = $6
	`

	result, err := r.db.DB.ExecContext(ctx, query,
		user.Name,
		user.Email,
		user.Provider,
		user.ProviderID,
		time.Now(),
		user.ID,
	)

	if err != nil {
		r.logger.Errorf("ユーザー更新エラー: %v", err)
		return fmt.Errorf("%w: %v", ErrUserUpdateFailed, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("更新結果確認エラー: %v", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// Delete はユーザーを削除する
func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Errorf("ユーザー削除エラー: %v", err)
		return fmt.Errorf("%w: %v", ErrUserDeleteFailed, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("削除結果確認エラー: %v", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// extractDisplayName はCognitoクレームから表示名を抽出する
func (r *UserRepositoryImpl) extractDisplayName(claims *auth.CognitoClaims) string {
	// 1. name フィールド
	if claims.Name != "" {
		return claims.Name
	}

	// 2. given_name + family_name
	if claims.GivenName != "" {
		name := claims.GivenName
		if claims.FamilyName != "" {
			name += " " + claims.FamilyName
		}
		return name
	}

	// 3. cognito:username
	if claims.CognitoUsername != "" {
		return claims.CognitoUsername
	}

	// 4. email（最後の手段）
	if claims.Email != "" {
		return claims.Email
	}

	// 5. ユーザーID（非常時）
	return "User"
}

// detectProviderType はCognitoクレームからプロバイダータイプを判定する
func (r *UserRepositoryImpl) detectProviderType(claims *auth.CognitoClaims) (provider string, providerID *string) {
	// Google SSO判定の複数パターン

	// パターン1: identities claim
	if len(claims.Identities) > 0 {
		for _, identity := range claims.Identities {
			if strings.Contains(identity.ProviderID, "google") ||
				strings.Contains(identity.ProviderID, "Google") {
				provider = "Google"
				providerID = &identity.UserID
				return
			}
		}
	}

	// パターン2: IdentityProvider フィールド
	if claims.IdentityProvider != "" {
		provider = "Google"
		if claims.Subject != "" {
			providerID = &claims.Subject
		}
		return
	}

	// パターン3: Picture（Googleプロフィール画像URL）
	if claims.Picture != "" && strings.Contains(claims.Picture, "googleusercontent.com") {
		provider = "Google"
		if claims.Subject != "" {
			providerID = &claims.Subject
		}
		return
	}

	// パターン4: Locale（Google特有のlocale形式）
	if claims.Locale != "" {
		provider = "Google"
		if claims.Subject != "" {
			providerID = &claims.Subject
		}
		return
	}

	// デフォルト: Cognito User Pool
	provider = "Cognito_UserPool"
	providerID = nil
	return
}

// UpdateProfile はユーザープロフィール（名前・メール）を更新する
func (r *UserRepositoryImpl) UpdateProfile(ctx context.Context, id uuid.UUID, name, email string) (*model.User, error) {
	// 現在のユーザー情報を取得
	user, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// プロフィール更新
	user.UpdateProfile(name, email)

	// データベース更新
	err = r.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ExistsByID はユーザーの存在確認を行う（軽量版）
func (r *UserRepositoryImpl) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

	var exists bool
	err := r.db.DB.GetContext(ctx, &exists, query, id)
	if err != nil {
		r.logger.Errorf("ユーザー存在確認エラー: %v", err)
		return false, err
	}

	return exists, nil
}

// GetUsersByProvider はプロバイダー別にユーザーを取得する（管理用）
func (r *UserRepositoryImpl) GetUsersByProvider(ctx context.Context, provider string, limit, offset int) ([]*model.User, error) {
	query := `
		SELECT id, name, email, provider, provider_id, created_at, updated_at
		FROM users
		WHERE provider = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var users []*model.User
	err := r.db.DB.SelectContext(ctx, &users, query, provider, limit, offset)
	if err != nil {
		r.logger.Errorf("プロバイダー別ユーザー取得エラー: %v", err)
		return nil, err
	}

	return users, nil
}
