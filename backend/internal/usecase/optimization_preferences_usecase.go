package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/entity"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/usecase/dto"
	"github.com/tsunakit99/selfpomodoro/internal/usecase/mapper"
	roundVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/round"
	sessionVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/session"
	userVO "github.com/tsunakit99/selfpomodoro/internal/domain/valueobject/user"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// OptimizationPreferencesUseCase はユーザー最適化設定に関するユースケースを定義するインターフェース
type OptimizationPreferencesUseCase interface {
	// GetOptimizationPreferences はユーザー最適化設定を取得する（デフォルト値フォールバック付き）
	GetOptimizationPreferences(ctx context.Context, userID uuid.UUID) (*dto.OptimizationPreferencesResponse, error)

	// CreateOptimizationPreferences は新しいユーザー最適化設定を作成する（PostConfirmation専用）
	CreateOptimizationPreferences(ctx context.Context, userID uuid.UUID, req *dto.CreateOptimizationPreferencesRequest) (*dto.OptimizationPreferencesResponse, error)

	// UpdateOptimizationPreferences はユーザー最適化設定を更新する
	UpdateOptimizationPreferences(ctx context.Context, userID uuid.UUID, req *dto.UpdateOptimizationPreferencesRequest) (*dto.OptimizationPreferencesResponse, error)

	// DeleteOptimizationPreferences はユーザー最適化設定を削除する
	DeleteOptimizationPreferences(ctx context.Context, userID uuid.UUID) error

	// GetOptimizationPreferencesForInternal は内部処理用にユーザー最適化設定を取得する（デフォルト値フォールバック付き）
	GetOptimizationPreferencesForInternal(ctx context.Context, userID userVO.UserID) (*entity.OptimizationPreferences, error)
}

// optimizationPreferencesUseCase はOptimizationPreferencesUseCaseインターフェースの実装（ドメイン強化版）
type optimizationPreferencesUseCase struct {
	optimizationPreferencesRepo repository.OptimizationPreferencesRepository
	mapper                      *mapper.OptimizationPreferencesMapper
	logger                      logger.Logger
}

// NewOptimizationPreferencesUseCase は新しいOptimizationPreferencesUseCaseインスタンスを作成する
func NewOptimizationPreferencesUseCase(optimizationPreferencesRepo repository.OptimizationPreferencesRepository, logger logger.Logger) OptimizationPreferencesUseCase {
	return &optimizationPreferencesUseCase{
		optimizationPreferencesRepo: optimizationPreferencesRepo,
		mapper:                      mapper.NewOptimizationPreferencesMapper(),
		logger:                      logger,
	}
}

// GetOptimizationPreferences はユーザー最適化設定を取得する（ドメイン強化版・デフォルト値フォールバック付き）
func (uc *optimizationPreferencesUseCase) GetOptimizationPreferences(ctx context.Context, userID uuid.UUID) (*dto.OptimizationPreferencesResponse, error) {
	// UUIDをValue Objectに変換
	userIDVO, err := userVO.NewUserIDFromString(userID.String())
	if err != nil {
		return nil, appErrors.NewInternalError(err)
	}

	if uc.optimizationPreferencesRepo == nil {
		uc.logger.Warn("OptimizationPreferencesRepository が初期化されていません。デフォルト設定を返します")
		// ✅ ドメインファクトリー使用：DynamoDBが利用できない場合のフォールバック
		defaultConfig, _ := entity.NewOptimizationPreferences(userIDVO)
		return uc.mapper.ToOptimizationPreferencesResponse(defaultConfig), nil
	}

	config, err := uc.optimizationPreferencesRepo.Get(ctx, userIDVO)
	if err != nil {
		uc.logger.Errorf("ユーザー最適化設定取得エラー、デフォルト設定を返します: %v", err)
		// ✅ ドメインファクトリー使用：設定が存在しない場合のデフォルト値フォールバック
		defaultConfig, _ := entity.NewOptimizationPreferences(userIDVO)

		// ✅ ドメインロジック活用：デフォルト値ログ出力
		uc.logger.Infof("デフォルト設定使用: work=%d分, break=%d分, rounds=%d, sessionBreak=%d分",
			defaultConfig.GetWorkTimeOrDefault().Minutes(),
			defaultConfig.GetBreakTimeOrDefault().Minutes(),
			defaultConfig.GetSessionRoundsOrDefault().Count(),
			defaultConfig.SessionBreakTime.Minutes())

		return uc.mapper.ToOptimizationPreferencesResponse(defaultConfig), nil
	}

	// ✅ ドメインロジック活用：取得した設定値の安全性確認
	uc.logger.Infof("ユーザー最適化設定取得成功: work=%d分, break=%d分, rounds=%d, sessionBreak=%d分",
		config.GetWorkTimeOrDefault().Minutes(),
		config.GetBreakTimeOrDefault().Minutes(),
		config.GetSessionRoundsOrDefault().Count(),
		config.SessionBreakTime.Minutes())

	return uc.mapper.ToOptimizationPreferencesResponse(config), nil
}

