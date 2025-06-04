package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/postgres"
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
	if err != nil && !errors.Is(err, postgres.ErrNoRoundsInSession) {
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
		return nil, err
	}

	return round.ToResponse(), nil
}

// GetAllRoundsBySessionID はセッションIDに紐づくすべてのラウンドを取得する
func (uc *roundUseCase) GetAllRoundsBySessionID(ctx context.Context, sessionID uuid.UUID) (*model.RoundsResponse, error) {
	rounds, err := uc.roundRepo.GetAllBySessionID(ctx, sessionID)
	if err != nil {
		uc.logger.Errorf("ラウンド一覧取得エラー: %v", err)
		return nil, err
	}

	// ラウンドをAPIレスポンス形式に変換
	roundResponses := make([]*model.RoundResponse, len(rounds))
	for i, round := range rounds {
		roundResponses[i] = round.ToResponse()
	}

	return &model.RoundsResponse{Rounds: roundResponses}, nil
}

// CompleteRound はラウンドを完了する(SQSメッセージ送信付き)
func (uc *roundUseCase) CompleteRound(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *model.RoundCompleteRequest) (*model.RoundResponse, error) {
	// ラウンドを取得して存在確認とセッションIDの取得
	round, err := uc.roundRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Errorf("ラウンド取得エラー: %v", err)
		return nil, err
	}

	// セッションを取得してユーザーIDを確認
	_, err = uc.sessionRepo.GetByID(ctx, round.SessionID, userID)
	if err != nil {
		uc.logger.Errorf("セッション取得エラー: %v", err)
		return nil, err
	}

	// ユーザーの作業時間と休憩時間を取得
	workTime, breakTime, err := uc.getUserWorkAndBreakTime(ctx, userID)
	if err != nil {
		uc.logger.Errorf("ユーザーの作業時間と休憩時間取得エラー: %v", err)
		return nil, err
	}

	// ラウンドを完了する
	if err := uc.roundRepo.Complete(ctx, id, req.FocusScore, workTime, breakTime); err != nil {
		uc.logger.Errorf("ラウンド完了エラー: %v", err)
		return nil, err
	}

	// ✅ 完了時のみSQS送信（集中度スコアが入力されている場合）
	if req.FocusScore != nil {
		uc.logger.Infof("ラウンド完了 - SQS最適化メッセージ送信開始: RoundID=%s, FocusScore=%d",
			id.String(), *req.FocusScore)

		// 同期でSQS送信
		uc.sendRoundOptimizationMessage(ctx, userID, id, *req.FocusScore)
	} else {
		uc.logger.Warn("集中度スコアが入力されていないため、SQS最適化メッセージは送信しません")
	}

	// 更新後のラウンドを取得
	updatedRound, err := uc.roundRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Errorf("完了後のラウンド取得エラー: %v", err)
		return nil, err
	}

	// TODO: 最適化アルゴリズムを実行（後で実装）

	// ラウンドをAPIレスポンス形式に変換
	return updatedRound.ToResponse(), nil
}

// AbortRound はラウンドを中止する（SQSメッセージは送信しない）
func (uc *roundUseCase) AbortRound(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.RoundResponse, error) {
	// ラウンドを取得して存在確認とセッションIDの取得
	round, err := uc.roundRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Errorf("ラウンド取得エラー: %v", err)
		return nil, err
	}

	// セッションを取得してユーザーIDでアクセス権限を確認
	_, err = uc.sessionRepo.GetByID(ctx, round.SessionID, userID)
	if err != nil {
		uc.logger.Errorf("セッション取得エラー: %v", err)
		return nil, err
	}

	// ラウンドを中止する
	if err := uc.roundRepo.AbortRound(ctx, id); err != nil {
		uc.logger.Errorf("ラウンドスキップエラー: %v", err)
		return nil, err
	}

	// ❌ 中止時はSQS送信なし（セッション全体が中止される想定）
	uc.logger.Infof("ラウンド中止完了: %s (SQS送信なし - セッション中止を想定)", id.String())

	// 更新後のラウンドを取得
	updatedRound, err := uc.roundRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Errorf("スキップ後のラウンド取得エラー: %v", err)
		return nil, err
	}

	// ラウンドをAPIレスポンス形式に変換
	return updatedRound.ToResponse(), nil
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

	// DynamoDBからユーザー設定を取得
	if uc.userConfigRepo != nil {
		uc.logger.Info("UserConfigRepository が利用可能です")
		userConfig, configErr := uc.userConfigRepo.GetOrCreateUserConfig(ctx, userID)
		if configErr != nil {
			// ユーザー設定が存在しない場合はデフォルト値を使用
			uc.logger.Warnf("ユーザー設定取得エラー、デフォルト値を使用します: %v", configErr)
		} else {
			// ユーザー設定が存在する場合はそれを使用
			uc.logger.Infof("DynamoDB設定値取得成功: workTime=%d, breakTime=%d", userConfig.RoundWorkTime, userConfig.RoundBreakTime)
			return userConfig.RoundWorkTime, userConfig.RoundBreakTime, nil
		}
	} else {
		uc.logger.Warn("UserConfigRepository が nil です")
	}

	// デフォルト値を使用
	defaultWorkTime := 25
	defaultBreakTime := 5
	uc.logger.Infof("デフォルト値を使用: workTime=%d, breakTime=%d", defaultWorkTime, defaultBreakTime)
	return defaultWorkTime, defaultBreakTime, nil
}
