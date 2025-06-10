package errors

import "errors"

var (
	ErrRoundNotFound       = errors.New("ラウンドが見つかりません")
	ErrRoundCreationFailed = errors.New("ラウンドの作成に失敗しました")
	ErrRoundUpdateFailed   = errors.New("ラウンドの更新に失敗しました")
	ErrNoRoundsInSession   = errors.New("セッションにラウンドが存在しません")
	ErrRoundAlreadyEnded   = errors.New("ラウンドは既に終了しています")
)
