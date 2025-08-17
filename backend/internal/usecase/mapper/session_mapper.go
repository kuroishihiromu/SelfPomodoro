package mapper

import (
	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	"github.com/tsunakit99/selfpomodoro/internal/usecase/dto"
)

// SessionMapper はセッション関連のマッピングを担当する
type SessionMapper struct{}

// NewSessionMapper は新しいSessionMapperを作成する
func NewSessionMapper() *SessionMapper {
	return &SessionMapper{}
}

// ToSessionResponse はドメインモデルからAPIレスポンス形式に変換する
func (m *SessionMapper) ToSessionResponse(session *entity.Session) *dto.SessionResponse {
	var averageFocus *float64
	if session.AverageFocus != nil {
		score := session.AverageFocus.Score()
		averageFocus = &score
	}
	
	var totalWorkMin *int
	if session.TotalWorkMin != nil {
		minutes := session.TotalWorkMin.Minutes()
		totalWorkMin = &minutes
	}
	
	var breakTime *int
	if session.BreakTime != nil {
		minutes := session.BreakTime.Minutes()
		breakTime = &minutes
	}
	
	return &dto.SessionResponse{
		ID:           session.ID.Value(),
		StartTime:    session.StartTime,
		EndTime:      session.EndTime,
		AverageFocus: averageFocus,
		TotalWorkMin: totalWorkMin,
		RoundCount: func() *int {
			if session.RoundCount == nil {
				return nil
			}
			count := session.RoundCount.Count()
			return &count
		}(),
		BreakTime: breakTime,
	}
}

// ToSessionsResponse はセッションのリストをレスポンス形式に変換する
func (m *SessionMapper) ToSessionsResponse(sessions []*entity.Session) *dto.SessionsResponse {
	responses := make([]*dto.SessionResponse, len(sessions))
	for i, session := range sessions {
		responses[i] = m.ToSessionResponse(session)
	}
	
	return &dto.SessionsResponse{
		Sessions: responses,
	}
}

// ToSessionResponseList はセッションのリストをレスポンス形式に変換する（複数パターン対応）
func (m *SessionMapper) ToSessionResponseList(sessions []*entity.Session) []*dto.SessionResponse {
	responses := make([]*dto.SessionResponse, len(sessions))
	for i, session := range sessions {
		responses[i] = m.ToSessionResponse(session)
	}
	return responses
}