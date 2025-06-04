package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/database"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/postgres"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/sqs"
	"github.com/tsunakit99/selfpomodoro/internal/usecase"
)

// RoundHandler はLambda用のラウンドハンドラー
type RoundHandler struct {
	roundUseCase usecase.RoundUseCase
	logger       logger.Logger
	validator    *validator.Validate
}

// handler はLambdaのエントリーポイント
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 依存関係の初期化
	cfg, err := config.Load()
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "設定読み込みエラー"), nil
	}

	appLogger, err := logger.NewLogger(cfg.LogLevel, cfg.Environment)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "ロガー初期化エラー"), nil
	}

	// PostgreSQL接続
	postgresDB, err := database.NewPostgresDB(cfg, appLogger)
	if err != nil {
		appLogger.Errorf("PostgreSQL接続エラー: %v", err)
		return errorResponse(http.StatusInternalServerError, "データベース接続エラー"), nil
	}
	defer postgresDB.Close()

	// DynamoDB接続（UserConfig用）
	var dynamoDB *database.DynamoDB
	dynamoDB, err = database.NewDynamoDB(cfg, appLogger)
	if err != nil {
		appLogger.Warnf("DynamoDB接続エラー: %v", err)
		appLogger.Warn("DynamoDBなしで続行します（デフォルト設定使用）")
		dynamoDB = nil
	}
	if dynamoDB != nil {
		defer dynamoDB.Close()
	}

	// SQS接続（最適化用）
	var sqsClient *sqs.SQSClient
	sqsClient, err = sqs.NewSQSClient(cfg, appLogger)
	if err != nil {
		appLogger.Warnf("SQS接続エラー: %v", err)
		appLogger.Warn("SQSなしで続行します（最適化機能無効）")
		sqsClient = nil
	}
	if sqsClient != nil {
		defer sqsClient.Close()

		// SQS接続確認（開発環境のみ）
		if cfg.Environment != "production" {
			if healthErr := sqsClient.HealthCheck(ctx); healthErr != nil {
				appLogger.Warnf("SQS接続確認失敗: %v", healthErr)
			}
		}
	}

	// リポジトリとユースケースの初期化
	repositoryFactory := repository.NewRepositoryFactory(postgresDB, dynamoDB, appLogger)
	roundUseCase := usecase.NewRoundUseCase(
		repositoryFactory.Round,
		repositoryFactory.Session,
		repositoryFactory.UserConfig, // DynamoDB UserConfig追加
		sqsClient,
		appLogger,
	)

	// ハンドラーの初期化
	roundHandler := &RoundHandler{
		roundUseCase: roundUseCase,
		logger:       appLogger,
		validator:    validator.New(),
	}

	// ユーザーID取得
	userID, err := getUserIDFromRequest(request)
	if err != nil {
		appLogger.Errorf("ユーザーID取得エラー: %v", err)
		return errorResponse(http.StatusUnauthorized, "認証エラー"), nil
	}

	// パスによるルーティング判定
	if strings.Contains(request.Path, "/sessions/") && strings.Contains(request.Path, "/rounds") {
		// /sessions/{session_id}/rounds パターン
		return roundHandler.handleSessionRounds(ctx, request, userID)
	} else if strings.Contains(request.Path, "/rounds/") {
		// /rounds/{round_id} パターン
		return roundHandler.handleIndividualRound(ctx, request, userID)
	}

	return errorResponse(http.StatusNotFound, "無効なパス"), nil
}

// getUserIDFromRequest はリクエストからユーザーIDを取得する
func getUserIDFromRequest(request events.APIGatewayProxyRequest) (uuid.UUID, error) {
	// 開発環境用の簡易認証（dev-tokenの場合）
	authHeader := request.Headers["Authorization"]
	if authHeader == "" {
		authHeader = request.Headers["authorization"] // 小文字の場合もチェック
	}

	if authHeader == "Bearer dev-token" {
		// 開発用固定ユーザーID
		return uuid.Parse("00000000-0000-0000-0000-000000000001")
	}

	return uuid.Nil, fmt.Errorf("認証情報が見つかりません")
}

