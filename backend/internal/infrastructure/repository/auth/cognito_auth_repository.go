// internal/infrastructure/repository/auth/cognito_auth_repository.go
package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/auth"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// CognitoAuthRepository はCognito認証を使用したAuthRepositoryの実装（新エラーハンドリング対応版）
type CognitoAuthRepository struct {
	jwtValidator *auth.CognitoJWTValidator
	config       *config.Config
	logger       logger.Logger
}

// NewCognitoAuthRepository は新しいCognitoAuthRepositoryを作成する
func NewCognitoAuthRepository(cfg *config.Config, logger logger.Logger) repository.AuthRepository {
	// JWT Validator初期化
	validatorConfig := &auth.CognitoJWTValidatorConfig{
		UserPoolID:   cfg.CognitoUserPoolID,
		ClientID:     cfg.CognitoClientID,
		Region:       cfg.AWSRegion,
		CacheTimeout: 30 * time.Minute,
		HTTPTimeout:  10 * time.Second,
	}

	jwtValidator := auth.NewCognitoJWTValidator(validatorConfig, logger)

	return &CognitoAuthRepository{
		jwtValidator: jwtValidator,
		config:       cfg,
		logger:       logger,
	}
}

// ValidateToken はトークンを検証してユーザーIDを返す（新エラーハンドリング対応版）
func (r *CognitoAuthRepository) ValidateToken(ctx context.Context, token string) (uuid.UUID, error) {
	// Bearer プレフィックスの処理
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		return uuid.Nil, appErrors.NewTokenNotFoundError()
	}

	// JWT検証
	claims, err := r.jwtValidator.ValidateJWT(token)
	if err != nil {
		r.logger.Errorf("JWT検証エラー: %v", err)

		// Infrastructure Error → Domain Error 変換（新エラーハンドリング）
		if appErrors.IsTokenExpiredError(err) {
			return uuid.Nil, appErrors.NewTokenExpiredError()
		}
		if appErrors.IsInvalidTokenError(err) {
			return uuid.Nil, appErrors.NewInvalidTokenError()
		}
		if appErrors.IsJWTError(err) {
			return uuid.Nil, appErrors.NewInvalidTokenError()
		}
		if appErrors.IsJWKSError(err) {
			return uuid.Nil, appErrors.NewInternalError(err)
		}
		if appErrors.IsHTTPError(err) {
			return uuid.Nil, appErrors.NewInternalError(err)
		}

		// その他のInfrastructure Error
		if appErrors.IsInfrastructureError(err) {
			return uuid.Nil, appErrors.NewInternalError(err)
		}

		// 不明なエラー
		return uuid.Nil, appErrors.NewUnauthorizedError("認証に失敗しました")
	}

	// ユーザーIDの取得
	userID, err := claims.GetUserID()
	if err != nil {
		r.logger.Errorf("ユーザーID取得エラー: %v", err)
		return uuid.Nil, appErrors.NewInvalidSubjectError()
	}

	r.logger.Debugf("JWT認証成功: UserID=%s, TokenUse=%s", userID.String()[:8]+"...", claims.TokenUse)
	return userID, nil
}

// ValidateTokenAndGetClaims はトークンを検証してクレーム情報を返す（新エラーハンドリング対応版）
func (r *CognitoAuthRepository) ValidateTokenAndGetClaims(ctx context.Context, token string) (*model.AuthClaims, error) {
	// Bearer プレフィックスの処理
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		return nil, appErrors.NewTokenNotFoundError()
	}

	// JWT検証
	cognitoClaims, err := r.jwtValidator.ValidateJWT(token)
	if err != nil {
		r.logger.Errorf("JWT検証エラー: %v", err)

		// Infrastructure Error → Domain Error 変換（新エラーハンドリング）
		if appErrors.IsTokenExpiredError(err) {
			return nil, appErrors.NewTokenExpiredError()
		}
		if appErrors.IsInvalidTokenError(err) {
			return nil, appErrors.NewInvalidTokenError()
		}
		if appErrors.IsJWTError(err) {
			return nil, appErrors.NewInvalidTokenError()
		}
		if appErrors.IsJWKSError(err) {
			return nil, appErrors.NewInternalError(err)
		}
		if appErrors.IsHTTPError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		// その他のInfrastructure Error
		if appErrors.IsInfrastructureError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		// 不明なエラー
		return nil, appErrors.NewUnauthorizedError("認証に失敗しました")
	}

	// ユーザーIDの取得
	userID, err := cognitoClaims.GetUserID()
	if err != nil {
		r.logger.Errorf("ユーザーID取得エラー: %v", err)
		return nil, appErrors.NewInvalidSubjectError()
	}

	// Domain ModelのAuthClaimsに変換
	authClaims := &model.AuthClaims{
		UserID:          userID,
		Email:           cognitoClaims.Email,
		Name:            cognitoClaims.GetDisplayName(),
		Provider:        cognitoClaims.GetProviderName(),
		CognitoUsername: cognitoClaims.CognitoUsername,
		EmailVerified:   cognitoClaims.EmailVerified,
		TokenUse:        cognitoClaims.TokenUse,
		ExpiresAt:       cognitoClaims.ExpiresAt,
		IssuedAt:        cognitoClaims.IssuedAt,
		Subject:         cognitoClaims.Subject,
		Issuer:          cognitoClaims.Issuer,
		Audience:        cognitoClaims.Audience,
	}

	r.logger.Debugf("JWT認証・クレーム取得成功: UserID=%s, Provider=%s",
		userID.String()[:8]+"...", authClaims.Provider)

	return authClaims, nil
}

// HealthCheck は認証サービスの接続確認を行う（新エラーハンドリング対応版）
func (r *CognitoAuthRepository) HealthCheck(ctx context.Context) error {
	// JWT Validatorのヘルスチェック
	err := r.jwtValidator.HealthCheck(ctx)
	if err != nil {
		r.logger.Errorf("Cognito認証サービス接続確認エラー: %v", err)

		// Infrastructure Error → Domain Error 変換（新エラーハンドリング）
		if appErrors.IsHTTPError(err) {
			return appErrors.NewExternalServiceUnavailableError("Cognito")
		}
		if appErrors.IsJWKSError(err) {
			return appErrors.NewExternalServiceError("Cognito", "jwks_check", err)
		}
		if appErrors.IsNetworkError(err) {
			return appErrors.NewExternalServiceTimeoutError("Cognito")
		}

		// その他のInfrastructure Error
		if appErrors.IsInfrastructureError(err) {
			return appErrors.NewInternalError(err)
		}

		return appErrors.NewInternalError(fmt.Errorf("Cognito認証サービス接続確認失敗: %w", err))
	}

	r.logger.Info("Cognito認証サービス接続確認成功")
	return nil
}
