package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// SQSClient はAWS SQSクライアントを管理する高機能クライアント
type SQSClient struct {
	client                 *sqs.Client
	roundOptimizationURL   string
	sessionOptimizationURL string
	logger                 logger.Logger
	maxRetries             int
	retryDelay             time.Duration
	messageTimeout         time.Duration
}

// SQSConfig はSQSクライアントの設定
type SQSConfig struct {
	MaxRetries     int           `default:"3"`
	RetryDelay     time.Duration `default:"1s"`
	MessageTimeout time.Duration `default:"30s"`
}

// NewSQSClient は新しいSQSクライアントを作成する（Lambda最適化版）
func NewSQSClient(cfg *config.Config, logger logger.Logger) (*SQSClient, error) {
	if cfg.SQSRoundOptimizationURL == "" || cfg.SQSSessionOptimizationURL == "" {
		return nil, fmt.Errorf("SQSキューURLが設定されていません")
	}

	// AWS設定の読み込み
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.AWSRegion),
	)
	if err != nil {
		return nil, fmt.Errorf("AWS設定読み込みエラー: %w", err)
	}

	// SQSクライアントの作成
	client := sqs.NewFromConfig(awsCfg)

	logger.Infof("SQSクライアント初期化成功: リージョン=%s", cfg.AWSRegion)
	logger.Infof("ラウンド最適化キュー: %s", cfg.SQSRoundOptimizationURL)
	logger.Infof("セッション最適化キュー: %s", cfg.SQSSessionOptimizationURL)

	return &SQSClient{
		client:                 client,
		roundOptimizationURL:   cfg.SQSRoundOptimizationURL,
		sessionOptimizationURL: cfg.SQSSessionOptimizationURL,
		logger:                 logger,
		maxRetries:             3,
		retryDelay:             1 * time.Second,
		messageTimeout:         30 * time.Second,
	}, nil
}

// SendRoundOptimizationMessage はラウンド最適化メッセージを送信する
func (s *SQSClient) SendRoundOptimizationMessage(ctx context.Context, message *model.RoundOptimizationMessage) error {
	if !message.IsValid() {
		return fmt.Errorf("無効なラウンド最適化メッセージ: %+v", message)
	}

	// JSON形式にシリアライズ
	messageBody, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("メッセージシリアライズエラー: %w", err)
	}

	s.logger.Infof("ラウンド最適化メッセージ送信開始: %s (サイズ: %d bytes)",
		message.ToLogString(), message.GetMessageSize())

	// リトライ付きでメッセージ送信
	err = s.sendMessageWithRetry(ctx, s.roundOptimizationURL, string(messageBody), message.MessageID)
	if err != nil {
		s.logger.Errorf("ラウンド最適化メッセージ送信失敗: %v", err)
		return err
	}

	s.logger.Infof("ラウンド最適化メッセージ送信成功: %s", message.ToLogString())
	return nil
}

// SendSessionOptimizationMessage はセッション最適化メッセージを送信する
func (s *SQSClient) SendSessionOptimizationMessage(ctx context.Context, message *model.SessionOptimizationMessage) error {
	if !message.IsValid() {
		return fmt.Errorf("無効なセッション最適化メッセージ: AvgFocusScore=%.2f, TotalWorkTime=%d",
			message.AvgFocusScore, message.TotalWorkTime)
	}

	// JSON形式にシリアライズ
	messageBody, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("メッセージシリアライズエラー: %w", err)
	}

	s.logger.Infof("セッション最適化メッセージ送信開始: %s (サイズ: %d bytes)",
		message.ToLogString(), message.GetMessageSize())

	// リトライ付きでメッセージ送信
	err = s.sendMessageWithRetry(ctx, s.sessionOptimizationURL, string(messageBody), message.MessageID)
	if err != nil {
		s.logger.Errorf("セッション最適化メッセージ送信失敗: %v", err)
		return err
	}

	s.logger.Infof("セッション最適化メッセージ送信成功: %s", message.ToLogString())
	return nil
}

// sendMessageWithRetry はリトライ機能付きでメッセージを送信する
func (s *SQSClient) sendMessageWithRetry(ctx context.Context, queueURL, messageBody, messageID string) error {
	var lastError error

	for attempt := 1; attempt <= s.maxRetries; attempt++ {
		// タイムアウト付きコンテキストを作成
		timeoutCtx, cancel := context.WithTimeout(ctx, s.messageTimeout)
		defer cancel()

		// SQSメッセージ送信
		input := &sqs.SendMessageInput{
			QueueUrl:    aws.String(queueURL),
			MessageBody: aws.String(messageBody),
			MessageAttributes: map[string]types.MessageAttributeValue{
				"MessageType": {
					DataType:    aws.String("String"),
					StringValue: aws.String("optimization"),
				},
				"Version": {
					DataType:    aws.String("String"),
					StringValue: aws.String("2.0"),
				},
				"Attempt": {
					DataType:    aws.String("Number"),
					StringValue: aws.String(fmt.Sprintf("%d", attempt)),
				},
				"MessageId": {
					DataType:    aws.String("String"),
					StringValue: aws.String(messageID),
				},
			},
		}

		_, err := s.client.SendMessage(timeoutCtx, input)
		if err == nil {
			// 送信成功
			if attempt > 1 {
				s.logger.Infof("SQSメッセージ送信成功（試行回数: %d/%d）", attempt, s.maxRetries)
			}
			return nil
		}

		lastError = err
		s.logger.Warnf("SQSメッセージ送信失敗（試行 %d/%d）: %v", attempt, s.maxRetries, err)

		// 最後の試行でない場合は待機
		if attempt < s.maxRetries {
			select {
			case <-time.After(s.retryDelay):
				// 次の試行まで待機
			case <-ctx.Done():
				return fmt.Errorf("コンテキストキャンセル: %w", ctx.Err())
			}
		}
	}

	return fmt.Errorf("SQSメッセージ送信失敗（%d回試行）: %w", s.maxRetries, lastError)
}

// GetQueueAttributes はキューの属性を取得する（デバッグ用）
func (s *SQSClient) GetQueueAttributes(ctx context.Context, queueURL string) (map[string]string, error) {
	input := &sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(queueURL),
		AttributeNames: []types.QueueAttributeName{
			types.QueueAttributeNameApproximateNumberOfMessages,
			types.QueueAttributeNameApproximateNumberOfMessagesNotVisible,
			types.QueueAttributeNameMessageRetentionPeriod,
		},
	}

	result, err := s.client.GetQueueAttributes(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("キュー属性取得エラー: %w", err)
	}

	return result.Attributes, nil
}

// HealthCheck はSQSクライアントの接続確認を行う
func (s *SQSClient) HealthCheck(ctx context.Context) error {
	// ラウンド最適化キューの確認
	_, err := s.GetQueueAttributes(ctx, s.roundOptimizationURL)
	if err != nil {
		return fmt.Errorf("ラウンド最適化キューへの接続確認失敗: %w", err)
	}

	// セッション最適化キューの確認
	_, err = s.GetQueueAttributes(ctx, s.sessionOptimizationURL)
	if err != nil {
		return fmt.Errorf("セッション最適化キューへの接続確認失敗: %w", err)
	}

	s.logger.Info("SQSクライアント接続確認成功")
	return nil
}

// Close はリソースをクリーンアップする（Lambda用）
func (s *SQSClient) Close() error {
	// SQSクライアントには明示的なCloseメソッドはないが、
	// 将来的な拡張のためにメソッドを定義
	s.logger.Debug("SQSクライアント接続をクリーンアップ")
	return nil
}
