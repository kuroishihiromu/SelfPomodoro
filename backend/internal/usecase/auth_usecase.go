package usecase

import "github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"

// AuthUseCase は認証に関するユースケースを定義するインターフェース
type AuthUseCase interface {
	// TODO: あとで拡張
}

// authUseCase は認証に関するユースケースの実装
type authUseCase struct {
	logger logger.Logger
}

func NewAuthUseCase(logger logger.Logger) AuthUseCase {
	return &authUseCase{
		logger: logger,
	}
}
