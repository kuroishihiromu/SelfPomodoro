package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// User はユーザーを表すドメインモデル（強化版）
type User struct {
	ID         uuid.UUID `db:"id" json:"id"`                   // Cognito sub
	Name       string    `db:"name" json:"name"`               // 表示名
	Email      string    `db:"email" json:"email"`             // メールアドレス
	Provider   string    `db:"provider" json:"provider"`       // "Cognito_UserPool", "Google"
	ProviderID *string   `db:"provider_id" json:"provider_id"` // Google sub等
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// CognitoUserParams はCognito PostConfirmation時のユーザー作成パラメータ
type CognitoUserParams struct {
	UserID     uuid.UUID
	Email      string
	Name       string
	GivenName  string
	FamilyName string
	// Google SSO関連（Cognito + Google統合時）
	Picture          string
	Locale           string
	IdentityProvider string
	// その他の属性
	EmailVerified bool
}

// UserCreationParams はユーザー作成時のパラメータ（後方互換性維持）
type UserCreationParams struct {
	UserID       uuid.UUID
	Name         string
	Email        string
	Provider     string
	ProviderID   *string
	IsGoogleUser bool
}

// ドメインルール：ユーザー作成ファクトリーメソッド

// NewUser は新しいユーザーを作成する（後方互換性維持）
func NewUser(params UserCreationParams) *User {
	now := time.Now()
	return &User{
		ID:         params.UserID,
		Name:       params.Name,
		Email:      params.Email,
		Provider:   params.Provider,
		ProviderID: params.ProviderID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// NewUserFromCognitoAttributes はCognito属性から新しいユーザーを作成する（メインファクトリー）
func NewUserFromCognitoAttributes(params CognitoUserParams) *User {
	displayName := determineDisplayName(params)
	provider, providerID := determineProvider(params)

	now := time.Now()
	return &User{
		ID:         params.UserID,
		Name:       displayName,
		Email:      params.Email,
		Provider:   provider,
		ProviderID: providerID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// ドメインルール：表示名決定ロジック（ビジネス優先順位）
func determineDisplayName(params CognitoUserParams) string {
	// 1. name フィールド（フルネーム）
	if params.Name != "" && !isEmailAddress(params.Name) {
		return params.Name
	}

	// 2. given_name + family_name の組み合わせ
	if params.GivenName != "" {
		name := params.GivenName
		if params.FamilyName != "" {
			name += " " + params.FamilyName
		}
		return name
	}

	// 3. メールアドレスのローカル部分（@より前）
	if params.Email != "" {
		if localPart := extractEmailLocalPart(params.Email); localPart != "" {
			return localPart
		}
		// フォールバック：メールアドレス全体
		return params.Email
	}

	// 4. 最終フォールバック
	return "ユーザー"
}

// ドメインルール：プロバイダー判定ロジック（Google SSO検出）
func determineProvider(params CognitoUserParams) (provider string, providerID *string) {
	// Google SSO判定の複数パターン

	// パターン1: identityProviderが設定されている（Federatedログイン）
	if params.IdentityProvider != "" {
		if strings.Contains(strings.ToLower(params.IdentityProvider), "google") {
			return "Google", stringPtr(params.UserID.String())
		}
	}

	// パターン2: picture フィールドがGoogleのドメインを含む
	if params.Picture != "" && strings.Contains(params.Picture, "googleusercontent.com") {
		return "Google", stringPtr(params.UserID.String())
	}

	// パターン3: locale フィールドが設定されている（Googleプロフィール特徴）
	if params.Locale != "" {
		return "Google", stringPtr(params.UserID.String())
	}

	// パターン4: given_name + family_name パターン（Google特有の構造）
	if params.GivenName != "" && params.FamilyName != "" && params.Name == "" {
		// このパターンはGoogle SSOで多い
		return "Google", stringPtr(params.UserID.String())
	}

	// デフォルト: Cognito User Pool
	return "Cognito_UserPool", nil
}

// ドメインルール：ユーザー検証・状態管理

// IsGoogleUser はGoogleユーザーかどうかを返す
func (u *User) IsGoogleUser() bool {
	return u.Provider == "Google"
}

// IsCognitoUser はCognitoユーザーかどうかを返す
func (u *User) IsCognitoUser() bool {
	return u.Provider == "Cognito_UserPool"
}

// ValidateEmail はメールアドレスの基本的な検証を行う
func (u *User) ValidateEmail() bool {
	return u.Email != "" && isEmailAddress(u.Email)
}

// ValidateName は名前の基本的な検証を行う
func (u *User) ValidateName() bool {
	return u.Name != "" && len(strings.TrimSpace(u.Name)) > 0
}

// IsValidForCreation はユーザー作成時の必須項目を検証する
func (u *User) IsValidForCreation() bool {
	return u.ID != uuid.Nil &&
		u.ValidateEmail() &&
		u.ValidateName() &&
		u.Provider != ""
}

// UpdateProfile はユーザープロフィールを更新する
func (u *User) UpdateProfile(name, email string) {
	if name != "" {
		u.Name = name
	}
	if email != "" {
		u.Email = email
	}
	u.UpdatedAt = time.Now()
}

// ドメインルール：プロバイダー別ビジネスロジック

// GetProviderDisplayName はプロバイダーの表示名を返す
func (u *User) GetProviderDisplayName() string {
	switch u.Provider {
	case "Google":
		return "Google"
	case "Cognito_UserPool":
		return "Email"
	default:
		return "不明"
	}
}

// CanChangeEmail はメールアドレス変更可能かを判定する
func (u *User) CanChangeEmail() bool {
	// Google ユーザーはメールアドレス変更不可（プロバイダー管理）
	return !u.IsGoogleUser()
}

// GetAvatarURL はアバターURLを生成する（将来の拡張用）
func (u *User) GetAvatarURL() string {
	if u.IsGoogleUser() {
		// Google ユーザーの場合はGravatarやGoogle Photos APIを使用
		// 現在は簡易実装
		return ""
	}
	// Cognito ユーザーの場合はGravatarを使用
	return ""
}

// ドメインルール：ヘルパー関数

// isEmailAddress はメールアドレス形式かを簡易チェック
func isEmailAddress(str string) bool {
	return strings.Contains(str, "@") && len(str) > 3
}

// extractEmailLocalPart はメールアドレスのローカル部分を抽出
func extractEmailLocalPart(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}
	return ""
}

// stringPtr は文字列のポインタを返すヘルパー
func stringPtr(s string) *string {
	return &s
}

// 既存のリクエスト・レスポンス構造体は維持

// CreateUserRequest はユーザー作成リクエストを表す
type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// UpdateUserRequest はユーザー更新リクエストを表す
type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty" validate:"omitempty"`
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
}

// UserResponse はユーザーのAPIレスポンスを表す
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse はドメインモデルからAPIレスポンス形式に変換する
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Provider:  u.Provider,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
