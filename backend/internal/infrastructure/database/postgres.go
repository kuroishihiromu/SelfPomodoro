package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQLドライバー
	"github.com/tsunakit99/selfpomodoro/internal/config"
	"github.com/tsunakit99/selfpomodoro/internal/infrastructure/logger"
)

// PostgresDB はPostgreSQLデータベース接続を管理する
type PostgresDB struct {
	DB     *sqlx.DB
	logger logger.Logger
}

// NewPostgresDB は新しいPostgreSQLデータベース接続を作成する（Lambda最適化版）
func NewPostgresDB(cfg *config.Config, logger logger.Logger) (*PostgresDB, error) {
	// 接続文字列の作成
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	// DBへの接続
	logger.Infof("PostgreSQLに接続: %s:%d/%s", cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("PostgreSQL接続エラー: %w", err)
	}

	// Lambda用接続設定（軽量化）
	db.SetMaxOpenConns(5)                  // Lambda用に少なめに設定
	db.SetMaxIdleConns(2)                  // アイドル接続も少なめ
	db.SetConnMaxLifetime(1 * time.Minute) // 短いライフタイム

	// 接続テスト
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("PostgreSQL Ping失敗: %w", err)
	}

	logger.Info("PostgreSQL接続成功")

	return &PostgresDB{
		DB:     db,
		logger: logger,
	}, nil
}

// Close はデータベース接続を閉じる
func (p *PostgresDB) Close() error {
	if p.DB != nil {
		p.logger.Info("PostgreSQL接続を閉じる")
		return p.DB.Close()
	}
	return nil
}

// ExecTx はトランザクション内で関数を実行する（既存機能維持）
func (p *PostgresDB) ExecTx(fn func(*sqlx.Tx) error) error {
	tx, err := p.DB.Beginx()
	if err != nil {
		return fmt.Errorf("トランザクション開始エラー: %w", err)
	}

	// パニックをキャッチしてロールバック
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			p.logger.Errorf("トランザクションでパニック発生: %v", r)
			panic(r) // パニックを再スロー
		}
	}()

	// 関数を実行
	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			p.logger.Errorf("トランザクションロールバックエラー: %v", rbErr)
		}
		return err
	}

	// コミット
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("トランザクションコミットエラー: %w", err)
	}

	return nil
}
