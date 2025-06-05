package auth

import (
	"errors"
	"fmt"
)

// 認証関連のエラー定義
var (
	// JWT検証エラー
	ErrInvalidToken     = errors.New("無効なJWTトークンです")
	ErrTokenExpired     = errors.New("JWTトークンの有効期限が切れています")
	ErrTokenNotFound    = errors.New("Authorizationヘッダーが見つかりません")
	ErrInvalidFormat    = errors.New("Authorizationヘッダーの形式が無効です")
	ErrInvalidIssuer    = errors.New("トークンの発行者が無効です")
	ErrInvalidAudience  = errors.New("トークンのaudienceが無効です")
	ErrInvalidSignature = errors.New("トークンの署名が無効です")

	// Cognito公開キー関連エラー
	ErrPublicKeyNotFound = errors.New("Cognito公開キーが見つかりません")
	ErrPublicKeyInvalid  = errors.New("Cognito公開キーが無効です")
	ErrJWKSFetchFailed   = errors.New("JWKS取得に失敗しました")

	// Claims関連エラー
	ErrMissingSubject  = errors.New("subjectクレームが見つかりません")
	ErrInvalidSubject  = errors.New("subjectクレームが無効です")
	ErrMissingTokenUse = errors.New("token_useクレームが見つかりません")
	ErrInvalidTokenUse = errors.New("token_useクレームが無効です")
)

// AuthError は認証関連のエラーを表す構造体
type AuthError struct {
	Type    string // エラータイプ
	Message string // エラーメッセージ
	Cause   error  // 元のエラー
}

// Error はエラーメッセージを返す
func (e *AuthError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("認証エラー [%s]: %s (原因: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("認証エラー [%s]: %s", e.Type, e.Message)
}

// Unwrap は元のエラーを返す
func (e *AuthError) Unwrap() error {
	return e.Cause
}

// NewAuthError は新しい認証エラーを作成する
func NewAuthError(errorType, message string, cause error) *AuthError {
	return &AuthError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
	}
}

// NewTokenValidationError はトークン検証エラーを作成する
func NewTokenValidationError(cause error) *AuthError {
	return NewAuthError("TOKEN_VALIDATION", "JWTトークンの検証に失敗しました", cause)
}

// NewPublicKeyError は公開キー取得エラーを作成する
func NewPublicKeyError(cause error) *AuthError {
	return NewAuthError("PUBLIC_KEY", "Cognito公開キーの取得に失敗しました", cause)
}

// NewClaimsError はクレーム検証エラーを作成する
func NewClaimsError(message string) *AuthError {
	return NewAuthError("CLAIMS", message, nil)
}

// IsTokenExpiredError はトークン有効期限エラーかどうかを判定する
func IsTokenExpiredError(err error) bool {
	if authErr, ok := err.(*AuthError); ok {
		return authErr.Type == "TOKEN_VALIDATION" &&
			(errors.Is(authErr.Cause, ErrTokenExpired) ||
				authErr.Cause != nil && errors.Is(authErr.Cause, ErrTokenExpired))
	}
	return errors.Is(err, ErrTokenExpired)
}

// IsInvalidTokenError は無効トークンエラーかどうかを判定する
func IsInvalidTokenError(err error) bool {
	if authErr, ok := err.(*AuthError); ok {
		return authErr.Type == "TOKEN_VALIDATION"
	}
	return errors.Is(err, ErrInvalidToken)
}
