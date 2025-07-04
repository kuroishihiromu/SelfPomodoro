package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// OptimizationUseCase は最適化に関するビジネスロジックを担当
type OptimizationUseCase interface {
	// 最適化履歴の取得
	GetRoundOptimizationHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*model.RoundOptimizationLog, error)
	GetSessionOptimizationHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*model.SessionOptimizationLog, error)
	
	// 最適化結果の処理
	ProcessRoundOptimizationResult(ctx context.Context, userID, roundID uuid.UUID, focusScore int, workTime, breakTime float64) error
	ProcessSessionOptimizationResult(ctx context.Context, userID, sessionID uuid.UUID, avgFocusScore float64, totalWorkTime int, roundCount int, breakTime float64) error
	
	// 最適化効果性の分析
	GetOptimizationEffectiveness(ctx context.Context, userID uuid.UUID, since time.Time) (*model.OptimizationEffectiveness, error)
	GetOptimizationSummary(ctx context.Context, userID uuid.UUID) (*model.OptimizationSummary, error)
	
	// UserConfigへの最適化結果適用
	ApplyOptimizationToUserConfig(ctx context.Context, userID uuid.UUID) error
}

// optimizationUseCase はOptimizationUseCaseの実装
type optimizationUseCase struct {
	optimizationRepo repository.OptimizationRepository
	userConfigRepo   repository.UserConfigRepository
	logger           logger.Logger
}

// NewOptimizationUseCase は新しい最適化UseCaseを作成する
func NewOptimizationUseCase(
	optimizationRepo repository.OptimizationRepository,
	userConfigRepo repository.UserConfigRepository,
	logger logger.Logger,
) OptimizationUseCase {
	return &optimizationUseCase{
		optimizationRepo: optimizationRepo,
		userConfigRepo:   userConfigRepo,
		logger:           logger,
	}
}

// GetRoundOptimizationHistory はラウンド最適化履歴を取得する
func (uc *optimizationUseCase) GetRoundOptimizationHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*model.RoundOptimizationLog, error) {
	uc.logger.Infof("ラウンド最適化履歴取得開始: UserID=%s, Limit=%d", userID.String()[:8]+"...", limit)
	
	if limit <= 0 || limit > 100 {
		limit = 50 // デフォルト制限
	}
	
	logs, err := uc.optimizationRepo.GetRoundOptimizationHistory(ctx, userID, limit)
	if err != nil {
		uc.logger.Errorf("ラウンド最適化履歴取得失敗: %v", err)
		return nil, err
	}
	
	uc.logger.Infof("ラウンド最適化履歴取得完了: UserID=%s, 件数=%d", userID.String()[:8]+"...", len(logs))
	return logs, nil
}

// GetSessionOptimizationHistory はセッション最適化履歴を取得する
func (uc *optimizationUseCase) GetSessionOptimizationHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*model.SessionOptimizationLog, error) {
	uc.logger.Infof("セッション最適化履歴取得開始: UserID=%s, Limit=%d", userID.String()[:8]+"...", limit)
	
	if limit <= 0 || limit > 100 {
		limit = 50 // デフォルト制限
	}
	
	logs, err := uc.optimizationRepo.GetSessionOptimizationHistory(ctx, userID, limit)
	if err != nil {
		uc.logger.Errorf("セッション最適化履歴取得失敗: %v", err)
		return nil, err
	}
	
	uc.logger.Infof("セッション最適化履歴取得完了: UserID=%s, 件数=%d", userID.String()[:8]+"...", len(logs))
	return logs, nil
}

// ProcessRoundOptimizationResult はラウンド最適化結果を処理する
func (uc *optimizationUseCase) ProcessRoundOptimizationResult(
	ctx context.Context, 
	userID, roundID uuid.UUID, 
	focusScore int, 
	workTime, breakTime float64,
) error {
	uc.logger.Infof("ラウンド最適化結果処理開始: UserID=%s, FocusScore=%d", userID.String()[:8]+"...", focusScore)
	
	// 入力値検証
	if focusScore < 0 || focusScore > 100 {
		return appErrors.NewValidationError("集中度スコアは0-100の範囲で入力してください")
	}
	
	if workTime < 5 || workTime > 120 {
		return appErrors.NewValidationError("作業時間は5-120分の範囲で入力してください")
	}
	
	if breakTime < 1 || breakTime > 60 {
		return appErrors.NewValidationError("休憩時間は1-60分の範囲で入力してください")
	}
	
	// 最適化ログを作成
	log := model.NewRoundOptimizationLog(userID, int(workTime), int(breakTime), focusScore)
	
	// データベースに保存
	err := uc.optimizationRepo.SaveRoundOptimizationLog(ctx, log)
	if err != nil {
		uc.logger.Errorf("ラウンド最適化結果保存失敗: %v", err)
		return err
	}
	
	uc.logger.Infof("ラウンド最適化結果処理完了: UserID=%s", userID.String()[:8]+"...")
	return nil
}

