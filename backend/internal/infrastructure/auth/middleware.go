package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// AuthMiddleware は認証処理を行うミドルウェア
type AuthMiddleware struct {
	jwtValidator *CognitoJWTValidator
	config       *config.Config
	logger       logger.Logger
}

// NewAuthMiddleware は新しい認証ミドルウェアを作成する
func NewAuthMiddleware(cfg *config.Config, logger logger.Logger) *AuthMiddleware {
	// JWT Validator の初期化
	validatorConfig := &CognitoJWTValidatorConfig{
		UserPoolID:   cfg.CognitoUserPoolID,
		ClientID:     cfg.CognitoClientID,
		Region:       cfg.AWSRegion,
		CacheTimeout: 30 * time.Minute, // 30分キャッシュ
		HTTPTimeout:  10 * time.Second,
	}

	jwtValidator := NewCognitoJWTValidator(validatorConfig, logger)

	return &AuthMiddleware{
		jwtValidator: jwtValidator,
		config:       cfg,
		logger:       logger,
	}
}

// AuthenticateRequest はリクエストを認証し、ユーザーIDを返す
func (m *AuthMiddleware) AuthenticateRequest(request events.APIGatewayProxyRequest) (uuid.UUID, error) {
	// Authorization ヘッダーの取得
	authHeader := m.getAuthorizationHeader(request)
	if authHeader == "" {
		m.logger.Warn("Authorizationヘッダーが見つかりません")
		return uuid.Nil, NewAuthError("MISSING_AUTH", "Authorizationヘッダーが必要です", ErrTokenNotFound)
	}

	// 開発環境での後方互換性（dev-token）
	if m.isDevelopmentEnvironment() && authHeader == "Bearer dev-token" {
		m.logger.Info("開発環境: dev-token認証を使用")
		return m.getDevTokenUserID()
	}

	// Cognito JWT認証
	return m.authenticateWithCognito(authHeader)
}

// getAuthorizationHeader は大文字小文字を考慮してAuthorizationヘッダーを取得する
func (m *AuthMiddleware) getAuthorizationHeader(request events.APIGatewayProxyRequest) string {
	// 標準的なケース
	if auth := request.Headers["Authorization"]; auth != "" {
		return auth
	}

	// 小文字のケース
	if auth := request.Headers["authorization"]; auth != "" {
		return auth
	}

	// その他のケースも考慮
	for key, value := range request.Headers {
		if strings.ToLower(key) == "authorization" {
			return value
		}
	}

	return ""
}

// isDevelopmentEnvironment は開発環境かどうかを判定する
func (m *AuthMiddleware) isDevelopmentEnvironment() bool {
	return m.config.Environment == "development" || m.config.Environment == "dev"
}

// getDevTokenUserID は開発用固定ユーザーIDを返す
func (m *AuthMiddleware) getDevTokenUserID() (uuid.UUID, error) {
	// 開発用固定ユーザーID
	devUserID := "00000000-0000-0000-0000-000000000001"
	return uuid.Parse(devUserID)
}

// authenticateWithCognito はCognito JWTを使用して認証する
func (m *AuthMiddleware) authenticateWithCognito(authHeader string) (uuid.UUID, error) {
	m.logger.Infof("Cognito JWT認証開始: ヘッダー長=%d", len(authHeader))

	// Bearer プレフィックスの確認
	if !strings.HasPrefix(authHeader, "Bearer ") {
		m.logger.Error("無効なAuthorizationヘッダー形式")
		return uuid.Nil, NewAuthError("INVALID_FORMAT", "Authorizationヘッダーは 'Bearer <token>' 形式である必要があります", ErrInvalidFormat)
	}

	// JWTトークンの抽出
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		m.logger.Error("空のJWTトークン")
		return uuid.Nil, NewAuthError("EMPTY_TOKEN", "JWTトークンが空です", ErrInvalidToken)
	}

	// JWT検証
	claims, err := m.jwtValidator.ValidateJWT(token)
	if err != nil {
		m.logger.Errorf("JWT検証エラー: %v", err)
		return uuid.Nil, err
	}

	// ユーザーIDの取得
	userID, err := claims.GetUserID()
	if err != nil {
		m.logger.Errorf("ユーザーID取得エラー: %v", err)
		return uuid.Nil, NewClaimsError("ユーザーIDの取得に失敗しました")
	}

	m.logger.Infof("Cognito JWT認証成功: UserID=%s, TokenUse=%s", userID.String()[:8]+"...", claims.TokenUse)
	return userID, nil
}

// GetUserIDFromRequest はリクエストからユーザーIDを取得する（既存関数との互換性）
func (m *AuthMiddleware) GetUserIDFromRequest(request events.APIGatewayProxyRequest) (uuid.UUID, error) {
	return m.AuthenticateRequest(request)
}

// ValidateTokenAndGetClaims はトークンを検証してクレームを返す（詳細情報が必要な場合）
func (m *AuthMiddleware) ValidateTokenAndGetClaims(request events.APIGatewayProxyRequest) (*CognitoClaims, error) {
	authHeader := m.getAuthorizationHeader(request)
	if authHeader == "" {
		return nil, NewAuthError("MISSING_AUTH", "Authorizationヘッダーが必要です", ErrTokenNotFound)
	}

	// dev-tokenの場合はダミークレームを返す
	if m.isDevelopmentEnvironment() && authHeader == "Bearer dev-token" {
		devUserID, _ := m.getDevTokenUserID()
		return &CognitoClaims{
			Subject:  devUserID.String(),
			TokenUse: "access",
			Email:    "dev@example.com",
			Name:     "Development User",
			Issuer:   "dev-environment",
		}, nil
	}

	// Cognito JWT検証
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, NewAuthError("INVALID_FORMAT", "Authorizationヘッダーは 'Bearer <token>' 形式である必要があります", ErrInvalidFormat)
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	return m.jwtValidator.ValidateJWT(token)
}

// HealthCheck は認証ミドルウェアの接続確認を行う
func (m *AuthMiddleware) HealthCheck() error {
	// 開発環境の場合はスキップ
	if m.isDevelopmentEnvironment() {
		m.logger.Info("開発環境: Cognito HealthCheck スキップ")
		return nil
	}

	// Cognito JWKS接続確認
	ctx := context.Background()
	err := m.jwtValidator.HealthCheck(ctx)
	if err != nil {
		return fmt.Errorf("Cognito接続確認失敗: %w", err)
	}

	m.logger.Info("認証ミドルウェア HealthCheck 成功")
	return nil
}
