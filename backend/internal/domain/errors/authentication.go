package errors

import (
	"errors"
	"net/http"
)

// 認証関連のエラー定義
var (
	ErrTokenExpired     = errors.New("JWTトークンの有効期限が切れています")
	ErrInvalidToken     = errors.New("無効なJWTトークンです")
	ErrTokenNotFound    = errors.New("Authorizationヘッダーが見つかりません")
	ErrInvalidFormat    = errors.New("Authorizationヘッダーの形式が無効です")
	ErrInvalidIssuer    = errors.New("トークンの発行者が無効です")
	ErrInvalidAudience  = errors.New("トークンのaudienceが無効です")
	ErrInvalidSignature = errors.New("トークンの署名が無効です")
	ErrMissingSubject   = errors.New("subjectクレームが見つかりません")
	ErrInvalidSubject   = errors.New("subjectクレームが無効です")
	ErrMissingTokenUse  = errors.New("token_useクレームが見つかりません")
	ErrInvalidTokenUse  = errors.New("token_useクレームが無効です")
)

// NewTokenExpiredError はトークン有効期限切れエラーを作成する
func NewTokenExpiredError() *AppError {
	return &AppError{
		Err:     ErrTokenExpired,
		Message: ErrTokenExpired.Error(),
		Status:  http.StatusUnauthorized,
	}
}

// NewInvalidTokenError は無効なトークンエラーを作成する
func NewInvalidTokenError() *AppError {
	return &AppError{
		Err:     ErrInvalidToken,
		Message: ErrInvalidToken.Error(),
		Status:  http.StatusUnauthorized,
	}
}

// NewTokenNotFoundError はトークン未検出エラーを作成する
func NewTokenNotFoundError() *AppError {
	return &AppError{
		Err:     ErrTokenNotFound,
		Message: ErrTokenNotFound.Error(),
		Status:  http.StatusUnauthorized,
	}
}

// NewInvalidFormatError は無効なフォーマットエラーを作成する
func NewInvalidFormatError() *AppError {
	return &AppError{
		Err:     ErrInvalidFormat,
		Message: ErrInvalidFormat.Error(),
		Status:  http.StatusBadRequest,
	}
}

// NewInvalidIssuerError は無効な発行者エラーを作成する
func NewInvalidIssuerError() *AppError {
	return &AppError{
		Err:     ErrInvalidIssuer,
		Message: ErrInvalidIssuer.Error(),
		Status:  http.StatusUnauthorized,
	}
}

// NewInvalidAudienceError は無効なオーディエンスエラーを作成する
func NewInvalidAudienceError() *AppError {
	return &AppError{
		Err:     ErrInvalidAudience,
		Message: ErrInvalidAudience.Error(),
		Status:  http.StatusUnauthorized,
	}
}

// NewInvalidSignatureError は無効な署名エラーを作成する
func NewInvalidSignatureError() *AppError {
	return &AppError{
		Err:     ErrInvalidSignature,
		Message: ErrInvalidSignature.Error(),
		Status:  http.StatusUnauthorized,
	}
}

// NewMissingSubjectError はsubjectクレーム未検出エラーを作成する
func NewMissingSubjectError() *AppError {
	return &AppError{
		Err:     ErrMissingSubject,
		Message: ErrMissingSubject.Error(),
		Status:  http.StatusUnauthorized,
	}
}

// NewInvalidSubjectError は無効なsubjectクレームエラーを作成する
func NewInvalidSubjectError() *AppError {
	return &AppError{
		Err:     ErrInvalidSubject,
		Message: ErrInvalidSubject.Error(),
		Status:  http.StatusUnauthorized,
	}
}

// NewMissingTokenUseError はtoken_useクレーム未検出エラーを作成する
func NewMissingTokenUseError() *AppError {
	return &AppError{
		Err:     ErrMissingTokenUse,
		Message: ErrMissingTokenUse.Error(),
		Status:  http.StatusUnauthorized,
	}
}

// NewInvalidTokenUseError は無効なtoken_useクレームエラーを作成する
func NewInvalidTokenUseError() *AppError {
	return &AppError{
		Err:     ErrInvalidTokenUse,
		Message: ErrInvalidTokenUse.Error(),
		Status:  http.StatusUnauthorized,
	}
}

// IsAuthenticationError は認証関連のエラーかどうかを判定する
func IsAuthenticationError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Status == http.StatusUnauthorized
	}
	return false
}