// handleSessionRounds はセッション関連のラウンド操作を処理
func (h *RoundHandler) handleSessionRounds(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	// session_idの取得
	sessionIDStr := request.PathParameters["session_id"]
	if sessionIDStr == "" {
		return errorResponse(http.StatusBadRequest, "セッションIDが指定されていません"), nil
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "無効なセッションID"), nil
	}

	switch request.HTTPMethod {
	case "GET":
		// セッションのラウンド一覧取得
		return h.handleGetRoundsBySession(ctx, sessionID)
	case "POST":
		// ラウンド開始（UserConfig考慮）
		return h.handleStartRound(ctx, sessionID, userID)
	default:
		return errorResponse(http.StatusMethodNotAllowed, "メソッドが許可されていません"), nil
	}
}

// handleIndividualRound は個別ラウンド操作を処理
func (h *RoundHandler) handleIndividualRound(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	// round_idの取得
	roundIDStr := request.PathParameters["round_id"]
	if roundIDStr == "" {
		return errorResponse(http.StatusBadRequest, "ラウンドIDが指定されていません"), nil
	}

	roundID, err := uuid.Parse(roundIDStr)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "無効なラウンドID"), nil
	}

	switch request.HTTPMethod {
	case "GET":
		// ラウンド取得
		return h.handleGetRound(ctx, roundID)
	case "PATCH":
		// ラウンド完了（DynamoDB設定値使用）
		return h.handleCompleteRound(ctx, request, roundID, userID)
	case "POST":
		if strings.Contains(request.Path, "/abort") {
			// ラウンド中止
			return h.handleAbortRound(ctx, roundID, userID)
		}
		return errorResponse(http.StatusNotFound, "無効なパス"), nil
	default:
		return errorResponse(http.StatusMethodNotAllowed, "メソッドが許可されていません"), nil
	}
}

// handleGetRoundsBySession はセッションのラウンド一覧取得を処理
func (h *RoundHandler) handleGetRoundsBySession(ctx context.Context, sessionID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	roundsResponse, err := h.roundUseCase.GetAllRoundsBySessionID(ctx, sessionID)
	if err != nil {
		h.logger.Errorf("ラウンド一覧取得エラー: %v", err)
		return errorResponse(http.StatusInternalServerError, "ラウンド一覧の取得に失敗しました"), nil
	}
	return successResponse(http.StatusOK, roundsResponse), nil
}

// handleStartRound はラウンド開始を処理（ログ強化版）
func (h *RoundHandler) handleStartRound(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	h.logger.Infof("ラウンド開始要求: セッションID=%s, ユーザーID=%s", sessionID.String(), userID.String())

	var req model.RoundCreateRequest
	roundResponse, err := h.roundUseCase.StartRound(ctx, sessionID, userID, &req)
	if err != nil {
		if isSessionNotFoundError(err) {
			return errorResponse(http.StatusNotFound, "セッションが見つかりません"), nil
		}
		if isSessionAccessDeniedError(err) {
			return errorResponse(http.StatusForbidden, "このセッションへのアクセス権限がありません"), nil
		}
		h.logger.Errorf("ラウンド開始エラー: %v", err)
		return errorResponse(http.StatusInternalServerError, "ラウンド開始に失敗しました"), nil
	}

	h.logger.Infof("ラウンド開始成功: ラウンドID=%s, ラウンド順序=%d", roundResponse.ID.String(), roundResponse.RoundOrder)
	return successResponse(http.StatusCreated, roundResponse), nil
}

// handleGetRound はラウンド取得を処理
func (h *RoundHandler) handleGetRound(ctx context.Context, roundID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	roundResponse, err := h.roundUseCase.GetRound(ctx, roundID)
	if err != nil {
		if isRoundNotFoundError(err) {
			return errorResponse(http.StatusNotFound, "ラウンドが見つかりません"), nil
		}
		h.logger.Errorf("ラウンド取得エラー: %v", err)
		return errorResponse(http.StatusInternalServerError, "ラウンド取得に失敗しました"), nil
	}
	return successResponse(http.StatusOK, roundResponse), nil
}

