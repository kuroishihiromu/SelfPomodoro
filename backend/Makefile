# Selfpomodoro Serverless Makefile

# 変数定義
GOOS := linux
GOARCH := amd64
CGO_ENABLED := 0

# Lambda関数一覧
FUNCTIONS := api-tasks api-sessions api-rounds api-statistics

# デフォルトのターゲット
.PHONY: help
help:
	@echo "使用可能なコマンド:"
	@echo "  make build-all      - 全Lambda関数をビルド"
	@echo "  make deploy-all     - 全Lambda関数をデプロイ"
	@echo "  make test-all       - 全APIのテストを実行"
	@echo "  make clean          - ビルド成果物を削除"
	@echo ""
	@echo "個別関数操作:"
	@echo "  make build-<function>   - 指定関数をビルド"
	@echo "  make deploy-<function>  - 指定関数をデプロイ"
	@echo ""
	@echo "その他:"
	@echo "  make migrate-up     - マイグレーションを適用"
	@echo "  make migrate-down   - マイグレーションをロールバック"
	@echo "  make logs-<function> - 指定関数のログを表示"

# 全Lambda関数ビルド
.PHONY: build-all
build-all: $(foreach func,$(FUNCTIONS),build-$(func))
	@echo "✅ 全Lambda関数のビルドが完了しました"

# 個別Lambda関数ビルド
.PHONY: build-api-tasks
build-api-tasks:
	@echo "🔨 api-tasks をビルド中..."
	@cd functions/api-tasks && \
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) \
	go build -ldflags="-s -w" -o bootstrap main.go && \
	zip bootstrap.zip bootstrap && \
	rm bootstrap
	@echo "✅ api-tasks ビルド完了"

.PHONY: build-api-sessions
build-api-sessions:
	@echo "🔨 api-sessions をビルド中..."
	@cd functions/api-sessions && \
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) \
	go build -ldflags="-s -w" -o bootstrap main.go && \
	zip bootstrap.zip bootstrap && \
	rm bootstrap
	@echo "✅ api-sessions ビルド完了"

.PHONY: build-api-rounds
build-api-rounds:
	@echo "🔨 api-rounds をビルド中..."
	@cd functions/api-rounds && \
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) \
	go build -ldflags="-s -w" -o bootstrap main.go && \
	zip bootstrap.zip bootstrap && \
	rm bootstrap
	@echo "✅ api-rounds ビルド完了"

.PHONY: build-api-statistics
build-api-statistics:
	@echo "🔨 api-statistics をビルド中..."
	@cd functions/api-statistics && \
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) \
	go build -ldflags="-s -w" -o bootstrap main.go && \
	zip bootstrap.zip bootstrap && \
	rm bootstrap
	@echo "✅ api-statistics ビルド完了"

# デプロイコマンド
.PHONY: deploy-all
deploy-all: build-all
	@echo "🚀 全Lambda関数をデプロイ中..."
	serverless deploy --aws-profile selfpomodoro-dev
	@echo "✅ 全Lambda関数のデプロイが完了しました"

.PHONY: deploy-function
deploy-function:
	@if [ -z "$(FUNC)" ]; then \
		echo "❌ エラー: FUNC変数を指定してください (例: make deploy-function FUNC=api-tasks)"; \
		exit 1; \
	fi
	@echo "🚀 $(FUNC) をデプロイ中..."
	serverless deploy function --function $(FUNC) --aws-profile selfpomodoro-dev
	@echo "✅ $(FUNC) のデプロイが完了しました"

# 個別関数デプロイ
.PHONY: deploy-api-tasks
deploy-api-tasks: build-api-tasks
	@echo "🚀 api-tasks をデプロイ中..."
	serverless deploy function --function api-tasks --aws-profile selfpomodoro-dev
	@echo "✅ api-tasks デプロイ完了"

.PHONY: deploy-api-sessions
deploy-api-sessions: build-api-sessions
	@echo "🚀 api-sessions をデプロイ中..."
	serverless deploy function --function api-sessions --aws-profile selfpomodoro-dev
	@echo "✅ api-sessions デプロイ完了"

.PHONY: deploy-api-rounds
deploy-api-rounds: build-api-rounds
	@echo "🚀 api-rounds をデプロイ中..."
	serverless deploy function --function api-rounds --aws-profile selfpomodoro-dev
	@echo "✅ api-rounds デプロイ完了"

