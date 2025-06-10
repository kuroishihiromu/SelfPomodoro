package errors

import "errors"

var (
	ErrSampleOptimizationDataNotFound     = errors.New("サンプル最適化データが見つかりません")
	ErrSampleOptimizationDataCreateFailed = errors.New("サンプル最適化データの作成に失敗しました")
	ErrSampleOptimizationDataUpdateFailed = errors.New("サンプル最適化データの更新に失敗しました")
	ErrSampleOptimizationDataDeleteFailed = errors.New("サンプル最適化データの削除に失敗しました")
)
