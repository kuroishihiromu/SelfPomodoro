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
 # migrationをするためのパッッケージインストール(Macのみ記載)
 brew install golang-migrate

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

スキーマ>public>テーブル>tasks で、tasksを右クリックすると以下のようなメニューが出てくるので、「すべての行」をクリックするとデータ見れます！

※ pgadminを日本語にする設定は以下のサイトから
https://qiita.com/sanapuuu/items/4e43f6ed0cf0a597efb5

![Image](https://github.com/user-attachments/assets/5e859e9b-8741-446b-9e67-d57e5b31219c)


## 2025/0518 セッション・ラウンド管理機能の動作確認

### 前提
---

※動作確認は、3の**典型的なワークフロー**を試すだけでも十分だと思います！

全部のエンドポイントが気になるなら1, 2に記載している他のやつも叩いてみてください！

---

まず、前回と同じ要領で、
 ```bash
 # backendディレクトリに移動
 cd backend

 # 一応前回のキャッシュが残ってないか確認のためサーバーダウン(ホストPCのDockerDesktopを開いてね！)
 make docker-down

 # サーバ起動
 make dev

 # migration実行
 make migrate-up
 ```

### 1. セッション管理機能

#### セッション開始

```bash
curl -X POST http://localhost:8080/api/v1/sessions \
  -H "Authorization: Bearer dev-token" \
  -H "Content-Type: application/json"

レスポンス例:
 {
  "id": "セッションID",
  "start_time": "2025-05-18T04:20:48.371273+09:00",
  "end_time": null,
  "average_focus": null,
  "total_work_min": null,
  "round_count": null,
  "break_time": null
  }
  ```

#### セッション一覧取得
```bash
curl -X GET http://localhost:8080/api/v1/sessions \
  -H "Authorization: Bearer dev-token" \
  -H "Content-Type: application/json"
  ```

#### 特定のセッション取得
```bash
curl -X GET http://localhost:8080/api/v1/sessions/{session_id} \
  -H "Authorization: Bearer dev-token" \
  -H "Content-Type: application/json"
  ```

#### セッション完了
```bash
curl -X PATCH http://localhost:8080/api/v1/sessions/{session_id}/complete \
  -H "Authorization: Bearer dev-token" \
  -H "Content-Type: application/json"
  ```


### 2. ラウンド管理機能
#### ラウンド開始
```bash
curl -X POST http://localhost:8080/api/v1/sessions/{session_id}/rounds \
  -H "Authorization: Bearer dev-token" \
  -H "Content-Type: application/json"

レスポンス例:
 {
  "id": "ラウンドID",
  "session_id": "セッションID",
  "round_order": 1,
  "start_time": "2025-05-18T04:25:12.371273+09:00",
  "end_time": null,
  "work_time": null,
  "break_time": null,
  "focus_score": null,
  "is_aborted": false
}
```

#### セッションの全ラウンド取得
```bash
curl -X GET http://localhost:8080/api/v1/sessions/{session_id}/rounds \
  -H "Authorization: Bearer dev-token" \
  -H "Content-Type: application/json"
  ```

#### 特定のラウンド取得
```bash
curl -X GET http://localhost:8080/api/v1/rounds/{round_id} \
  -H "Authorization: Bearer dev-token" \
  -H "Content-Type: application/json"
  ```

#### ラウンド完了
```bash
curl -X PATCH http://localhost:8080/api/v1/rounds/{round_id}/complete \
  -H "Authorization: Bearer dev-token" \
  -H "Content-Type: application/json" \
  -d '{"focus_score": 85}'
  ```
#### ラウンド中断
```bash
curl -X POST http://localhost:8080/api/v1/rounds/{round_id}/abort \
  -H "Authorization: Bearer dev-token" \
  -H "Content-Type: application/json"
  ```

### 3. 典型的なワークフロー
以下は時系列順のセッション・ラウンドのワークフローの例です!

フローで気になるところとかあったら修正するんでなんでも言ってください！

#### セッション開始
```bash
curl -X POST http://localhost:8080/api/v1/sessions \
  -H "Authorization: Bearer dev-token" \
  -H "Content-Type: application/json"
```

#### 最初のラウンド開始
```bash
curl -X POST http://localhost:8080/api/v1/sessions/$SESSION_ID/rounds \
  -H "Authorization: Bearer dev-token" \
  -H "Content-Type: application/json"
```

#### ラウンド完了（集中度スコア入力）
```bash
curl -X PATCH http://localhost:8080/api/v1/rounds/$ROUND_ID/complete \
  -H "Authorization: Bearer dev-token" \
  -H "Content-Type: application/json" \
  -d '{"focus_score": 90}'
  ```

#### 2回目のラウンド開始
```bash
curl -X POST http://localhost:8080/api/v1/sessions/$SESSION_ID/rounds \
  -H "Authorization: Bearer dev-token" \
  -H "Content-Type: application/json"
  ```

#### 3回目のラウンド開始
```bash
curl -X POST http://localhost:8080/api/v1/sessions/$SESSION_ID/rounds \
  -H "Authorization: Bearer dev-token" \
  -H "Content-Type: application/json"
```

#### 3回目のラウンド完了
```bash
curl -X PATCH http://localhost:8080/api/v1/rounds/$ROUND_ID/complete \
  -H "Authorization: Bearer dev-token" \
  -H "Content-Type: application/json" \
  -d '{"focus_score": 95}'
  ```

#### セッション完了（統計計算）
```bash
curl -X PATCH http://localhost:8080/api/v1/sessions/$SESSION_ID/complete \
  -H "Authorization: Bearer dev-token" \
  -H "Content-Type: application/json"
  ```

#### セッションの最終結果確認
```bash
curl -X GET http://localhost:8080/api/v1/sessions/$SESSION_ID \
  -H "Authorization: Bearer dev-token" \
  -H "Content-Type: application/json"
  ```

4. 実装上の注意点

- ラウンドの順序は自動的に管理される（クライアントが指定する必要はない）
- 進行中のラウンドがある場合、新しいラウンドを開始することはできない
- セッション完了時にはラウンドの統計情報（平均集中度、総作業時間など）が自動計算される
- 集中度スコアは0〜100の整数値で入力
- 現在の実装では作業時間と休憩時間はバックエンドでデフォルト値（25分/5分）に設定されている(Dynamoと最適化との連携がまだなため)

現状
- セッションは完了してなくても新しいセッション開始のリクエストは通る
- ラウンドが中断した後は、同セッション内でラウンド開始リクエストが通ってしまう

↓

より詳細なフローの認識をすり合わせて最終的な処理の実装方針を決めたい！





