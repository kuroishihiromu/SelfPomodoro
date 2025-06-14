package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/config"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// AuthResult は認証結果を表す構造体
type AuthResult struct {
	UserID uuid.UUID      // 認証されたユーザーID
	Claims *CognitoClaims // JWTクレーム情報
}

// AuthMiddleware は認証処理のみを行うミドルウェア（Pure・UseCase依存なし・新エラーハンドリング対応版）
type AuthMiddleware struct {
	jwtValidator *CognitoJWTValidator
	config       *config.Config
	logger       logger.Logger
}

// NewAuthMiddleware は新しい認証ミドルウェアを作成する（Pure版・新エラーハンドリング対応）
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

// Authenticate はリクエストを認証し、AuthResultを返す（Pure認証・User作成なし・新エラーハンドリング対応版）
func (m *AuthMiddleware) Authenticate(request events.APIGatewayProxyRequest) (*AuthResult, error) {
	// Authorization ヘッダーの取得
	authHeader := m.getAuthorizationHeader(request)
	if authHeader == "" {
		m.logger.Warn("Authorizationヘッダーが見つかりません")
		return nil, appErrors.NewTokenNotFoundError() // 新エラー構造使用
	}

	// 開発環境での後方互換性（dev-token）
	if m.isDevelopmentEnvironment() && authHeader == "Bearer dev-token" {
		m.logger.Info("開発環境: dev-token認証を使用")
		return m.getDevTokenAuthResult()
	}

	// Cognito JWT認証
	return m.authenticateWithCognito(authHeader)
}

// GetUserIDFromRequest はリクエストからユーザーIDを取得する（既存互換性・新エラーハンドリング対応版）
func (m *AuthMiddleware) GetUserIDFromRequest(request events.APIGatewayProxyRequest) (uuid.UUID, error) {
	authResult, err := m.Authenticate(request)
	if err != nil {
		return uuid.Nil, err // 新エラー構造がそのまま返される
	}
	return authResult.UserID, nil
}

// ValidateTokenAndGetClaims はトークンを検証してクレームを返す（新エラーハンドリング対応版）
func (m *AuthMiddleware) ValidateTokenAndGetClaims(request events.APIGatewayProxyRequest) (*CognitoClaims, error) {
	authResult, err := m.Authenticate(request)
	if err != nil {
		return nil, err // 新エラー構造がそのまま返される
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

// getDevTokenAuthResult は開発用のAuthResultを返す（新エラーハンドリング対応版）
func (m *AuthMiddleware) getDevTokenAuthResult() (*AuthResult, error) {
	devUserID := "00000000-0000-0000-0000-000000000001"
	userID, err := uuid.Parse(devUserID)
	if err != nil {
		return nil, appErrors.NewInternalError(fmt.Errorf("dev-token ユーザーID解析エラー: %w", err))
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

// authenticateWithCognito はCognito JWTを使用して認証する（新エラーハンドリング対応版）
func (m *AuthMiddleware) authenticateWithCognito(authHeader string) (*AuthResult, error) {
	m.logger.Infof("Cognito JWT認証開始: ヘッダー長=%d", len(authHeader))

	// Bearer プレフィックスの確認
	if !strings.HasPrefix(authHeader, "Bearer ") {
		m.logger.Error("無効なAuthorizationヘッダー形式")
		return nil, appErrors.NewInvalidFormatError() // 新エラー構造使用
	}

	// JWTトークンの抽出
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		m.logger.Error("空のJWTトークン")
		return nil, appErrors.NewInvalidTokenError() // 新エラー構造使用
	}

	// JWT検証（Infrastructure Error → Domain Error 変換）
	claims, err := m.jwtValidator.ValidateJWT(token)
	if err != nil {
		m.logger.Errorf("JWT検証エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if appErrors.IsJWTError(err) {
			// JWT関連のInfrastructure Errorを適切なDomain Errorに変換
			if errors.Is(err, appErrors.ErrJWTTokenExpired) {
				return nil, appErrors.NewTokenExpiredError()
			}
			if errors.Is(err, appErrors.ErrJWTTokenInvalid) ||
				errors.Is(err, appErrors.ErrJWTParsingFailed) ||
				errors.Is(err, appErrors.ErrJWTSignatureInvalid) {
				return nil, appErrors.NewInvalidTokenError()
			}
		}

		if appErrors.IsJWKSError(err) {
			// JWKS関連エラーは一般的に認証エラーとして扱う
			return nil, appErrors.NewUnauthorizedError("認証サービスエラー")
		}

		if appErrors.IsHTTPError(err) {
			// HTTP関連エラーは一時的な認証サービス障害として扱う
			return nil, appErrors.NewUnauthorizedError("認証サービス接続エラー")
		}

		// その他のエラーは一般的な認証失敗として扱う
		return nil, appErrors.NewUnauthorizedError("認証に失敗しました")
	}

	// ユーザーIDの取得
	userID, err := claims.GetUserID()
	if err != nil {
		m.logger.Errorf("ユーザーID取得エラー: %v", err)

		// ClaimsのエラーをDomain Errorに変換
		if errors.Is(err, appErrors.ErrInvalidSubject) {
			return nil, appErrors.NewInvalidSubjectError()
		}
		if errors.Is(err, appErrors.ErrMissingSubject) {
			return nil, appErrors.NewMissingSubjectError()
		}

		return nil, appErrors.NewInvalidTokenError()
	}

	m.logger.Infof("Cognito JWT認証成功: UserID=%s, TokenUse=%s", userID.String()[:8]+"...", claims.TokenUse)

	return &AuthResult{
		UserID: userID,
		Claims: claims,
	}, nil
}

// HealthCheck は認証ミドルウェアの接続確認を行う（新エラーハンドリング対応版）
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
		// Infrastructure Error → Domain Error 変換
		if appErrors.IsHTTPError(err) {
			return appErrors.NewInternalError(fmt.Errorf("Cognito接続確認失敗: 通信エラー"))
		}
		if appErrors.IsJWKSError(err) {
			return appErrors.NewInternalError(fmt.Errorf("Cognito接続確認失敗: JWKS取得エラー"))
		}

		return appErrors.NewInternalError(fmt.Errorf("Cognito接続確認失敗: %w", err))
	}

	m.logger.Info("認証ミドルウェア HealthCheck 成功")
	return nil
}
