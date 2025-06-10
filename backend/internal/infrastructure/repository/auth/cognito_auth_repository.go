// internal/infrastructure/auth/cognito_auth_repository.go
package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/domain/errors"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/auth"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	infraError "github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository/auth/errors"
)

// CognitoAuthRepository はCognito認証を使用したAuthRepositoryの実装
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

// ValidateToken はトークンを検証してユーザーIDを返す
func (r *CognitoAuthRepository) ValidateToken(ctx context.Context, token string) (uuid.UUID, error) {
	// Bearer プレフィックスの処理
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		return uuid.Nil, errors.NewUnauthorizedError("空のトークンです")
	}

	// JWT検証
	claims, err := r.jwtValidator.ValidateJWT(token)
	if err != nil {
		r.logger.Errorf("JWT検証エラー: %v", err)
		// Infrastructureエラーを適切なDomainエラーに変換
		if infraError.IsTokenExpiredError(err) {
			return uuid.Nil, errors.NewUnauthorizedError("トークンの有効期限が切れています")
		}
		if infraError.IsInvalidTokenError(err) {
			return uuid.Nil, errors.NewUnauthorizedError("無効なトークンです")
		}
		return uuid.Nil, errors.NewUnauthorizedError("認証に失敗しました")
	}

	// ユーザーIDの取得
	userID, err := claims.GetUserID()
	if err != nil {
		r.logger.Errorf("ユーザーID取得エラー: %v", err)
		return uuid.Nil, errors.NewUnauthorizedError("ユーザーIDの取得に失敗しました")
	}

	return userID, nil
}

// ValidateTokenAndGetClaims はトークンを検証してクレーム情報を返す
func (r *CognitoAuthRepository) ValidateTokenAndGetClaims(ctx context.Context, token string) (*model.AuthClaims, error) {
	// Bearer プレフィックスの処理
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		return nil, errors.NewUnauthorizedError("空のトークンです")
	}

	// JWT検証
	cognitoClaims, err := r.jwtValidator.ValidateJWT(token)
	if err != nil {
		r.logger.Errorf("JWT検証エラー: %v", err)
		if infraError.IsTokenExpiredError(err) {
			return nil, errors.NewUnauthorizedError("トークンの有効期限が切れています")
		}
		if infraError.IsInvalidTokenError(err) {
			return nil, errors.NewUnauthorizedError("無効なトークンです")
		}
		return nil, errors.NewUnauthorizedError("認証に失敗しました")
	}

	// ユーザーIDの取得
	userID, err := cognitoClaims.GetUserID()
	if err != nil {
		return nil, errors.NewUnauthorizedError("ユーザーIDの取得に失敗しました")
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

	return authClaims, nil
}

// // RefreshToken はリフレッシュトークンを使って新しいアクセストークンを取得する
// func (r *CognitoAuthRepository) RefreshToken(ctx context.Context, refreshToken string) (*model.AuthTokenResponse, error) {
// 	// TODO: Cognito AdminInitiateAuth APIを使用してリフレッシュトークンを処理
// 	// 現在は未実装
// 	return nil, errors.NewNotImplementedError("リフレッシュトークン機能は未実装です")
// }

// // RevokeToken はトークンを無効化する
// func (r *CognitoAuthRepository) RevokeToken(ctx context.Context, token string) error {
// 	// TODO: Cognito GlobalSignOut APIを使用してトークンを無効化
// 	// 現在は未実装
// 	return errors.NewNotImplementedError("トークン無効化機能は未実装です")
// }

// HealthCheck は認証サービスの接続確認を行う
func (r *CognitoAuthRepository) HealthCheck(ctx context.Context) error {
	// JWT Validatorのヘルスチェック
	err := r.jwtValidator.HealthCheck(ctx)
	if err != nil {
		return fmt.Errorf("Cognito認証サービス接続確認失敗: %w", err)
	}

	r.logger.Info("Cognito認証サービス接続確認成功")
	return nil
}
