package errors

import "errors"

var (
	ErrUserNotFound       = errors.New("ユーザーが見つかりません")
	ErrUserCreationFailed = errors.New("ユーザーの作成に失敗しました")
	ErrUserUpdateFailed   = errors.New("ユーザーの更新に失敗しました")
	ErrUserDeleteFailed   = errors.New("ユーザーの削除に失敗しました")
	ErrEmailAlreadyExists = errors.New("メールアドレスが既に使用されています")
)
