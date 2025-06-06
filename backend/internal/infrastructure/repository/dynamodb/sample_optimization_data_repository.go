package dynamodb

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// SampleOptimizationDataRepositoryImpl はDynamoDBを使用したSampleOptimizationDataRepositoryの実装
type SampleOptimizationDataRepositoryImpl struct {
	client                   *dynamodb.Client
	roundOptimizationTable   string
	sessionOptimizationTable string
	logger                   logger.Logger
}

// 固定サンプルデータ定義
var (
	// 10日間のサンプルパターン（日数, セッション数, 基準集中度）
	sampleDayPatterns = []struct {
		day              int
		sessions         int
		baseFocus        float64
		roundsPerSession int
	}{
		{0, 1, 65.0, 3}, // 1日目: 1セッション, 集中度65, 3ラウンド
		{1, 2, 68.0, 3}, // 2日目: 2セッション, 集中度68, 3ラウンド
		{2, 1, 62.0, 2}, // 3日目: 調子悪い日
		{3, 2, 72.0, 3}, // 4日目: 改善
		{4, 3, 75.0, 4}, // 5日目: 好調
		{5, 2, 78.0, 3}, // 6日目: 安定
		{6, 1, 70.0, 3}, // 7日目: 休息日
		{7, 3, 82.0, 4}, // 8日目: とても好調
		{8, 2, 85.0, 3}, // 9日目: 維持
		{9, 2, 88.0, 4}, // 10日目: 最高調
	}

	// ラウンドパターン（ラウンド順, 作業時間, 休憩時間, 集中度修正値）
	roundPatterns = []struct {
		order       int
		workTime    int
		breakTime   int
		focusAdjust int
	}{
		{0, 25, 5, +5},  // 1ラウンド目: 集中しやすい
		{1, 25, 5, 0},   // 2ラウンド目: 標準
		{2, 25, 10, -3}, // 3ラウンド目: 少し疲れ
		{3, 20, 15, -5}, // 4ラウンド目: 短めで調整
	}
)

// NewSampleOptimizationDataRepository は新しいSampleOptimizationDataRepositoryImplを作成する
func NewSampleOptimizationDataRepository(client *dynamodb.Client, cfg *config.Config, logger logger.Logger) repository.SampleOptimizationDataRepository {
	return &SampleOptimizationDataRepositoryImpl{
		client:                   client,
		roundOptimizationTable:   cfg.DynamoRoundOptimizationTable,
		sessionOptimizationTable: cfg.DynamoSessionOptimizationTable,
		logger:                   logger,
	}
}

// CreateSampleOptimizationData はサンプル最適化データを作成する（固定パターン）
func (r *SampleOptimizationDataRepositoryImpl) CreateSampleOptimizationData(ctx context.Context, userID uuid.UUID) error {
	r.logger.Infof("サンプル最適化データ作成開始: UserID=%s", userID.String()[:8]+"...")

	// 固定パターンからサンプルデータを生成
	roundLogs, sessionLogs := r.generateSampleData(userID)

	// ラウンド最適化ログを一括作成
	err := r.CreateRoundOptimizationLogs(ctx, roundLogs)
	if err != nil {
		return fmt.Errorf("ラウンド最適化ログ作成失敗: %w", err)
	}

	// セッション最適化ログを一括作成
	err = r.CreateSessionOptimizationLogs(ctx, sessionLogs)
	if err != nil {
		return fmt.Errorf("セッション最適化ログ作成失敗: %w", err)
	}

	r.logger.Infof("サンプル最適化データ作成完了: ラウンド=%d件, セッション=%d件",
		len(roundLogs), len(sessionLogs))
	return nil
}

// generateSampleData は固定パターンからサンプルデータを生成する
func (r *SampleOptimizationDataRepositoryImpl) generateSampleData(userID uuid.UUID) ([]*model.RoundOptimizationLog, []*model.SessionOptimizationLog) {
	var roundLogs []*model.RoundOptimizationLog
	var sessionLogs []*model.SessionOptimizationLog

	baseTime := time.Now().AddDate(0, 0, -10) // 10日前から開始

	// 各日のパターンを処理
	for _, dayPattern := range sampleDayPatterns {
		currentDate := baseTime.AddDate(0, 0, dayPattern.day)

		// 1日のセッション数分ループ
		for sessionNum := 0; sessionNum < dayPattern.sessions; sessionNum++ {
			// セッション開始時刻（9時、13時、17時）
			sessionStartTime := currentDate.Add(time.Duration(9+sessionNum*4) * time.Hour)

			// セッション最適化ログ作成
			sessionLog := r.generateSessionLog(userID, sessionStartTime, dayPattern)
			sessionLogs = append(sessionLogs, sessionLog)

			// このセッションのラウンド最適化ログ作成
			for roundNum := 0; roundNum < dayPattern.roundsPerSession; roundNum++ {
				if roundNum >= len(roundPatterns) {
					break // パターン数を超えた場合は終了
				}

				roundStartTime := sessionStartTime.Add(time.Duration(roundNum*30) * time.Minute)
				roundLog := r.generateRoundLog(userID, roundStartTime, dayPattern, roundNum)
				roundLogs = append(roundLogs, roundLog)
			}
		}
	}

	return roundLogs, sessionLogs
}

