package errors

import (
	"errors"
	"net/http"
)

// ラウンド関連のエラー定義
var (
	ErrRoundNotFound     = errors.New("ラウンドが見つかりません")
	ErrRoundAlreadyEnded = errors.New("ラウンドは既に終了しています")
	ErrInvalidRoundState = errors.New("ラウンドの状態が無効です")
	ErrRoundInProgress   = errors.New("既に進行中のラウンドが存在します")
	ErrInvalidRoundType  = errors.New("無効なラウンドタイプです")
)

// NewRoundNotFoundError はラウンド未検出エラーを作成する
func NewRoundNotFoundError() *AppError {
	return &AppError{
		Err:     ErrRoundNotFound,
		Message: ErrRoundNotFound.Error(),
		Status:  http.StatusNotFound,
	}
}

// NewRoundAlreadyEndedError はラウンド終了済みエラーを作成する
func NewRoundAlreadyEndedError() *AppError {
	return &AppError{
		Err:     ErrRoundAlreadyEnded,
		Message: ErrRoundAlreadyEnded.Error(),
		Status:  http.StatusBadRequest,
	}
}

// NewInvalidRoundStateError は無効なラウンド状態エラーを作成する
func NewInvalidRoundStateError() *AppError {
	return &AppError{
		Err:     ErrInvalidRoundState,
		Message: ErrInvalidRoundState.Error(),
		Status:  http.StatusBadRequest,
	}
}

// NewRoundInProgressError は進行中ラウンド存在エラーを作成する
func NewRoundInProgressError() *AppError {
	return &AppError{
		Err:     ErrRoundInProgress,
		Message: ErrRoundInProgress.Error(),
		Status:  http.StatusConflict,
	}
}

// NewInvalidRoundTypeError は無効なラウンドタイプエラーを作成する
func NewInvalidRoundTypeError() *AppError {
	return &AppError{
		Err:     ErrInvalidRoundType,
		Message: ErrInvalidRoundType.Error(),
		Status:  http.StatusBadRequest,
	}
}

// IsRoundError はラウンド関連のエラーかどうかを判定する
func IsRoundError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Err == ErrRoundNotFound ||
			appErr.Err == ErrRoundAlreadyEnded ||
			appErr.Err == ErrInvalidRoundState ||
			appErr.Err == ErrRoundInProgress ||
			appErr.Err == ErrInvalidRoundType
	}
	return false
}
