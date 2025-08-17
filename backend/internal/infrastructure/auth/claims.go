package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
)

// CognitoClaims はCognito JWTのクレームを表す構造体
type CognitoClaims struct {
	// 標準クレーム
	Subject   string `json:"sub"` // ユーザーID（Cognito User Pool内のユニークID）
	Issuer    string `json:"iss"` // 発行者（Cognito User Pool）
	Audience  string `json:"aud"` // オーディエンス（Client ID）
	ExpiresAt int64  `json:"exp"` // 有効期限（Unix timestamp）
	IssuedAt  int64  `json:"iat"` // 発行時刻（Unix timestamp）

	// Cognito固有クレーム
	TokenUse        string `json:"token_use"`                // "access" or "id"
	CognitoUsername string `json:"cognito:username"`         // Cognitoユーザー名
	Email           string `json:"email,omitempty"`          // メールアドレス
	EmailVerified   bool   `json:"email_verified,omitempty"` // メール確認済み
	GivenName       string `json:"given_name,omitempty"`     // 名前
	FamilyName      string `json:"family_name,omitempty"`    // 姓
	Name            string `json:"name,omitempty"`           // フルネーム

	// Google SSO関連クレーム（Cognito + Google IDプロバイダー統合時）
	IdentityProvider string                    `json:"identities,omitempty"` // IDプロバイダー情報
	Picture          string                    `json:"picture,omitempty"`    // プロフィール画像URL
	Locale           string                    `json:"locale,omitempty"`     // ロケール情報
	Identities       []CognitoIdentityProvider `json:"-"`                    // 詳細IDプロバイダー情報

	// カスタムクレーム（将来の拡張用）
	CustomAttributes map[string]interface{} `json:"-"` // カスタム属性

	jwt.RegisteredClaims
}

// CognitoIdentityProvider はCognitoのIDプロバイダー情報を表す
type CognitoIdentityProvider struct {
	UserID      string `json:"userId"`
	ProviderID  string `json:"providerName"`
	IsPrimary   bool   `json:"primary"`
	DateCreated string `json:"dateCreated"`
}

// NewCognitoClaims は新しいCognitoクレームを作成する
func NewCognitoClaims() *CognitoClaims {
	return &CognitoClaims{
		CustomAttributes: make(map[string]interface{}),
	}
}

// IsValid はクレームの有効性を検証する
func (c *CognitoClaims) IsValid() error {
	now := time.Now()

	// 必須フィールドの確認
	if c.Subject == "" {
		return appErrors.ErrMissingSubject
	}

	if c.TokenUse == "" {
		return appErrors.ErrMissingTokenUse
	}

	// token_useの値をチェック（accessまたはidのみ許可）
	if c.TokenUse != "access" && c.TokenUse != "id" {
		return appErrors.ErrInvalidTokenUse
	}

	// 有効期限の確認
	if c.ExpiresAt != 0 && now.Unix() > c.ExpiresAt {
		return appErrors.ErrTokenExpired
	}

	// subjectがUUID形式かどうかチェック
	if _, err := uuid.Parse(c.Subject); err != nil {
		return appErrors.ErrInvalidSubject
	}

	return nil
}

// GetUserID はユーザーIDをUUID形式で返す
func (c *CognitoClaims) GetUserID() (uuid.UUID, error) {
	return uuid.Parse(c.Subject)
}

// IsExpired はトークンが期限切れかどうかを返す
func (c *CognitoClaims) IsExpired() bool {
	if c.ExpiresAt == 0 {
		return false
	}
	return time.Now().Unix() > c.ExpiresAt
}

// IsAccessToken はアクセストークンかどうかを返す
func (c *CognitoClaims) IsAccessToken() bool {
	return c.TokenUse == "access"
}

// IsIDToken はIDトークンかどうかを返す
func (c *CognitoClaims) IsIDToken() bool {
	return c.TokenUse == "id"
}

