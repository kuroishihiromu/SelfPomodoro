package auth

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// GetUserIDFromAPIGatewayContext はAPI Gateway Authorizerから認証済みユーザーIDを取得
func GetUserIDFromAPIGatewayContext(request events.APIGatewayProxyRequest, logger logger.Logger) (uuid.UUID, error) {
	// API Gateway Authorizerで認証済みのクレーム情報を取得
	claims, ok := request.RequestContext.Authorizer["claims"].(map[string]interface{})
	if !ok {
		logger.Error("認証コンテキストが見つかりません")
		return uuid.Nil, fmt.Errorf("認証エラー")
	}

	// ユーザーID（sub）を取得
	subClaim, ok := claims["sub"].(string)
	if !ok || subClaim == "" {
		logger.Error("subクレームが見つかりません")
		return uuid.Nil, fmt.Errorf("無効な認証情報")
	}

	// UUIDに変換
	userID, err := uuid.Parse(subClaim)
	if err != nil {
		logger.Errorf("無効なユーザーID: %v", err)
		return uuid.Nil, fmt.Errorf("無効なユーザーID")
	}

	return userID, nil
}