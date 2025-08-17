package service

import (
	"context"

	roundVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/round"
	sessionVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/session"
	statisticsVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/statistics"
	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
)

// MessagingService はメッセージ送信の抽象化インターフェース
// UseCase層がインフラ層に直接依存することを避けるためのドメインサービス
type MessagingService interface {
	// SendSessionOptimizationMessage はセッション最適化メッセージを送信する
	SendSessionOptimizationMessage(ctx context.Context, data SessionOptimizationData) error
	
	// SendRoundOptimizationMessage はラウンド最適化メッセージを送信する
	SendRoundOptimizationMessage(ctx context.Context, data RoundOptimizationData) error
}

// SessionOptimizationData はセッション最適化メッセージのドメインデータ
// Value Objectsを使用してタイプセーフティを確保
type SessionOptimizationData struct {
	UserID        userVO.UserID                   // ユーザーID（Value Object）
	SessionID     sessionVO.SessionID             // セッションID（Value Object）
	AvgFocusScore statisticsVO.FocusScore         // 平均集中度スコア（Value Object）
	TotalWorkTime statisticsVO.WorkTime           // 合計作業時間（Value Object）
}

// IsValid はデータの有効性を検証する
func (data SessionOptimizationData) IsValid() bool {
	return !data.UserID.IsEmpty() &&
		!data.SessionID.IsEmpty()
	// FocusScoreとWorkTimeは作成時に検証済みのため、追加検証不要
}

// RoundOptimizationData はラウンド最適化メッセージのドメインデータ
// Value Objectsを使用してタイプセーフティを確保
type RoundOptimizationData struct {
	UserID      userVO.UserID               // ユーザーID（Value Object）
	SessionID   sessionVO.SessionID         // セッションID（Value Object）
	RoundID     roundVO.RoundID             // ラウンドID（Value Object）
	FocusScore  statisticsVO.FocusScore     // 集中度スコア（Value Object）
	WorkTime    statisticsVO.WorkTime       // 作業時間（Value Object）
}

// IsValid はデータの有効性を検証する
func (data RoundOptimizationData) IsValid() bool {
	return !data.UserID.IsEmpty() &&
		!data.SessionID.IsEmpty() &&
		!data.RoundID.IsEmpty()
	// FocusScoreとWorkTimeは作成時に検証済みのため、追加検証不要
}