// CreateOptimizationPreferences は新しいユーザー最適化設定を作成する（ドメイン強化版）
func (uc *optimizationPreferencesUseCase) CreateOptimizationPreferences(ctx context.Context, userID uuid.UUID, req *dto.CreateOptimizationPreferencesRequest) (*dto.OptimizationPreferencesResponse, error) {
	// UUIDをValue Objectに変換
	userIDVO, err := userVO.NewUserIDFromString(userID.String())
	if err != nil {
		return nil, appErrors.NewInternalError(err)
	}

	if uc.optimizationPreferencesRepo == nil {
		return nil, appErrors.NewInternalError(errors.New("ユーザー最適化設定機能は現在利用できません"))
	}

	// ✅ ドメインファクトリー使用：リクエストから新しい設定を作成
	config, err := entity.NewOptimizationPreferences(userIDVO)
	if err != nil {
		return nil, appErrors.NewInternalError(err)
	}

	// Value Objectsを作成してリクエストデータを適用
	if req.RoundWorkTime != nil {
		workTime, err := roundVO.NewWorkTime(*req.RoundWorkTime)
		if err != nil {
			return nil, appErrors.NewValidationError(err.Error())
		}
		config.RoundWorkTime = workTime
	}

	if req.RoundBreakTime != nil {
		breakTime, err := roundVO.NewBreakTime(*req.RoundBreakTime)
		if err != nil {
			return nil, appErrors.NewValidationError(err.Error())
		}
		config.RoundBreakTime = breakTime
	}

	if req.SessionRounds != nil {
		sessionRounds, err := sessionVO.NewSessionRounds(*req.SessionRounds)
		if err != nil {
			return nil, appErrors.NewValidationError(err.Error())
		}
		config.SessionRounds = sessionRounds
	}

	if req.SessionBreakTime != nil {
		sessionBreakTime, err := sessionVO.NewBreakTime(*req.SessionBreakTime)
		if err != nil {
			return nil, appErrors.NewValidationError(err.Error())
		}
		config.SessionBreakTime = sessionBreakTime
	}

	// ✅ ドメインロジック活用：設定の有効性をチェック
	if !config.IsValid() {
		uc.logger.Error("ユーザー最適化設定バリデーションエラー: 無効な設定値")
		return nil, appErrors.NewValidationError("無効な設定値です")
	}

	// リポジトリに保存
	if err := uc.optimizationPreferencesRepo.Create(ctx, config); err != nil {
		uc.logger.Errorf("ユーザー最適化設定作成エラー: %v", err)
		if errors.Is(err, appErrors.ErrUniqueConstraint) {
			return nil, appErrors.NewUserConfigCreateFailedError()
		}
		return nil, appErrors.NewInternalError(err)
	}

	uc.logger.Infof("ユーザー最適化設定作成成功: UserID=%s, work=%d分, break=%d分, rounds=%d",
		userID.String()[:8]+"...", config.RoundWorkTime.Minutes(), config.RoundBreakTime.Minutes(), config.SessionRounds.Count())

	return uc.mapper.ToOptimizationPreferencesResponse(config), nil
}

// UpdateOptimizationPreferences はユーザー最適化設定を更新する（ドメイン強化版）
func (uc *optimizationPreferencesUseCase) UpdateOptimizationPreferences(ctx context.Context, userID uuid.UUID, req *dto.UpdateOptimizationPreferencesRequest) (*dto.OptimizationPreferencesResponse, error) {
	// UUIDをValue Objectに変換
	userIDVO, err := userVO.NewUserIDFromString(userID.String())
	if err != nil {
		return nil, appErrors.NewInternalError(err)
	}

	if uc.optimizationPreferencesRepo == nil {
		return nil, appErrors.NewInternalError(errors.New("ユーザー最適化設定機能は現在利用できません"))
	}

	// 既存の設定を取得（存在前提）
	config, err := uc.optimizationPreferencesRepo.Get(ctx, userIDVO)
	if err != nil {
		uc.logger.Errorf("ユーザー最適化設定取得エラー: %v", err)
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return nil, appErrors.NewUserConfigNotFoundError()
		}
		return nil, appErrors.NewInternalError(err)
	}

	// ✅ ドメインロジック活用：設定を部分更新（Value Objects使用）
	updateCount := 0

	if req.RoundWorkTime != nil {
		workTime, err := roundVO.NewWorkTime(*req.RoundWorkTime)
		if err != nil {
			return nil, appErrors.NewValidationError(err.Error())
		}
		oldValue := config.RoundWorkTime.Minutes()
		config.RoundWorkTime = workTime
		uc.logger.Debugf("作業時間更新: %d分 → %d分", oldValue, *req.RoundWorkTime)
		updateCount++
	}

	if req.RoundBreakTime != nil {
		breakTime, err := roundVO.NewBreakTime(*req.RoundBreakTime)
		if err != nil {
			return nil, appErrors.NewValidationError(err.Error())
		}
		oldValue := config.RoundBreakTime.Minutes()
		config.RoundBreakTime = breakTime
		uc.logger.Debugf("休憩時間更新: %d分 → %d分", oldValue, *req.RoundBreakTime)
		updateCount++
	}

	if req.SessionRounds != nil {
		sessionRounds, err := sessionVO.NewSessionRounds(*req.SessionRounds)
		if err != nil {
			return nil, appErrors.NewValidationError(err.Error())
		}
		oldValue := config.SessionRounds.Count()
		config.SessionRounds = sessionRounds
		uc.logger.Debugf("セッションラウンド数更新: %d回 → %d回", oldValue, *req.SessionRounds)
		updateCount++
	}

	if req.SessionBreakTime != nil {
		sessionBreakTime, err := sessionVO.NewBreakTime(*req.SessionBreakTime)
		if err != nil {
			return nil, appErrors.NewValidationError(err.Error())
		}
		oldValue := config.SessionBreakTime.Minutes()
		config.SessionBreakTime = sessionBreakTime
		uc.logger.Debugf("セッション長休憩時間更新: %d分 → %d分", oldValue, *req.SessionBreakTime)
		updateCount++
	}

	// ✅ ドメインロジック活用：更新後の設定の有効性をチェック
	if !config.IsValid() {
		uc.logger.Error("ユーザー最適化設定更新バリデーションエラー: 無効な設定値")
		return nil, appErrors.NewValidationError("無効な設定値です")
	}

	// リポジトリに保存
	if err := uc.optimizationPreferencesRepo.Update(ctx, config); err != nil {
		uc.logger.Errorf("ユーザー最適化設定更新エラー: %v", err)
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return nil, appErrors.NewUserConfigUpdateFailedError()
		}
		return nil, appErrors.NewInternalError(err)
	}

	uc.logger.Infof("ユーザー最適化設定更新成功: UserID=%s, work=%d分, break=%d分, rounds=%d, 更新項目数=%d",
		userID.String()[:8]+"...", config.RoundWorkTime.Minutes(), config.RoundBreakTime.Minutes(), config.SessionRounds.Count(), updateCount)

	return uc.mapper.ToOptimizationPreferencesResponse(config), nil
}

