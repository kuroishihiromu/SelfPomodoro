package errors

import (
	"errors"
	"net/http"
)

// ユーザー設定関連のエラー定義
var (
	ErrUserConfigNotFound     = errors.New("ユーザー設定が見つかりません")
	ErrUserConfigCreateFailed = errors.New("ユーザー設定の作成に失敗しました")
	ErrUserConfigUpdateFailed = errors.New("ユーザー設定の更新に失敗しました")
	ErrUserConfigDeleteFailed = errors.New("ユーザー設定の削除に失敗しました")
	ErrInvalidUserConfig      = errors.New("無効なユーザー設定です")
)

// NewUserConfigNotFoundError はユーザー設定未検出エラーを作成する
func NewUserConfigNotFoundError() *AppError {
	return &AppError{
		Err:     ErrUserConfigNotFound,
		Message: ErrUserConfigNotFound.Error(),
		Status:  http.StatusNotFound,
	}
}

// NewUserConfigCreateFailedError はユーザー設定作成失敗エラーを作成する
func NewUserConfigCreateFailedError() *AppError {
	return &AppError{
		Err:     ErrUserConfigCreateFailed,
		Message: ErrUserConfigCreateFailed.Error(),
		Status:  http.StatusInternalServerError,
	}
}

// NewUserConfigUpdateFailedError はユーザー設定更新失敗エラーを作成する
func NewUserConfigUpdateFailedError() *AppError {
	return &AppError{
		Err:     ErrUserConfigUpdateFailed,
		Message: ErrUserConfigUpdateFailed.Error(),
		Status:  http.StatusInternalServerError,
	}
}

// NewUserConfigDeleteFailedError はユーザー設定削除失敗エラーを作成する
func NewUserConfigDeleteFailedError() *AppError {
	return &AppError{
		Err:     ErrUserConfigDeleteFailed,
		Message: ErrUserConfigDeleteFailed.Error(),
		Status:  http.StatusInternalServerError,
	}
}

// NewInvalidUserConfigError は無効なユーザー設定エラーを作成する
func NewInvalidUserConfigError() *AppError {
	return &AppError{
		Err:     ErrInvalidUserConfig,
		Message: ErrInvalidUserConfig.Error(),
		Status:  http.StatusBadRequest,
	}
}

// IsUserConfigError はユーザー設定関連のエラーかどうかを判定する
func IsUserConfigError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Err == ErrUserConfigNotFound ||
			appErr.Err == ErrUserConfigCreateFailed ||
			appErr.Err == ErrUserConfigUpdateFailed ||
			appErr.Err == ErrUserConfigDeleteFailed ||
			appErr.Err == ErrInvalidUserConfig
	}
	return false
}
