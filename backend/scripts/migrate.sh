#!/usr/bin/env bash

# エラーが発生した場合に終了
set -e

# 環境変数の読み込み
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

# データベース接続文字列
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"

# マイグレーションディレクトリ
MIGRATIONS_DIR="./migrations/sql"

# コマンドラインパラメータの処理
CMD=$1
VERSION=$2

# マイグレーションコマンドの実行
case $CMD in
  up)
    echo "マイグレーションを実行: $MIGRATIONS_DIR -> $DB_URL"
    migrate -path $MIGRATIONS_DIR -database $DB_URL up
    ;;
  down)
    echo "マイグレーションをロールバック: $MIGRATIONS_DIR -> $DB_URL"
    migrate -path $MIGRATIONS_DIR -database $DB_URL down ${VERSION:-1}
    ;;
  create)
    if [ -z "$VERSION" ]; then
      echo "エラー: マイグレーション名を指定してください"
      echo "使用法: $0 create <マイグレーション名>"
      exit 1
    fi
    echo "マイグレーションファイルを作成: $VERSION"
    migrate create -ext sql -dir $MIGRATIONS_DIR -seq $VERSION
    ;;
  force)
    if [ -z "$VERSION" ]; then
      echo "エラー: バージョン番号を指定してください"
      echo "使用法: $0 force <バージョン>"
      exit 1
    fi
    echo "マイグレーションバージョンを強制設定: $VERSION"
    migrate -path $MIGRATIONS_DIR -database $DB_URL force $VERSION
    ;;
  version)
    echo "現在のマイグレーションバージョンを確認:"
    migrate -path $MIGRATIONS_DIR -database $DB_URL version
    ;;
  *)
    echo "使用法: $0 {up|down|create|force|version} [引数]"
    echo "  up                  : 全てのマイグレーションを適用"
    echo "  down [N]            : N個のマイグレーションをロールバック（デフォルト: 1）"
    echo "  create <name>       : 新しいマイグレーションファイルを作成"
    echo "  force <version>     : マイグレーションバージョンを強制設定"
    echo "  version             : 現在のマイグレーションバージョンを表示"
    exit 1
    ;;
esac

echo "完了"