// IsGoogleUser はGoogleユーザーかどうかを返す
func (c *CognitoClaims) IsGoogleUser() bool {
	return c.IdentityProvider != "" || len(c.Identities) > 0
}

// GetProviderName はプロバイダー名を返す
func (c *CognitoClaims) GetProviderName() string {
	if c.IsGoogleUser() {
		return "Google"
	}
	return "Cognito_UserPool"
}

// GetDisplayName は表示名を返す（優先順位: name > given_name + family_name > cognito:username > email）
func (c *CognitoClaims) GetDisplayName() string {
	if c.Name != "" {
		return c.Name
	}

	if c.GivenName != "" {
		name := c.GivenName
		if c.FamilyName != "" {
			name += " " + c.FamilyName
		}
		return name
	}

	if c.CognitoUsername != "" {
		return c.CognitoUsername
	}

	return c.Email
}

// jwt.Claims インターフェース実装

// GetExpirationTime は有効期限をjwt.NumericDate形式で返す
func (c *CognitoClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	if c.ExpiresAt == 0 {
		return nil, nil
	}
	return jwt.NewNumericDate(time.Unix(c.ExpiresAt, 0)), nil
}

// GetIssuedAt は発行時刻をjwt.NumericDate形式で返す
func (c *CognitoClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	if c.IssuedAt == 0 {
		return nil, nil
	}
	return jwt.NewNumericDate(time.Unix(c.IssuedAt, 0)), nil
}

// GetNotBefore は有効開始時刻をjwt.NumericDate形式で返す
func (c *CognitoClaims) GetNotBefore() (*jwt.NumericDate, error) {
	// Cognitoでは通常 nbf クレームは使用されない
	return nil, nil
}

// GetIssuer は発行者を返す
func (c *CognitoClaims) GetIssuer() (string, error) {
	return c.Issuer, nil
}

// GetSubject はサブジェクトを返す
func (c *CognitoClaims) GetSubject() (string, error) {
	return c.Subject, nil
}

// GetAudience はオーディエンスを返す
func (c *CognitoClaims) GetAudience() (jwt.ClaimStrings, error) {
	if c.Audience == "" {
		return nil, nil
	}
	return jwt.ClaimStrings{c.Audience}, nil
}

// 便利メソッド（後方互換性）

// GetExpiresAtTime は有効期限をtime.Time形式で返す
func (c *CognitoClaims) GetExpiresAtTime() time.Time {
	if c.ExpiresAt == 0 {
		return time.Time{}
	}
	return time.Unix(c.ExpiresAt, 0)
}

// GetIssuedAtTime は発行時刻をtime.Time形式で返す
func (c *CognitoClaims) GetIssuedAtTime() time.Time {
	if c.IssuedAt == 0 {
		return time.Time{}
	}
	return time.Unix(c.IssuedAt, 0)
}

// SetCustomAttribute はカスタム属性を設定する
func (c *CognitoClaims) SetCustomAttribute(key string, value interface{}) {
	if c.CustomAttributes == nil {
		c.CustomAttributes = make(map[string]interface{})
	}
	c.CustomAttributes[key] = value
}

// GetCustomAttribute はカスタム属性を取得する
func (c *CognitoClaims) GetCustomAttribute(key string) (interface{}, bool) {
	if c.CustomAttributes == nil {
		return nil, false
	}
	value, exists := c.CustomAttributes[key]
	return value, exists
}

// ToLogString はログ出力用の文字列を返す（セキュリティ考慮）
func (c *CognitoClaims) ToLogString() string {
	userIDShort := c.Subject
	if len(userIDShort) > 8 {
		userIDShort = userIDShort[:8] + "..."
	}

	return fmt.Sprintf("CognitoClaims[Sub=%s, TokenUse=%s, Email=%s, Provider=%s, Expired=%t]",
		userIDShort, c.TokenUse, c.Email, c.GetProviderName(), c.IsExpired())
}
