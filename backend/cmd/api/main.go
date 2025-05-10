package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/tsunakit99/selfpomodoro/internal/app"
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/database"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// ビルド情報（ビルド時に設定される）
var (
	Version     = "開発版"
	Commit      = "unknown"
	BuildTime   = "unknown"
	Environment = "development"
)

func main() {
	// 設定の読み込み
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("設定読み込みエラー: %v\n", err)
		os.Exit(1)
	}

	// 環境設定がない場合はビルド時の環境を使用
	if cfg.Environment == "" {
		cfg.Environment = Environment
	}

	// ロガーの初期化
	appLogger, err := logger.NewLogger(cfg.LogLevel, cfg.Environment)
	if err != nil {
		fmt.Printf("ロガー初期化エラー: %v\n", err)
		os.Exit(1)
	}

	// ビルド情報をログに記録
	appLogger.Infof("ポモドーロAPIサーバー 起動中... バージョン: %s, コミット: %s, ビルド時間: %s, 環境: %s",
		Version, Commit, BuildTime, cfg.Environment)

	// PostgreSQL接続
	postgresDB, err := database.NewPostgresDB(cfg, appLogger)
	if err != nil {
		appLogger.Fatalf("PostgreSQL初期化エラー: %v", err)
	}
	defer postgresDB.Close()

	// DynamoDB接続
	// dynamoDB, err := database.NewDynamoDB(cfg, appLogger)
	// if err != nil {
	// 	appLogger.Fatalf("DynamoDB初期化エラー: %v", err)
	// }

	// ここで各リポジトリ、ユースケース、ハンドラーを初期化
	// TODO

	// サーバー初期化
	server := app.NewServer(cfg, appLogger)

	// ルーティング設定
	// TODO

	// ヘルスチェックエンドポイント
	server.GetEcho().GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":  "OK",
			"version": Version,
		})
	})

	// サーバー起動
	if err := server.Start(); err != nil {
		appLogger.Fatalf("サーバー起動エラー: %v", err)
	}
}
