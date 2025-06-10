package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/config"
	domainErrors "github.com/tsunakit99/selfpomodoro/internal/domain/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// AuthResult は認証結果を表す構造体
type AuthResult struct {
	UserID uuid.UUID      // 認証されたユーザーID
	Claims *CognitoClaims // JWTクレーム情報
}

// AuthMiddleware は認証処理のみを行うミドルウェア（Pure・UseCase依存なし）
type AuthMiddleware struct {
	jwtValidator *CognitoJWTValidator
	config       *config.Config
	logger       logger.Logger
}

// NewAuthMiddleware は新しい認証ミドルウェアを作成する（Pure版）
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

// Authenticate はリクエストを認証し、AuthResultを返す（Pure認証・User作成なし）
func (m *AuthMiddleware) Authenticate(request events.APIGatewayProxyRequest) (*AuthResult, error) {
	// Authorization ヘッダーの取得
	authHeader := m.getAuthorizationHeader(request)
	if authHeader == "" {
		m.logger.Warn("Authorizationヘッダーが見つかりません")
		return nil, domainErrors.NewTokenNotFoundError()
	}

	// 開発環境での後方互換性（dev-token）
	if m.isDevelopmentEnvironment() && authHeader == "Bearer dev-token" {
		m.logger.Info("開発環境: dev-token認証を使用")
		return m.getDevTokenAuthResult()
	}

	// Cognito JWT認証
	return m.authenticateWithCognito(authHeader)
}

// GetUserIDFromRequest はリクエストからユーザーIDを取得する（既存互換性）
func (m *AuthMiddleware) GetUserIDFromRequest(request events.APIGatewayProxyRequest) (uuid.UUID, error) {
	authResult, err := m.Authenticate(request)
	if err != nil {
		return uuid.Nil, err
	}
	return authResult.UserID, nil
}

// ValidateTokenAndGetClaims はトークンを検証してクレームを返す
func (m *AuthMiddleware) ValidateTokenAndGetClaims(request events.APIGatewayProxyRequest) (*CognitoClaims, error) {
	authResult, err := m.Authenticate(request)
	if err != nil {
		return nil, err
	}
	return authResult.Claims, nil
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

// getDevTokenAuthResult は開発用のAuthResultを返す
func (m *AuthMiddleware) getDevTokenAuthResult() (*AuthResult, error) {
	devUserID := "00000000-0000-0000-0000-000000000001"
	userID, err := uuid.Parse(devUserID)
	if err != nil {
		return nil, fmt.Errorf("dev-token ユーザーID解析エラー: %w", err)
	}

	// 開発用ダミークレーム
	claims := &CognitoClaims{
		Subject:         devUserID,
		TokenUse:        "access",
		Email:           "dev@example.com",
		Name:            "Development User",
		CognitoUsername: "dev-user",
		Issuer:          "dev-environment",
	}

	return &AuthResult{
		UserID: userID,
		Claims: claims,
	}, nil
}

// authenticateWithCognito はCognito JWTを使用して認証する
func (m *AuthMiddleware) authenticateWithCognito(authHeader string) (*AuthResult, error) {
	m.logger.Infof("Cognito JWT認証開始: ヘッダー長=%d", len(authHeader))

	// Bearer プレフィックスの確認
	if !strings.HasPrefix(authHeader, "Bearer ") {
		m.logger.Error("無効なAuthorizationヘッダー形式")
		return nil, domainErrors.NewInvalidFormatError()
	}

	// JWTトークンの抽出
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		m.logger.Error("空のJWTトークン")
		return nil, domainErrors.NewInvalidTokenError()
	}

	// JWT検証
	claims, err := m.jwtValidator.ValidateJWT(token)
	if err != nil {
		m.logger.Errorf("JWT検証エラー: %v", err)
		return nil, err
	}

	// ユーザーIDの取得
	userID, err := claims.GetUserID()
	if err != nil {
		m.logger.Errorf("ユーザーID取得エラー: %v", err)
		return nil, domainErrors.NewInvalidSubjectError()
	}

	m.logger.Infof("Cognito JWT認証成功: UserID=%s, TokenUse=%s", userID.String()[:8]+"...", claims.TokenUse)

	return &AuthResult{
		UserID: userID,
		Claims: claims,
	}, nil
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