// generateSessionLog はセッション最適化ログを生成する
func (r *SampleOptimizationDataRepositoryImpl) generateSessionLog(userID uuid.UUID, timestamp time.Time, dayPattern struct {
	day              int
	sessions         int
	baseFocus        float64
	roundsPerSession int
}) *model.SessionOptimizationLog {

	// セッション後の設定値（次回セッション用）
	roundCount := dayPattern.roundsPerSession
	breakTime := 15 // 基本15分長休憩

	// 後半はより長い休憩
	if dayPattern.day > 5 {
		breakTime = 20
	}

	totalWorkTime := roundCount * 25 // 基本25分×ラウンド数

	return model.NewSessionOptimizationLogWithTime(
		userID,
		timestamp,
		roundCount,
		breakTime,
		dayPattern.baseFocus,
		totalWorkTime,
	)
}

// generateRoundLog はラウンド最適化ログを生成する
func (r *SampleOptimizationDataRepositoryImpl) generateRoundLog(userID uuid.UUID, timestamp time.Time, dayPattern struct {
	day              int
	sessions         int
	baseFocus        float64
	roundsPerSession int
}, roundNum int) *model.RoundOptimizationLog {

	pattern := roundPatterns[roundNum]

	// 作業時間・休憩時間（次回ラウンド用の最適化値）
	workTime := pattern.workTime
	breakTime := pattern.breakTime

	// 経験による微調整（後半の日では最適化が進む）
	if dayPattern.day > 5 {
		if roundNum == 0 {
			workTime = 30 // 最初のラウンドは30分に延長
		} else if roundNum >= 3 {
			workTime = 20 // 後半は短めに調整
		}
	}

	// 集中度スコア（このラウンドの実績値）
	focusScore := int(dayPattern.baseFocus) + pattern.focusAdjust

	// 値の範囲チェック
	if focusScore > 100 {
		focusScore = 100
	}
	if focusScore < 30 {
		focusScore = 30
	}

	return model.NewRoundOptimizationLogWithTime(
		userID,
		timestamp,
		workTime,
		breakTime,
		focusScore,
	)
}

// CreateRoundOptimizationLogs はラウンド最適化ログを一括作成する
func (r *SampleOptimizationDataRepositoryImpl) CreateRoundOptimizationLogs(ctx context.Context, logs []*model.RoundOptimizationLog) error {
	if len(logs) == 0 {
		return nil
	}

	r.logger.Infof("ラウンド最適化ログ一括作成開始: %d件", len(logs))

	// BatchWriteItemは最大25件なので、分割して処理
	batchSize := 25
	for i := 0; i < len(logs); i += batchSize {
		end := i + batchSize
		if end > len(logs) {
			end = len(logs)
		}

		batch := logs[i:end]
		err := r.batchWriteRoundLogs(ctx, batch)
		if err != nil {
			return fmt.Errorf("ラウンドログバッチ作成失敗: %w", err)
		}
	}

	r.logger.Infof("ラウンド最適化ログ一括作成完了: %d件", len(logs))
	return nil
}

// CreateSessionOptimizationLogs はセッション最適化ログを一括作成する
func (r *SampleOptimizationDataRepositoryImpl) CreateSessionOptimizationLogs(ctx context.Context, logs []*model.SessionOptimizationLog) error {
	if len(logs) == 0 {
		return nil
	}

	r.logger.Infof("セッション最適化ログ一括作成開始: %d件", len(logs))

	// BatchWriteItemは最大25件なので、分割して処理
	batchSize := 25
	for i := 0; i < len(logs); i += batchSize {
		end := i + batchSize
		if end > len(logs) {
			end = len(logs)
		}

		batch := logs[i:end]
		err := r.batchWriteSessionLogs(ctx, batch)
		if err != nil {
			return fmt.Errorf("セッションログバッチ作成失敗: %w", err)
		}
	}

	r.logger.Infof("セッション最適化ログ一括作成完了: %d件", len(logs))
	return nil
}

// batchWriteRoundLogs はラウンドログをバッチ書き込みする
func (r *SampleOptimizationDataRepositoryImpl) batchWriteRoundLogs(ctx context.Context, logs []*model.RoundOptimizationLog) error {
	var writeRequests []types.WriteRequest

	for _, log := range logs {
		item := map[string]types.AttributeValue{
			"user_id":     &types.AttributeValueMemberS{Value: log.UserID},
			"timestamp":   &types.AttributeValueMemberS{Value: log.Timestamp},
			"work_time":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", log.WorkTime)},
			"break_time":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", log.BreakTime)},
			"focus_score": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", log.FocusScore)},
			"created_at":  &types.AttributeValueMemberS{Value: log.CreatedAt},
		}

		writeRequests = append(writeRequests, types.WriteRequest{
			PutRequest: &types.PutRequest{Item: item},
		})
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			r.roundOptimizationTable: writeRequests,
		},
	}

	_, err := r.client.BatchWriteItem(ctx, input)
	return err
}

// batchWriteSessionLogs はセッションログをバッチ書き込みする
func (r *SampleOptimizationDataRepositoryImpl) batchWriteSessionLogs(ctx context.Context, logs []*model.SessionOptimizationLog) error {
	var writeRequests []types.WriteRequest

	for _, log := range logs {
		item := map[string]types.AttributeValue{
			"user_id":         &types.AttributeValueMemberS{Value: log.UserID},
			"timestamp":       &types.AttributeValueMemberS{Value: log.Timestamp},
			"round_count":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", log.RoundCount)},
			"break_time":      &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", log.BreakTime)},
			"avg_focus_score": &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", log.AvgFocusScore)},
			"total_work_time": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", log.TotalWorkTime)},
			"created_at":      &types.AttributeValueMemberS{Value: log.CreatedAt},
		}

		writeRequests = append(writeRequests, types.WriteRequest{
			PutRequest: &types.PutRequest{Item: item},
		})
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			r.sessionOptimizationTable: writeRequests,
		},
	}

	_, err := r.client.BatchWriteItem(ctx, input)
	return err
}
