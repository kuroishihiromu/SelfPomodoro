package message

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

// GetMessageSize はメッセージの正確なサイズを返す（バイト）
func (msg *SessionOptimizationMessage) GetMessageSize() int {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		// エラー時は概算値を返す
		return len(msg.MessageID) + len(msg.MessageType) + len(msg.Version) +
			len(msg.UserID) + len(msg.SessionID) + 80
	}
	return len(jsonData)
}

// ToJSON はメッセージをJSON形式にシリアライズする
func (msg *SessionOptimizationMessage) ToJSON() ([]byte, error) {
	return json.Marshal(msg)
}

// GetRoutingKey はSQSルーティング用のキーを返す
func (msg *SessionOptimizationMessage) GetRoutingKey() string {
	return fmt.Sprintf("optimization.session.%s", msg.UserID)
}

// GetQueueName はキューナメを返す
func (msg *SessionOptimizationMessage) GetQueueName() string {
	return "session-optimization-queue"
}

// ToLogString はログ出力用の文字列を返す（セキュリティ考慮）
func (msg *SessionOptimizationMessage) ToLogString() string {
	userIDShort := msg.UserID
	if len(userIDShort) > 8 {
		userIDShort = userIDShort[:8] + "..."
	}
	
	sessionIDShort := msg.SessionID
	if len(sessionIDShort) > 8 {
		sessionIDShort = sessionIDShort[:8] + "..."
	}
	
	messageIDShort := msg.MessageID
	if len(messageIDShort) > 8 {
		messageIDShort = messageIDShort[:8] + "..."
	}

	return fmt.Sprintf("SessionOptimization[ID=%s, UserID=%s, SessionID=%s, AvgFocus=%.1f, TotalWork=%dmin]",
		messageIDShort, userIDShort, sessionIDShort, msg.AvgFocusScore, msg.TotalWorkTime)
}