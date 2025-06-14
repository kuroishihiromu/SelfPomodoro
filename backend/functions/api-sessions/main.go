package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/container"
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

// SessionHandler はDI Container使用版のセッションハンドラー
type SessionHandler struct {
	useCases *usecase.UseCases
	logger   logger.Logger
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
	sessionHandler := &SessionHandler{
		useCases: useCases,
		logger:   logger,
	}

	// 4. 認証・User存在確認（統一処理）
	userID, err := sessionHandler.authenticateAndValidateUser(ctx, request)
	if err != nil {
		return sessionHandler.handleError(err), nil
	}

	// 5. 操作ルーティング
	return sessionHandler.routeOperation(ctx, request, userID)
}

// authenticateAndValidateUser は認証・User存在確認の統一処理
func (h *SessionHandler) authenticateAndValidateUser(ctx context.Context, request events.APIGatewayProxyRequest) (uuid.UUID, error) {
	return h.useCases.Auth.AuthenticateAndValidateUser(ctx, request)
}

// routeOperation は操作ルーティング
func (h *SessionHandler) routeOperation(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	// パスパラメータからセッションIDを取得
	sessionIDStr := request.PathParameters["session_id"]

	switch request.HTTPMethod {
	case "GET":
		if sessionIDStr == "" {
			return h.handleGetSessions(ctx, userID)
		}
		sessionID, err := uuid.Parse(sessionIDStr)
		if err != nil {
			return createErrorResponse(http.StatusBadRequest, "INVALID_SESSION_ID", "無効なセッションID"), nil
		}
		return h.handleGetSession(ctx, sessionID, userID)

	case "POST":
		return h.handleStartSession(ctx, userID)

	case "PATCH":
		if sessionIDStr == "" {
			return createErrorResponse(http.StatusBadRequest, "MISSING_SESSION_ID", "セッションIDが指定されていません"), nil
		}
		if !strings.Contains(request.Path, "/complete") {
			return createErrorResponse(http.StatusNotFound, "NOT_FOUND", "無効なパス"), nil
		}
		sessionID, err := uuid.Parse(sessionIDStr)
		if err != nil {
			return createErrorResponse(http.StatusBadRequest, "INVALID_SESSION_ID", "無効なセッションID"), nil
		}
		return h.handleCompleteSession(ctx, sessionID, userID)

	case "DELETE":
		if sessionIDStr == "" {
			return createErrorResponse(http.StatusBadRequest, "MISSING_SESSION_ID", "セッションIDが指定されていません"), nil
		}
		sessionID, err := uuid.Parse(sessionIDStr)
		if err != nil {
			return createErrorResponse(http.StatusBadRequest, "INVALID_SESSION_ID", "無効なセッションID"), nil
		}
		return h.handleDeleteSession(ctx, sessionID, userID)

	default:
		return createErrorResponse(http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "メソッドが許可されていません"), nil
	}
}

// handleGetSessions はセッション一覧取得を処理
func (h *SessionHandler) handleGetSessions(ctx context.Context, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	sessionsResponse, err := h.useCases.Session.GetAllSessions(ctx, userID)
	if err != nil {
		h.logger.Errorf("セッション一覧取得エラー: %v", err)
		return h.handleError(err), nil
	}
	return createSuccessResponse(http.StatusOK, sessionsResponse), nil
}

// handleGetSession は個別セッション取得を処理
func (h *SessionHandler) handleGetSession(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	sessionResponse, err := h.useCases.Session.GetSession(ctx, sessionID, userID)
	if err != nil {
		h.logger.Errorf("セッション取得エラー: %v", err)
		return h.handleError(err), nil
	}
	return createSuccessResponse(http.StatusOK, sessionResponse), nil
}

// handleStartSession はセッション開始を処理
func (h *SessionHandler) handleStartSession(ctx context.Context, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	h.logger.Infof("セッション開始要求: ユーザーID=%s", userID.String())

	sessionResponse, err := h.useCases.Session.StartSession(ctx, userID)
	if err != nil {
		h.logger.Errorf("セッション開始エラー: %v", err)
		return h.handleError(err), nil
	}

	h.logger.Infof("セッション開始成功: セッションID=%s", sessionResponse.ID.String())
	return createSuccessResponse(http.StatusCreated, sessionResponse), nil
}

// handleCompleteSession はセッション完了を処理
func (h *SessionHandler) handleCompleteSession(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	h.logger.Infof("セッション完了要求: セッションID=%s", sessionID.String())

	sessionResponse, err := h.useCases.Session.CompleteSession(ctx, sessionID, userID)
	if err != nil {
		h.logger.Errorf("セッション完了エラー: %v", err)
		return h.handleError(err), nil
	}

	// SQS送信ログ出力
	if sessionResponse.RoundCount != nil && *sessionResponse.RoundCount > 0 {
		h.logger.Infof("セッション完了成功: SessionID=%s, RoundCount=%d, AvgFocus=%.1f, TotalWork=%dmin (SQS送信済み)",
			sessionID.String(), *sessionResponse.RoundCount, *sessionResponse.AverageFocus, *sessionResponse.TotalWorkMin)
	} else {
		h.logger.Infof("セッション完了成功: SessionID=%s, RoundCount=0 (SQS送信なし)",
			sessionID.String())
	}

	return createSuccessResponse(http.StatusOK, sessionResponse), nil
}

// handleDeleteSession はセッション削除を処理
func (h *SessionHandler) handleDeleteSession(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	err := h.useCases.Session.DeleteSession(ctx, sessionID, userID)
	if err != nil {
		h.logger.Errorf("セッション削除エラー: %v", err)
		return h.handleError(err), nil
	}

	return createSuccessResponse(http.StatusOK, map[string]string{"message": "セッションが削除されました"}), nil
}

// handleError はエラーを統一処理（error_mapper.go使用版）
func (h *SessionHandler) handleError(err error) events.APIGatewayProxyResponse {
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
