package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/database"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/postgres"
	"github.com/tsunakit99/selfpomodoro/internal/usecase"
)

// SessionHandler はLambda用のセッションハンドラー
type SessionHandler struct {
	sessionUseCase usecase.SessionUseCase
	logger         logger.Logger
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

	// リポジトリとユースケースの初期化
	repositoryFactory := repository.NewRepositoryFactory(postgresDB, dynamoDB, appLogger)
	sessionUseCase := usecase.NewSessionUseCase(
		repositoryFactory.Session,
		repositoryFactory.Round,
		repositoryFactory.UserConfig, // DynamoDB UserConfig追加
		appLogger,
	)

	// ハンドラーの初期化
	sessionHandler := &SessionHandler{
		sessionUseCase: sessionUseCase,
		logger:         appLogger,
	}

	// ユーザーID取得
	userID, err := getUserIDFromRequest(request)
	if err != nil {
		appLogger.Errorf("ユーザーID取得エラー: %v", err)
		return errorResponse(http.StatusUnauthorized, "認証エラー"), nil
	}

	// ルーティング
	switch request.HTTPMethod {
	case "GET":
		return sessionHandler.handleGetSessions(ctx, request, userID)
	case "POST":
		return sessionHandler.handleStartSession(ctx, userID)
	case "PATCH":
		return sessionHandler.handleCompleteSession(ctx, request, userID)
	case "DELETE":
		return sessionHandler.handleDeleteSession(ctx, request, userID)
	default:
		return errorResponse(http.StatusMethodNotAllowed, "メソッドが許可されていません"), nil
	}
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

	// 本番環境用Cognito認証（将来実装）
	// if claims, exists := request.RequestContext.Authorizer["claims"]; exists {
	//     claimsMap := claims.(map[string]interface{})
	//     if sub, ok := claimsMap["sub"].(string); ok {
	//         return uuid.Parse(sub)
	//     }
	// }

	return uuid.Nil, fmt.Errorf("認証情報が見つかりません")
}

// handleGetSessions はセッション取得を処理（一覧または個別）
func (h *SessionHandler) handleGetSessions(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	// パスパラメータから session_id を確認
	sessionIDStr := request.PathParameters["session_id"]

	if sessionIDStr == "" {
		// セッション一覧取得
		sessionsResponse, err := h.sessionUseCase.GetAllSessions(ctx, userID)
		if err != nil {
			h.logger.Errorf("セッション一覧取得エラー: %v", err)
			return errorResponse(http.StatusInternalServerError, "セッション一覧の取得に失敗しました"), nil
		}
		return successResponse(http.StatusOK, sessionsResponse), nil
	} else {
		// 個別セッション取得
		sessionID, err := uuid.Parse(sessionIDStr)
		if err != nil {
			return errorResponse(http.StatusBadRequest, "無効なセッションID"), nil
		}

		sessionResponse, err := h.sessionUseCase.GetSession(ctx, sessionID, userID)
		if err != nil {
			if isSessionNotFoundError(err) {
				return errorResponse(http.StatusNotFound, "セッションが見つかりません"), nil
			}
			if isSessionAccessDeniedError(err) {
				return errorResponse(http.StatusForbidden, "このセッションへのアクセス権限がありません"), nil
			}
			h.logger.Errorf("セッション取得エラー: %v", err)
			return errorResponse(http.StatusInternalServerError, "セッション取得に失敗しました"), nil
		}
		return successResponse(http.StatusOK, sessionResponse), nil
	}
}

// handleStartSession はセッション開始を処理（UserConfig確認付き）
func (h *SessionHandler) handleStartSession(ctx context.Context, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	h.logger.Infof("セッション開始要求: ユーザーID=%s", userID.String())

	sessionResponse, err := h.sessionUseCase.StartSession(ctx, userID)
	if err != nil {
		h.logger.Errorf("セッション開始エラー: %v", err)
		return errorResponse(http.StatusInternalServerError, "セッション開始に失敗しました"), nil
	}

	h.logger.Infof("セッション開始成功: セッションID=%s", sessionResponse.ID.String())
	return successResponse(http.StatusCreated, sessionResponse), nil
}

// handleCompleteSession はセッション完了を処理
func (h *SessionHandler) handleCompleteSession(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	// パスパラメータからセッションIDを取得
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

	sessionResponse, err := h.sessionUseCase.CompleteSession(ctx, sessionID, userID)
	if err != nil {
		if isSessionNotFoundError(err) {
			return errorResponse(http.StatusNotFound, "セッションが見つかりません"), nil
		}
		if isSessionAccessDeniedError(err) {
			return errorResponse(http.StatusForbidden, "このセッションへのアクセス権限がありません"), nil
		}
		h.logger.Errorf("セッション完了エラー: %v", err)
		return errorResponse(http.StatusInternalServerError, "セッション完了に失敗しました"), nil
	}

	h.logger.Infof("セッション完了成功: セッションID=%s", sessionID.String())
	return successResponse(http.StatusOK, sessionResponse), nil
}

// handleDeleteSession はセッション削除を処理
func (h *SessionHandler) handleDeleteSession(ctx context.Context, request events.APIGatewayProxyRequest, userID uuid.UUID) (events.APIGatewayProxyResponse, error) {
	sessionIDStr := request.PathParameters["session_id"]
	if sessionIDStr == "" {
		return errorResponse(http.StatusBadRequest, "セッションIDが指定されていません"), nil
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "無効なセッションID"), nil
	}

	err = h.sessionUseCase.DeleteSession(ctx, sessionID, userID)
	if err != nil {
		if isSessionNotFoundError(err) {
			return errorResponse(http.StatusNotFound, "セッションが見つかりません"), nil
		}
		if isSessionAccessDeniedError(err) {
			return errorResponse(http.StatusForbidden, "このセッションへのアクセス権限がありません"), nil
		}
		h.logger.Errorf("セッション削除エラー: %v", err)
		return errorResponse(http.StatusInternalServerError, "セッション削除に失敗しました"), nil
	}

	return successResponse(http.StatusOK, map[string]string{"message": "セッションが削除されました"}), nil
}

// ヘルパー関数
func isSessionNotFoundError(err error) bool {
	return strings.Contains(err.Error(), postgres.ErrSessionNotFound.Error())
}

func isSessionAccessDeniedError(err error) bool {
	return strings.Contains(err.Error(), postgres.ErrSessionAccessDenied.Error())
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
