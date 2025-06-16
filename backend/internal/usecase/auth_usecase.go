package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/domain/model"
	"github.com/tsunakit99/selfpomodoro/internal/domain/repository"
	appErrors "github.com/tsunakit99/selfpomodoro/internal/errors"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// AuthUseCase は認証に関するユースケースを定義するインターフェース
type AuthUseCase interface {
	// AuthenticateRequest はAPIリクエストを認証してユーザーIDを返す
	AuthenticateRequest(ctx context.Context, request events.APIGatewayProxyRequest) (uuid.UUID, error)

	// AuthenticateAndValidateUser は認証とUser存在確認を統合して行う（統一パターン）
	AuthenticateAndValidateUser(ctx context.Context, request events.APIGatewayProxyRequest) (uuid.UUID, error)

	// AuthenticateAndGetUser はAPIリクエストを認証してユーザー情報を返す
	AuthenticateAndGetUser(ctx context.Context, request events.APIGatewayProxyRequest) (*model.AuthUser, error)

	// ValidateToken はトークンの有効性を検証する
	ValidateToken(ctx context.Context, token string) (*model.AuthClaims, error)

	// CheckAuthHealth は認証サービスの接続確認を行う
	CheckAuthHealth(ctx context.Context) error
}

// authUseCase はAuthUseCaseの実装（新エラーハンドリング完全対応版）
type authUseCase struct {
	authRepo repository.AuthRepository
	userRepo repository.UserRepository // User存在確認用
	config   *config.Config
	logger   logger.Logger
}

// NewAuthUseCase は新しいAuthUseCaseを作成する
func NewAuthUseCase(authRepo repository.AuthRepository, userRepo repository.UserRepository, config *config.Config, logger logger.Logger) AuthUseCase {
	return &authUseCase{
		authRepo: authRepo,
		userRepo: userRepo,
		config:   config,
		logger:   logger,
	}
}

// AuthenticateRequest はAPIリクエストを認証してユーザーIDを返す（新エラーハンドリング対応版）
func (uc *authUseCase) AuthenticateRequest(ctx context.Context, request events.APIGatewayProxyRequest) (uuid.UUID, error) {
	// Authorization ヘッダー取得
	authHeader := uc.getAuthorizationHeader(request)
	if authHeader == "" {
		return uuid.Nil, appErrors.NewTokenNotFoundError()
	}

	// 開発環境での後方互換性（dev-token）
	if uc.isDevelopmentEnvironment() && authHeader == "Bearer dev-token" {
		devUserID, err := uuid.Parse("00000000-0000-0000-0000-000000000001")
		if err != nil {
			return uuid.Nil, appErrors.NewInternalError(err)
		}
		uc.logger.Info("開発環境: dev-token認証を使用")
		return devUserID, nil
	}

	// Repository経由で認証
	userID, err := uc.authRepo.ValidateToken(ctx, authHeader)
	if err != nil {
		uc.logger.Errorf("認証失敗: %v", err)

		// Infrastructure Error → Domain Error 変換
		if appErrors.IsJWTError(err) {
			if appErrors.IsTokenExpiredError(err) {
				return uuid.Nil, appErrors.NewTokenExpiredError()
			}
			if appErrors.IsInvalidTokenError(err) {
				return uuid.Nil, appErrors.NewInvalidTokenError()
			}
			return uuid.Nil, appErrors.NewUnauthorizedError("JWT認証エラー")
		}
		if appErrors.IsJWKSError(err) {
			return uuid.Nil, appErrors.NewUnauthorizedError("公開キー取得エラー")
		}
		if appErrors.IsHTTPError(err) {
			return uuid.Nil, appErrors.NewInternalError(err)
		}

		return uuid.Nil, appErrors.NewUnauthorizedError("認証に失敗しました")
	}

	uc.logger.Infof("認証成功: UserID=%s", userID.String()[:8]+"...")
	return userID, nil
}

// AuthenticateAndValidateUser は認証とUser存在確認を統合して行う（新エラーハンドリング対応版）
func (uc *authUseCase) AuthenticateAndValidateUser(ctx context.Context, request events.APIGatewayProxyRequest) (uuid.UUID, error) {
	// 1. 認証処理
	userID, err := uc.AuthenticateRequest(ctx, request)
	if err != nil {
		return uuid.Nil, err // Domain Errorがそのまま返される
	}

	// 2. User存在確認（PostConfirmation前提）
	exists, err := uc.userRepo.ExistsByID(ctx, userID)
	if err != nil {
		uc.logger.Errorf("User存在確認エラー: %v", err)

		// Infrastructure Error → Domain Error 変換
		if appErrors.IsDatabaseError(err) {
			return uuid.Nil, appErrors.NewInternalError(err)
		}

		return uuid.Nil, appErrors.NewInternalError(err)
	}
	if !exists {
		uc.logger.Warnf("存在しないユーザーのアクセス: %s", userID.String())
		return uuid.Nil, appErrors.NewUserNotFoundError()
	}

	uc.logger.Infof("認証・User存在確認成功: UserID=%s", userID.String()[:8]+"...")
	return userID, nil
}

