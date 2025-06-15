package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SessionOptimizationMessage はセッション最適化用のSQSメッセージを表す
type SessionOptimizationMessage struct {
	MessageID     string    `json:"message_id" validate:"required"`                    // メッセージ一意ID
	MessageType   string    `json:"message_type" validate:"required"`                  // "session_optimization"
	Timestamp     time.Time `json:"timestamp" validate:"required"`                     // メッセージ作成時刻
	Version       string    `json:"version" validate:"required"`                       // メッセージ形式バージョン
	UserID        string    `json:"user_id" validate:"required"`                       // ユーザーID
	SessionID     string    `json:"session_id" validate:"required"`                    // セッションID
	AvgFocusScore float64   `json:"avg_focus_score" validate:"required,min=0,max=100"` // 平均集中度スコア
	TotalWorkTime int       `json:"total_work_time" validate:"required,min=0"`         // 合計作業時間（分）
}

// NewSessionOptimizationMessage は新しいセッション最適化メッセージを作成する
func NewSessionOptimizationMessage(userID, sessionID uuid.UUID, avgFocusScore float64, totalWorkTime int) *SessionOptimizationMessage {
	return &SessionOptimizationMessage{
		MessageID:     uuid.New().String(),
		MessageType:   "session_optimization",
		Timestamp:     time.Now(),
		Version:       "2.0",
		UserID:        userID.String(),
		SessionID:     sessionID.String(),
		AvgFocusScore: avgFocusScore,
		TotalWorkTime: totalWorkTime,
	}
}

// IsValid はメッセージの有効性をチェックする
func (msg *SessionOptimizationMessage) IsValid() bool {
	return msg.MessageID != "" &&
		msg.MessageType == "session_optimization" &&
		msg.Version == "2.0" &&
		msg.UserID != "" &&
		msg.SessionID != "" &&
		msg.AvgFocusScore >= 0 && msg.AvgFocusScore <= 100 &&
		msg.TotalWorkTime >= 0
}

// GetMessageSize はメッセージのおおよそのサイズを返す（バイト）
func (msg *SessionOptimizationMessage) GetMessageSize() int {
	// 実際のJSONマーシャルしてサイズを計算
	jsonData, err := json.Marshal(msg)
	if err != nil {
		// エラー時は概算値を返す
		return len(msg.MessageID) + len(msg.MessageType) + len(msg.Version) +
			len(msg.UserID) + len(msg.SessionID) + 80
	}
	return len(jsonData)
}

// ToLogString はログ出力用の文字列を返す
func (msg *SessionOptimizationMessage) ToLogString() string {
	return fmt.Sprintf("SessionOptimization[ID=%s, UserID=%s, SessionID=%s, AvgFocus=%.1f, TotalWork=%dmin]",
		msg.MessageID[:8], msg.UserID[:8], msg.SessionID[:8], msg.AvgFocusScore, msg.TotalWorkTime)
}