// ProcessSessionOptimizationResult はセッション最適化結果を処理する
func (uc *optimizationUseCase) ProcessSessionOptimizationResult(
	ctx context.Context,
	userID, sessionID uuid.UUID,
	avgFocusScore float64,
	totalWorkTime int,
	roundCount int,
	breakTime float64,
) error {
	uc.logger.Infof("セッション最適化結果処理開始: UserID=%s, AvgFocusScore=%.1f", userID.String()[:8]+"...", avgFocusScore)
	
	// 入力値検証
	if avgFocusScore < 0 || avgFocusScore > 100 {
		return appErrors.NewValidationError("平均集中度スコアは0-100の範囲で入力してください")
	}
	
	if totalWorkTime < 30 || totalWorkTime > 480 {
		return appErrors.NewValidationError("合計作業時間は30-480分の範囲で入力してください")
	}
	
	if roundCount < 1 || roundCount > 20 {
		return appErrors.NewValidationError("ラウンド数は1-20の範囲で入力してください")
	}
	
	if breakTime < 5 || breakTime > 120 {
		return appErrors.NewValidationError("休憩時間は5-120分の範囲で入力してください")
	}
	
	// 最適化ログを作成
	log := model.NewSessionOptimizationLog(userID, roundCount, int(breakTime), avgFocusScore, totalWorkTime)
	
	// データベースに保存
	err := uc.optimizationRepo.SaveSessionOptimizationLog(ctx, log)
	if err != nil {
		uc.logger.Errorf("セッション最適化結果保存失敗: %v", err)
		return err
	}
	
	uc.logger.Infof("セッション最適化結果処理完了: UserID=%s", userID.String()[:8]+"...")
	return nil
}

// GetOptimizationEffectiveness は最適化の効果性を取得する
func (uc *optimizationUseCase) GetOptimizationEffectiveness(ctx context.Context, userID uuid.UUID, since time.Time) (*model.OptimizationEffectiveness, error) {
	uc.logger.Infof("最適化効果性取得開始: UserID=%s, Since=%s", userID.String()[:8]+"...", since.Format("2006-01-02"))
	
	effectiveness, err := uc.optimizationRepo.GetOptimizationEffectiveness(ctx, userID, since)
	if err != nil {
		uc.logger.Errorf("最適化効果性取得失敗: %v", err)
		return nil, err
	}
	
	uc.logger.Infof("最適化効果性取得完了: UserID=%s", userID.String()[:8]+"...")
	return effectiveness, nil
}

// GetOptimizationSummary は最適化サマリーを取得する
func (uc *optimizationUseCase) GetOptimizationSummary(ctx context.Context, userID uuid.UUID) (*model.OptimizationSummary, error) {
	uc.logger.Infof("最適化サマリー取得開始: UserID=%s", userID.String()[:8]+"...")
	
	// 最適化ログ数を取得
	roundCount, sessionCount, err := uc.optimizationRepo.CountOptimizationLogs(ctx, userID)
	if err != nil {
		uc.logger.Errorf("最適化ログカウント失敗: %v", err)
		return nil, err
	}
	
	// 最新の最適化結果を取得
	var latestRoundResult *model.RoundOptimizationLog
	var latestSessionResult *model.SessionOptimizationLog
	
	if roundCount > 0 {
		latestRoundResult, _ = uc.optimizationRepo.GetLatestRoundOptimizationResult(ctx, userID)
	}
	
	if sessionCount > 0 {
		latestSessionResult, _ = uc.optimizationRepo.GetLatestSessionOptimizationResult(ctx, userID)
	}
	
	// サマリーを作成
	summary := &model.OptimizationSummary{
		UserID:                userID,
		TotalRoundOptimizations:   roundCount,
		TotalSessionOptimizations: sessionCount,
		LatestRoundResult:        latestRoundResult,
		LatestSessionResult:      latestSessionResult,
		IsOptimizationActive:     roundCount > 0 || sessionCount > 0,
		CreatedAt:                time.Now(),
	}
	
	// 最後の最適化時刻を設定
	if latestRoundResult != nil {
		if timestamp, err := time.Parse(time.RFC3339, latestRoundResult.Timestamp); err == nil {
			summary.LastOptimizedAt = timestamp
		}
	}
	
	if latestSessionResult != nil {
		if timestamp, err := time.Parse(time.RFC3339, latestSessionResult.Timestamp); err == nil {
			if summary.LastOptimizedAt.IsZero() || timestamp.After(summary.LastOptimizedAt) {
				summary.LastOptimizedAt = timestamp
			}
		}
	}
	
	uc.logger.Infof("最適化サマリー取得完了: UserID=%s, Round=%d, Session=%d", 
		userID.String()[:8]+"...", roundCount, sessionCount)
	
	return summary, nil
}

