package main

import (
	"fmt"
	"os"

	"github.com/tsunakit99/selfpomodoro/internal/app"
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/database"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/repository"
	"github.com/tsunakit99/selfpomodoro/internal/interface/handler"
	"github.com/tsunakit99/selfpomodoro/internal/usecase"
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
	RepositoryFactory := repository.NewRepositoryFactory(postgresDB, nil, appLogger)
	useCases := usecase.NewUseCases(RepositoryFactory.Task, RepositoryFactory.Session, RepositoryFactory.Round, appLogger)
	handlers := handler.NewHandlers(useCases, appLogger)

	// サーバー初期化
	server := app.NewServer(cfg, appLogger)

	// ルーティング設定
	server.SetupRouter(handlers)

	// サーバー起動
	if err := server.Start(); err != nil {
		appLogger.Fatalf("サーバー起動エラー: %v", err)
	}
}