// handleCompleteRound はラウンド完了を処理（DynamoDB設定値使用版, SQS最適化メッセージ送信付き）
func (h *RoundHandler) handleCompleteRound(ctx context.Context, request events.APIGatewayProxyRequest, roundID uuid.UUID, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	if !strings.Contains(request.Path, "/complete") {
		return errorResponse(http.StatusNotFound, "無効なパス"), nil
	}

	var req model.RoundCompleteRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return errorResponse(http.StatusBadRequest, "無効なリクエスト形式"), nil
	}

	if err := h.validator.Struct(req); err != nil {
		return errorResponse(http.StatusBadRequest, "集中度スコアは0から100の間である必要があります"), nil
	}

	h.logger.Infof("ラウンド完了要求: ラウンドID=%s, 集中度スコア=%v", roundID.String(), req.FocusScore)

	roundResponse, err := h.roundUseCase.CompleteRound(ctx, roundID, userID, &req)
	if err != nil {
		if isRoundNotFoundError(err) {
			return errorResponse(http.StatusNotFound, "ラウンドが見つかりません"), nil
		}
		h.logger.Errorf("ラウンド完了エラー: %v", err)
		return errorResponse(http.StatusInternalServerError, "ラウンド完了に失敗しました"), nil
	}

	// SQS送信の有無をログ出力
	if req.FocusScore != nil {
		h.logger.Infof("ラウンド完了成功: ラウンドID=%s, 作業時間=%d分, 休憩時間=%d分 (SQS最適化メッセージ送信済み)",
			roundID.String(), *roundResponse.WorkTime, *roundResponse.BreakTime)
	} else {
		h.logger.Infof("ラウンド完了成功: ラウンドID=%s, 作業時間=%d分, 休憩時間=%d分 (集中度スコア未入力・SQS送信なし)",
			roundID.String(), *roundResponse.WorkTime, *roundResponse.BreakTime)
	}

	return successResponse(http.StatusOK, roundResponse), nil
}

// handleAbortRound はラウンド中止を処理（SQS最適化メッセージは送信しない）
func (h *RoundHandler) handleAbortRound(ctx context.Context, roundID uuid.UUID, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	h.logger.Infof("ラウンド中止要求: ラウンドID=%s", roundID.String())

	roundResponse, err := h.roundUseCase.AbortRound(ctx, roundID, userID)
	if err != nil {
		if isRoundNotFoundError(err) {
			return errorResponse(http.StatusNotFound, "ラウンドが見つかりません"), nil
		}
		h.logger.Errorf("ラウンド中止エラー: %v", err)
		return errorResponse(http.StatusInternalServerError, "ラウンド中止に失敗しました"), nil
	}

	h.logger.Infof("ラウンド中止成功: ラウンドID=%s (SQS送信なし)", roundID.String())
	return successResponse(http.StatusOK, roundResponse), nil
}

// ヘルパー関数
func isSessionNotFoundError(err error) bool {
	return strings.Contains(err.Error(), postgres.ErrSessionNotFound.Error())
}

func isSessionAccessDeniedError(err error) bool {
	return strings.Contains(err.Error(), postgres.ErrSessionAccessDenied.Error())
}

func isRoundNotFoundError(err error) bool {
	return strings.Contains(err.Error(), postgres.ErrRoundNotFound.Error())
}

func successResponse(statusCode int, data interface{}) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(data)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type,Authorization",
			"Access-Control-Allow-Methods": "GET,POST,PATCH,DELETE,OPTIONS",
		},
		Body: string(body),
	}
}

func errorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(map[string]string{"error": message})
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type,Authorization",
			"Access-Control-Allow-Methods": "GET,POST,PATCH,DELETE,OPTIONS",
		},
		Body: string(body),
	}
}

func main() {
	lambda.Start(handler)
}
