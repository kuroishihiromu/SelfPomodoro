package mapper

import (
	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	"github.com/tsunakit99/selfpomodoro/internal/usecase/dto"
)

// OptimizationPreferencesMapper は最適化設定関連のマッピングを担当する
type OptimizationPreferencesMapper struct{}

// NewOptimizationPreferencesMapper は新しいOptimizationPreferencesMapperを作成する
func NewOptimizationPreferencesMapper() *OptimizationPreferencesMapper {
	return &OptimizationPreferencesMapper{}
}

// ToOptimizationPreferencesResponse はドメインモデルからAPIレスポンス形式に変換する
func (m *OptimizationPreferencesMapper) ToOptimizationPreferencesResponse(prefs *entity.OptimizationPreferences) *dto.OptimizationPreferencesResponse {
	return &dto.OptimizationPreferencesResponse{
		UserID:           prefs.UserID.Value(),
		RoundWorkTime:    prefs.RoundWorkTime.Minutes(),
		RoundBreakTime:   prefs.RoundBreakTime.Minutes(),
		SessionRounds:    prefs.SessionRounds.Count(),
		SessionBreakTime: prefs.SessionBreakTime.Minutes(),
		CreatedAt:        prefs.CreatedAt,
		UpdatedAt:        prefs.UpdatedAt,
	}
}