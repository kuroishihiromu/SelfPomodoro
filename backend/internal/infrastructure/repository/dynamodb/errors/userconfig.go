package errors

import "errors"

var (
	ErrUserConfigNotFound     = errors.New("ユーザー設定が見つかりません")
	ErrUserConfigCreateFailed = errors.New("ユーザー設定の作成に失敗しました")
	ErrUserConfigUpdateFailed = errors.New("ユーザー設定の更新に失敗しました")
	ErrUserConfigDeleteFailed = errors.New("ユーザー設定の削除に失敗しました")
)
