package errors

import (
	"errors"
	"net/http"
)

// セッション関連のエラー定義
var (
	ErrSessionNotFound     = errors.New("セッションが見つかりません")
	ErrSessionAlreadyEnded = errors.New("セッションは既に終了しています")
	ErrInvalidSessionState = errors.New("セッションの状態が無効です")
	ErrSessionInProgress   = errors.New("既に進行中のセッションが存在します")
)

// NewSessionNotFoundError はセッション未検出エラーを作成する
func NewSessionNotFoundError() *AppError {
	return &AppError{
		Err:     ErrSessionNotFound,
		Message: ErrSessionNotFound.Error(),
		Status:  http.StatusNotFound,
	}
}

// NewSessionAlreadyEndedError はセッション終了済みエラーを作成する
func NewSessionAlreadyEndedError() *AppError {
	return &AppError{
		Err:     ErrSessionAlreadyEnded,
		Message: ErrSessionAlreadyEnded.Error(),
		Status:  http.StatusBadRequest,
	}
}

// NewInvalidSessionStateError は無効なセッション状態エラーを作成する
func NewInvalidSessionStateError() *AppError {
	return &AppError{
		Err:     ErrInvalidSessionState,
		Message: ErrInvalidSessionState.Error(),
		Status:  http.StatusBadRequest,
	}
}

// NewSessionInProgressError は進行中セッション存在エラーを作成する
func NewSessionInProgressError() *AppError {
	return &AppError{
		Err:     ErrSessionInProgress,
		Message: ErrSessionInProgress.Error(),
		Status:  http.StatusConflict,
	}
}

// IsSessionError はセッション関連のエラーかどうかを判定する
func IsSessionError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Err == ErrSessionNotFound ||
			appErr.Err == ErrSessionAlreadyEnded ||
			appErr.Err == ErrInvalidSessionState ||
			appErr.Err == ErrSessionInProgress
	}
	return false
}
