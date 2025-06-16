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
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	httpError "github.com/tsunakit99/selfpomodoro/internal/handler"
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
		return createErrorResponse(http.StatusInternalServerError, "INTERNAL_ERROR", "サービス初期化エラー"), nil
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
	return roundHandler.routeOperation(ctx, request, userID)
}

// authenticateAndValidateUser は認証・User存在確認の統一処理
func (h *RoundHandler) authenticateAndValidateUser(ctx context.Context, request events.APIGatewayProxyRequest) (uuid.UUID, error) {
	return h.useCases.Auth.AuthenticateAndValidateUser(ctx, request)
}

// routeOperation は操作ルーティング
func (h *RoundHandler) routeOperation(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	// パスパラメータからセッションIDとラウンドIDを取得
	sessionIDStr := request.PathParameters["session_id"]
	roundIDStr := request.PathParameters["round_id"]

	// セッション関連のラウンド操作
	if strings.Contains(request.Path, "/sessions/") && strings.Contains(request.Path, "/rounds") {
		if sessionIDStr == "" {
			return createErrorResponse(http.StatusBadRequest, "MISSING_SESSION_ID", "セッションIDが指定されていません"), nil
		}
		sessionID, err := uuid.Parse(sessionIDStr)
		if err != nil {
			return createErrorResponse(http.StatusBadRequest, "INVALID_SESSION_ID", "無効なセッションID"), nil
		}

		switch request.HTTPMethod {
		case "GET":
			return h.handleGetRoundsBySession(ctx, sessionID)
		case "POST":
			return h.handleStartRound(ctx, sessionID, userID)
		default:
			return createErrorResponse(http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "メソッドが許可されていません"), nil
		}
	}

	// 個別ラウンド操作
	if roundIDStr != "" {
		roundID, err := uuid.Parse(roundIDStr)
		if err != nil {
			return createErrorResponse(http.StatusBadRequest, "INVALID_ROUND_ID", "無効なラウンドID"), nil
		}

		switch request.HTTPMethod {
		case "GET":
			return h.handleGetRound(ctx, roundID)
		case "PATCH":
			if !strings.Contains(request.Path, "/complete") {
				return createErrorResponse(http.StatusNotFound, "NOT_FOUND", "無効なパス"), nil
			}
			return h.handleCompleteRound(ctx, request, roundID, userID)
		case "POST":
			if !strings.Contains(request.Path, "/abort") {
				return createErrorResponse(http.StatusNotFound, "NOT_FOUND", "無効なパス"), nil
			}
			return h.handleAbortRound(ctx, roundID, userID)
		default:
			return createErrorResponse(http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "メソッドが許可されていません"), nil
		}
	}

	return createErrorResponse(http.StatusNotFound, "NOT_FOUND", "無効なパス"), nil
}

// handleGetRoundsBySession はセッションのラウンド一覧取得を処理
func (h *RoundHandler) handleGetRoundsBySession(ctx context.Context, sessionID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	roundsResponse, err := h.useCases.Round.GetAllRoundsBySessionID(ctx, sessionID)
	if err != nil {
		h.logger.Errorf("ラウンド一覧取得エラー: %v", err)
		return h.handleError(err), nil
	}
	return createSuccessResponse(http.StatusOK, roundsResponse), nil
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

	h.logger.Infof("ラウンド開始成功: ラウンドID=%s", roundResponse.ID.String())
	return createSuccessResponse(http.StatusCreated, roundResponse), nil
}

// handleGetRound はラウンド取得を処理
func (h *RoundHandler) handleGetRound(ctx context.Context, roundID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	roundResponse, err := h.useCases.Round.GetRound(ctx, roundID)
	if err != nil {
		h.logger.Errorf("ラウンド取得エラー: %v", err)
		return h.handleError(err), nil
	}
	return createSuccessResponse(http.StatusOK, roundResponse), nil
}

// handleCompleteRound はラウンド完了を処理
func (h *RoundHandler) handleCompleteRound(ctx context.Context, request events.APIGatewayProxyRequest, roundID uuid.UUID, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	var req model.RoundCompleteRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return createErrorResponse(http.StatusBadRequest, "INVALID_REQUEST_FORMAT", "無効なリクエスト形式"), nil
	}

	if req.FocusScore == nil {
		return createErrorResponse(http.StatusBadRequest, "VALIDATION_ERROR", "集中度スコアは必須です"), nil
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

	return createSuccessResponse(http.StatusOK, roundResponse), nil
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
	return createSuccessResponse(http.StatusOK, roundResponse), nil
}

// handleError はエラーを統一処理（error_mapper.go使用版）
func (h *RoundHandler) handleError(err error) events.APIGatewayProxyResponse {
	// error_mapper.goを使用してHTTPエラーにマッピング
	httpErr := httpError.MapErrorToHTTP(err)

	// ログ出力（サーバーエラーのみ詳細ログ）
	if httpError.IsServerError(err) {
		h.logger.Errorf("サーバーエラー: %v", err)
	} else if httpError.IsClientError(err) {
		h.logger.Warnf("クライアントエラー: %s - %s", httpErr.Code, httpErr.Message)
	}

	// 統一されたエラーレスポンス作成
	return createErrorResponse(httpErr.StatusCode, httpErr.Code, httpErr.Message)
}

// createSuccessResponse は成功レスポンスを作成
func createSuccessResponse(statusCode int, data interface{}) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(data)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    getCORSHeaders(),
		Body:       string(body),
	}
}

// createErrorResponse はエラーレスポンスを作成（統一フォーマット）
func createErrorResponse(statusCode int, code, message string) events.APIGatewayProxyResponse {
	// エラーレスポンスの統一フォーマット
	errorBody := map[string]interface{}{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	}

	body, _ := json.Marshal(errorBody)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    getCORSHeaders(),
		Body:       string(body),
	}
}

// getCORSHeaders はCORSヘッダーを取得
func getCORSHeaders() map[string]string {
	return map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type,Authorization",
		"Access-Control-Allow-Methods": "GET,POST,PATCH,DELETE,OPTIONS",
	}
}

func main() {
	lambda.Start(handler)
}
