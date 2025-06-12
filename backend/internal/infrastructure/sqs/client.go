package sqs

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// SQSClient はAWS SQSクライアントを管理する高機能クライアント（新エラーハンドリング対応版）
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

// NewSQSClient は新しいSQSクライアントを作成する（新エラーハンドリング対応版）
func NewSQSClient(cfg *config.Config, logger logger.Logger) (*SQSClient, error) {
	// 設定値検証
	if cfg.SQSRoundOptimizationURL == "" || cfg.SQSSessionOptimizationURL == "" {
		logger.Error("SQSキューURLが設定されていません")
		return nil, appErrors.NewConfigMissingError("SQS_QUEUE_URL")
	}

	if cfg.AWSRegion == "" {
		logger.Error("AWSリージョンが設定されていません")
		return nil, appErrors.NewConfigMissingError("AWS_REGION")
	}

	// AWS設定の読み込み
	logger.Infof("SQS初期化開始: リージョン=%s", cfg.AWSRegion)
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.AWSRegion),
	)
	if err != nil {
		logger.Errorf("AWS設定読み込みエラー: %v", err)
		return nil, appErrors.NewConfigError("AWS_CONFIG", err)
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

// SendRoundOptimizationMessage はラウンド最適化メッセージを送信する（新エラーハンドリング対応版）
func (s *SQSClient) SendRoundOptimizationMessage(ctx context.Context, message *model.RoundOptimizationMessage) error {
	// メッセージバリデーション
	if !message.IsValid() {
		s.logger.Errorf("無効なラウンド最適化メッセージ: %+v", message)
		return appErrors.NewSQSError("validate", appErrors.ErrSQSMessageInvalid)
	}

	// JSON形式にシリアライズ
	messageBody, err := json.Marshal(message)
	if err != nil {
		s.logger.Errorf("メッセージシリアライズエラー: %v", err)
		return appErrors.NewSQSError("serialize", err)
	}

	s.logger.Infof("ラウンド最適化メッセージ送信開始: %s (サイズ: %d bytes)",
		message.ToLogString(), message.GetMessageSize())

	// リトライ付きでメッセージ送信
	err = s.sendMessageWithRetry(ctx, s.roundOptimizationURL, string(messageBody), message.MessageID, "round_optimization")
	if err != nil {
		s.logger.Errorf("ラウンド最適化メッセージ送信失敗: %v", err)
		return err // sendMessageWithRetryが適切なSQSエラーを返す
	}

	s.logger.Infof("ラウンド最適化メッセージ送信成功: %s", message.ToLogString())
	return nil
}

// SendSessionOptimizationMessage はセッション最適化メッセージを送信する（新エラーハンドリング対応版）
func (s *SQSClient) SendSessionOptimizationMessage(ctx context.Context, message *model.SessionOptimizationMessage) error {
	// メッセージバリデーション
	if !message.IsValid() {
		s.logger.Errorf("無効なセッション最適化メッセージ: AvgFocusScore=%.2f, TotalWorkTime=%d",
			message.AvgFocusScore, message.TotalWorkTime)
		return appErrors.NewSQSError("validate", appErrors.ErrSQSMessageInvalid)
	}

	// JSON形式にシリアライズ
	messageBody, err := json.Marshal(message)
	if err != nil {
		s.logger.Errorf("メッセージシリアライズエラー: %v", err)
		return appErrors.NewSQSError("serialize", err)
	}

	s.logger.Infof("セッション最適化メッセージ送信開始: %s (サイズ: %d bytes)",
		message.ToLogString(), message.GetMessageSize())

	// リトライ付きでメッセージ送信
	err = s.sendMessageWithRetry(ctx, s.sessionOptimizationURL, string(messageBody), message.MessageID, "session_optimization")
	if err != nil {
		s.logger.Errorf("セッション最適化メッセージ送信失敗: %v", err)
		return err // sendMessageWithRetryが適切なSQSエラーを返す
	}

	s.logger.Infof("セッション最適化メッセージ送信成功: %s", message.ToLogString())
	return nil
}

// sendMessageWithRetry はリトライ機能付きでメッセージを送信する（新エラーハンドリング対応版）
func (s *SQSClient) sendMessageWithRetry(ctx context.Context, queueURL, messageBody, messageID, messageType string) error {
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
					StringValue: aws.String(messageType),
				},
				"Version": {
					DataType:    aws.String("String"),
					StringValue: aws.String("2.0"),
				},
				"Attempt": {
					DataType:    aws.String("String"),
					StringValue: aws.String(string(rune(attempt + '0'))), // 数値を文字列に変換
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

		// エラータイプ判定
		if s.isTimeoutError(timeoutCtx, err) {
			s.logger.Warnf("SQS送信タイムアウト（試行 %d/%d）", attempt, s.maxRetries)
			// タイムアウトの場合は次の試行まで待機
		} else if s.isRetryableError(err) {
			s.logger.Warnf("SQS送信一時的エラー（試行 %d/%d）", attempt, s.maxRetries)
			// リトライ可能エラーの場合は次の試行まで待機
		} else {
			// リトライ不可能エラーの場合は即座に終了
			s.logger.Errorf("SQS送信致命的エラー: %v", err)
			return appErrors.NewSQSSendError(err)
		}

		// 最後の試行でない場合は待機
		if attempt < s.maxRetries {
			select {
			case <-time.After(s.retryDelay):
				// 次の試行まで待機
			case <-ctx.Done():
				s.logger.Error("SQS送信コンテキストキャンセル")
				return appErrors.NewSQSError("context", ctx.Err())
			}
		}
	}

	// 最大試行回数に達した場合
	if s.isTimeoutError(context.Background(), lastError) {
		return appErrors.NewSQSTimeoutError(lastError)
	}

	return appErrors.NewSQSError("retry_exhausted", lastError)
}

