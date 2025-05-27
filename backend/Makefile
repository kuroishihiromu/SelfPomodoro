# 複数Lambda関数対応 Makefile

.PHONY: build deploy clean build-health build-api-tasks

# 全関数ビルド
build: build-health build-api-tasks

# health関数ビルド
build-health:
	@echo "Building health function..."
	GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
	zip bootstrap.zip bootstrap
	rm bootstrap

# api-tasks関数ビルド
build-api-tasks:
	@echo "Building api-tasks function..."
	cd functions/api-tasks && GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
	cd functions/api-tasks && zip bootstrap.zip bootstrap
	cd functions/api-tasks && rm bootstrap

# デプロイ
deploy: build
	@echo "Deploying to AWS..."
	npx serverless deploy --stage dev --verbose

# クリーンアップ
clean:
	@echo "Cleaning up build artifacts..."
	rm -f bootstrap.zip
	rm -f functions/*/bootstrap.zip
	rm -f functions/*/bootstrap

# 開発用：health関数のみビルド・デプロイ
deploy-health: build-health
	npx serverless deploy function --function health --stage dev

# 開発用：api-tasks関数のみビルド・デプロイ
deploy-api-tasks: build-api-tasks
	npx serverless deploy function --function api-tasks --stage dev
