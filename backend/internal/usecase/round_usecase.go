package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	domainErrors "github.com/tsunakit99/selfpomodoro/internal/domain/errors"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	rounderrors "github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/postgres/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/sqs"
)

// RoundUseCase はラウンドに関するユースケースを定義するインターフェース
type RoundUseCase interface {
	// StartRound は新しいラウンドを開始する
	StartRound(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, req *model.RoundCreateRequest) (*model.RoundResponse, error)

	// GetRound は指定されたIDのラウンドを取得する
	GetRound(ctx context.Context, id uuid.UUID) (*model.RoundResponse, error)

	// GetAllRoundsBySessionID は指定されたセッションIDのラウンドを取得する
	GetAllRoundsBySessionID(ctx context.Context, sessionID uuid.UUID) (*model.RoundsResponse, error)

	// CompleteRound はラウンドを完了する(SQSメッセージ送信付き)
	CompleteRound(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *model.RoundCompleteRequest) (*model.RoundResponse, error)

	// AbortRound はラウンドを中止する（SQSメッセージは送信しない）
	AbortRound(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.RoundResponse, error)
}

// roundUseCase はラウンドに関するユースケースの実装
type roundUseCase struct {
	roundRepo      repository.RoundRepository
	sessionRepo    repository.SessionRepository
	userConfigRepo repository.UserConfigRepository
	sqsClient      *sqs.SQSClient
	logger         logger.Logger
}

// NewRoundUseCase は新しいラウンドユースケースを作成する
func NewRoundUseCase(roundrepo repository.RoundRepository, sessionRepo repository.SessionRepository, userConfigRepo repository.UserConfigRepository, sqsClient *sqs.SQSClient, logger logger.Logger) RoundUseCase {
	return &roundUseCase{
		roundRepo:      roundrepo,
		sessionRepo:    sessionRepo,
		userConfigRepo: userConfigRepo,
		sqsClient:      sqsClient,
		logger:         logger,
	}
}

// StartRound は新しいラウンドを開始する
func (uc *roundUseCase) StartRound(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID, req *model.RoundCreateRequest) (*model.RoundResponse, error) {
	// セッションが存在するか確認
	session, err := uc.sessionRepo.GetByID(ctx, sessionID, userID)
	if err != nil {
		uc.logger.Errorf("セッション取得エラー: %v", err)
		return nil, err
	}

	// 最後のラウンドを取得して、完了しているか確認
	lastRound, err := uc.roundRepo.GetLastRoundBySessionID(ctx, session.ID)
	if err != nil && !errors.Is(err, rounderrors.ErrNoRoundsInSession) {
		uc.logger.Errorf("最後のラウンド取得エラー: %v", err)
		return nil, err
	}

	// 最後のラウンドが存在して、まだ完了していない場合はエラー
	if lastRound != nil && lastRound.EndTime == nil {
		uc.logger.Errorf("進行中のラウンドがあります: %v", lastRound.ID)
		return nil, errors.New("進行中のラウンドがあります。新しいラウンドを開始する前に現在のラウンドを完了またはスキップしてください。")
	}

	// 新しいラウンド順序を計算
	var newRoundOrder int
	if lastRound == nil {
		newRoundOrder = 1 // 最初のラウンド
	} else {
		newRoundOrder = lastRound.RoundOrder + 1 // 次のラウンド
	}

	// 新しいラウンドを作成
	round := model.NewRound(session.ID, newRoundOrder)

	// ラウンドをデータベースに保存
	if err = uc.roundRepo.Create(ctx, round); err != nil {
		uc.logger.Errorf("ラウンド作成エラー: %v", err)
		return nil, err
	}

	// ラウンドをレスポンスに変換
	return round.ToResponse(), nil
}

// GetRound は指定されたIDのラウンドを取得する
func (uc *roundUseCase) GetRound(ctx context.Context, id uuid.UUID) (*model.RoundResponse, error) {
	round, err := uc.roundRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Errorf("ラウンド取得エラー: %v", err)
		if errors.Is(err, rounderrors.ErrRoundNotFound) {
			return nil, domainErrors.NewRoundNotFoundError()
		}
		return nil, domainErrors.NewInternalServerError()
	}
	return round.ToResponse(), nil
}

// GetAllRoundsBySessionID はセッションIDに紐づくすべてのラウンドを取得する
func (uc *roundUseCase) GetAllRoundsBySessionID(ctx context.Context, sessionID uuid.UUID) (*model.RoundsResponse, error) {
	rounds, err := uc.roundRepo.GetAllBySessionID(ctx, sessionID)
	if err != nil {
		uc.logger.Errorf("ラウンド一覧取得エラー: %v", err)
		return nil, domainErrors.NewInternalServerError()
	}

	// ラウンドをAPIレスポンス形式に変換
	roundResponses := make([]*model.RoundResponse, len(rounds))
	for i, round := range rounds {
		roundResponses[i] = round.ToResponse()
	}

	return &model.RoundsResponse{
		Rounds: roundResponses,
	}, nil
}

