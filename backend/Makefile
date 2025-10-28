# Selfpomodoro Serverless Makefile

# 変数定義
GOOS := linux
GOARCH := amd64
CGO_ENABLED := 0

# Lambda関数一覧
FUNCTIONS := api-tasks api-sessions api-rounds api-statistics api-optimization-preferences post-confirmation optimization-round optimization-session

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
	@echo "Cognito Trigger:"
	@echo "  make build-post-confirmation  - PostConfirmation関数をビルド"
	@echo "  make deploy-post-confirmation - PostConfirmation関数をデプロイ"
	@echo ""
	@echo "🎯 最適化システム監視:"
	@echo "  make check-optimization    - 最適化システム動作確認"
	@echo "  make optimization-logs     - 最適化Lambda関数ログ表示"
	@echo "  make optimization-status   - 最適化システム総合ステータス"
	@echo "  make deploy-opt-session IMAGE_TAG=<tag> - セッション最適化Lambdaを指定イメージで更新"
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

.PHONY: build-api-optimization-preferences
build-api-optimization-preferences:
	@echo "🔨 api-optimization-preferences をビルド中..."
	@cd functions/api-optimization-preferences && \
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) \
	go build -ldflags="-s -w" -o bootstrap main.go && \
	zip bootstrap.zip bootstrap && \
	rm bootstrap
	@echo "✅ api-optimization-preferences ビルド完了"

# PostConfirmation Lambda関数ビルド
.PHONY: build-post-confirmation
build-post-confirmation:
	@echo "🔨 post-confirmation をビルド中..."
	@cd functions/cognito-triggers/post-confirmation && \
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) \
	go build -ldflags="-s -w" -o bootstrap main.go && \
	zip bootstrap.zip bootstrap && \
	rm bootstrap
	@echo "✅ post-confirmation ビルド完了"

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

.PHONY: deploy-api-optimization-preferences
deploy-api-optimization-preferences: build-api-optimization-preferences
	@echo "🚀 api-optimization-preferences をデプロイ中..."
	serverless deploy function --function api-optimization-preferences --aws-profile selfpomodoro-dev
	@echo "✅ api-optimization-preferences デプロイ完了"

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

# PostConfirmation関数デプロイ
.PHONY: deploy-post-confirmation
deploy-post-confirmation: build-post-confirmation
	@echo "🚀 post-confirmation をデプロイ中..."
	serverless deploy function --function post-confirmation --aws-profile selfpomodoro-dev
	@echo "✅ post-confirmation デプロイ完了"

# テスト実行
.PHONY: test-all
test-all:
	@echo "🧪 全APIテストを実行中..."
	@chmod +x test-all-apis.sh
	@./test-all-apis.sh

# PostConfirmation テスト
.PHONY: test-post-confirmation
test-post-confirmation:
	@echo "🧪 PostConfirmation Triggerテストを実行中..."
	@echo "テストユーザーでCognito ユーザー確認を実行..."
	@read -p "テストユーザー名を入力してください: " username; \
	aws cognito-idp admin-confirm-sign-up \
		--user-pool-id ap-northeast-1_UltsO4poH \
		--username $$username \
		--profile selfpomodoro-dev

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

.PHONY: logs-post-confirmation
logs-post-confirmation:
	@echo "📝 post-confirmation のログを表示..."
	aws logs tail /aws/lambda/selfpomodoro-serverless-dev-post-confirmation --follow --profile selfpomodoro-dev


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

# PostConfirmation開発用
.PHONY: dev-cognito
dev-cognito:
	@echo "🔐 Cognito PostConfirmation 開発モード..."
	@echo "1. PostConfirmation関数ビルド"
	@make build-post-confirmation
	@echo "2. PostConfirmation関数デプロイ"
	@make deploy-post-confirmation
	@echo "3. テスト準備完了"
	@echo "📝 テスト実行: make test-post-confirmation"

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
	@for func in $(FUNCTIONS); do echo "  - $$func"; done
	@echo ""
	@echo "🌐 API Gateway:"
	@echo "  Base URL: https://rtc4ooslz6.execute-api.ap-northeast-1.amazonaws.com/dev"
	@echo ""
	@echo "🔐 Cognito:"
	@echo "  User Pool ID: ap-northeast-1_UltsO4poH"
	@echo "  Client ID: 2dbjjcvvon3trtt5imonmpnq9g"

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

