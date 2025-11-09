# SelfPomodoro API エンドポイント仕様書

このドキュメントは、SelfPomodoro バックエンドAPIの全エンドポイントの仕様を記載しています。

## 目次

- [認証](#認証)
- [タスク管理API](#タスク管理api)
- [セッション管理API](#セッション管理api)
- [ラウンド管理API](#ラウンド管理api)
- [統計API](#統計api)
- [最適化設定API](#最適化設定api)
- [ヘルスチェック](#ヘルスチェック)
- [エラーレスポンス](#エラーレスポンス)

---

## 認証

全HTTPエンドポイント（ヘルスチェックを除く）で以下の認証が必須です。

**ヘッダー:**
```
Authorization: Bearer {Cognito_JWT_Token}
```

ユーザーIDはAWS Cognito User PoolsのJWTトークンから自動的に抽出されます。

**ベースURL:**
- 開発環境: `https://{api-gateway-id}.execute-api.ap-northeast-1.amazonaws.com/dev`
- 本番環境: `https://{api-gateway-id}.execute-api.ap-northeast-1.amazonaws.com/prod`

---

## タスク管理API

### 1. タスク一覧取得

認証ユーザーの全タスクを取得します。

**エンドポイント:** `GET /api/v1/tasks`

**リクエスト**
- ヘッダー: `Authorization: Bearer {token}`
- ボディ: なし

**レスポンス (200 OK)**
```json
{
  "tasks": [
    {
      "id": "uuid",
      "detail": "タスクの詳細",
      "is_completed": false,
      "created_at": "2025-01-15T10:00:00Z",
      "updated_at": "2025-01-15T10:00:00Z"
    }
  ]
}
```

**実装場所:** `functions/api-tasks/main.go:87-94`

---

### 2. タスク作成

新しいタスクを作成します。

**エンドポイント:** `POST /api/v1/tasks`

**リクエスト**
```json
{
  "detail": "新しいタスク" // 必須
}
```

**レスポンス (201 Created)**
```json
{
  "id": "uuid",
  "detail": "新しいタスク",
  "is_completed": false,
  "created_at": "2025-01-15T10:00:00Z",
  "updated_at": "2025-01-15T10:00:00Z"
}
```

**バリデーション:**
- `detail`: 必須フィールド

**実装場所:** `functions/api-tasks/main.go:97-113`

---

### 3. タスク更新

特定のタスクの詳細情報を更新します。

**エンドポイント:** `PATCH /api/v1/tasks/{task_id}/edit`

**パスパラメータ:**
- `task_id`: タスクのUUID（必須）

**リクエスト**
```json
{
  "detail": "更新されたタスク" // 必須
}
```

**レスポンス (200 OK)**
```json
{
  "id": "uuid",
  "detail": "更新されたタスク",
  "is_completed": false,
  "created_at": "2025-01-15T10:00:00Z",
  "updated_at": "2025-01-15T10:30:00Z"
}
```

**バリデーション:**
- `detail`: 必須フィールド

**実装場所:** `functions/api-tasks/main.go:136-152`

---

### 4. タスク完了状態切り替え

タスクの完了/未完了状態をトグルします。

**エンドポイント:** `PATCH /api/v1/tasks/{task_id}/toggle`

**パスパラメータ:**
- `task_id`: タスクのUUID（必須）

**リクエスト**
- ボディ: なし

**レスポンス (200 OK)**
```json
{
  "id": "uuid",
  "detail": "タスクの詳細",
  "is_completed": true,
  "created_at": "2025-01-15T10:00:00Z",
  "updated_at": "2025-01-15T11:00:00Z"
}
```

**実装場所:** `functions/api-tasks/main.go:127-134`

---

### 5. タスク削除

特定のタスクを削除します。

**エンドポイント:** `DELETE /api/v1/tasks/{task_id}`

**パスパラメータ:**
- `task_id`: タスクのUUID（必須）

**リクエスト**
- ボディ: なし

**レスポンス (200 OK)**
```json
{
  "message": "タスクが削除されました"
}
```

**実装場所:** `functions/api-tasks/main.go:155-173`

---

## セッション管理API

### 6. セッション開始

新しいポモドーロセッションを開始します。

**エンドポイント:** `POST /api/v1/sessions`

**リクエスト**
- ボディ: なし

**レスポンス (201 Created)**
```json
{
  "id": "uuid",
  "start_time": "2025-01-15T10:00:00Z",
  "end_time": null,
  "average_focus": null,
  "total_work_min": null,
  "round_count": null,
  "break_time": null
}
```

**実装場所:** `functions/api-sessions/main.go:94-105`

---

### 7. セッション完了

セッションを完了し、統計情報を計算します。

**エンドポイント:** `PATCH /api/v1/sessions/{session_id}/complete`

**パスパラメータ:**
- `session_id`: セッションのUUID（必須）

**リクエスト**
- ボディ: なし

**レスポンス (200 OK)**
```json
{
  "id": "uuid",
  "start_time": "2025-01-15T10:00:00Z",
  "end_time": "2025-01-15T12:00:00Z",
  "average_focus": 75.5,
  "total_work_min": 100,
  "round_count": 4,
  "break_time": 20
}
```

**処理内容:**
- セッションの終了時刻を記録
- 平均集中度を計算（全ラウンドの集中スコアの平均）
- 総作業時間を計算
- ラウンド数をカウント
- ラウンド数が1以上の場合、SQSメッセージを送信して最適化処理を起動

**実装場所:** `functions/api-sessions/main.go:108-127`

---

## ラウンド管理API

### 8. ラウンド開始

セッション内で新しいラウンド（作業または休憩）を開始します。

**エンドポイント:** `POST /api/v1/sessions/{session_id}/rounds`

**パスパラメータ:**
- `session_id`: セッションのUUID（必須）

**リクエスト**
```json
{}
```

**レスポンス (201 Created)**
```json
{
  "id": "uuid",
  "session_id": "uuid",
  "round_order": 1,
  "start_time": "2025-01-15T10:00:00Z",
  "end_time": null,
  "work_time": null,
  "break_time": null,
  "focus_score": null
}
```

**実装場所:** `functions/api-rounds/main.go:116-128`

---

### 9. ラウンド完了

集中スコアと共にラウンドを完了します。

**エンドポイント:** `PATCH /api/v1/rounds/{round_id}/complete`

**パスパラメータ:**
- `round_id`: ラウンドのUUID（必須）

**リクエスト**
```json
{
  "focus_score": 80,  // 0-100の整数（必須）
  "work_time": 25,    // 作業時間（分）（オプション、min: 1）
  "break_time": 5     // 休憩時間（分）（オプション、min: 1）
}
```

**レスポンス (200 OK)**
```json
{
  "id": "uuid",
  "session_id": "uuid",
  "round_order": 1,
  "start_time": "2025-01-15T10:00:00Z",
  "end_time": "2025-01-15T10:30:00Z",
  "work_time": 25,
  "break_time": 5,
  "focus_score": 80
}
```

**バリデーション:**
- `focus_score`: 必須、0-100の整数
- `work_time`: オプション、1以上の整数
- `break_time`: オプション、1以上の整数

**処理内容:**
- ラウンドの終了時刻を記録
- 集中スコアが指定されている場合、SQSメッセージを送信して最適化処理を起動

**実装場所:** `functions/api-rounds/main.go:131-168`

---

## 統計API

### 10. 集中度トレンド取得

特定の日付から始まる集中度トレンドデータを取得します。

**エンドポイント:** `GET /api/v1/statistics/focus-trend/{date}`

**パスパラメータ:**
- `date`: YYYYMMDD形式の開始日（例: 20250115）

**リクエスト**
- ボディ: なし

**レスポンス (200 OK)**
```json
[
  {
    "date": "2025-01-15",
    "focus_score": 75.5
  },
  {
    "date": "2025-01-16",
    "focus_score": 82.3
  }
]
```

**パラメータ形式:**
- 日付は8桁の数字（YYYYMMDD）
- 例: `20250115` → 2025年1月15日

**実装場所:** `functions/api-statistics/main.go:84-116`

---

### 11. 集中度ヒートマップ取得

集中度のヒートマップデータを取得します。週、月、またはカスタム期間を指定可能です。

**エンドポイント:** `GET /api/v1/statistics/focus-heatmap`

**クエリパラメータ:**
- `period`: 期間の種類（"week" | "month" | "custom"）
- `start_date`: 開始日（YYYY-MM-DD形式、periodがcustomの場合必須）
- `end_date`: 終了日（YYYY-MM-DD形式、periodがcustomの場合必須）

**リクエスト例**
```
GET /api/v1/statistics/focus-heatmap?period=week
GET /api/v1/statistics/focus-heatmap?period=month
GET /api/v1/statistics/focus-heatmap?period=custom&start_date=2025-01-01&end_date=2025-01-15
```

**レスポンス (200 OK)**
```json
[
  {
    "date": "2025-01-15",
    "hour": 10,
    "focus_score": 78.5
  },
  {
    "date": "2025-01-15",
    "hour": 14,
    "focus_score": 85.2
  }
]
```

**期間の指定:**
- `week`: 今週のデータ
- `month`: 今月のデータ
- `custom`: カスタム期間（start_dateとend_dateが必須）

**実装場所:** `functions/api-statistics/main.go:119-157`

---

## 最適化設定API

### 12. 最適化設定取得

ユーザーの最適化設定（作業時間、休憩時間、セッションラウンド数など）を取得します。

**エンドポイント:** `GET /api/v1/optimization-preferences`

**リクエスト**
- ボディ: なし

**レスポンス (200 OK)**
```json
{
  "user_id": "uuid",
  "round_work_time": 25,      // 作業時間（分）
  "round_break_time": 5,      // 休憩時間（分）
  "session_rounds": 4,        // セッションあたりのラウンド数
  "session_break_time": 15,   // セッション休憩時間（分）
  "created_at": "2025-01-15T10:00:00Z",
  "updated_at": "2025-01-15T10:00:00Z"
}
```

**実装場所:** `functions/api-optimization-preferences/main.go:81-90`

---

### 13. 最適化設定更新

ユーザーの最適化設定を更新します。

**エンドポイント:** `PUT /api/v1/optimization-preferences`

**リクエスト**
```json
{
  "round_work_time": 30,      // 1-180分（オプション）
  "round_break_time": 10,     // 1-60分（オプション）
  "session_rounds": 5,        // 1-20ラウンド（オプション）
  "session_break_time": 20    // 1-180分（オプション）
}
```

**レスポンス (200 OK)**
```json
{
  "user_id": "uuid",
  "round_work_time": 30,
  "round_break_time": 10,
  "session_rounds": 5,
  "session_break_time": 20,
  "created_at": "2025-01-15T10:00:00Z",
  "updated_at": "2025-01-15T12:00:00Z"
}
```

**バリデーション:**
- `round_work_time`: 1-180分
- `round_break_time`: 1-60分
- `session_rounds`: 1-20ラウンド
- `session_break_time`: 1-180分
- 全てのフィールドはオプション（指定されたフィールドのみ更新）

**実装場所:** `functions/api-optimization-preferences/main.go:93-115`

---

## ヘルスチェック

### 15. ヘルスチェック

サービスの稼働状態を確認します。

**エンドポイント:** `GET /`

**リクエスト**
- 認証不要
- ボディ: なし

**レスポンス (200 OK)**
```json
{
  "status": "OK",
  "version": "serverless-1.0.0",
  "service": "selfpomodoro-api"
}
```

**実装場所:** `functions/health/main.go:11-30`

---

## エラーレスポンス

全エンドポイント共通のエラーレスポンス形式です。

**エラーレスポンス形式**
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "エラーメッセージ"
  }
}
```

### HTTPステータスコードとエラーコード

#### 400 Bad Request
クライアントからの不正なリクエスト

**エラーコード:**
- `INVALID_REQUEST_FORMAT`: 無効なリクエスト形式
- `VALIDATION_ERROR`: バリデーションエラー
- `MISSING_TASK_ID`: タスクIDが指定されていません
- `MISSING_SESSION_ID`: セッションIDが指定されていません
- `MISSING_DATE`: 日付が指定されていません
- `INVALID_TASK_ID`: 無効なタスクID
- `INVALID_SESSION_ID`: 無効なセッションID
- `INVALID_ROUND_ID`: 無効なラウンドID
- `INVALID_DATE_FORMAT`: 無効な日付形式
- `INVALID_REQUEST`: 無効なリクエストボディ

#### 401 Unauthorized
認証が必要または認証情報が無効

**エラーコード:**
- `UNAUTHORIZED`: 認証が必要です

#### 403 Forbidden
認証済みだがアクセス権限がない

**エラーコード:**
- `FORBIDDEN`: アクセスが拒否されました

#### 404 Not Found
リソースが見つからない

**エラーコード:**
- `NOT_FOUND`: リソースが見つかりません
- `TASK_NOT_FOUND`: タスクが見つかりません
- `SESSION_NOT_FOUND`: セッションが見つかりません
- `ROUND_NOT_FOUND`: ラウンドが見つかりません

#### 405 Method Not Allowed
許可されていないHTTPメソッド

**エラーコード:**
- `METHOD_NOT_ALLOWED`: メソッドが許可されていません

#### 500 Internal Server Error
サーバー内部エラー

**エラーコード:**
- `INTERNAL_ERROR`: サーバー内部エラー

---

## 技術仕様

### インフラストラクチャ
- **ランタイム:** Go (provided.al2)
- **デプロイ:** AWS Lambda + API Gateway (Serverless Framework)
- **データベース:** DynamoDB（unified tableデザイン）
- **メッセージング:** AWS SQS（非同期最適化処理用）
- **リージョン:** ap-northeast-1（東京）

### CORS設定
全エンドポイントで以下のCORSヘッダーが設定されています：

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Headers: Content-Type,Authorization
Access-Control-Allow-Methods: GET,POST,PATCH,DELETE,OPTIONS
```

### 認証方式
- AWS Cognito User Pools
- JWT Bearer Token

### レスポンス形式
- Content-Type: `application/json`
- 文字エンコーディング: UTF-8

---

## 特別な機能

### Cognito PostConfirmationトリガー

**タイプ:** AWS Cognitoトリガー（HTTP APIではない）

**トリガーイベント:** ユーザーがサインアップ確認を完了した際

**処理内容:**
1. データベースにUserレコードを作成
2. デフォルト値でUserConfigを作成
3. 10日分のサンプル最適化データを生成

**実装場所:** `functions/cognito-triggers/post-confirmation/main.go:30-59`

---

## 非同期処理

### SQSメッセージング

以下の操作でSQSメッセージが送信され、非同期で最適化処理が実行されます：

1. **セッション完了時:**
   - 条件: ラウンド数が1以上
   - メッセージ内容: セッション統計データ

2. **ラウンド完了時:**
   - 条件: 集中スコアが指定されている
   - メッセージ内容: ラウンドの作業時間、休憩時間、集中スコア

---

## バージョン履歴

- **v1.0.0** (2025-01-15): 初版リリース
  - タスク管理API
  - セッション管理API
  - ラウンド管理API
  - 統計API
  - 最適化設定API
  - ヘルスチェック

---

## サポート

問題が発生した場合は、以下を確認してください：

1. 認証トークンが有効であることを確認
2. リクエストボディが正しいJSON形式であることを確認
3. 必須パラメータが全て含まれていることを確認
4. エラーレスポンスのエラーコードとメッセージを確認

---

**Last Updated:** 2025-10-16
