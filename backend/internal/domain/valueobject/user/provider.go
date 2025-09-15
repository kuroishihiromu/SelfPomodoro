package user

import (
	"errors"
	"strings"
)

// Provider は認証プロバイダーを表すValue Object
type Provider struct {
	value string
}

// 認証プロバイダーの定数
const (
	ProviderCognito = "Cognito_UserPool"
	ProviderGoogle  = "Google"
)

// NewProvider は新しいProviderを作成する
func NewProvider(provider string) (Provider, error) {
	trimmed := strings.TrimSpace(provider)
	
	if trimmed == "" {
		return Provider{}, errors.New("プロバイダーは必須です")
	}
	
	// 有効なプロバイダーかチェック
	if !isValidProvider(trimmed) {
		return Provider{}, errors.New("無効なプロバイダーです")
	}
	
	return Provider{value: trimmed}, nil
}

// NewCognitoProvider はCognitoプロバイダーを作成する
func NewCognitoProvider() Provider {
	return Provider{value: ProviderCognito}
}

// NewGoogleProvider はGoogleプロバイダーを作成する
func NewGoogleProvider() Provider {
	return Provider{value: ProviderGoogle}
}

// MustNewProvider はパニックを起こす可能性があるProvider作成（テスト用）
func MustNewProvider(provider string) Provider {
	p, err := NewProvider(provider)
	if err != nil {
		panic(err)
	}
	return p
}

// Value はプロバイダーの値を返す
func (p Provider) Value() string {
	return p.value
}

// String はプロバイダーの文字列表現を返す
func (p Provider) String() string {
	return p.value
}

// IsEmpty はプロバイダーが空かどうかを判定する
func (p Provider) IsEmpty() bool {
	return p.value == ""
}

// IsCognito はCognitoプロバイダーかどうかを判定する
func (p Provider) IsCognito() bool {
	return p.value == ProviderCognito
}

// IsGoogle はGoogleプロバイダーかどうかを判定する
func (p Provider) IsGoogle() bool {
	return p.value == ProviderGoogle
}

// IsThirdParty はサードパーティープロバイダーかどうかを判定する
func (p Provider) IsThirdParty() bool {
	return p.IsGoogle() // 現在はGoogleのみ
}

// CanChangeEmail はメールアドレス変更可能かを判定する
func (p Provider) CanChangeEmail() bool {
	// サードパーティープロバイダーはメールアドレス変更不可
	return !p.IsThirdParty()
}

// GetDisplayName はプロバイダーの表示名を返す
func (p Provider) GetDisplayName() string {
	switch p.value {
	case ProviderGoogle:
		return "Google"
	case ProviderCognito:
		return "Email"
	default:
		return "不明"
	}
}

// Equals は他のProviderと等しいかを判定する
func (p Provider) Equals(other Provider) bool {
	return p.value == other.value
}

// RequiresProviderID はプロバイダーIDが必要かを判定する
func (p Provider) RequiresProviderID() bool {
	return p.IsThirdParty()
}

// GetAuthenticationMethod は認証方法を返す
func (p Provider) GetAuthenticationMethod() string {
	switch p.value {
	case ProviderGoogle:
		return "OAuth 2.0"
	case ProviderCognito:
		return "Username/Password"
	default:
		return "Unknown"
	}
}

// isValidProvider は有効なプロバイダーかをチェックする
func isValidProvider(provider string) bool {
	validProviders := []string{
		ProviderCognito,
		ProviderGoogle,
	}
	
	for _, valid := range validProviders {
		if provider == valid {
			return true
		}
	}
	
	return false
}

// GetAllValidProviders は全ての有効なプロバイダーを返す
func GetAllValidProviders() []string {
	return []string{
		ProviderCognito,
		ProviderGoogle,
	}
}