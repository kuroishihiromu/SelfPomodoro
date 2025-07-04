# DynamoDB テーブル設計書（AWSコンソール用）

## Single Table Design概要

Self PomodoroアプリケーションのDynamoDB統合では、Single Table Designを採用し、すべてのデータエンティティを1つのテーブルに格納します。

## AWSコンソールでのテーブル作成

### 基本設定
- **テーブル名**: `selfpomodoro_unified_table_dev`
- **パーティションキー**: `PK` (String)
- **ソートキー**: `SK` (String)
- **設定**: デフォルト設定を使用

## エンティティ設計

### 1. ユーザー設定 (UserConfig)
```
PK: USER#{cognito_sub}
SK: CONFIG
```

**属性:**
- `user_id`: string (Cognito sub)
- `round_work_time`: number
- `round_break_time`: number
- `session_rounds`: number
- `session_break_time`: number
- `created_at`: string (ISO 8601)
- `updated_at`: string (ISO 8601)

### 2. タスク (Task)
```
PK: USER#{cognito_sub}
SK: TASK#{task_id}
```

**属性:**
- `user_id`: string
- `task_id`: string (UUID)
- `detail`: string
- `is_completed`: boolean
- `created_at`: string (ISO 8601)
- `updated_at`: string (ISO 8601)

### 3. セッション (Session)
```
PK: USER#{cognito_sub}
SK: SESSION#{date}#{session_id}
```

**属性:**
- `user_id`: string
- `session_id`: string (UUID)
- `date`: string (YYYY-MM-DD)
- `start_time`: string (ISO 8601)
- `end_time`: string (ISO 8601) - nullable
- `average_focus`: number - nullable
- `total_work_min`: number - nullable
- `round_count`: number - nullable
- `break_time`: number - nullable
- `created_at`: string (ISO 8601)
- `updated_at`: string (ISO 8601)

### 4. ラウンド (Round)
```
PK: USER#{cognito_sub}
SK: ROUND#{date}#{session_id}#{round_order}
```

**属性:**
- `user_id`: string
- `round_id`: string (UUID)
- `session_id`: string (UUID)
- `date`: string (YYYY-MM-DD)
- `round_order`: number (3桁ゼロパディング)
- `start_time`: string (ISO 8601)
- `end_time`: string (ISO 8601) - nullable
- `work_time`: number - nullable
- `break_time`: number - nullable
- `focus_score`: number - nullable
- `is_aborted`: boolean
- `created_at`: string (ISO 8601)
- `updated_at`: string (ISO 8601)

### 5. 日別統計 (Aggregated Daily Stats)
```
PK: USER#{cognito_sub}
SK: STATS_DAILY#{date}
```

**属性:**
- `user_id`: string
- `date`: string (YYYY-MM-DD)
- `total_rounds`: number
- `avg_focus_score`: number
- `total_work_min`: number
- `total_break_min`: number
- `session_count`: number
- `updated_at`: string (ISO 8601)

### 6. 時間別統計 (Aggregated Hourly Stats)
```
PK: USER#{cognito_sub}
SK: STATS_HOURLY#{date}#{hour}
```

**属性:**
- `user_id`: string
- `date`: string (YYYY-MM-DD)
- `hour`: number (2桁ゼロパディング)
- `total_rounds`: number
- `avg_focus_score`: number
- `total_work_min`: number
- `total_break_min`: number
- `updated_at`: string (ISO 8601)

### 7. 週別統計 (Aggregated Weekly Stats)
```
PK: USER#{cognito_sub}
SK: STATS_WEEKLY#{week_start}
```

**属性:**
- `user_id`: string
- `week_start`: string (YYYY-MM-DD, Monday)
- `week_end`: string (YYYY-MM-DD, Sunday)
- `total_rounds`: number
- `avg_focus_score`: number
- `total_work_min`: number
- `total_break_min`: number
- `session_count`: number
- `days_active`: number
- `updated_at`: string (ISO 8601)

### 8. 最適化ログ (Round Optimization)
```
PK: USER#{cognito_sub}
SK: OPTIMIZATION_ROUND#{timestamp}
```

### 9. 最適化ログ (Session Optimization)
```
PK: USER#{cognito_sub}
SK: OPTIMIZATION_SESSION#{timestamp}
```

## インデックス設計

### Global Secondary Index (GSI)

#### GSI1: 日付別クエリ用
- **Partition Key**: `date`
- **Sort Key**: `user_id`
- **用途**: 特定日のすべてのユーザーのセッション・ラウンドデータ取得

#### GSI2: エンティティタイプ別クエリ用
- **Partition Key**: `entity_type` (SKから抽出)
- **Sort Key**: `created_at`
- **用途**: 特定タイプのエンティティの一覧取得

