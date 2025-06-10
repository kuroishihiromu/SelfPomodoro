package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/container"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/usecase"
)

// Global container (Lambda再利用最適化)
var globalContainer container.Container

// init はLambda init phaseで実行
func init() {
	globalContainer = container.NewLambdaContainer()
}

// PostConfirmationHandler はDI Container使用版のPostConfirmationハンドラー
type PostConfirmationHandler struct {
	useCases *usecase.UseCases
	logger   logger.Logger
}

// handler はCognito PostConfirmation Triggerのエントリーポイント（DI Container版）
func handler(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	// 1. Container初期化
	if err := globalContainer.Initialize(ctx); err != nil {
		// PostConfirmationでは詳細エラーを返す（CognitoがLogで確認）
		return event, fmt.Errorf("サービス初期化エラー: %w", err)
	}

	// 2. Dependencies取得（Infrastructure依存なし！）
	useCases := globalContainer.GetUseCases()
	logger := globalContainer.GetLogger()

	// 3. Handler初期化
	postConfirmationHandler := &PostConfirmationHandler{
		useCases: useCases,
		logger:   logger,
	}

	logger.Infof("PostConfirmation Trigger開始: UserPoolID=%s, Username=%s",
		event.UserPoolID, event.UserName)

	// 4. PostConfirmation処理実行
	err := postConfirmationHandler.processPostConfirmation(ctx, event)
	if err != nil {
		logger.Errorf("PostConfirmation処理エラー: %v", err)
		return event, err
	}

	logger.Infof("PostConfirmation Trigger完了: Username=%s", event.UserName)
	return event, nil
}

// processPostConfirmation はPostConfirmation後の初期化処理を実行する
func (h *PostConfirmationHandler) processPostConfirmation(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) error {
	// Cognito属性から情報取得
	userAttributes := event.Request.UserAttributes

	// ユーザーIDの取得（Cognito sub）
	userIDStr, exists := userAttributes["sub"]
	if !exists || userIDStr == "" {
		return fmt.Errorf("subクレームが見つかりません")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fmt.Errorf("無効なユーザーID: %w", err)
	}

	// PostConfirmationパラメータの作成
	params := usecase.PostConfirmationParams{
		UserID:     userID,
		Email:      userAttributes["email"],
		Name:       userAttributes["name"],
		GivenName:  userAttributes["given_name"],
		FamilyName: userAttributes["family_name"],
		Provider:   "Cognito_UserPool",
	}

	h.logger.Infof("ユーザーオンボーディング実行: UserID=%s, Email=%s, Name=%s",
		userID.String()[:8]+"...", params.Email, params.Name)

	// OnboardingUseCaseに完全委譲（DI Container経由）
	err = h.useCases.Onboarding.CompletePostConfirmationSetup(ctx, params)
	if err != nil {
		return fmt.Errorf("オンボーディング処理失敗: %w", err)
	}

	h.logger.Infof("ユーザーオンボーディング完了: UserID=%s", userID.String()[:8]+"...")
	return nil
}

func main() {
	lambda.Start(handler)
}
