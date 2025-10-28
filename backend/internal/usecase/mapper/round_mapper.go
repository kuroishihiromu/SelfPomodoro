package mapper

import (
	"time"

	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	"github.com/tsunakit99/selfpomodoro/internal/usecase/dto"
)

// RoundMapper はラウンド関連のマッピングを担当する
type RoundMapper struct{}

// NewRoundMapper は新しいRoundMapperを作成する
func NewRoundMapper() *RoundMapper {
	return &RoundMapper{}
}

// ToRoundResponse はドメインモデルからAPIレスポンス形式に変換する
func (m *RoundMapper) ToRoundResponse(round *entity.Round) *dto.RoundResponse {
	startTime := round.StartTime.UTC().Truncate(time.Second)

	var endTime *time.Time
	if round.EndTime != nil {
		end := round.EndTime.UTC().Truncate(time.Second)
		endTime = &end
	}

	var workTime *int
	if round.WorkTime != nil {
		minutes := round.WorkTime.Minutes()
		workTime = &minutes
	}

	var breakTime *int
	if round.BreakTime != nil {
		minutes := round.BreakTime.Minutes()
		breakTime = &minutes
	}

	var focusScore *int
	if round.FocusScore != nil {
		score := round.FocusScore.Value()
		focusScore = &score
	}

	return &dto.RoundResponse{
		ID:         round.ID.Value(),
		SessionID:  round.SessionID.Value(),
		RoundOrder: round.RoundOrder.Order(),
		StartTime:  startTime,
		EndTime:    endTime,
		WorkTime:   workTime,
		BreakTime:  breakTime,
		FocusScore: focusScore,
	}
}

// ToRoundsResponse はラウンドのリストをレスポンス形式に変換する
func (m *RoundMapper) ToRoundsResponse(rounds []*entity.Round) *dto.RoundsResponse {
	responses := make([]*dto.RoundResponse, len(rounds))
	for i, round := range rounds {
		responses[i] = m.ToRoundResponse(round)
	}

	return &dto.RoundsResponse{
		Rounds: responses,
	}
}

// ToRoundResponseList はラウンドのリストをレスポンス形式に変換する（複数パターン対応）
func (m *RoundMapper) ToRoundResponseList(rounds []*entity.Round) []*dto.RoundResponse {
	responses := make([]*dto.RoundResponse, len(rounds))
	for i, round := range rounds {
		responses[i] = m.ToRoundResponse(round)
	}
	return responses
}
