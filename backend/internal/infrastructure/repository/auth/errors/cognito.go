package errors

import (
	"errors"
	"fmt"
)

// Cognito認証インフラエラー定義
var (
	// JWT検証エラー
	ErrJWTTokenNotFound    = errors.New("JWTトークンが見つかりません")
	ErrJWTTokenExpired     = errors.New("JWTトークンの有効期限が切れています")
	ErrJWTTokenInvalid     = errors.New("無効なJWTトークンです")
	ErrJWTParsingFailed    = errors.New("JWTトークンの解析に失敗しました")
	ErrJWTSignatureInvalid = errors.New("JWTトークンの署名が無効です")

	// JWT Claims エラー
	ErrJWTInvalidIssuer   = errors.New("JWTトークンの発行者が無効です")
	ErrJWTInvalidAudience = errors.New("JWTトークンのaudienceが無効です")
	ErrJWTMissingSubject  = errors.New("JWTトークンのsubjectが見つかりません")
	ErrJWTInvalidSubject  = errors.New("JWTトークンのsubjectが無効です")
	ErrJWTMissingTokenUse = errors.New("JWTトークンのtoken_useが見つかりません")
	ErrJWTInvalidTokenUse = errors.New("JWTトークンのtoken_useが無効です")

	// JWKS関連エラー
	ErrJWKSFetchFailed   = errors.New("JWKS取得に失敗しました")
	ErrJWKSDecodeFailed  = errors.New("JWKSデコードに失敗しました")
	ErrPublicKeyNotFound = errors.New("公開キーが見つかりません")
	ErrPublicKeyInvalid  = errors.New("公開キーが無効です")
	ErrKeyIDNotFound     = errors.New("指定されたKey IDが見つかりません")

	// HTTP関連エラー
	ErrHTTPRequestFailed   = errors.New("HTTPリクエストに失敗しました")
	ErrHTTPResponseInvalid = errors.New("HTTPレスポンスが無効です")

	// キャッシュ関連エラー
	ErrCacheExpired  = errors.New("キャッシュの有効期限が切れています")
	ErrCacheNotFound = errors.New("キャッシュが見つかりません")
)

// CognitoAuthError はCognito認証関連のインフラエラーを表す構造体
type CognitoAuthError struct {
	Type    string // エラータイプ
	Message string // エラーメッセージ
	Cause   error  // 元のエラー
}

// Error はエラーメッセージを返す
func (e *CognitoAuthError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("Cognito認証エラー [%s]: %s (原因: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("Cognito認証エラー [%s]: %s", e.Type, e.Message)
}

// Unwrap は元のエラーを返す
func (e *CognitoAuthError) Unwrap() error {
	return e.Cause
}

// Is はエラーの種類を判定する
func (e *CognitoAuthError) Is(target error) bool {
	return errors.Is(e.Cause, target)
}

// NewCognitoAuthError は新しいCognito認証エラーを作成する
func NewCognitoAuthError(errorType, message string, cause error) *CognitoAuthError {
	return &CognitoAuthError{
		Type:    errorType,
		Message: message,
		Cause:   cause,
	}
}

// 具体的なエラー作成関数
func NewJWTParsingError(cause error) *CognitoAuthError {
	return NewCognitoAuthError("JWT_PARSING", "JWTトークンの解析に失敗しました", cause)
}

func NewPublicKeyError(cause error) *CognitoAuthError {
	return NewCognitoAuthError("PUBLIC_KEY", "公開キーの取得または処理に失敗しました", cause)
}

func NewJWKSError(cause error) *CognitoAuthError {
	return NewCognitoAuthError("JWKS", "JWKS取得に失敗しました", cause)
}

func NewHTTPError(cause error) *CognitoAuthError {
	return NewCognitoAuthError("HTTP", "HTTP通信に失敗しました", cause)
}

// エラー判定用ヘルパー関数
func IsTokenExpiredError(err error) bool {
	var cognitoErr *CognitoAuthError
	if errors.As(err, &cognitoErr) {
		return errors.Is(cognitoErr.Cause, ErrJWTTokenExpired)
	}
	return errors.Is(err, ErrJWTTokenExpired)
}

func IsInvalidTokenError(err error) bool {
	var cognitoErr *CognitoAuthError
	if errors.As(err, &cognitoErr) {
		return errors.Is(cognitoErr.Cause, ErrJWTTokenInvalid) ||
			errors.Is(cognitoErr.Cause, ErrJWTParsingFailed) ||
			errors.Is(cognitoErr.Cause, ErrJWTSignatureInvalid)
	}
	return errors.Is(err, ErrJWTTokenInvalid) ||
		errors.Is(err, ErrJWTParsingFailed) ||
		errors.Is(err, ErrJWTSignatureInvalid)
}

func IsJWKSError(err error) bool {
	var cognitoErr *CognitoAuthError
	if errors.As(err, &cognitoErr) {
		return cognitoErr.Type == "JWKS" || cognitoErr.Type == "PUBLIC_KEY"
	}
	return false
}

func IsHTTPError(err error) bool {
	var cognitoErr *CognitoAuthError
	if errors.As(err, &cognitoErr) {
		return cognitoErr.Type == "HTTP"
	}
	return false
}
