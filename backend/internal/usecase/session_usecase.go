package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/domain/service"
	sessionVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/session"
	statisticsVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/statistics"
	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/usecase/dto"
	"github.com/tsunakit99/selfpomodoro/internal/usecase/mapper"
)

// SessionUseCase はセッションに関するユースケースを定義するインターフェース
type SessionUseCase interface {
	// StartSession は新しいセッションを開始する
	StartSession(ctx context.Context, userID uuid.UUID) (*dto.SessionResponse, error)

	// CompleteSession はセッションを完了する（SQSメッセージ送信付き）
	CompleteSession(ctx context.Context, id, userID uuid.UUID) (*dto.SessionResponse, error)

	// GetActiveSession はアクティブなセッションを取得する
	GetActiveSession(ctx context.Context, userID uuid.UUID) (*dto.SessionResponse, error)

	// GetRecentSessions は最近のセッション履歴を取得する
	GetRecentSessions(ctx context.Context, userID uuid.UUID, limit int) ([]*dto.SessionResponse, error)
}

// sessionUseCase はSessionUseCaseインターフェースの実装（新エラーハンドリング対応版）
type sessionUseCase struct {
	sessionRepo                 repository.SessionRepository
	optimizationPreferencesRepo repository.OptimizationPreferencesRepository
	messagingService            service.MessagingService
	sessionMapper               *mapper.SessionMapper
	logger                      logger.Logger
}

// NewSessionUseCase は新しいSessionUseCaseインスタンスを作成する
func NewSessionUseCase(
	sessionRepo repository.SessionRepository,
	optimizationPreferencesRepo repository.OptimizationPreferencesRepository,
	messagingService service.MessagingService,
	logger logger.Logger,
) SessionUseCase {
	return &sessionUseCase{
		sessionRepo:                 sessionRepo,
		optimizationPreferencesRepo: optimizationPreferencesRepo,
		messagingService:            messagingService,
		sessionMapper:               mapper.NewSessionMapper(),
		logger:                      logger,
	}
}

