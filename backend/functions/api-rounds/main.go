package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/container"
	domainErrors "github.com/tsunakit99/selfpomodoro/internal/domain/errors"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/usecase"
)

// Global container (Lambda再利用最適化)
var globalContainer container.Container

// init はLambda init phaseで実行
func init() {
	globalContainer = container.NewLambdaContainer()
}

// RoundHandler はDI Container使用版のラウンドハンドラー
type RoundHandler struct {
	useCases  *usecase.UseCases
	logger    logger.Logger
	validator *validator.Validate
}

// handler はLambdaのエントリーポイント（DI Container版）
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 1. Container初期化
	if err := globalContainer.Initialize(ctx); err != nil {
		return errorResponse(http.StatusInternalServerError, "サービス初期化エラー"), nil
	}

	// 2. Dependencies取得（Infrastructure依存なし！）
	useCases := globalContainer.GetUseCases()
	logger := globalContainer.GetLogger()

	// 3. Handler初期化
	roundHandler := &RoundHandler{
		useCases:  useCases,
		logger:    logger,
		validator: validator.New(),
	}

	// 4. 認証・User存在確認（統一処理）
	userID, err := roundHandler.authenticateAndValidateUser(ctx, request)
	if err != nil {
		return roundHandler.handleError(err), nil
	}

	// 5. パスによるルーティング判定
	if strings.Contains(request.Path, "/sessions/") && strings.Contains(request.Path, "/rounds") {
		// /sessions/{session_id}/rounds パターン
		return roundHandler.handleSessionRounds(ctx, request, userID)
	} else if strings.Contains(request.Path, "/rounds/") {
		// /rounds/{round_id} パターン
		return roundHandler.handleIndividualRound(ctx, request, userID)
	}

	return errorResponse(http.StatusNotFound, "無効なパス"), nil
}

// authenticateAndValidateUser は認証・User存在確認の統一処理
func (h *RoundHandler) authenticateAndValidateUser(ctx context.Context, request events.APIGatewayProxyRequest) (uuid.UUID, error) {
	return h.useCases.Auth.AuthenticateAndValidateUser(ctx, request)
}

// handleSessionRounds はセッション関連のラウンド操作を処理
func (h *RoundHandler) handleSessionRounds(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
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
		// ラウンド開始
		return h.handleStartRound(ctx, sessionID, userID)
	default:
		return errorResponse(http.StatusMethodNotAllowed, "メソッドが許可されていません"), nil
	}
}

// handleIndividualRound は個別ラウンド操作を処理
func (h *RoundHandler) handleIndividualRound(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
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
		// ラウンド完了
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
	roundsResponse, err := h.useCases.Round.GetAllRoundsBySessionID(ctx, sessionID)
	if err != nil {
		h.logger.Errorf("ラウンド一覧取得エラー: %v", err)
		return h.handleError(err), nil
	}
	return successResponse(http.StatusOK, roundsResponse), nil
}

// handleStartRound はラウンド開始を処理
func (h *RoundHandler) handleStartRound(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	h.logger.Infof("ラウンド開始要求: セッションID=%s, ユーザーID=%s", sessionID.String(), userID.String())

	var req model.RoundCreateRequest
	roundResponse, err := h.useCases.Round.StartRound(ctx, sessionID, userID, &req)
	if err != nil {
		h.logger.Errorf("ラウンド開始エラー: %v", err)
		return h.handleError(err), nil
	}

	h.logger.Infof("ラウンド開始成功: ラウンドID=%s, ラウンド順序=%d", roundResponse.ID.String(), roundResponse.RoundOrder)
	return successResponse(http.StatusCreated, roundResponse), nil
}

// handleGetRound はラウンド取得を処理
func (h *RoundHandler) handleGetRound(ctx context.Context, roundID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	roundResponse, err := h.useCases.Round.GetRound(ctx, roundID)
	if err != nil {
		h.logger.Errorf("ラウンド取得エラー: %v", err)
		return h.handleError(err), nil
	}
	return successResponse(http.StatusOK, roundResponse), nil
}

// handleCompleteRound はラウンド完了を処理
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

	roundResponse, err := h.useCases.Round.CompleteRound(ctx, roundID, userID, &req)
	if err != nil {
		h.logger.Errorf("ラウンド完了エラー: %v", err)
		return h.handleError(err), nil
	}

	// SQS送信ログ出力
	if req.FocusScore != nil {
		h.logger.Infof("ラウンド完了成功: ラウンドID=%s, 作業時間=%d分, 休憩時間=%d分 (SQS送信済み)",
			roundID.String(), *roundResponse.WorkTime, *roundResponse.BreakTime)
	} else {
		h.logger.Infof("ラウンド完了成功: ラウンドID=%s, 作業時間=%d分, 休憩時間=%d分 (SQS送信なし)",
			roundID.String(), *roundResponse.WorkTime, *roundResponse.BreakTime)
	}

	return successResponse(http.StatusOK, roundResponse), nil
}

// handleAbortRound はラウンド中止を処理
func (h *RoundHandler) handleAbortRound(ctx context.Context, roundID uuid.UUID, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	h.logger.Infof("ラウンド中止要求: ラウンドID=%s", roundID.String())

	roundResponse, err := h.useCases.Round.AbortRound(ctx, roundID, userID)
	if err != nil {
		h.logger.Errorf("ラウンド中止エラー: %v", err)
		return h.handleError(err), nil
	}

	h.logger.Infof("ラウンド中止成功: ラウンドID=%s (SQS送信なし)", roundID.String())
	return successResponse(http.StatusOK, roundResponse), nil
}

// handleError はドメインエラーを統一処理
func (h *RoundHandler) handleError(err error) events.APIGatewayProxyResponse {
	if appErr, ok := err.(*domainErrors.AppError); ok {
		return errorResponse(appErr.Status, appErr.Error())
	}

	h.logger.Errorf("予期しないエラータイプ: %T, %v", err, err)
	return errorResponse(http.StatusInternalServerError, "内部エラーが発生しました")
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