// DeleteOptimizationPreferences はユーザー最適化設定を削除する
func (uc *optimizationPreferencesUseCase) DeleteOptimizationPreferences(ctx context.Context, userID uuid.UUID) error {
	// UUIDをValue Objectに変換
	userIDVO, err := userVO.NewUserIDFromString(userID.String())
	if err != nil {
		return appErrors.NewInternalError(err)
	}

	if uc.optimizationPreferencesRepo == nil {
		return appErrors.NewInternalError(errors.New("ユーザー最適化設定機能は現在利用できません"))
	}

	if err := uc.optimizationPreferencesRepo.Delete(ctx, userIDVO); err != nil {
		uc.logger.Errorf("ユーザー最適化設定削除エラー: %v", err)
		if errors.Is(err, appErrors.ErrRecordNotFound) {
			return appErrors.NewUserConfigNotFoundError()
		}
		return appErrors.NewInternalError(err)
	}

	uc.logger.Infof("ユーザー最適化設定削除成功: UserID=%s", userID.String()[:8]+"...")
	return nil
}

// GetOptimizationPreferencesForInternal は内部処理用にユーザー最適化設定を取得する（ドメイン強化版・デフォルト値フォールバック）
func (uc *optimizationPreferencesUseCase) GetOptimizationPreferencesForInternal(ctx context.Context, userID userVO.UserID) (*entity.OptimizationPreferences, error) {
	if uc.optimizationPreferencesRepo == nil {
		uc.logger.Warn("OptimizationPreferencesRepository が初期化されていません。デフォルト設定を返します")
		// ✅ ドメインファクトリー使用
		defaultConfig, _ := entity.NewOptimizationPreferences(userID)
		return defaultConfig, nil
	}

	config, err := uc.optimizationPreferencesRepo.Get(ctx, userID)
	if err != nil {
		uc.logger.Warnf("内部処理用ユーザー最適化設定取得エラー、デフォルト設定にフォールバックします: %v", err)
		// ✅ ドメインファクトリー使用：エラーの場合もデフォルト設定でフォールバック（処理は継続）
		defaultConfig, _ := entity.NewOptimizationPreferences(userID)

		// ✅ ドメインロジック活用：内部処理用のベース値ログ出力
		uc.logger.Infof("内部処理用デフォルト設定: work=%d分, break=%d分, rounds=%d, sessionBreak=%d分",
			defaultConfig.GetWorkTimeOrDefault().Minutes(),
			defaultConfig.GetBreakTimeOrDefault().Minutes(),
			defaultConfig.GetSessionRoundsOrDefault().Count(),
			defaultConfig.SessionBreakTime.Minutes())

		return defaultConfig, nil
	}

	// ✅ ドメインロジック活用：内部処理用設定確認ログ
	uc.logger.Infof("内部処理用設定取得成功: work=%d分, break=%d分, rounds=%d, sessionBreak=%d分",
		config.GetWorkTimeOrDefault().Minutes(),
		config.GetBreakTimeOrDefault().Minutes(),
		config.GetSessionRoundsOrDefault().Count(),
		config.SessionBreakTime.Minutes())

	return config, nil
}