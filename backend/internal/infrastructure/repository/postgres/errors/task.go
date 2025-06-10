package errors

import "errors"

var (
	ErrTaskNotFound       = errors.New("タスクが見つかりません")
	ErrTaskAccessDenied   = errors.New("このタスクへのアクセス権限がありません")
	ErrTaskCreationFailed = errors.New("タスクの作成に失敗しました")
	ErrTaskUpdateFailed   = errors.New("タスクの更新に失敗しました")
	ErrTaskDeleteFailed   = errors.New("タスクの削除に失敗しました")
	ErrTaskAlreadyExists  = errors.New("タスクが既に存在します")
)
