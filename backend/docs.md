# Self Pomodoro アプリケーション仕様書

## 概要

Self Pomodoroは、ポモドーロテクニックをベースとしたタスク管理・時間管理アプリケーションです。ユーザーがタスクを作成し、ポモドーロセッションを通じて集中的に作業を行い、統計データに基づいて最適化されたワークフローを提供します。

## システムアーキテクチャ

### インフラ構成
- **プラットフォーム**: AWS Serverless
- **API**: AWS Lambda + API Gateway
- **認証**: AWS Cognito
- **データベース**: 
  - PostgreSQL（メインデータ）
  - DynamoDB（ユーザー設定・最適化データ）
- **非同期処理**: AWS SQS
- **ログ**: AWS CloudWatch
- **デプロイ**: Serverless Framework

### Lambda関数構成
- `api-tasks`: タスク管理API
- `api-sessions`: セッション管理API
- `api-rounds`: ラウンド管理API
- `api-statistics`: 統計データAPI
- `post-confirmation`: Cognito PostConfirmation Trigger

## 機能仕様

### 1. タスク管理機能

#### 1.1 データモデル
```go
type Task struct {
    ID          uuid.UUID
    UserID      uuid.UUID
    Detail      string
    IsCompleted bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

#### 1.2 API エンドポイント
- `GET /api/v1/tasks` - タスク一覧取得
- `POST /api/v1/tasks` - タスク作成
- `PATCH /api/v1/tasks/{task_id}/edit` - タスク更新
- `PATCH /api/v1/tasks/{task_id}/toggle` - 完了状態切り替え
- `DELETE /api/v1/tasks/{task_id}` - タスク削除

#### 1.3 バリデーション
- タスク詳細は必須項目
- ユーザーIDによるアクセス権限制御

### 2. セッション管理機能

#### 2.1 データモデル
```go
type Session struct {
    ID           uuid.UUID
    UserID       uuid.UUID
    StartTime    time.Time
    EndTime      *time.Time
    AverageFocus *float64
    TotalWorkMin *int
    RoundCount   *int
    BreakTime    *int
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

#### 2.2 API エンドポイント
- `GET /api/v1/sessions` - セッション一覧取得
- `POST /api/v1/sessions` - セッション開始
- `GET /api/v1/sessions/{session_id}` - セッション取得
- `PATCH /api/v1/sessions/{session_id}/complete` - セッション完了
- `DELETE /api/v1/sessions/{session_id}` - セッション削除

#### 2.3 ビジネスロジック
- セッション開始時に新しいUUIDを生成
- 完了時に統計データ（平均集中度、総作業時間、ラウンド数）を計算
- 最適化メッセージ送信判定とSQS投稿

### 3. ラウンド管理機能

#### 3.1 データモデル
```go
type Round struct {
    ID         uuid.UUID
    SessionID  uuid.UUID
    UserID     uuid.UUID
    RoundType  string // "work" or "break"
    StartTime  time.Time
    EndTime    *time.Time
    FocusScore *int
    CreatedAt  time.Time
    UpdatedAt  time.Time
}
```

#### 3.2 API エンドポイント
- `GET /api/v1/sessions/{session_id}/rounds` - ラウンド一覧取得
- `POST /api/v1/sessions/{session_id}/rounds` - ラウンド開始
- `GET /api/v1/rounds/{round_id}` - ラウンド取得
- `PATCH /api/v1/rounds/{round_id}/complete` - ラウンド完了
- `POST /api/v1/rounds/{round_id}/abort` - ラウンド中止

#### 3.3 ラウンドタイプ
- **work**: 作業ラウンド（ポモドーロ）
- **break**: 休憩ラウンド（短休憩・長休憩）

### 4. ユーザー設定機能

#### 4.1 データモデル
```go
type UserConfig struct {
    UserID           uuid.UUID
    RoundWorkTime    int // 1-120分
    RoundBreakTime   int // 1-60分
    SessionRounds    int // 1-10ラウンド
    SessionBreakTime int // 5-120分
    CreatedAt        time.Time
    UpdatedAt        time.Time
}
```

#### 4.2 デフォルト値とフォールバック機能
- **作業時間**: 25分
- **休憩時間**: 5分
- **セッションラウンド数**: 4回
- **セッション長休憩時間**: 15分

##### 安全なフォールバック機能
- ユーザー設定が取得できない場合、自動的にデフォルト値を使用
- DynamoDBが利用できない環境でもアプリケーションが動作継続
- `GetWorkTimeOrDefault()`, `GetBreakTimeOrDefault()` などの安全なアクセサメソッド

#### 4.3 バリデーション制約
- **作業時間**: 1-120分
- **休憩時間**: 1-60分
- **セッションラウンド数**: 1-10回
- **セッション長休憩時間**: 5-120分

#### 4.4 API エンドポイント
```
GET /api/v1/user-config     - ユーザー設定取得
POST /api/v1/user-config    - ユーザー設定作成
PUT /api/v1/user-config     - ユーザー設定更新
```

#### 4.5 ドメインロジック
- **設定の有効性チェック**: `IsValid()`, `ValidateSettings()`
- **最適化データベース値提供**: `GetOptimizationBaseValues()`
- **設定更新**: `UpdateSettings()` でドメインルール適用

### 5. オンボーディング機能

#### 5.1 PostConfirmation Trigger
AWS Cognito PostConfirmation Triggerを使用して、ユーザー登録完了時に自動的に初期設定を実行。

#### 5.2 オンボーディングフロー
```
1. Cognito PostConfirmation Trigger 発生
2. User作成（ドメインファクトリー使用）
3. UserConfig作成（デフォルト値設定）
4. サンプル最適化データ作成
```

#### 5.3 実行ステップ詳細

##### ステップ1: User作成
- Cognito属性から `CognitoUserParams` を生成
- `NewUserFromCognitoAttributes()` でドメインモデル作成
- `IsValidForCreation()` で作成前バリデーション
- プロバイダー情報（Google、Cognito UserPool等）の保存

##### ステップ2: UserConfig作成
- `NewDefaultUserConfig()` でデフォルト設定作成
- DynamoDBへの保存（失敗時も処理続行）
- 冪等性の確保（既存設定がある場合はスキップ）

##### ステップ3: サンプル最適化データ作成
- 過去10日分のサンプルデータを自動生成
- ラウンド最適化ログ・セッション最適化ログを作成
- リアルな使用パターンを模擬したデータ

#### 5.4 エラーハンドリング
- **User作成失敗**: オンボーディング全体が失敗
- **UserConfig作成失敗**: 警告ログ出力、処理継続（デフォルト値フォールバック）
- **サンプルデータ作成失敗**: 警告ログ出力、処理継続

### 6. 最適化サイクル機能

#### 6.1 最適化の仕組み
Self Pomodoroは、ユーザーの実際の作業パフォーマンスを学習し、個人に最適化された設定を提供します。

#### 6.2 最適化メッセージ

##### ラウンド最適化メッセージ
```go
type RoundOptimizationMessage struct {
    MessageID   string    // メッセージ一意ID
    MessageType string    // "round_optimization"
    Timestamp   time.Time // メッセージ作成時刻
    Version     string    // "2.0"
    UserID      string    // ユーザーID
    RoundID     string    // ラウンドID
    FocusScore  int       // 集中度スコア (0-100)
}
```

##### セッション最適化メッセージ
```go
type SessionOptimizationMessage struct {
    MessageID     string    // メッセージ一意ID
    MessageType   string    // "session_optimization"
    Timestamp     time.Time // メッセージ作成時刻
    Version       string    // "2.0"
    UserID        string    // ユーザーID
    SessionID     string    // セッションID
    AvgFocusScore float64   // 平均集中度スコア
    TotalWorkTime int       // 合計作業時間（分）
}
```

#### 6.3 最適化サイクルのフロー

##### ラウンド最適化サイクル
```
1. ラウンド完了時に集中度スコアを記録
2. 集中度スコアが有効な場合、SQSメッセージ送信
3. 最適化処理（別システム）でラウンド設定を調整
4. 次回ラウンド開始時に最適化された設定を適用
```

##### セッション最適化サイクル
```
1. セッション完了時に統計データを計算
2. 平均集中度・総作業時間を算出
3. SQSメッセージ送信（条件: 完了ラウンド存在）
4. 最適化処理でセッション設定を調整
5. 次回セッション開始時に最適化された設定を適用
```

#### 6.4 SQS メッセージ送信機能

##### 送信条件
- **ラウンド最適化**: 集中度スコアが設定されている場合
- **セッション最適化**: セッションが完了し、1つ以上のラウンドが存在する場合

##### 送信設定
- **最大リトライ回数**: 3回
- **リトライ間隔**: 1秒
- **メッセージタイムアウト**: 30秒
- **メッセージ属性**: MessageType, Version, Attempt, MessageId

##### エラーハンドリング
- タイムアウトエラー: リトライ実行
- 一時的エラー: リトライ実行
- 致命的エラー: 即座に終了

#### 6.5 最適化データの保存

##### DynamoDBテーブル
- `selfpomodoro_round_optimization_logs_dev`: ラウンド最適化ログ
- `selfpomodoro_session_optimization_logs_dev`: セッション最適化ログ

##### データ構造
```
RoundOptimizationLog:
  - user_id: ユーザーID
  - timestamp: ラウンド完了時刻
  - work_time: 次回推奨作業時間
  - break_time: 次回推奨休憩時間
  - focus_score: 実績集中度スコア

SessionOptimizationLog:
  - user_id: ユーザーID
  - timestamp: セッション完了時刻
  - round_count: 次回推奨ラウンド数
  - break_time: 次回推奨長休憩時間
  - avg_focus_score: 平均集中度
  - total_work_time: 総作業時間
```

#### 6.6 サンプルデータ生成機能
新規ユーザーに対して、過去10日分のリアルなサンプル最適化データを生成。

##### 生成パターン
- **初期集中度**: 60-80%
- **学習効果**: 日数が進むにつれて集中度向上
- **作業時間調整**: 実績に基づく微調整
- **休憩時間最適化**: 疲労度に応じた調整

##### データ特徴
- 日により異なる集中度パターン
- 1日3-5セッションの現実的な使用パターン
- セッションあたり2-6ラウンド
- 時間帯による集中度の変化を模擬

## 認証・認可

### 認証方式
- AWS Cognito JWT認証
- Bearer Token形式
- 開発環境では`dev-token`による簡易認証サポート

### 認可制御
- 全APIでユーザーIDによるリソースアクセス制御
- セッション、ラウンド、タスクは作成者のみがアクセス可能

## データベース設計

### PostgreSQL テーブル
- `users` - ユーザー情報
- `tasks` - タスク情報
- `sessions` - セッション情報
- `rounds` - ラウンド情報
- `statistics` - 統計情報（実装中）

### DynamoDB テーブル
- `selfpomodoro_user_configs_dev` - ユーザー設定
- `selfpomodoro_round_optimization_logs_dev` - ラウンド最適化ログ
- `selfpomodoro_session_optimization_logs_dev` - セッション最適化ログ

## エラーハンドリング

### エラー分類
1. **Domain Errors**: ビジネスロジックエラー
2. **Infrastructure Errors**: 技術的エラー
3. **Validation Errors**: 入力検証エラー

### HTTPステータスコード
- `200 OK`: 正常処理
- `201 Created`: リソース作成成功
- `400 Bad Request`: 入力エラー・バリデーションエラー
- `401 Unauthorized`: 認証エラー
- `403 Forbidden`: 認可エラー
- `404 Not Found`: リソースが見つからない
- `409 Conflict`: リソース競合
- `500 Internal Server Error`: システムエラー

## API仕様

### 共通仕様
- **Base URL**: `/api/v1`
- **認証**: `Authorization: Bearer <JWT_TOKEN>`
- **Content-Type**: `application/json`
- **CORS**: 有効

### レスポンス形式
```json
{
  "data": {...},
  "message": "操作が完了しました"
}
```

### エラーレスポンス形式
```json
{
  "code": "ERROR_CODE",
  "message": "エラーメッセージ",
  "status_code": 400
}
```

## 環境設定

### 環境変数
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - PostgreSQL接続情報
- `COGNITO_USER_POOL_ID`, `COGNITO_CLIENT_ID` - Cognito設定
- `DYNAMO_USER_CONFIG_TABLE` - DynamoDB設定
- `SQS_ROUND_OPTIMIZATION_URL`, `SQS_SESSION_OPTIMIZATION_URL` - SQS設定
- `LOG_LEVEL`, `ENVIRONMENT` - アプリケーション設定

### デプロイ設定
- **Stage**: dev
- **Region**: ap-northeast-1
- **Runtime**: provided.al2 (Go 1.24.3)
- **Profile**: selfpomodoro-dev

## 開発・運用

### ビルド・デプロイ
```bash
# 全関数ビルド
make build-all

# 全関数デプロイ
make deploy-all

# 個別関数デプロイ
make deploy-api-tasks
```

### テスト
```bash
# 全APIテスト実行
make test-all

# PostConfirmationテスト
make test-post-confirmation
```

### ログ確認
```bash
# 各関数のログ確認
make logs-api-tasks
make logs-api-sessions
make logs-api-rounds
```

## 今後の拡張予定

1. **統計機能の強化**: より詳細な分析機能
2. **チーム機能**: 複数ユーザーでのセッション共有
3. **モバイルアプリ**: iOS/Android対応
4. **通知機能**: セッション開始・終了通知
5. **レポート機能**: 週次・月次レポート生成

---

*最終更新: 2025年6月23日*
