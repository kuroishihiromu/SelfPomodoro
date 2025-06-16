package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RoundOptimizationMessage はラウンド最適化用のSQSメッセージを表す
type RoundOptimizationMessage struct {
	MessageID   string    `json:"message_id" validate:"required"`                // メッセージ一意ID
	MessageType string    `json:"message_type" validate:"required"`              // "round_optimization"
	Timestamp   time.Time `json:"timestamp" validate:"required"`                 // メッセージ作成時刻
	Version     string    `json:"version" validate:"required"`                   // メッセージ形式バージョン
	UserID      string    `json:"user_id" validate:"required"`                   // ユーザーID
	RoundID     string    `json:"round_id" validate:"required"`                  // ラウンドID
	FocusScore  int       `json:"focus_score" validate:"required,min=0,max=100"` // 集中度スコア
}

// NewRoundOptimizationMessage は新しいラウンド最適化メッセージを作成する
func NewRoundOptimizationMessage(userID, roundID uuid.UUID, focusScore int) *RoundOptimizationMessage {
	return &RoundOptimizationMessage{
		MessageID:   uuid.New().String(),
		MessageType: "round_optimization",
		Timestamp:   time.Now(),
		Version:     "2.0",
		UserID:      userID.String(),
		RoundID:     roundID.String(),
		FocusScore:  focusScore,
	}
}

// IsValid はメッセージの有効性をチェックする
func (msg *RoundOptimizationMessage) IsValid() bool {
	return msg.MessageID != "" &&
		msg.MessageType == "round_optimization" &&
		msg.Version == "2.0" &&
		msg.UserID != "" &&
		msg.RoundID != "" &&
		msg.FocusScore >= 0 && msg.FocusScore <= 100
}

// GetMessageSize はメッセージの正確なサイズを返す（バイト）
func (msg *RoundOptimizationMessage) GetMessageSize() int {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		// エラー時は概算値を返す
		return len(msg.MessageID) + len(msg.MessageType) + len(msg.Version) +
			len(msg.UserID) + len(msg.RoundID) + 100
	}
	return len(jsonData)
}

// ToLogString はログ出力用の文字列を返す
func (msg *RoundOptimizationMessage) ToLogString() string {
	return fmt.Sprintf("RoundOptimization[ID=%s, UserID=%s, RoundID=%s, FocusScore=%d]",
		msg.MessageID[:8], msg.UserID[:8], msg.RoundID[:8], msg.FocusScore)
}