// AuthenticateAndGetUser はAPIリクエストを認証してユーザー情報を返す（新エラーハンドリング対応版）
func (uc *authUseCase) AuthenticateAndGetUser(ctx context.Context, request events.APIGatewayProxyRequest) (*model.AuthUser, error) {
	authHeader := uc.getAuthorizationHeader(request)
	if authHeader == "" {
		return nil, appErrors.NewTokenNotFoundError()
	}

	// 開発環境での後方互換性
	if uc.isDevelopmentEnvironment() && authHeader == "Bearer dev-token" {
		devUserID, _ := uuid.Parse("00000000-0000-0000-0000-000000000001")
		return &model.AuthUser{
			UserID:   devUserID,
			Email:    "dev@example.com",
			Name:     "Development User",
			Provider: "dev-environment",
		}, nil
	}

	// Repository経由でクレーム取得
	claims, err := uc.authRepo.ValidateTokenAndGetClaims(ctx, authHeader)
	if err != nil {
		uc.logger.Errorf("クレーム取得エラー: %v", err)

		// Infrastructure Error → Domain Error 変換（AuthenticateRequestと同様）
		if appErrors.IsJWTError(err) {
			if appErrors.IsTokenExpiredError(err) {
				return nil, appErrors.NewTokenExpiredError()
			}
			if appErrors.IsInvalidTokenError(err) {
				return nil, appErrors.NewInvalidTokenError()
			}
			return nil, appErrors.NewUnauthorizedError("JWT認証エラー")
		}
		if appErrors.IsJWKSError(err) {
			return nil, appErrors.NewUnauthorizedError("公開キー取得エラー")
		}
		if appErrors.IsHTTPError(err) {
			return nil, appErrors.NewInternalError(err)
		}

		return nil, appErrors.NewUnauthorizedError("認証に失敗しました")
	}

	authUser := &model.AuthUser{
		UserID:   claims.UserID,
		Email:    claims.Email,
		Name:     claims.Name,
		Provider: claims.Provider,
		Claims:   claims,
	}

	return authUser, nil
}

// ValidateToken はトークンの有効性を検証する（新エラーハンドリング対応版）
func (uc *authUseCase) ValidateToken(ctx context.Context, token string) (*model.AuthClaims, error) {
	claims, err := uc.authRepo.ValidateTokenAndGetClaims(ctx, token)
	if err != nil {
		uc.logger.Errorf("トークン検証失敗: %v", err)

		// Infrastructure Error → Domain Error 変換
		if appErrors.IsJWTError(err) {
			if appErrors.IsTokenExpiredError(err) {
				return nil, appErrors.NewTokenExpiredError()
			}
			if appErrors.IsInvalidTokenError(err) {
				return nil, appErrors.NewInvalidTokenError()
			}
			return nil, appErrors.NewUnauthorizedError("JWT検証エラー")
		}
		if appErrors.IsJWKSError(err) {
			return nil, appErrors.NewUnauthorizedError("公開キー取得エラー")
		}

		return nil, appErrors.NewInternalError(err)
	}
	return claims, nil
}

// CheckAuthHealth は認証サービスの接続確認を行う（新エラーハンドリング対応版）
func (uc *authUseCase) CheckAuthHealth(ctx context.Context) error {
	if uc.authRepo == nil {
		return appErrors.NewInternalError(fmt.Errorf("AuthRepositoryが初期化されていません"))
	}

	// 開発環境の場合はスキップ
	if uc.isDevelopmentEnvironment() {
		uc.logger.Info("開発環境: 認証ヘルスチェックをスキップ")
		return nil
	}

	// Repository経由でヘルスチェック
	err := uc.authRepo.HealthCheck(ctx)
	if err != nil {
		uc.logger.Errorf("認証サービス接続確認失敗: %v", err)

		// Infrastructure Error → Domain Error 変換
		if appErrors.IsHTTPError(err) || appErrors.IsJWKSError(err) {
			return appErrors.NewInternalError(fmt.Errorf("認証サービス接続確認失敗: %w", err))
		}

		return appErrors.NewInternalError(fmt.Errorf("認証サービス接続確認失敗: %w", err))
	}

	uc.logger.Info("認証サービス接続確認成功")
	return nil
}

// getAuthorizationHeader はAuthorizationヘッダーを取得する
func (uc *authUseCase) getAuthorizationHeader(request events.APIGatewayProxyRequest) string {
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
func (uc *authUseCase) isDevelopmentEnvironment() bool {
	return uc.config.Environment == "development" || uc.config.Environment == "dev"
}
