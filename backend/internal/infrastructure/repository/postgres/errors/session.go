package errors

import "errors"

var (
	ErrSessionNotFound       = errors.New("セッションが見つかりません")
	ErrSessionAccessDenied   = errors.New("このセッションへのアクセス権限がありません")
	ErrSessionCreationFailed = errors.New("セッションの作成に失敗しました")
	ErrSessionUpdateFailed   = errors.New("セッションの更新に失敗しました")
	ErrSessionDeleteFailed   = errors.New("セッションの削除に失敗しました")
	ErrSessionInProgress     = errors.New("既に進行中のセッションが存在します")
	ErrSessionAlreadyEnded   = errors.New("セッションは既に終了しています")
)
