#!/usr/bin/env bash

# エラーが発生した場合に終了
set -e

# 環境変数を読み込む
ENV=${1:-development}

if [ -f .env.$ENV ]; then
  export $(grep -v '^#' .env.$ENV | xargs)
fi

# ビルド情報
VERSION=$(git describe --tags --always --dirty || echo "unknown")
COMMIT=$(git rev-parse --short HEAD || echo "unknown")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

echo "ビルド環境: $ENV"
echo "バージョン: $VERSION"
echo "コミット: $COMMIT"
echo "ビルド時間: $BUILD_TIME"

# ビルドディレクトリの作成
mkdir -p ./bin

# ビルドフラグの設定
BUILD_FLAGS="-X 'main.Version=$VERSION' -X 'main.Commit=$COMMIT' -X 'main.BuildTime=$BUILD_TIME' -X 'main.Environment=$ENV'"

# ビルドの実行
echo "APIサーバーをビルド中..."
GOOS=linux GOARCH=amd64 go build -ldflags "$BUILD_FLAGS" -o ./bin/api ./cmd/api/main.go

echo "ビルド完了: ./bin/api"