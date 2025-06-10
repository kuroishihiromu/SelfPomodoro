package errors

import (
	"errors"
	"net/http"
)

// ユーザー関連のエラー定義
var (
	ErrUserNotFound       = errors.New("ユーザーが見つかりません")
	ErrUserCreationFailed = errors.New("ユーザーの作成に失敗しました")
	ErrUserUpdateFailed   = errors.New("ユーザーの更新に失敗しました")
	ErrUserDeleteFailed   = errors.New("ユーザーの削除に失敗しました")
	ErrEmailAlreadyExists = errors.New("メールアドレスが既に使用されています")
	ErrInvalidUserData    = errors.New("無効なユーザーデータです")
)

// NewUserNotFoundError はユーザー未検出エラーを作成する
func NewUserNotFoundError() *AppError {
	return &AppError{
		Err:     ErrUserNotFound,
		Message: ErrUserNotFound.Error(),
		Status:  http.StatusNotFound,
	}
}

// NewUserCreationFailedError はユーザー作成失敗エラーを作成する
func NewUserCreationFailedError() *AppError {
	return &AppError{
		Err:     ErrUserCreationFailed,
		Message: ErrUserCreationFailed.Error(),
		Status:  http.StatusInternalServerError,
	}
}

// NewUserUpdateFailedError はユーザー更新失敗エラーを作成する
func NewUserUpdateFailedError() *AppError {
	return &AppError{
		Err:     ErrUserUpdateFailed,
		Message: ErrUserUpdateFailed.Error(),
		Status:  http.StatusInternalServerError,
	}
}

// NewUserDeleteFailedError はユーザー削除失敗エラーを作成する
func NewUserDeleteFailedError() *AppError {
	return &AppError{
		Err:     ErrUserDeleteFailed,
		Message: ErrUserDeleteFailed.Error(),
		Status:  http.StatusInternalServerError,
	}
}

// NewEmailAlreadyExistsError はメールアドレス重複エラーを作成する
func NewEmailAlreadyExistsError() *AppError {
	return &AppError{
		Err:     ErrEmailAlreadyExists,
		Message: ErrEmailAlreadyExists.Error(),
		Status:  http.StatusConflict,
	}
}

// NewInvalidUserDataError は無効なユーザーデータエラーを作成する
func NewInvalidUserDataError() *AppError {
	return &AppError{
		Err:     ErrInvalidUserData,
		Message: ErrInvalidUserData.Error(),
		Status:  http.StatusBadRequest,
	}
}

// IsUserError はユーザー関連のエラーかどうかを判定する
func IsUserError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Err == ErrUserNotFound ||
			appErr.Err == ErrUserCreationFailed ||
			appErr.Err == ErrUserUpdateFailed ||
			appErr.Err == ErrUserDeleteFailed ||
			appErr.Err == ErrEmailAlreadyExists ||
			appErr.Err == ErrInvalidUserData
	}
	return false
}
