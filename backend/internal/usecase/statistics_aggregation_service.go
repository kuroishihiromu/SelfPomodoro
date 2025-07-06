package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	dynamo "github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/dynamodb"
)

// StatisticsAggregationService は統計データの事前集約を行うサービス
type StatisticsAggregationService interface {
	// UpdateStatisticsOnRoundComplete ラウンド完了時に統計データを更新する
	UpdateStatisticsOnRoundComplete(ctx context.Context, userID uuid.UUID, round *entity.Round) error
}

// statisticsAggregationService は統計データ事前集約サービスの実装
type statisticsAggregationService struct {
	dynamoStatsRepo *dynamo.StatisticsRepositoryImpl
	logger          logger.Logger
}

// NewStatisticsAggregationService は新しい統計集約サービスを作成する
func NewStatisticsAggregationService(dynamoStatsRepo *dynamo.StatisticsRepositoryImpl, logger logger.Logger) StatisticsAggregationService {
	return &statisticsAggregationService{
		dynamoStatsRepo: dynamoStatsRepo,
		logger:          logger,
	}
}

// UpdateStatisticsOnRoundComplete はラウンド完了時に統計データを更新する
func (s *statisticsAggregationService) UpdateStatisticsOnRoundComplete(ctx context.Context, userID uuid.UUID, round *entity.Round) error {
	// 集中度スコアのないラウンドは統計から除外
	if round.FocusScore == nil {
		s.logger.Debugf("統計更新スキップ: ラウンドID=%s (集中度スコアなし)", round.ID.String())
		return nil
	}

	// ラウンドデータをDynamoDB形式に変換
	dynamoRound := s.convertRoundToDynamoDBFormat(userID.String(), round)

	// 日付と時間を取得
	date := round.StartTime.Format("2006-01-02")
	hour := round.StartTime.Hour()

	// 日別統計を更新
	if err := s.dynamoStatsRepo.UpdateDailyStats(ctx, userID.String(), date, dynamoRound); err != nil {
		s.logger.Errorf("日別統計更新エラー: %v", err)
		return err
	}

	// 時間別統計を更新
	if err := s.dynamoStatsRepo.UpdateHourlyStats(ctx, userID.String(), date, hour, dynamoRound); err != nil {
		s.logger.Errorf("時間別統計更新エラー: %v", err)
		return err
	}

	// 週別統計を更新
	if err := s.dynamoStatsRepo.UpdateWeeklyStats(ctx, userID.String(), date, dynamoRound); err != nil {
		s.logger.Errorf("週別統計更新エラー: %v", err)
		return err
	}

	s.logger.Infof("統計データ更新成功: UserID=%s, Date=%s, Hour=%d, FocusScore=%d",
		userID.String(), date, hour, *round.FocusScore)

	return nil
}

// convertRoundToDynamoDBFormat はRoundモデルをDynamoDB形式に変換する（統計更新用）
func (s *statisticsAggregationService) convertRoundToDynamoDBFormat(_ string, round *entity.Round) *entity.Round {
	// 統計処理では直接Roundモデルを使用するため、変換は不要
	return round
}

// StatisticsAggregationServiceFactory は統計集約サービスのファクトリ
type StatisticsAggregationServiceFactory struct {
	dynamoStatsRepo *dynamo.StatisticsRepositoryImpl
	logger          logger.Logger
}

// NewStatisticsAggregationServiceFactory は新しいファクトリを作成する
func NewStatisticsAggregationServiceFactory(dynamoStatsRepo *dynamo.StatisticsRepositoryImpl, logger logger.Logger) *StatisticsAggregationServiceFactory {
	return &StatisticsAggregationServiceFactory{
		dynamoStatsRepo: dynamoStatsRepo,
		logger:          logger,
	}
}

// CreateService は統計集約サービスを作成する
func (f *StatisticsAggregationServiceFactory) CreateService() StatisticsAggregationService {
	if f.dynamoStatsRepo == nil {
		f.logger.Warn("DynamoDB統計リポジトリが利用できないため、統計集約サービスは無効化されます")
		return &noOpStatisticsAggregationService{logger: f.logger}
	}

	return NewStatisticsAggregationService(f.dynamoStatsRepo, f.logger)
}

// noOpStatisticsAggregationService は何もしない統計集約サービス（フォールバック用）
type noOpStatisticsAggregationService struct {
	logger logger.Logger
}

// UpdateStatisticsOnRoundComplete は何もしない実装
func (s *noOpStatisticsAggregationService) UpdateStatisticsOnRoundComplete(ctx context.Context, userID uuid.UUID, round *entity.Round) error {
	s.logger.Debug("統計集約サービスが無効化されているため、統計更新をスキップします")
	return nil
}