// CompleteRound はラウンドを完了する(SQSメッセージ送信付き)
func (uc *roundUseCase) CompleteRound(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *model.RoundCompleteRequest) (*model.RoundResponse, error) {
	// ラウンドの存在確認
	_, err := uc.roundRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Errorf("ラウンド取得エラー: %v", err)
		if errors.Is(err, rounderrors.ErrRoundNotFound) {
			return nil, domainErrors.NewRoundNotFoundError()
		}
		return nil, domainErrors.NewInternalServerError()
	}

	// ユーザーの作業時間と休憩時間を取得
	workTime, breakTime, err := uc.getUserWorkAndBreakTime(ctx, userID)
	if err != nil {
		uc.logger.Errorf("ユーザー設定取得エラー: %v", err)
		return nil, domainErrors.NewInternalServerError()
	}

	// ラウンドを完了
	err = uc.roundRepo.Complete(ctx, id, req.FocusScore, workTime, breakTime)
	if err != nil {
		uc.logger.Errorf("ラウンド完了エラー: %v", err)
		if errors.Is(err, rounderrors.ErrRoundAlreadyEnded) {
			return nil, domainErrors.NewRoundAlreadyEndedError()
		}
		return nil, domainErrors.NewInternalServerError()
	}

	// 集中度スコアが設定されている場合、SQSメッセージを送信
	if req.FocusScore != nil {
		uc.logger.Infof("ラウンド完了 - SQS最適化メッセージ送信開始: RoundID=%s, FocusScore=%d",
			id.String(), *req.FocusScore)

		// 同期でSQS送信
		uc.sendRoundOptimizationMessage(ctx, userID, id, *req.FocusScore)
	} else {
		uc.logger.Warn("集中度スコアが入力されていないため、SQS最適化メッセージは送信しません")
	}

	// 完了したラウンドを取得して返す
	completedRound, err := uc.roundRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Errorf("完了ラウンド取得エラー: %v", err)
		return nil, domainErrors.NewInternalServerError()
	}

	return completedRound.ToResponse(), nil
}

// AbortRound はラウンドを中止する（SQSメッセージは送信しない）
func (uc *roundUseCase) AbortRound(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.RoundResponse, error) {
	// ラウンドの存在確認
	round, err := uc.roundRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Errorf("ラウンド取得エラー: %v", err)
		if errors.Is(err, rounderrors.ErrRoundNotFound) {
			return nil, domainErrors.NewRoundNotFoundError()
		}
		return nil, domainErrors.NewInternalServerError()
	}

	// セッションを取得してユーザーIDでアクセス権限を確認
	_, err = uc.sessionRepo.GetByID(ctx, round.SessionID, userID)
	if err != nil {
		uc.logger.Errorf("セッション取得エラー: %v", err)
		return nil, err
	}

	// ラウンドを中止する
	if err := uc.roundRepo.AbortRound(ctx, id); err != nil {
		uc.logger.Errorf("ラウンド中止エラー: %v", err)
		if errors.Is(err, rounderrors.ErrRoundAlreadyEnded) {
			return nil, domainErrors.NewRoundAlreadyEndedError()
		}
		return nil, domainErrors.NewInternalServerError()
	}

	// ❌ 中止時はSQS送信なし（セッション全体が中止される想定）
	uc.logger.Infof("ラウンド中止完了: %s (SQS送信なし - セッション中止を想定)", id.String())

	// 中止したラウンドを取得して返す
	abortedRound, err := uc.roundRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Errorf("中止ラウンド取得エラー: %v", err)
		return nil, domainErrors.NewInternalServerError()
	}

	return abortedRound.ToResponse(), nil
}

// sendRoundOptimizationMessage はラウンド最適化メッセージをSQSに送信する（同期）
func (uc *roundUseCase) sendRoundOptimizationMessage(ctx context.Context, userID, roundID uuid.UUID, focusScore int) {
	if uc.sqsClient == nil {
		uc.logger.Warn("SQSクライアントが初期化されていません。最適化メッセージは送信されません。")
		return
	}

	// 最小限のメッセージを作成
	message := model.NewRoundOptimizationMessage(userID, roundID, focusScore)

	// SQS送信実行
	err := uc.sqsClient.SendRoundOptimizationMessage(ctx, message)
	if err != nil {
		uc.logger.Errorf("ラウンド最適化メッセージ送信エラー: %v", err)
		return
	}

	uc.logger.Infof("ラウンド最適化メッセージ送信成功: %s", message.ToLogString())
}

// getUserWorkAndBreakTime はユーザーの作業時間と休憩時間を取得する
func (uc *roundUseCase) getUserWorkAndBreakTime(ctx context.Context, userID uuid.UUID) (workTime int, breakTime int, err error) {
	uc.logger.Infof("getUserWorkAndBreakTime 開始: userID=%s", userID.String())

	// UserConfigを取得（GetOrCreateではない）
	if uc.userConfigRepo != nil {
		userConfig, configErr := uc.userConfigRepo.GetUserConfig(ctx, userID)
		if configErr != nil {
			// UserConfig取得失敗時はデフォルト値 + 警告ログ
			uc.logger.Warnf("UserConfig取得エラー、デフォルト値を使用します: %v", configErr)
		} else {
			// UserConfig取得成功
			uc.logger.Infof("UserConfig取得成功: workTime=%d, breakTime=%d",
				userConfig.RoundWorkTime, userConfig.RoundBreakTime)
			return userConfig.RoundWorkTime, userConfig.RoundBreakTime, nil
		}
	} else {
		uc.logger.Warn("UserConfigRepository が nil です")
	}

	// フォールバック：デフォルト値
	defaultWorkTime := 25
	defaultBreakTime := 5
	uc.logger.Infof("デフォルト値を使用: workTime=%d, breakTime=%d", defaultWorkTime, defaultBreakTime)
	return defaultWorkTime, defaultBreakTime, nil
}