// StartSession は新しいセッションを開始する（新エラーハンドリング対応版）
func (uc *sessionUseCase) StartSession(ctx context.Context, userID uuid.UUID) (*dto.SessionResponse, error) {
	// 入力値検証
	if userID == uuid.Nil {
		uc.logger.Error("StartSession: 無効なUserIDが指定されました")
		return nil, appErrors.NewValidationError("ユーザーIDが無効です")
	}

	// UUIDをValue Objectに変換
	userIDVO, err := userVO.NewUserIDFromString(userID.String())
	if err != nil {
		uc.logger.Errorf("StartSession: UserID Value Object変換エラー: %v", err)
		return nil, appErrors.NewInternalError(err)
	}

	// 既存のアクティブセッションをチェック
	existingSession, err := uc.sessionRepo.GetActiveSession(ctx, userIDVO)
	if err != nil && !errors.Is(err, appErrors.ErrRecordNotFound) {
		uc.logger.Errorf("StartSession: アクティブセッション確認エラー: %v", err)
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}
		return nil, appErrors.NewInternalError(err)
	}

	// アクティブなセッションが既に存在する場合
	if existingSession != nil {
		uc.logger.Warnf("StartSession: 既にアクティブなセッションが存在します: %s", existingSession.ID.String())
		return nil, appErrors.NewSessionInProgressError()
	}

	// ✅ ドメインロジック活用：OptimizationPreferences安全取得（デフォルト値フォールバック）
	optimizationPreferences := uc.getOptimizationPreferencesWithFallback(ctx, userIDVO)

	// ユーザー設定確認ログ（MaxRoundsはRoundUseCase側で使用）
	uc.logger.Infof("セッション開始 - ユーザー設定確認完了: work=%d分, break=%d分, rounds=%d",
		optimizationPreferences.GetWorkTimeOrDefault().Minutes(), optimizationPreferences.GetBreakTimeOrDefault().Minutes(), optimizationPreferences.GetSessionRoundsOrDefault().Count())

	// ✅ ドメインファクトリー使用（シンプルなセッション作成）
	session := entity.NewSession(userIDVO)

	// DBにセッションを保存
	if err := uc.sessionRepo.CreateSession(ctx, session); err != nil {
		uc.logger.Errorf("セッション開始エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrUniqueConstraint) {
			return nil, appErrors.NewSessionInProgressError()
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	uc.logger.Infof("セッション開始成功: %s", session.ID.String())
	return uc.sessionMapper.ToSessionResponse(session), nil
}

// CompleteSession はセッションを完了する（新エラーハンドリング対応版）
func (uc *sessionUseCase) CompleteSession(ctx context.Context, id, userID uuid.UUID) (*dto.SessionResponse, error) {
	// 入力値検証
	if id == uuid.Nil {
		uc.logger.Error("CompleteSession: 無効なSessionIDが指定されました")
		return nil, appErrors.NewValidationError("セッションIDが無効です")
	}
	if userID == uuid.Nil {
		uc.logger.Error("CompleteSession: 無効なUserIDが指定されました")
		return nil, appErrors.NewValidationError("ユーザーIDが無効です")
	}

	// UUIDをValue Objectに変換
	sessionIDVO, err := sessionVO.NewSessionIDFromString(id.String())
	if err != nil {
		uc.logger.Errorf("CompleteSession: SessionID Value Object変換エラー: %v", err)
		return nil, appErrors.NewInternalError(err)
	}
	userIDVO, err := userVO.NewUserIDFromString(userID.String())
	if err != nil {
		uc.logger.Errorf("CompleteSession: UserID Value Object変換エラー: %v", err)
		return nil, appErrors.NewInternalError(err)
	}

	// セッション取得（ラウンドも含めて取得）
	session, err := uc.sessionRepo.GetSessionWithRounds(ctx, sessionIDVO, userIDVO)
	if err != nil {
		uc.logger.Errorf("セッション取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return nil, appErrors.NewSessionNotFoundError()
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// ✅ ドメインロジック活用：状態チェック
	if session.IsCompleted() {
		uc.logger.Errorf("セッションは既に完了しています: %s", id.String())
		return nil, appErrors.NewSessionAlreadyEndedError()
	}

	// ✅ DDD Aggregate: セッション取得時に既にラウンドが読み込まれている
	// SessionWithRoundsで取得済みのため、追加の読み込みは不要
	uc.logger.Infof("セッションラウンド数確認: %d個のラウンドが読み込まれました", len(session.GetRounds()))

	// ✅ DDD Aggregate: 統計計算とセッション完了（Aggregate内部で完結）
	session.CompleteWithRounds()

	// 統計ログ出力
	stats := session.CalculateStatistics()
	uc.logger.Infof("セッション統計計算完了: 平均集中度=%.1f, 総作業時間=%d分, ラウンド数=%d, 休憩時間=%d分",
		stats.AverageFocus(), stats.TotalWorkMin(), stats.RoundCount(), stats.BreakTime())

	// データベース更新（Session集約として更新）
	err = uc.sessionRepo.CompleteSession(ctx, session)
	if err != nil {
		uc.logger.Errorf("セッション完了永続化エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return nil, appErrors.NewSessionNotFoundError()
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// ✅ ドメインロジック活用：最適化メッセージ送信判定
	if session.ShouldSendOptimizationMessage() {
		avgFocus, totalWork, _ := session.GetOptimizationMessageData()
		uc.logger.Infof("セッション完了 - SQS最適化メッセージ送信開始: SessionID=%s, RoundCount=%d, AvgFocus=%.1f, TotalWork=%dmin",
			id.String(), stats.RoundCount(), avgFocus, totalWork)

		// Value Objectsに変換してメッセージ送信
		focusScoreVO, err := statisticsVO.NewFocusScore(avgFocus)
		if err != nil {
			uc.logger.Errorf("FocusScore変換エラー: %v", err)
		} else {
			workTimeVO, err := statisticsVO.NewWorkTime(totalWork)
			if err != nil {
				uc.logger.Errorf("WorkTime変換エラー: %v", err)
			} else {
				// ドメインサービス経由でメッセージ送信
				sessionData := service.SessionOptimizationData{
					UserID:        userIDVO,
					SessionID:     sessionIDVO,
					AvgFocusScore: focusScoreVO,
					TotalWorkTime: workTimeVO,
				}

				if err := uc.messagingService.SendSessionOptimizationMessage(ctx, sessionData); err != nil {
					uc.logger.Errorf("SQSメッセージ送信エラー: %v", err)
				}
			}
		}
	} else {
		uc.logger.Warn("完了したラウンドが存在しないため、SQS最適化メッセージは送信しません")
	}

	// 完了したセッションを取得
	updatedSession, err := uc.sessionRepo.GetSessionWithRounds(ctx, sessionIDVO, userIDVO)
	if err != nil {
		uc.logger.Errorf("完了後のセッション取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return nil, appErrors.NewSessionNotFoundError()
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// ✅ ドメインロジック活用：セッション品質評価ログ
	quality := updatedSession.GetSessionQuality()
	efficiency := updatedSession.GetEfficiency()
	uc.logger.Infof("セッション完了成功: SessionID=%s, 品質=%s, 効率=%.1f%%, 平均集中度=%.1f, 総作業時間=%d分, ラウンド数=%d",
		id.String(), quality, efficiency, stats.AverageFocus(), stats.TotalWorkMin(), stats.RoundCount())

	return uc.sessionMapper.ToSessionResponse(updatedSession), nil
}

// GetActiveSession はアクティブなセッションを取得する（Session集約完全対応版）
func (uc *sessionUseCase) GetActiveSession(ctx context.Context, userID uuid.UUID) (*dto.SessionResponse, error) {
	// UUIDをValue Objectに変換
	userIDVO, err := userVO.NewUserIDFromString(userID.String())
	if err != nil {
		return nil, appErrors.NewInternalError(err)
	}

	// アクティブなセッションを取得
	session, err := uc.sessionRepo.GetActiveSession(ctx, userIDVO)
	if err != nil {
		uc.logger.Errorf("アクティブセッション取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return nil, appErrors.NewSessionNotFoundError()
		}
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	uc.logger.Infof("アクティブセッション取得成功: SessionID=%s, 開始時刻=%s",
		session.ID.String(), session.StartTime.Format("2006-01-02 15:04:05"))

	return uc.sessionMapper.ToSessionResponse(session), nil
}

// GetRecentSessions は最近のセッション履歴を取得する（Session集約完全対応版）
func (uc *sessionUseCase) GetRecentSessions(ctx context.Context, userID uuid.UUID, limit int) ([]*dto.SessionResponse, error) {
	// UUIDをValue Objectに変換
	userIDVO, err := userVO.NewUserIDFromString(userID.String())
	if err != nil {
		return nil, appErrors.NewInternalError(err)
	}

	// パラメータ検証
	if limit <= 0 || limit > 100 {
		limit = 20 // デフォルト値
	}

	// 最近のセッション履歴を取得
	sessions, err := uc.sessionRepo.GetRecentSessions(ctx, userIDVO, limit)
	if err != nil {
		uc.logger.Errorf("セッション履歴取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if appErrors.IsDatabaseError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewInternalError(err)
	}

	// レスポンス形式に変換
	responses := uc.sessionMapper.ToSessionResponseList(sessions)

	uc.logger.Infof("セッション履歴取得成功: UserID=%s, 取得件数=%d, リクエスト上限=%d",
		userID.String()[:8]+"...", len(responses), limit)

	return responses, nil
}

// ✅ ドメインロジック活用：OptimizationPreferences安全取得（フォールバック）
func (uc *sessionUseCase) getOptimizationPreferencesWithFallback(ctx context.Context, userID userVO.UserID) *entity.OptimizationPreferences {
	if uc.optimizationPreferencesRepo == nil {
		uc.logger.Warn("OptimizationPreferencesRepository が nil です - デフォルト設定使用")
		defaultPrefs, _ := entity.NewOptimizationPreferences(userID)
		return defaultPrefs
	}

	optimizationPreferences, err := uc.optimizationPreferencesRepo.GetOrCreateDefault(ctx, userID)
	if err != nil {
		uc.logger.Warnf("OptimizationPreferences取得エラー、デフォルト設定を使用: %v", err)
		defaultPrefs, _ := entity.NewOptimizationPreferences(userID)
		return defaultPrefs
	}

	return optimizationPreferences
}