// ApplyOptimizationToUserConfig は最適化結果をUserConfigに適用する
func (uc *optimizationUseCase) ApplyOptimizationToUserConfig(ctx context.Context, userID uuid.UUID) error {
	uc.logger.Infof("最適化結果のUserConfig適用開始: UserID=%s", userID.String()[:8]+"...")
	
	// 現在のUserConfigを取得
	currentConfig, err := uc.userConfigRepo.GetUserConfig(ctx, userID)
	if err != nil {
		uc.logger.Errorf("UserConfig取得失敗: %v", err)
		return err
	}
	
	// 最新のラウンド最適化結果を取得
	latestRoundResult, err := uc.optimizationRepo.GetLatestRoundOptimizationResult(ctx, userID)
	if err != nil && !appErrors.IsNotFoundError(err) {
		uc.logger.Errorf("最新ラウンド最適化結果取得失敗: %v", err)
		return err
	}
	
	// 最新のセッション最適化結果を取得
	latestSessionResult, err := uc.optimizationRepo.GetLatestSessionOptimizationResult(ctx, userID)
	if err != nil && !appErrors.IsNotFoundError(err) {
		uc.logger.Errorf("最新セッション最適化結果取得失敗: %v", err)
		return err
	}
	
	// 更新フラグ
	configUpdated := false
	
	// ラウンド最適化結果を適用
	if latestRoundResult != nil {
		// 段階的適用（急激な変化を避ける）
		newWorkTime := uc.calculateGradualChange(currentConfig.RoundWorkTime, latestRoundResult.WorkTime, 0.3)
		newBreakTime := uc.calculateGradualChange(currentConfig.RoundBreakTime, latestRoundResult.BreakTime, 0.3)
		
		if newWorkTime != currentConfig.RoundWorkTime || newBreakTime != currentConfig.RoundBreakTime {
			currentConfig.RoundWorkTime = newWorkTime
			currentConfig.RoundBreakTime = newBreakTime
			configUpdated = true
			
			uc.logger.Infof("ラウンド最適化適用: Work=%d分, Break=%d分", newWorkTime, newBreakTime)
		}
	}
	
	// セッション最適化結果を適用
	if latestSessionResult != nil {
		// 段階的適用
		newRounds := uc.calculateGradualChange(currentConfig.SessionRounds, latestSessionResult.RoundCount, 0.5)
		newBreakTime := uc.calculateGradualChange(currentConfig.SessionBreakTime, latestSessionResult.BreakTime, 0.3)
		
		if newRounds != currentConfig.SessionRounds || newBreakTime != currentConfig.SessionBreakTime {
			currentConfig.SessionRounds = newRounds
			currentConfig.SessionBreakTime = newBreakTime
			configUpdated = true
			
			uc.logger.Infof("セッション最適化適用: Rounds=%d, Break=%d分", newRounds, newBreakTime)
		}
	}
	
	// 変更がある場合のみ保存
	if configUpdated {
		err = uc.userConfigRepo.UpdateUserConfig(ctx, currentConfig)
		if err != nil {
			uc.logger.Errorf("UserConfig更新失敗: %v", err)
			return err
		}
		
		uc.logger.Infof("最適化結果のUserConfig適用完了: UserID=%s", userID.String()[:8]+"...")
	} else {
		uc.logger.Infof("最適化結果の変更なし: UserID=%s", userID.String()[:8]+"...")
	}
	
	return nil
}

// calculateGradualChange は段階的な変更を計算する
func (uc *optimizationUseCase) calculateGradualChange(current, target int, changeRate float64) int {
	if current == target {
		return current
	}
	
	diff := float64(target - current)
	change := diff * changeRate
	
	// 最小1の変更は適用
	if change > 0 && change < 1 {
		change = 1
	} else if change < 0 && change > -1 {
		change = -1
	}
	
	newValue := current + int(change)
	
	// 合理的な範囲内に制限
	if newValue < 1 {
		newValue = 1
	} else if newValue > 240 { // 最大4時間
		newValue = 240
	}
	
	return newValue
}