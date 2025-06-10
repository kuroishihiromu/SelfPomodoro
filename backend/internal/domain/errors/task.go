package errors

import (
	"errors"
	"net/http"
)

// タスク関連のエラー定義
var (
	ErrTaskNotFound     = errors.New("タスクが見つかりません")
	ErrTaskAlreadyDone  = errors.New("タスクは既に完了しています")
	ErrInvalidTaskState = errors.New("タスクの状態が無効です")
	ErrTaskInProgress   = errors.New("既に進行中のタスクが存在します")
	ErrInvalidTaskType  = errors.New("無効なタスクタイプです")
)

// NewTaskNotFoundError はタスク未検出エラーを作成する
func NewTaskNotFoundError() *AppError {
	return &AppError{
		Err:     ErrTaskNotFound,
		Message: ErrTaskNotFound.Error(),
		Status:  http.StatusNotFound,
	}
}

// NewTaskAlreadyDoneError はタスク完了済みエラーを作成する
func NewTaskAlreadyDoneError() *AppError {
	return &AppError{
		Err:     ErrTaskAlreadyDone,
		Message: ErrTaskAlreadyDone.Error(),
		Status:  http.StatusBadRequest,
	}
}

// NewInvalidTaskStateError は無効なタスク状態エラーを作成する
func NewInvalidTaskStateError() *AppError {
	return &AppError{
		Err:     ErrInvalidTaskState,
		Message: ErrInvalidTaskState.Error(),
		Status:  http.StatusBadRequest,
	}
}

// NewTaskInProgressError は進行中タスク存在エラーを作成する
func NewTaskInProgressError() *AppError {
	return &AppError{
		Err:     ErrTaskInProgress,
		Message: ErrTaskInProgress.Error(),
		Status:  http.StatusConflict,
	}
}

// NewInvalidTaskTypeError は無効なタスクタイプエラーを作成する
func NewInvalidTaskTypeError() *AppError {
	return &AppError{
		Err:     ErrInvalidTaskType,
		Message: ErrInvalidTaskType.Error(),
		Status:  http.StatusBadRequest,
	}
}

// IsTaskError はタスク関連のエラーかどうかを判定する
func IsTaskError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Err == ErrTaskNotFound ||
			appErr.Err == ErrTaskAlreadyDone ||
			appErr.Err == ErrInvalidTaskState ||
			appErr.Err == ErrTaskInProgress ||
			appErr.Err == ErrInvalidTaskType
	}
	return false
}
