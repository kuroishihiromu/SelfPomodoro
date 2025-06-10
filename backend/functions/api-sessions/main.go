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
	domainErrors "github.com/tsunakit99/selfpomodoro/internal/domain/errors"
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
		return errorResponse(http.StatusInternalServerError, "サービス初期化エラー"), nil
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
	switch request.HTTPMethod {
	case "GET":
		return h.handleGetSessions(ctx, request, userID)
	case "POST":
		return h.handleStartSession(ctx, userID)
	case "PATCH":
		return h.handleCompleteSession(ctx, request, userID)
	case "DELETE":
		return h.handleDeleteSession(ctx, request, userID)
	default:
		return errorResponse(http.StatusMethodNotAllowed, "メソッドが許可されていません"), nil
	}
}

// handleGetSessions はセッション取得処理（一覧または個別）
func (h *SessionHandler) handleGetSessions(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	sessionIDStr := request.PathParameters["session_id"]

	if sessionIDStr == "" {
		// セッション一覧取得
		sessionsResponse, err := h.useCases.Session.GetAllSessions(ctx, userID)
		if err != nil {
			h.logger.Errorf("セッション一覧取得エラー: %v", err)
			return h.handleError(err), nil
		}
		return successResponse(http.StatusOK, sessionsResponse), nil
	} else {
		// 個別セッション取得
		sessionID, err := uuid.Parse(sessionIDStr)
		if err != nil {
			return errorResponse(http.StatusBadRequest, "無効なセッションID"), nil
		}

		sessionResponse, err := h.useCases.Session.GetSession(ctx, sessionID, userID)
		if err != nil {
			h.logger.Errorf("セッション取得エラー: %v", err)
			return h.handleError(err), nil
		}
		return successResponse(http.StatusOK, sessionResponse), nil
	}
}

// handleStartSession はセッション開始処理
func (h *SessionHandler) handleStartSession(ctx context.Context, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	h.logger.Infof("セッション開始要求: ユーザーID=%s", userID.String())

	sessionResponse, err := h.useCases.Session.StartSession(ctx, userID)
	if err != nil {
		h.logger.Errorf("セッション開始エラー: %v", err)
		return h.handleError(err), nil
	}

	h.logger.Infof("セッション開始成功: セッションID=%s", sessionResponse.ID.String())
	return successResponse(http.StatusCreated, sessionResponse), nil
}

// handleCompleteSession はセッション完了処理
func (h *SessionHandler) handleCompleteSession(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	sessionIDStr := request.PathParameters["session_id"]
	if sessionIDStr == "" {
		return errorResponse(http.StatusBadRequest, "セッションIDが指定されていません"), nil
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "無効なセッションID"), nil
	}

	// パスが /complete かどうか確認
	if !strings.Contains(request.Path, "/complete") {
		return errorResponse(http.StatusNotFound, "無効なパス"), nil
	}

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

	return successResponse(http.StatusOK, sessionResponse), nil
}

// handleDeleteSession はセッション削除処理
func (h *SessionHandler) handleDeleteSession(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	sessionIDStr := request.PathParameters["session_id"]
	if sessionIDStr == "" {
		return errorResponse(http.StatusBadRequest, "セッションIDが指定されていません"), nil
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "無効なセッションID"), nil
	}

	err = h.useCases.Session.DeleteSession(ctx, sessionID, userID)
	if err != nil {
		h.logger.Errorf("セッション削除エラー: %v", err)
		return h.handleError(err), nil
	}

	return successResponse(http.StatusOK, map[string]string{"message": "セッションが削除されました"}), nil
}

// handleError はドメインエラーを統一処理
func (h *SessionHandler) handleError(err error) events.APIGatewayProxyResponse {
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
