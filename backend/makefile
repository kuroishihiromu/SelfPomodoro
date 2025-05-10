# PgAdminコンテナのログを表示
.PHONY: pgadmin-logs
pgadmin-logs:
	@echo "PgAdminコンテナのログを表示しています..."
	@docker logs -f $(PG_ADMIN_CONTAINER)# 変数定義
APP_NAME=api
DOCKER_COMPOSE=docker-compose
GO_BUILD=go build
GO_RUN=go run
GO_TEST=go test
GO_MIGRATE=migrate
PG_CONTAINER=pomodoro-postgres
PG_ADMIN_CONTAINER=pomodoro-pgadmin

# デフォルトのターゲット
.PHONY: help
help:
	@echo "使用可能なコマンド:"
	@echo "  make setup      - 開発環境を初期化 (Docker起動、マイグレーション実行)"
	@echo "  make dev        - 開発サーバーを起動 (Docker + Go)"
	@echo "  make run        - Go APIサーバーのみ起動"
	@echo "  make docker-up  - Dockerコンテナを起動"
	@echo "  make docker-build - Dockerイメージをビルド"
	@echo "  make docker-down - Dockerコンテナを停止"
	@echo "  make docker-restart - Dockerコンテナを再起動"
	@echo "  make build      - APIサーバーをビルド"
	@echo "  make test       - テストを実行"
	@echo "  make migrate-up - マイグレーションを適用"
	@echo "  make migrate-down - マイグレーションを1つロールバック"
	@echo "  make clean      - ビルド成果物を削除"
	@echo "  make logs       - Dockerコンテナのログを表示"
	@echo "  make pg-console - PostgreSQLコンソールに接続"

# 開発環境のセットアップ
.PHONY: setup
setup: docker-up migrate-up
	@echo "開発環境のセットアップが完了しました"

# 開発サーバーを起動 (バックグラウンドでDockerを起動 + Goサーバー)
.PHONY: dev
dev: docker-up
	@echo "Goサーバーを起動しています..."
	@$(GO_RUN) cmd/api/main.go

# Dockerコンテナを起動
.PHONY: docker-up
docker-up:
	@echo "Dockerコンテナを起動しています..."
	@$(DOCKER_COMPOSE) up -d --build
	@echo "Dockerコンテナが起動しました"
	@echo "PostgreSQL: localhost:5432"
	@echo "PgAdmin: http://localhost:8081"

# Dockerコンテナを停止
.PHONY: docker-down
docker-down:
	@echo "Dockerコンテナを停止しています..."
	@$(DOCKER_COMPOSE) down
	@echo "Dockerコンテナが停止しました"

# Goサーバーのみ起動
.PHONY: run
run:
	@echo "Goサーバーを起動しています..."
	@$(GO_RUN) cmd/api/main.go

# アプリケーションをビルド
.PHONY: build
build:
	@echo "アプリケーションをビルドしています..."
	@mkdir -p bin
	@$(GO_BUILD) -o bin/$(APP_NAME) cmd/api/main.go
	@echo "ビルド完了: bin/$(APP_NAME)"

# マイグレーションを適用
.PHONY: migrate-up
migrate-up:
	@echo "マイグレーションを適用しています..."
	@./scripts/migrate.sh up
	@echo "マイグレーション完了"

# マイグレーションをロールバック
.PHONY: migrate-down
migrate-down:
	@echo "マイグレーションをロールバックしています..."
	@./scripts/migrate.sh down
	@echo "ロールバック完了"

# 新しいマイグレーションファイルを作成
.PHONY: migrate-create
migrate-create:
	@read -p "マイグレーション名を入力してください: " name; \
	./scripts/migrate.sh create $$name

# テストを実行
.PHONY: test
test:
	@echo "テストを実行しています..."
	@$(GO_TEST) ./...

# ビルド成果物を削除
.PHONY: clean
clean:
	@echo "ビルド成果物を削除しています..."
	@rm -rf bin
	@echo "クリーンアップ完了"

# Dockerイメージをビルド
.PHONY: docker-build
docker-build:
	@echo "Dockerイメージをビルドしています..."
	@$(DOCKER_COMPOSE) build
	@echo "Dockerイメージのビルドが完了しました"

# Dockerコンテナを再起動
.PHONY: docker-restart
docker-restart: docker-down docker-up
	@echo "Dockerコンテナを再起動しました"

# PostgreSQLコンソールに接続
.PHONY: pg-console
pg-console:
	@echo "PostgreSQLコンソールに接続しています..."
	@docker exec -it $(PG_CONTAINER) psql -U postgres -d pomodoro

# Dockerコンテナのログを表示
.PHONY: logs
logs:
	@echo "Dockerコンテナのログを表示しています..."
	@$(DOCKER_COMPOSE) logs -f