# 🎯 最適化Lambda関数（Python）
.PHONY: build-optimization-round
build-optimization-round:
	@echo "🔨 optimization-round をビルド中..."
	@cd functions/optimization-round && \
	pip install -r requirements.txt -t . --no-deps && \
	zip -r optimization-round.zip . -x "*.pyc" "__pycache__/*" "*.git*" && \
	echo "✅ optimization-round ビルド完了"

.PHONY: build-optimization-session
build-optimization-session:
	@echo "🔨 optimization-session をビルド中..."
	@cd functions/optimization-session && \
	pip install -r requirements.txt -t . --no-deps && \
	zip -r optimization-session.zip . -x "*.pyc" "__pycache__/*" "*.git*" && \
	echo "✅ optimization-session ビルド完了"

.PHONY: deploy-optimization-round
deploy-optimization-round:
	@echo "🚀 optimization-round をデプロイ中..."
	@npx serverless deploy function --function optimization-round --stage dev
	@echo "✅ optimization-round デプロイ完了"

.PHONY: deploy-optimization-session
deploy-optimization-session:
	@echo "🚀 optimization-session をデプロイ中..."
	@npx serverless deploy function --function optimization-session --stage dev
	@echo "✅ optimization-session デプロイ完了"

# 最適化関数統合操作
.PHONY: build-optimization-all
build-optimization-all: build-optimization-round build-optimization-session
	@echo "✅ 全最適化Lambda関数のビルドが完了しました"

.PHONY: deploy-optimization-all
deploy-optimization-all: deploy-optimization-round deploy-optimization-session
	@echo "✅ 全最適化Lambda関数のデプロイが完了しました"

.PHONY: logs-optimization-round
logs-optimization-round:
	@echo "📋 optimization-round のログを表示中..."
	@npx serverless logs --function optimization-round --stage dev --tail

.PHONY: logs-optimization-session
logs-optimization-session:
	@echo "📋 optimization-session のログを表示中..."
	@npx serverless logs --function optimization-session --stage dev --tail

# 🎯 最適化動作確認コマンド
.PHONY: check-optimization
check-optimization:
	@echo "📊 最適化システム動作確認を開始..."
	@echo ""
	@echo "🔄 1. SQSキュー状況確認:"
	@aws sqs get-queue-attributes \
		--queue-url https://sqs.ap-northeast-1.amazonaws.com/545009834501/selfpomodoro-round-optimization-dev \
		--attribute-names ApproximateNumberOfMessages \
		--profile selfpomodoro-dev \
		--query 'Attributes.ApproximateNumberOfMessages' \
		--output text | \
		awk '{print "   Round最適化キュー: " $$1 " メッセージ待機中"}'
	@aws sqs get-queue-attributes \
		--queue-url https://sqs.ap-northeast-1.amazonaws.com/545009834501/selfpomodoro-session-optimization-dev \
		--attribute-names ApproximateNumberOfMessages \
		--profile selfpomodoro-dev \
		--query 'Attributes.ApproximateNumberOfMessages' \
		--output text | \
		awk '{print "   Session最適化キュー: " $$1 " メッセージ待機中"}'
	@echo ""
	@echo "📈 2. 最適化ログ確認（直近5件）:"
	@aws dynamodb scan \
		--table-name selfpomodoro_unified_table_dev \
		--filter-expression "begins_with(SK, :sk_prefix)" \
		--expression-attribute-values '{":sk_prefix": {"S": "OPTIMIZATION_LOG#"}}' \
		--projection-expression "SK, work_time, break_time, focus_score, updated_at" \
		--max-items 5 \
		--profile selfpomodoro-dev \
		--output table
	@echo ""
	@echo "⚙️  3. UserConfig最新状況:"
	@aws dynamodb get-item \
		--table-name selfpomodoro_unified_table_dev \
		--key '{"PK": {"S": "USER#67448aa8-0061-7040-0113-b1d35fcf76f6"}, "SK": {"S": "CONFIG"}}' \
		--projection-expression "round_work_time, round_break_time, session_break_time, session_rounds, updated_at" \
		--profile selfpomodoro-dev \
		--output table
	@echo ""
	@echo "✅ 最適化システム動作確認完了"

