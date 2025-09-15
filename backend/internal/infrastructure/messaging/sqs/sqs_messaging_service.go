package sqs

import (
	"context"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/service"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/messaging/sqs/message"
)

// SQSMessagingService はMessagingServiceのSQS実装
// ドメインサービスインターフェースを実装してDDDの層分離を実現
type SQSMessagingService struct {
	sqsClient *SQSClient
	logger    logger.Logger
}

// NewSQSMessagingService は新しいSQSメッセージングサービスを作成する
func NewSQSMessagingService(sqsClient *SQSClient, logger logger.Logger) service.MessagingService {
	return &SQSMessagingService{
		sqsClient: sqsClient,
		logger:    logger,
	}
}

// SendSessionOptimizationMessage はセッション最適化メッセージをSQSに送信する
func (s *SQSMessagingService) SendSessionOptimizationMessage(ctx context.Context, data service.SessionOptimizationData) error {
	// データ検証
	if !data.IsValid() {
		s.logger.Errorf("無効なセッション最適化データ: UserID=%s, SessionID=%s, AvgFocus=%s, TotalWork=%s",
			data.UserID.String(), data.SessionID.String(), data.AvgFocusScore.String(), data.TotalWorkTime.String())
		return nil // 無効データは静かに無視（セッション完了を阻害しない）
	}

	// SQSクライアントが利用できない場合
	if s.sqsClient == nil {
		s.logger.Warn("SQSクライアントが初期化されていません。最適化メッセージは送信されません。")
		return nil
	}

	// Value Objects → UUIDに変換してからSQSメッセージ作成
	userUUID, err := uuid.Parse(data.UserID.String())
	if err != nil {
		s.logger.Errorf("UserID変換エラー: %v", err)
		return nil // 変換エラーは静かに無視
	}

	sessionUUID, err := uuid.Parse(data.SessionID.String())
	if err != nil {
		s.logger.Errorf("SessionID変換エラー: %v", err)
		return nil // 変換エラーは静かに無視
	}

	// SQSメッセージを作成
	sqsMessage := message.NewSessionOptimizationMessage(
		userUUID,
		sessionUUID,
		data.AvgFocusScore.Value(),
		data.TotalWorkTime.Minutes(),
	)

	s.logger.Infof("セッション最適化メッセージ作成: UserID=%s, SessionID=%s, AvgFocus=%s, TotalWork=%s",
		data.UserID.String()[:8]+"...", data.SessionID.String()[:8]+"...", 
		data.AvgFocusScore.String(), data.TotalWorkTime.String())

	// SQSに送信
	err = s.sqsClient.SendSessionOptimizationMessage(ctx, sqsMessage)
	if err != nil {
		s.logger.Errorf("セッション最適化メッセージ送信エラー: %v", err)
		return nil // 送信エラーもセッション完了を阻害しない
	}

	s.logger.Infof("セッション最適化メッセージ送信成功: %s", sqsMessage.ToLogString())
	return nil
}

// SendRoundOptimizationMessage はラウンド最適化メッセージをSQSに送信する
func (s *SQSMessagingService) SendRoundOptimizationMessage(ctx context.Context, data service.RoundOptimizationData) error {
	// データ検証
	if !data.IsValid() {
		s.logger.Errorf("無効なラウンド最適化データ: UserID=%s, SessionID=%s, RoundID=%s, FocusScore=%s, WorkTime=%s",
			data.UserID.String(), data.SessionID.String(), data.RoundID.String(), 
			data.FocusScore.String(), data.WorkTime.String())
		return nil // 無効データは静かに無視（ラウンド完了を阻害しない）
	}

	// SQSクライアントが利用できない場合
	if s.sqsClient == nil {
		s.logger.Warn("SQSクライアントが初期化されていません。ラウンド最適化メッセージは送信されません。")
		return nil
	}

	// Value Objects → UUIDに変換してからSQSメッセージ作成
	userUUID, err := uuid.Parse(data.UserID.String())
	if err != nil {
		s.logger.Errorf("UserID変換エラー: %v", err)
		return nil // 変換エラーは静かに無視
	}

	// SessionIDは現在RoundOptimizationMessageで使用しないが、将来の拡張のため保持
	_, err = uuid.Parse(data.SessionID.String())
	if err != nil {
		s.logger.Errorf("SessionID変換エラー: %v", err)
		return nil // 変換エラーは静かに無視
	}

	roundUUID, err := uuid.Parse(data.RoundID.String())
	if err != nil {
		s.logger.Errorf("RoundID変換エラー: %v", err)
		return nil // 変換エラーは静かに無視
	}

	// SQSメッセージを作成
	sqsMessage := message.NewRoundOptimizationMessage(
		userUUID,
		roundUUID,
		int(data.FocusScore.Value()),
	)

	s.logger.Infof("ラウンド最適化メッセージ作成: UserID=%s, SessionID=%s, RoundID=%s, FocusScore=%s, WorkTime=%s",
		data.UserID.String()[:8]+"...", data.SessionID.String()[:8]+"...", data.RoundID.String()[:8]+"...",
		data.FocusScore.String(), data.WorkTime.String())

	// SQSに送信
	err = s.sqsClient.SendRoundOptimizationMessage(ctx, sqsMessage)
	if err != nil {
		s.logger.Errorf("ラウンド最適化メッセージ送信エラー: %v", err)
		return nil // 送信エラーもラウンド完了を阻害しない
	}

	s.logger.Infof("ラウンド最適化メッセージ送信成功: %s", sqsMessage.ToLogString())
	return nil
}