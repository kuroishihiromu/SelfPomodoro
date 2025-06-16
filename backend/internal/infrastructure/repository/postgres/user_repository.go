package postgres

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/auth"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/database"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// UserRepositoryImpl はUserRepositoryインターフェースの実装（新エラーハンドリング対応版）
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

// GetByID はIDによってユーザーを取得する（新エラーハンドリング対応版）
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
			r.logger.Debugf("ユーザーが見つかりません: %s", id.String())
			return nil, appErrors.ErrRecordNotFound // Infrastructure Error
		}
		r.logger.Errorf("ユーザー取得エラー: %v", err)
		return nil, appErrors.NewDatabaseQueryError(err)
	}

	return &user, nil
}

// GetByEmail はメールアドレスによってユーザーを取得する（新エラーハンドリング対応版）
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
			r.logger.Debugf("ユーザー（Email）が見つかりません: %s", email)
			return nil, appErrors.ErrRecordNotFound // Infrastructure Error
		}
		r.logger.Errorf("ユーザー（Email）取得エラー: %v", err)
		return nil, appErrors.NewDatabaseQueryError(err)
	}

	return &user, nil
}

// Create は新しいユーザーを作成する（新エラーハンドリング対応版）
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
		// PostgreSQL固有のエラーハンドリング
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if pqErr.Constraint == "users_pkey" {
					r.logger.Errorf("ユーザーID重複: %v", err)
					return appErrors.NewUniqueConstraintError(err)
				}
				if pqErr.Constraint == "users_email_key" {
					r.logger.Errorf("メールアドレス重複: %v", err)
					return appErrors.NewUniqueConstraintError(err)
				}
				r.logger.Errorf("ユーザー作成一意制約違反: %v", err)
				return appErrors.NewUniqueConstraintError(err)
			case "23503": // foreign_key_violation
				r.logger.Errorf("ユーザー作成外部キー制約違反: %v", err)
				return appErrors.NewDatabaseError("create_user_fk", err)
			}
		}
		r.logger.Errorf("ユーザー作成エラー: %v", err)
		return appErrors.NewDatabaseError("create_user", err)
	}

	return nil
}

// Update はユーザー情報を更新する（新エラーハンドリング対応版）
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
		// PostgreSQL固有のエラーハンドリング
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if pqErr.Constraint == "users_email_key" {
					r.logger.Errorf("メールアドレス重複（更新）: %v", err)
					return appErrors.NewUniqueConstraintError(err)
				}
				r.logger.Errorf("ユーザー更新一意制約違反: %v", err)
				return appErrors.NewUniqueConstraintError(err)
			}
		}
		r.logger.Errorf("ユーザー更新エラー: %v", err)
		return appErrors.NewDatabaseError("update_user", err)
	}

	// 更新された行数の確認
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Errorf("ユーザー更新結果確認エラー: %v", err)
		return appErrors.NewDatabaseError("update_user_check", err)
	}

	if rowsAffected == 0 {
		r.logger.Warnf("ユーザー更新対象なし: %s", user.ID.String())
		return appErrors.ErrRecordNotFound
	}

	return nil
}

// Delete はユーザーを削除する（新エラーハンドリング対応版）
func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Errorf("ユーザー削除エラー: %v", err)
		return appErrors.NewDatabaseError("delete_user", err)
	}

	// 削除された行数の確認
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Errorf("ユーザー削除結果確認エラー: %v", err)
		return appErrors.NewDatabaseError("delete_user_check", err)
	}

	if rowsAffected == 0 {
		r.logger.Warnf("ユーザー削除対象なし: %s", id.String())
		return appErrors.ErrRecordNotFound
	}

	return nil
}

// extractDisplayName はCognitoクレームから表示名を抽出する（ヘルパーメソッド維持）
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

// detectProviderType はCognitoクレームからプロバイダータイプを判定する（ヘルパーメソッド維持）
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

// UpdateProfile はユーザープロフィール（名前・メール）を更新する（新エラーハンドリング対応版）
func (r *UserRepositoryImpl) UpdateProfile(ctx context.Context, id uuid.UUID, name, email string) (*model.User, error) {
	// 現在のユーザー情報を取得
	user, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err // GetByIDのInfrastructure Errorをそのまま返す
	}

	// プロフィール更新
	user.UpdateProfile(name, email)

	// データベース更新
	err = r.Update(ctx, user)
	if err != nil {
		return nil, err // UpdateのInfrastructure Errorをそのまま返す
	}

	return user, nil
}

// ExistsByID はユーザーの存在確認を行う（軽量版）（新エラーハンドリング対応版）
func (r *UserRepositoryImpl) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

	var exists bool
	err := r.db.DB.GetContext(ctx, &exists, query, id)
	if err != nil {
		r.logger.Errorf("ユーザー存在確認エラー: %v", err)
		return false, appErrors.NewDatabaseQueryError(err)
	}

	return exists, nil
}

// GetUsersByProvider はプロバイダー別にユーザーを取得する（管理用）（新エラーハンドリング対応版）
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
		return nil, appErrors.NewDatabaseQueryError(err)
	}

	return users, nil
}