// isTimeoutError はタイムアウトエラーかどうかを判定する
func (s *SQSClient) isTimeoutError(ctx context.Context, err error) bool {
	if ctx.Err() == context.DeadlineExceeded {
		return true
	}
	// AWS SDK特有のタイムアウトエラーパターンをチェック
	errStr := err.Error()
	return strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "deadline") ||
		strings.Contains(errStr, "context deadline exceeded")
}

// isRetryableError はリトライ可能エラーかどうかを判定する
func (s *SQSClient) isRetryableError(err error) bool {
	errStr := err.Error()
	// AWS SDK特有のリトライ可能エラーパターンをチェック
	return strings.Contains(errStr, "throttling") ||
		strings.Contains(errStr, "service unavailable") ||
		strings.Contains(errStr, "internal error") ||
		strings.Contains(errStr, "connection") ||
		strings.Contains(errStr, "network")
}

// GetQueueAttributes はキューの属性を取得する（新エラーハンドリング対応版）
func (s *SQSClient) GetQueueAttributes(ctx context.Context, queueURL string) (map[string]string, error) {
	if queueURL == "" {
		return nil, appErrors.NewSQSError("validate", appErrors.ErrSQSMessageInvalid)
	}

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
		s.logger.Errorf("キュー属性取得エラー: %v", err)
		return nil, appErrors.NewSQSError("get_attributes", err)
	}

	return result.Attributes, nil
}

// HealthCheck はSQSクライアントの接続確認を行う（新エラーハンドリング対応版）
func (s *SQSClient) HealthCheck(ctx context.Context) error {
	// ラウンド最適化キューの確認
	_, err := s.GetQueueAttributes(ctx, s.roundOptimizationURL)
	if err != nil {
		s.logger.Errorf("ラウンド最適化キューへの接続確認失敗: %v", err)
		return appErrors.NewSQSConnectionError(err)
	}

	// セッション最適化キューの確認
	_, err = s.GetQueueAttributes(ctx, s.sessionOptimizationURL)
	if err != nil {
		s.logger.Errorf("セッション最適化キューへの接続確認失敗: %v", err)
		return appErrors.NewSQSConnectionError(err)
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

// =================================
// ヘルパーメソッド（デバッグ・監視用）
// =================================

// GetClientInfo はクライアント情報を返す（デバッグ用）
func (s *SQSClient) GetClientInfo() map[string]interface{} {
	return map[string]interface{}{
		"round_optimization_url":   s.roundOptimizationURL,
		"session_optimization_url": s.sessionOptimizationURL,
		"max_retries":              s.maxRetries,
		"retry_delay":              s.retryDelay.String(),
		"message_timeout":          s.messageTimeout.String(),
	}
}

// IsConfigured は設定が正常かどうかを確認する
func (s *SQSClient) IsConfigured() bool {
	return s.client != nil &&
		s.roundOptimizationURL != "" &&
		s.sessionOptimizationURL != ""
}

// =================================
// エラー統計取得（監視用）
// =================================

// GetConnectionStatus は接続状態を確認する（軽量版）
func (s *SQSClient) GetConnectionStatus(ctx context.Context) map[string]bool {
	status := map[string]bool{
		"round_queue_accessible":   false,
		"session_queue_accessible": false,
		"overall_healthy":          false,
	}

	// ラウンドキューの確認
	if _, err := s.GetQueueAttributes(ctx, s.roundOptimizationURL); err == nil {
		status["round_queue_accessible"] = true
	}

	// セッションキューの確認
	if _, err := s.GetQueueAttributes(ctx, s.sessionOptimizationURL); err == nil {
		status["session_queue_accessible"] = true
	}

	// 全体的な健全性
	status["overall_healthy"] = status["round_queue_accessible"] && status["session_queue_accessible"]

	return status
}
