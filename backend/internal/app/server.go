package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// Server はEchoサーバーのラッパー
type Server struct {
	echo   *echo.Echo
	config *config.Config
	logger logger.Logger
}

// NewServer は新しいEchoサーバーを作成する
func NewServer(cfg *config.Config, logger logger.Logger) *Server {
	e := echo.New()

	// エラーハンドラーのカスタマイズ
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := "内部サーバーエラー"

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = fmt.Sprintf("%v", he.Message)
		}

		// クライアントにはエラーメッセージのみ返す
		if err := c.JSON(code, map[string]string{"error": message}); err != nil {
			logger.Errorf("JSONレスポンス生成エラー: %v", err)
		}

		// エラーをログに記録
		logger.Errorf("HTTPエラー: %d - %s", code, message)
	}

	// 基本的なミドルウェアを設定
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// リクエストロギングミドルウェア
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} | ${id} | ${method} ${uri} | ${status} | ${latency_human}\n",
	}))

	return &Server{
		echo:   e,
		config: cfg,
		logger: logger,
	}
}

// Start はサーバーを起動する
func (s *Server) Start() error {
	// サーバーの設定
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.ServerPort),
		ReadTimeout:  s.config.ServerReadTimeout,
		WriteTimeout: s.config.ServerWriteTimeout,
		IdleTimeout:  s.config.ServerIdleTimeout,
	}

	// シグナル処理のためのチャネル
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// 非同期でサーバーを起動
	go func() {
		s.logger.Infof("サーバーをポート %d で起動", s.config.ServerPort)
		if err := s.echo.StartServer(server); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("サーバー起動エラー: %v", err)
		}
	}()

	// シグナルを待機
	<-quit
	s.logger.Info("シャットダウンシグナルを受信")

	// シャットダウンのタイムアウト設定
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// グレースフルシャットダウン
	if err := s.echo.Shutdown(ctx); err != nil {
		s.logger.Errorf("サーバーシャットダウンエラー: %v", err)
		return err
	}

	s.logger.Info("サーバーは正常にシャットダウンしました")
	return nil
}

// GetEcho はEchoインスタンスを返す
func (s *Server) GetEcho() *echo.Echo {
	return s.echo
}