## クエリパターン

### 1. ユーザー設定取得
```
Query: PK = USER#{user_id} AND SK = CONFIG
```

### 2. ユーザーのタスク一覧
```
Query: PK = USER#{user_id} AND begins_with(SK, "TASK#")
```

### 3. 期間内セッション一覧
```
Query: PK = USER#{user_id} AND SK BETWEEN "SESSION#{start_date}" AND "SESSION#{end_date}#zzz"
```

### 4. 期間内ラウンド一覧
```
Query: PK = USER#{user_id} AND SK BETWEEN "ROUND#{start_date}" AND "ROUND#{end_date}#zzz"
```

### 5. 日別統計取得
```
Query: PK = USER#{user_id} AND SK BETWEEN "STATS_DAILY#{start_date}" AND "STATS_DAILY#{end_date}"
```

### 6. 時間別統計取得
```
Query: PK = USER#{user_id} AND SK BETWEEN "STATS_HOURLY#{start_date}#00" AND "STATS_HOURLY#{end_date}#23"
```

### 7. 週別統計取得
```
Query: PK = USER#{user_id} AND SK BETWEEN "STATS_WEEKLY#{start_week}" AND "STATS_WEEKLY#{end_week}"
```

## 容量設計

### 読み込み容量 (RCU)
- **統計クエリ**: 1日あたり約100回 × 期間内日数
- **セッション/ラウンドクエリ**: 1日あたり約50回
- **設定読み込み**: 1日あたり約30回

**推定RCU**: 20-30 RCU (開発環境)

### 書き込み容量 (WCU)
- **ラウンド完了**: 1日あたり約20回
- **統計更新**: 1日あたり約20回 (日別・時間別)
- **セッション更新**: 1日あたり約5回

**推定WCU**: 10-15 WCU (開発環境)

## 移行戦略

### Phase 1: 統計データのDynamoDB移行
1. DynamoDB統計リポジトリの実装
2. 事前集約処理の実装
3. 既存PostgreSQLデータの移行スクリプト作成

### Phase 2: 二重書き込み
1. PostgreSQL + DynamoDB同時書き込み
2. データ整合性の検証
3. 性能監視

### Phase 3: 読み取り切り替え
1. 統計クエリをDynamoDBに切り替え
2. PostgreSQL統計データの読み取り停止

### Phase 4: PostgreSQL統計機能削除
1. PostgreSQL統計関連テーブル・機能の削除
2. DynamoDBへの完全移行完了

### Phase 5: 全エンティティ移行
1. タスク、セッション、ラウンドのDynamoDB移行
2. PostgreSQL完全廃止

## 監視・運用

### CloudWatch メトリクス
- `ConsumedReadCapacityUnits`
- `ConsumedWriteCapacityUnits`
- `ThrottledRequests`
- `SuccessfulRequestLatency`

### アラート設定
- RCU/WCU使用率 > 80%
- スロットリング発生
- レイテンシ > 100ms

## セキュリティ

### IAM ポリシー
- Lambda実行ロールに最小権限の付与
- テーブル単位でのアクセス制御
- VPC エンドポイントの使用推奨

### 暗号化
- 保存時暗号化: AWS KMS
- 転送中暗号化: TLS 1.2+

## バックアップ・災害復旧

### Point-in-Time Recovery (PITR)
- 35日間の自動バックアップ
- 秒単位での復旧ポイント

### クロスリージョンレプリケーション
- プロダクション環境での災害復旧対応
- 複数AZでの冗長化

## AWSコンソールでの設定手順

### 1. テーブル作成
```
1. DynamoDBコンソールを開く
2. 「テーブルの作成」をクリック
3. テーブル名: selfpomodoro_unified_table_dev
4. パーティションキー: PK (文字列)
5. ソートキー: SK (文字列)
6. 「テーブルの作成」をクリック
```

### 2. GSI1の作成（メールアドレス検索用）
```
1. 作成したテーブルを選択
2. 「インデックス」タブをクリック
3. 「インデックスの作成」をクリック
4. インデックス名: GSI1
5. パーティションキー: email (文字列)
6. 射影: すべての属性
7. 「インデックスの作成」をクリック
```

### 3. GSI2の作成（プロバイダー別検索用）
```
1. 「インデックスの作成」をクリック
2. インデックス名: GSI2
3. パーティションキー: provider (文字列)
4. ソートキー: created_at (文字列)
5. 射影: すべての属性
6. 「インデックスの作成」をクリック
```

### 4. 環境変数の設定
アプリケーションで新しいテーブルを使用するため、以下の環境変数を設定：
```bash
export USE_DYNAMODB_FOR_STATS=true
export DYNAMO_USER_CONFIG_TABLE=selfpomodoro_unified_table_dev
```