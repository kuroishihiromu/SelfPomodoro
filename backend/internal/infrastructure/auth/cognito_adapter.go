package auth

import (
	"strings"

	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
)

// CognitoAdapter はCognito固有の処理をドメインエンティティに変換するアダプター
type CognitoAdapter struct{}

// NewCognitoAdapter は新しいCognitoAdapterを作成する
func NewCognitoAdapter() *CognitoAdapter {
	return &CognitoAdapter{}
}

// ConvertToUser はCognitoUserCreationParamsからUserエンティティを作成する
func (a *CognitoAdapter) ConvertToUser(params CognitoUserCreationParams) (*entity.User, error) {
	displayName := a.determineDisplayName(params)
	provider, providerID := a.determineProvider(params)

	// UserName Value Object を作成
	userName, err := userVO.NewUserName(displayName)
	if err != nil {
		return nil, err
	}

	// EmailAddress Value Object を作成
	emailAddress, err := userVO.NewEmailAddress(params.Email)
	if err != nil {
		return nil, err
	}

	// Provider Value Object を作成
	userProvider, err := userVO.NewProvider(provider)
	if err != nil {
		return nil, err
	}

	user := entity.NewUser(entity.UserCreationParams{
		UserID:     userVO.NewUserIDFromUUID(params.UserID),
		Name:       userName,
		Email:      emailAddress,
		Provider:   userProvider,
		ProviderID: providerID,
	})

	return user, nil
}

// determineDisplayName は表示名決定ロジック（ビジネス優先順位）
func (a *CognitoAdapter) determineDisplayName(params CognitoUserCreationParams) string {
	// 1. name フィールド（フルネーム）
	if params.Name != "" && !a.isEmailAddress(params.Name) {
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
		if localPart := a.extractEmailLocalPart(params.Email); localPart != "" {
			return localPart
		}
		// フォールバック：メールアドレス全体
		return params.Email
	}

	// 4. 最終フォールバック
	return "ユーザー"
}

// determineProvider はプロバイダー判定ロジック（Google SSO検出）
func (a *CognitoAdapter) determineProvider(params CognitoUserCreationParams) (provider string, providerID *string) {
	// Google SSO判定の複数パターン

	// パターン1: identityProviderが設定されている（Federatedログイン）
	if params.IdentityProvider != "" {
		if strings.Contains(strings.ToLower(params.IdentityProvider), "google") {
			return "Google", a.stringPtr(params.UserID.String())
		}
	}

	// パターン2: picture フィールドがGoogleのドメインを含む
	if params.Picture != "" && strings.Contains(params.Picture, "googleusercontent.com") {
		return "Google", a.stringPtr(params.UserID.String())
	}

	// パターン3: locale フィールドが設定されている（Googleプロフィール特徴）
	if params.Locale != "" {
		return "Google", a.stringPtr(params.UserID.String())
	}

	// パターン4: given_name + family_name パターン（Google特有の構造）
	if params.GivenName != "" && params.FamilyName != "" && params.Name == "" {
		// このパターンはGoogle SSOで多い
		return "Google", a.stringPtr(params.UserID.String())
	}

	// デフォルト: Cognito User Pool
	return "Cognito_UserPool", nil
}

// ヘルパー関数

// isEmailAddress はメールアドレス形式かを簡易チェック
func (a *CognitoAdapter) isEmailAddress(str string) bool {
	return strings.Contains(str, "@") && len(str) > 3
}

// extractEmailLocalPart はメールアドレスのローカル部分を抽出
func (a *CognitoAdapter) extractEmailLocalPart(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}
	return ""
}

// stringPtr は文字列のポインタを返すヘルパー
func (a *CognitoAdapter) stringPtr(s string) *string {
	return &s
}