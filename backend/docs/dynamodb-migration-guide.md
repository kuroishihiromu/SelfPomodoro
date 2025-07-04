# DynamoDB統合・移行ガイド

## 概要

Self Pomodoroアプリケーションでは、PostgreSQLからDynamoDBへの段階的移行を実装しました。統計データの高速化とサーバーレス環境での最適化を目的としています。

## 実装済み機能

### 1. ドメイン層の拡張

**事前集約統計モデル:**
- `AggregatedDailyStats`: 日別統計データ
- `AggregatedHourlyStats`: 時間別統計データ
- `DynamoDBSession`, `DynamoDBRound`, `DynamoDBTask`: DynamoDB用データモデル

**ファイル:**
- `/internal/domain/model/aggregated_statistics.go`

### 2. インフラ層の実装

**DynamoDB統計リポジトリ:**
- Single Table Designによる効率的なデータ格納
- 事前集約データによる高速クエリ
- PostgreSQL統計リポジトリと同じインターフェースを実装

**ファイル:**
- `/internal/infrastructure/repository/dynamodb/statistics_repository.go`

### 3. 段階移行対応

**RepositoryFactory拡張:**
- `USE_DYNAMODB_FOR_STATS`環境変数による切り替え
- PostgreSQL/DynamoDB統計リポジトリの自動選択

**設定追加:**
- `/internal/config/config.go`にUseDynamoDBForStatsフラグ追加

### 4. 事前集約処理

**統計集約サービス:**
- ラウンド完了時のリアルタイム統計更新
- 日別・時間別統計の自動集約
- エラー時のグレースフルな処理継続

**ファイル:**
- `/internal/usecase/statistics_aggregation_service.go`
- `/internal/usecase/round_usecase.go` (統計更新統合)
- `/internal/usecase/usecases.go` (依存注入)

## 移行手順

### Phase 1: 開発環境でのテスト

1. **環境変数設定:**
```bash
export USE_DYNAMODB_FOR_STATS=false  # PostgreSQL継続使用
```

2. **DynamoDBテーブル作成:**
```bash
# 開発環境のDynamoDBテーブルを作成
# テーブル名: selfpomodoro_user_configs_dev (既存テーブルを統一テーブルとして使用)
```

3. **アプリケーションデプロイ:**
```bash
make build-all
make deploy-all
```

### Phase 2: DynamoDB統計有効化

1. **環境変数変更:**
```bash
export USE_DYNAMODB_FOR_STATS=true  # DynamoDB統計有効化
```

2. **アプリケーション再デプロイ:**
```bash
make deploy-all
```

3. **動作確認:**
- ラウンド完了時の統計データ更新
- 統計API (`/api/v1/statistics/*`) の動作
- ログでの統計集約サービス有効化確認

### Phase 3: データ移行 (将来実装)

1. **移行スクリプト作成:**
- PostgreSQL既存統計データのDynamoDB移行
- データ整合性チェック

2. **二重書き込み実装:**
- PostgreSQL + DynamoDB同時更新
- データ同期確認

3. **完全移行:**
- PostgreSQL統計機能の削除

## Single Table Design構造

### Primary Key構成
```
PK: USER#{cognito_sub}
SK: エンティティタイプ#識別子
```

### エンティティ例
```
USER#123 | CONFIG                    → ユーザー設定
USER#123 | STATS_DAILY#2025-01-15    → 日別統計
USER#123 | STATS_HOURLY#2025-01-15#14 → 時間別統計
USER#123 | SESSION#2025-01-15#uuid   → セッション
USER#123 | ROUND#2025-01-15#uuid#001 → ラウンド
USER#123 | TASK#uuid                 → タスク
```

## パフォーマンス最適化

### 事前集約のメリット
- **クエリ性能**: 複雑なJOINやGROUP BYが不要
- **レスポンス時間**: サブ秒レベルの応答時間
- **コスト効率**: RCU/WCUの最適化

### 統計更新フロー
```
1. ラウンド完了 (RoundUseCase.CompleteRound)
   ↓
2. データベース永続化 (PostgreSQL/RoundRepository)
   ↓
3. 統計集約サービス実行
   ↓
4. DynamoDB日別・時間別統計更新
   ↓
5. SQS最適化メッセージ送信
```

## 監視・トラブルシューティング

### ログ確認ポイント

**統計集約サービス有効化:**
```
[INFO] DynamoDB統計集約サービスを有効化しました
```

**統計更新成功:**
```
[INFO] 統計データ更新成功: UserID=xxx, Date=2025-01-15, Hour=14, FocusScore=85
```

**統計更新エラー (処理継続):**
```
[WARN] 統計データ更新エラー（処理続行）: DynamoDB operation error
```

### 主要メトリクス

**CloudWatch Logs:**
- Lambda関数ログで統計更新状況を確認
- エラー率とレスポンス時間の監視

**DynamoDB:**
- ConsumedReadCapacityUnits
- ConsumedWriteCapacityUnits
- ThrottledRequests

### トラブルシューティング

**統計データが更新されない:**
1. `USE_DYNAMODB_FOR_STATS=true`の設定確認
2. DynamoDBテーブルの存在確認
3. Lambda実行ロールのDynamoDB権限確認

**パフォーマンス問題:**
1. DynamoDBのRCU/WCU設定確認
2. ホットパーティションの回避
3. Query条件の最適化

## API使用方法

### 統計データ取得 (変更なし)

**集中度トレンド:**
```http
GET /api/v1/statistics/focus-trend?period=week
Authorization: Bearer {jwt_token}
```

**集中度ヒートマップ:**
```http
GET /api/v1/statistics/focus-heatmap?period=month
Authorization: Bearer {jwt_token}
```

### 期待される改善効果

**レスポンス時間:**
- PostgreSQL: 500-1000ms
- DynamoDB: 50-200ms

**スケーラビリティ:**
- PostgreSQL: 同時接続数制限
- DynamoDB: 事実上無制限

**運用コスト:**
- PostgreSQL: 固定インスタンス費用
- DynamoDB: 使用量ベース課金

## 開発・運用ベストプラクティス

### 開発時の注意点

1. **ドメインインターフェース準拠:**
   - `StatisticsRepository`インターフェースの実装
   - 既存APIとの互換性維持

2. **エラーハンドリング:**
   - 統計更新エラー時の処理継続
   - グレースフルなフォールバック

3. **ログ出力:**
   - 統計更新の成功/失敗ログ
   - パフォーマンス監視用メトリクス

### 運用時のベストプラクティス

1. **容量計画:**
   - 日次・週次の使用量監視
   - 自動スケーリング設定

2. **バックアップ:**
   - Point-in-Time Recovery有効化
   - 定期的なデータ整合性チェック

3. **セキュリティ:**
   - 最小権限のIAMロール
   - VPCエンドポイント使用

## 今後の拡張計画

### Phase 4: 全エンティティ移行

1. **タスク管理のDynamoDB移行**
2. **セッション・ラウンドデータ移行**
3. **PostgreSQL完全廃止**

### Phase 5: 高度な機能

1. **DynamoDB Streams活用**
   - リアルタイム分析
   - 外部システム連携

2. **Global Secondary Index追加**
   - 多軸分析対応
   - 管理者向けダッシュボード

3. **Multi-Region対応**
   - 災害復旧対応
   - グローバル展開