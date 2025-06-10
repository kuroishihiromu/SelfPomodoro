package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// 共通エラー定義
var (
	ErrInternalServer = errors.New("内部サーバーエラーが発生しました")
	ErrNotFound       = errors.New("リソースが見つかりません")
	ErrBadRequest     = errors.New("不正なリクエストです")
	ErrUnauthorized   = errors.New("認証が必要です")
	ErrForbidden      = errors.New("アクセス権限がありません")
	ErrConflict       = errors.New("リソースが既に存在します")
	ErrValidation     = errors.New("バリデーションエラーが発生しました")
)

// AppError はアプリケーションエラーを表す
type AppError struct {
	Err     error
	Message string
	Status  int
	Cause   error
}

// Error はエラーメッセージを返す
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap は元のエラーを返す
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewInternalError は内部サーバーエラーを作成する
func NewInternalError(cause error) *AppError {
	return &AppError{
		Err:     ErrInternalServer,
		Message: ErrInternalServer.Error(),
		Status:  http.StatusInternalServerError,
		Cause:   cause,
	}
}

// NewNotFoundError はリソースが見つからないエラーを作成する
func NewNotFoundError(message string) *AppError {
	if message == "" {
		message = ErrNotFound.Error()
	}
	return &AppError{
		Err:     ErrNotFound,
		Message: message,
		Status:  http.StatusNotFound,
	}
}

// NewBadRequestError は不正なリクエストエラーを作成する
func NewBadRequestError(message string) *AppError {
	if message == "" {
		message = ErrBadRequest.Error()
	}
	return &AppError{
		Err:     ErrBadRequest,
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

// NewUnauthorizedError は認証エラーを作成する
func NewUnauthorizedError(message string) *AppError {
	if message == "" {
		message = ErrUnauthorized.Error()
	}
	return &AppError{
		Err:     ErrUnauthorized,
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

// NewForbiddenError はアクセス権限エラーを作成する
func NewForbiddenError(message string) *AppError {
	if message == "" {
		message = ErrForbidden.Error()
	}
	return &AppError{
		Err:     ErrForbidden,
		Message: message,
		Status:  http.StatusForbidden,
	}
}

// NewConflictError はリソース競合エラーを作成する
func NewConflictError(message string) *AppError {
	if message == "" {
		message = ErrConflict.Error()
	}
	return &AppError{
		Err:     ErrConflict,
		Message: message,
		Status:  http.StatusConflict,
	}
}

// NewInternalServerError は内部サーバーエラーを作成する
func NewInternalServerError() *AppError {
	return &AppError{
		Err:     ErrInternalServer,
		Message: ErrInternalServer.Error(),
		Status:  http.StatusInternalServerError,
	}
}

// NewValidationError はバリデーションエラーを作成する
func NewValidationError(message string) *AppError {
	if message == "" {
		message = ErrValidation.Error()
	}
	return &AppError{
		Err:     ErrValidation,
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

// Is はエラーの種類を判定する
func Is(err, target error) bool {
	return errors.Is(err, target)
}
