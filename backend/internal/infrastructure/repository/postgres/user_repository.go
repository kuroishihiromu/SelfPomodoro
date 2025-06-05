package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

// GetOrCreateUser はユーザーを取得し、存在しない場合はCognitoクレームから作成する（最重要メソッド）
func (r *UserRepositoryImpl) GetOrCreateUser(ctx context.Context, userID uuid.UUID, claims *auth.CognitoClaims) (*model.User, error) {
	r.logger.Infof("GetOrCreateUser開始: UserID=%s, Email=%s, Provider=%s",
		userID.String()[:8]+"...", claims.Email, claims.GetProviderName())

	// まずユーザーの取得を試行
	user, err := r.GetByID(ctx, userID)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, fmt.Errorf("ユーザー取得エラー: %w", err)
	}

	// ユーザーが見つかった場合はそのまま返す
	if user != nil {
		r.logger.Infof("既存ユーザー取得成功: %s (%s)", user.Name, user.Provider)
		return user, nil
	}

	// ユーザーが見つからない場合は新規作成
	r.logger.Infof("ユーザーが見つかりません。新規作成します: UserID=%s", userID.String()[:8]+"...")

	newUser := model.NewUserFromCognito(userID, claims)

	// 作成を試行（競合状態を考慮）
	err = r.Create(ctx, newUser)
	if err != nil {
		// 別のリクエストが同時に作成した可能性があるため、再度取得を試行
		if errors.Is(err, ErrUserCreationFailed) ||
			(err != nil && fmt.Sprintf("%v", err) == "ユーザーID重複") {
			r.logger.Warnf("ユーザー作成競合が発生、再取得を試行: %v", err)

			existingUser, getErr := r.GetByID(ctx, userID)
			if getErr == nil {
				r.logger.Infof("競合後の再取得成功: %s", existingUser.Name)
				return existingUser, nil
			}
		}
		return nil, fmt.Errorf("ユーザー作成失敗: %w", err)
	}

	r.logger.Infof("新規ユーザー作成成功: %s (%s)", newUser.Name, newUser.Provider)
	return newUser, nil
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