.PHONY: optimization-logs
optimization-logs:
	@echo "📋 最適化Lambda関数の最新ログを確認..."
	@echo ""
	@echo "🔄 Round最適化ログ（直近10件）:"
	@aws logs describe-log-streams \
		--log-group-name "/aws/lambda/selfpomodoro-serverless-dev-optimization-round-v2" \
		--order-by LastEventTime \
		--descending \
		--max-items 1 \
		--profile selfpomodoro-dev \
		--query 'logStreams[0].logStreamName' \
		--output text | \
		xargs -I {} aws logs get-log-events \
		--log-group-name "/aws/lambda/selfpomodoro-serverless-dev-optimization-round-v2" \
		--log-stream-name {} \
		--profile selfpomodoro-dev \
		--query 'events[-10:].message' \
		--output text
	@echo ""
	@echo "🔄 Session最適化ログ（直近10件）:"
	@aws logs describe-log-streams \
		--log-group-name "/aws/lambda/selfpomodoro-serverless-dev-optimization-session-v2" \
		--order-by LastEventTime \
		--descending \
		--max-items 1 \
		--profile selfpomodoro-dev \
		--query 'logStreams[0].logStreamName' \
		--output text | \
		xargs -I {} aws logs get-log-events \
		--log-group-name "/aws/lambda/selfpomodoro-serverless-dev-optimization-session-v2" \
		--log-stream-name {} \
		--profile selfpomodoro-dev \
		--query 'events[-10:].message' \
		--output text || echo "   Session最適化ログは未実行です"

.PHONY: optimization-status
optimization-status:
	@echo "🎯 最適化システム総合ステータス"
	@echo "=================================="
	@echo ""
	@echo "📊 最適化実行回数:"
	@aws dynamodb scan \
		--table-name selfpomodoro_unified_table_dev \
		--filter-expression "begins_with(SK, :sk_prefix)" \
		--expression-attribute-values '{":sk_prefix": {"S": "OPTIMIZATION_LOG#ROUND#"}}' \
		--select COUNT \
		--profile selfpomodoro-dev \
		--query 'Count' \
		--output text | \
		awk '{print "   Round最適化: " $$1 " 回実行"}'
	@aws dynamodb scan \
		--table-name selfpomodoro_unified_table_dev \
		--filter-expression "begins_with(SK, :sk_prefix)" \
		--expression-attribute-values '{":sk_prefix": {"S": "OPTIMIZATION_LOG#SESSION#"}}' \
		--select COUNT \
		--profile selfpomodoro-dev \
		--query 'Count' \
		--output text | \
		awk '{print "   Session最適化: " $$1 " 回実行"}'
	@echo ""
	@echo "📈 最新の最適化結果:"
	@aws dynamodb scan \
		--table-name selfpomodoro_unified_table_dev \
		--filter-expression "begins_with(SK, :sk_prefix)" \
		--expression-attribute-values '{":sk_prefix": {"S": "OPTIMIZATION_LOG#"}}' \
		--projection-expression "SK, work_time, break_time, focus_score" \
		--max-items 1 \
		--profile selfpomodoro-dev \
		--query 'Items[0]' \
		--output table
	@echo ""
	@echo "⚙️  現在のUserConfig:"
	@aws dynamodb get-item \
		--table-name selfpomodoro_unified_table_dev \
		--key '{"PK": {"S": "USER#67448aa8-0061-7040-0113-b1d35fcf76f6"}, "SK": {"S": "CONFIG"}}' \
		--projection-expression "round_work_time, round_break_time, session_break_time, session_rounds" \
		--profile selfpomodoro-dev \
		--query 'Item' \
		--output table

.PHONY: deploy-opt-session
deploy-opt-session:
	@if [ -z "$(IMAGE_TAG)" ]; then \
		echo "❌ エラー: IMAGE_TAG を指定してください (例: make deploy-opt-session IMAGE_TAG=v2)"; \
		exit 1; \
	fi
	@echo "🚀 optimization-session-v2 をイメージタグ $(IMAGE_TAG) で更新中..."
	aws lambda update-function-code \
		--function-name selfpomodoro-serverless-dev-optimization-session-v2 \
		--image-uri 545009834501.dkr.ecr.ap-northeast-1.amazonaws.com/selfpomodoro-optimization-session:$(IMAGE_TAG) \
		--profile selfpomodoro-dev \
		--region ap-northeast-1 >/dev/null
	@echo "✅ optimization-session-v2 を更新しました (image tag: $(IMAGE_TAG))"
