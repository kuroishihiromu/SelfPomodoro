# バックエンド

このディレクトリは、ポモドーロ最適化アプリ用のバックエンドAPIサーバーやAWSリソースとのコネクションを構築する。

## 機能

- ユーザー認証（Cognito JWT認証）
- タスク管理（CRUD操作）
- セッション管理
- ラウンド管理
- 集中度の記録
- セッションとラウンドの最適化
- 時間帯別の集中度ヒートマップ
- 集中度評価の推移グラフ

## 技術スタック

- **言語**: Go 1.24.3
- **フレームワーク**: Echo
- **データベース**: 
  - PostgreSQL (RDS) - 主要データ
  - DynamoDB - 最適化ログと設定
- **認証**: AWS Cognito
- **メッセージング**: AWS SNS/SQS
- **サーバーレス**: AWS Lambda
- **インフラ**: AWS (EC2, ALB, VPC, etc.)
- **IaC**: Terraform
- **CI/CD**: GitHub Actions

## 開発環境のセットアップ

### 前提条件のインストール

#### 1. Goのインストール

**Windows:**

1. [Go公式サイト](https://golang.org/dl/) から Windows 用のインストーラーをダウンロード
2. ダウンロードしたインストーラーを実行し、画面の指示に従う
3. インストール完了後、コマンドプロンプトを開き、`go version` を実行して1.24.3であることを確認

**macOS:**

Homebrewを使用する:
```bash
brew install go
```
ダウンロード後に `go version` を実行して1.24.3であることを確認してください！


**Linux:**

Ubuntuの場合:
```bash
sudo apt update
sudo apt install golang-go
```

または、公式バイナリを使用:
```bash
wget https://golang.org/dl/go1.24.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.3.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile
```
ダウンロード後に `go version` を実行して1.24.3であることを確認してください！

#### 2. PostgreSQLに関して

今回は簡易性や環境統一の観点からPostgres(ローカルDB)はDocker環境で構築します。
したがって、自分のPCにDocker Desktopをインストールする必要があります。(まだインストールしていない人のみ)

https://www.docker.com/ja-jp/products/docker-desktop/ から自身のOSにあったインストーラをダウンロードして、インストールをしてください！


### プロジェクトのセットアップと実行

#### 1. 依存関係のインストール

```bash
go mod download
```

#### 2. 環境設定

```bash
# .env.example をコピーして .env を作成
cp .env.example .env
```

#### 3. GoサーバとDocker起動
**Docker Desktopのアプリケーションを立ち上げた後**に、以下のコマンドを実行
```bash
# 内部で go run cmd/api/main.go と docker compose up -d --build　が走っています
make dev
```

以下のようなログがコンソールに出力されれば、Goサーバの起動、Postgresとpgadminの起動、そしてGoサーバとPostres間の通信が成功しています！
```bash
Dockerコンテナが起動しました
PostgreSQL: localhost:5432
PgAdmin: http://localhost:8081
Goサーバーを起動しています...
設定ファイル読み込み: .env
2025-05-10T16:12:47.543848+09:00        INFO    api/main.go:43  ポモドーロAPIサーバー 起動中... バージョン: 開発版, コミット: unknown, ビルド時間: unknown, 環境: development
2025-05-10T16:12:47.544599+09:00        INFO    database/postgres.go:28 PostgreSQLに接続: localhost:5432/pomodoro
2025-05-10T16:12:47.552867+09:00        INFO    database/postgres.go:44 PostgreSQL接続成功
2025-05-10T16:12:47.553807+09:00        INFO    app/server.go:84        サーバーをポート 8080 で起動

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.13.3
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:8080
```


- PostgreSQL→http://localhost:5432
- PgAdmin→http://localhost:8081
- Goサーバ(Echo)→http://localhost:8080

で起動しています！

#### 4. 動作確認
**Goサーバ:**

```bash
http://localhost:8080/health
```
結果

```bash
{"status":"OK","version":"開発版"}
```

**pgadmin(DBをGUIで閲覧できるサービス):**

```bash
http://localhost:8081
```

上記のエンドポイントにアクセスして、
 - Email →　admin@example.com
 - password → admin
 
 でログインできたらOK!

 ## 2025/05/14 タスク処理の動作確認
 ### 手順
 ```bash
 # backendディレクトリに移動
 cd backend

# マイグレーションスクリプトの実行権限付与
 chmod +x scripts/migrate.sh

 # 一応前回のキャッシュが残ってないか確認のためサーバーダウン(ホストPCのDockerDesktopを開いてね！)
 make docker-down

 # サーバ起動
 make dev

 # migration実行
 make migrate-up

 # Task動作確認(今回は開発環境用にダミーのトークンを使用)
 curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer dev-token" \
  -H "Content-Type: application/json" \
  -d '{"detail": "テストタスク"}'

  #以下のようなレスポンンスが来たら成功(idは一意生成なので違くてOK)
  # 今回は固定値のuser_idに基づいてタスクを作成
  {"id":"70d7757d-b1b7-467c-a5d4-da09f2f6f0e1","detail":"テストタスク","is_completed":false,"created_at":"2025-05-14T18:04:43.476769+09:00","updated_at":"2025-05-14T18:04:43.476769+09:00"}
```

### 他のエンドポイントも試してみるといいかも
 - 全部のタスクをとってくる
 ```bash
 curl -X GET http://localhost:8080/api/v1/tasks \
 -H "Authorization: Bearer dev-token" \
 -H "Content-Type: application/json" \
 ```

 - タスクの詳細を編集する
 ```bash
 curl -X PATCH http://localhost:8080/api/v1/tasks/idのとこに出てきた文字列/edit \
-H "Authorization: Bearer dev-token" \
-H "Content-Type: application/json" \
-d '{"detail": "テストタスク更新"}'
 ```

 - タスクの完了状態を切り替える
 ```bash
 curl -X PATCH http://localhost:8080/api/v1/tasks/idのとこに出てきた文字列/toggle \
-H "Authorization: Bearer dev-token" \
-H "Content-Type: application/json" \
 ```

 - タスクを削除する
 ```bash
 curl -X DELETE http://localhost:8080/api/v1/tasks/idのとこに出てきた文字列 \
-H "Authorization: Bearer dev-token" \
-H "Content-Type: application/json" \
```

※後々ユニットテストも実装予定です

### pgadminでの確認
作成したタスクをGUIで確認しよう！
```bash
# 以下にアクセス
localhost:8081
```
以下のような画面が出てくるはず
![Image](https://github.com/user-attachments/assets/67fd26e0-6d8b-49c2-8f5e-28446daef3ca)

Local Docker PostgreSQLをクリックするとパスワードを求められるので、
**postgres**
と入力してください。


