package repository

import (
	"context"

	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	roundVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/round"
	sessionVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/session"
	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
)

// SessionRepository はSession集約（Session + Round）の永続化を担当するリポジトリ
// DDD原則: 1つの集約ルートに対して1つのリポジトリ
type SessionRepository interface {
	// Session集約ルートの操作
	CreateSession(ctx context.Context, session *entity.Session) error
	GetSession(ctx context.Context, sessionID sessionVO.SessionID, userID userVO.UserID) (*entity.Session, error)
	GetSessionWithRounds(ctx context.Context, sessionID sessionVO.SessionID, userID userVO.UserID) (*entity.Session, error)
	CompleteSession(ctx context.Context, session *entity.Session) error
	
	// Session集約内のRound操作（Sessionを通じて管理）
	AddRoundToSession(ctx context.Context, sessionID sessionVO.SessionID, userID userVO.UserID, round *entity.Round) error
	CompleteRound(ctx context.Context, sessionID sessionVO.SessionID, userID userVO.UserID, roundID roundVO.RoundID, focusScore *int, workTime, breakTime int) error
	GetRoundsBySession(ctx context.Context, sessionID sessionVO.SessionID, userID userVO.UserID) ([]*entity.Round, error)
	GetRoundByID(ctx context.Context, roundID roundVO.RoundID, userID userVO.UserID) (*entity.Round, error)
	
	// Session集約のビジネスロジック支援
	GetActiveSession(ctx context.Context, userID userVO.UserID) (*entity.Session, error)
	GetRecentSessions(ctx context.Context, userID userVO.UserID, limit int) ([]*entity.Session, error)
}