.PHONY: deploy-api-statistics
deploy-api-statistics: build-api-statistics
	@echo "🚀 api-statistics をデプロイ中..."
	serverless deploy function --function api-statistics --aws-profile selfpomodoro-dev
	@echo "✅ api-statistics デプロイ完了"

# テスト実行
.PHONY: test-all
test-all:
	@echo "🧪 全APIテストを実行中..."
	@chmod +x test-all-apis.sh
	@./test-all-apis.sh

# ログ確認
.PHONY: logs-api-tasks
logs-api-tasks:
	@echo "📝 api-tasks のログを表示..."
	aws logs tail /aws/lambda/selfpomodoro-serverless-dev-api-tasks --follow --profile selfpomodoro-dev

.PHONY: logs-api-sessions
logs-api-sessions:
	@echo "📝 api-sessions のログを表示..."
	aws logs tail /aws/lambda/selfpomodoro-serverless-dev-api-sessions --follow --profile selfpomodoro-dev

.PHONY: logs-api-rounds
logs-api-rounds:
	@echo "📝 api-rounds のログを表示..."
	aws logs tail /aws/lambda/selfpomodoro-serverless-dev-api-rounds --follow --profile selfpomodoro-dev

.PHONY: logs-api-statistics
logs-api-statistics:
	@echo "📝 api-statistics のログを表示..."
	aws logs tail /aws/lambda/selfpomodoro-serverless-dev-api-statistics --follow --profile selfpomodoro-dev

# マイグレーション（既存機能維持）
.PHONY: migrate-up
migrate-up:
	@echo "マイグレーションを実行: ./migrations/sql"
	@if [ -f .env ]; then \
		export $$(grep -v '^#' .env | xargs); \
		migrate -path ./migrations/sql -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSL_MODE" up; \
	else \
		echo "❌ .env ファイルが見つかりません"; \
		exit 1; \
	fi

.PHONY: migrate-down
migrate-down:
	@echo "マイグレーションをロールバック: ./migrations/sql"
	@if [ -f .env ]; then \
		export $$(grep -v '^#' .env | xargs); \
		migrate -path ./migrations/sql -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSL_MODE" down 1; \
	else \
		echo "❌ .env ファイルが見つかりません"; \
		exit 1; \
	fi

# 新しいマイグレーションファイルを作成
.PHONY: migrate-create
migrate-create:
	@read -p "マイグレーション名を入力してください: " name; \
	./scripts/migrate.sh create $$name

# クリーンアップ
.PHONY: clean
clean:
	@echo "🧹 ビルド成果物を削除中..."
	@find . -name "bootstrap.zip" -delete
	@find . -name "bootstrap" -delete
	@echo "✅ クリーンアップ完了"

# 開発用コマンド
.PHONY: dev
dev:
	@echo "🏃 開発モードでの迅速デプロイ..."
	@echo "1. 全関数ビルド"
	@make build-all
	@echo "2. デプロイ"
	@make deploy-all
	@echo "3. テスト実行"
	@sleep 10  # デプロイ完了待機
	@make test-all

# 情報表示
.PHONY: info
info:
	@echo "📊 現在の設定情報:"
	@echo "  プロジェクト: selfpomodoro-serverless"
	@echo "  ステージ: dev"
	@echo "  リージョン: ap-northeast-1"
	@echo "  プロファイル: selfpomodoro-dev"
	@echo ""
	@echo "📁 Lambda関数:"
	@for func in $(FUNCTIONS); do echo "  - $func"; done
	@echo ""
	@echo "🌐 API Gateway:"
	@echo "  Base URL: https://f0vvgsmjj7.execute-api.ap-northeast-1.amazonaws.com/dev"

# エラー時のヘルプ表示
.PHONY: error-help
error-help:
	@echo "❌ よくあるエラーと解決方法:"
	@echo ""
	@echo "1. ビルドエラー:"
	@echo "   - go.mod の依存関係を確認: go mod tidy"
	@echo "   - internal パッケージのパスを確認"
	@echo ""
	@echo "2. デプロイエラー:"
	@echo "   - AWS認証情報を確認: aws configure list --profile selfpomodoro-dev"
	@echo "   - serverless.yml の設定を確認"
	@echo ""
	@echo "3. テストエラー:"
	@echo "   - .env ファイルの設定を確認"
	@echo "   - API Gateway のエンドポイントURLを確認"
