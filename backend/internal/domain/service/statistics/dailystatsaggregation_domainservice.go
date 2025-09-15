package statistics

import (
	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
)

// DailyStatsAggregationDomainService は日別統計集約のドメインサービス
type DailyStatsAggregationDomainService struct{}

// NewDailyStatsAggregationDomainService は新しいドメインサービスを作成する
func NewDailyStatsAggregationDomainService() *DailyStatsAggregationDomainService {
	return &DailyStatsAggregationDomainService{}
}

// AggregateRoundsIntoDailyStats は複数のRoundデータを日別統計に集約する
func (s *DailyStatsAggregationDomainService) AggregateRoundsIntoDailyStats(userID userVO.UserID, date string, rounds []*entity.Round) *entity.DailyStatistics {
	dailyStats := entity.NewDailyStatistics(userID, date)
	
	for _, round := range rounds {
		if round != nil {
			dailyStats.UpdateWithRound(round)
		}
	}
	
	return dailyStats
}

// ShouldUpdateStatistics は統計を更新すべきかを判定する（複数集約の協調）
func (s *DailyStatsAggregationDomainService) ShouldUpdateStatistics(existingStats *entity.DailyStatistics, newRounds []*entity.Round) bool {
	// 新しいラウンドデータがある場合は更新
	if len(newRounds) > 0 {
		return true
	}
	
	// 既存統計がない場合は更新不要
	if existingStats == nil {
		return false
	}
	
	// データが存在する場合のみ更新対象
	return existingStats.HasData()
}