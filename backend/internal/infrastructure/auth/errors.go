package auth

import (
	"errors"
	"fmt"
)

// 技術的エラーのみを残す
var (
	// JWKS関連エラー（外部API技術エラー）
	ErrJWKSFetchFailed   = errors.New("JWKS取得に失敗しました")
	ErrPublicKeyNotFound = errors.New("Cognito公開キーが見つかりません")
	ErrPublicKeyInvalid  = errors.New("Cognito公開キーが無効です")

	// JWT解析関連エラー（技術エラー）
	ErrJWTParsingFailed = errors.New("JWT解析に失敗しました")
	ErrSignatureInvalid = errors.New("JWT署名が無効です")
)

// AuthError は認証関連の技術的エラーを表す構造体
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

// NewJWKSError はJWKS取得エラーを作成する
func NewJWKSError(cause error) *AuthError {
	return NewAuthError("JWKS", "JWKS取得に失敗しました", cause)
}

// NewPublicKeyError は公開キー取得エラーを作成する
func NewPublicKeyError(cause error) *AuthError {
	return NewAuthError("PUBLIC_KEY", "Cognito公開キーの取得に失敗しました", cause)
}

// NewJWTParsingError はJWT解析エラーを作成する
func NewJWTParsingError(cause error) *AuthError {
	return NewAuthError("JWT_PARSING", "JWT解析に失敗しました", cause)
}
