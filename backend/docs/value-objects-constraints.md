# Value Objects 制約・ビジネスルール一覧

## 概要

このドキュメントは、SelfPomodoroプロジェクトで実装されている全Value Objectsの制約（バリデーション、範囲制限、ビジネスルール）を詳細にまとめたものです。

---

## 1. ユーザー関連（User）

### UserID
**パッケージ**: `internal/domain/valueobject/user/user_id.go`

**制約・バリデーション**:
- UUID形式必須
- 無効なUUID形式はエラー

**ビジネスルール**:
- `IsEmpty()`: 空ID判定（uuid.Nil）
- `Value()`: UUID値を取得
- `String()`: 文字列表現を取得

---

### EmailAddress
**パッケージ**: `internal/domain/valueobject/user/email_address.go`

**制約・バリデーション**:
- **必須項目**: 空文字列不可
- **正規表現**: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
- **最大長**: 254文字
- **前後の空白**: 自動除去
- **大文字変換**: 自動的に小文字に変換

**ビジネスルール**:
- `Domain()`: ドメイン部分を取得
- `LocalPart()`: ローカル部分を取得
- `IsGmailAddress()`: Gmail判定
- `IsEmpty()`: 空判定

---

### UserName
**パッケージ**: `internal/domain/valueobject/user/user_name.go`

**制約・バリデーション**:
- **必須項目**: 空文字列不可
- **文字数**: 1-50文字（UTF-8文字数）
- **制御文字禁止**: タブ、改行、復帰以外の制御文字
- **前後の空白**: 自動除去

**ビジネスルール**:
- `IsShort()`: 5文字以下判定
- `IsLong()`: 20文字以上判定
- `ContainsEmail()`: メールアドレス様文字列判定
- `IsValid()`: 有効性検証

---

### Provider
**パッケージ**: `internal/domain/valueobject/user/provider.go`

**制約・バリデーション**:
- **必須項目**: 空文字列不可
- **有効値**: `Cognito_UserPool`, `Google`

**ビジネスルール**:
- `IsCognito()`, `IsGoogle()`: プロバイダー種別判定
- `IsThirdParty()`: サードパーティー判定（現在はGoogleのみ）
- `CanChangeEmail()`: メール変更可否（サードパーティーは不可）
- `RequiresProviderID()`: プロバイダーID必要性判定
- `GetAuthenticationMethod()`: 認証方式取得（OAuth 2.0/Username-Password）
- `GetDisplayName()`: 表示名取得

---

## 2. タスク関連（Task）

### TaskID
**パッケージ**: `internal/domain/valueobject/task/task_id.go`

**制約・バリデーション**:
- UUID形式必須

**ビジネスルール**:
- `IsNil()`: 空ID判定
- `Value()`: UUID値を取得

---

### TaskDetail
**パッケージ**: `internal/domain/valueobject/task/task_detail.go`

**制約・バリデーション**:
- **文字数**: 1-500文字（UTF-8文字数）
- **空文字列不可**: 必須項目
- **改行文字正規化**: CRLF→LF
- **連続改行制限**: 最大2個まで

**定数**:
```go
MinDetailLength = 1
MaxDetailLength = 500
MaxConsecutiveNewlines = 2
```

**ビジネスルール**:
- `IsShort()`: 10文字以下判定
- `IsLong()`: 100文字以上判定
- `HasMultipleLines()`: 複数行判定
- `GetPreview(maxLength)`: プレビュー文字列生成
- `ContainsKeywords()`: キーワード検索（大小文字不問）

---

### TaskStatus
**パッケージ**: `internal/domain/valueobject/task/task_status.go`

**制約・バリデーション**:
- **有効値**: `StatusIncomplete`(0), `StatusCompleted`(1)
- **文字列変換**: 
  - 未完了: `incomplete`, `false`, `0`
  - 完了: `completed`, `true`, `1`

**ビジネスルール**:
- `Toggle()`: 状態切り替え
- `CanTransitionTo()`: 遷移可能性判定（現在は全ての遷移を許可）
- `Compare()`: 順序比較（未完了 < 完了）
- `IsCompleted()`, `IsIncomplete()`: 状態判定

---

## 3. 統計関連（Statistics）

### Date
**パッケージ**: `internal/domain/valueobject/statistics/date.go`

**制約・バリデーション**:
- **フォーマット**: YYYY-MM-DD（ISO 8601）
- `time.Parse()`で検証

**ビジネスルール**:
- 各種日付判定: `IsToday()`, `IsYesterday()`, `IsFuture()`, `IsPast()`
- 曜日判定: `IsWeekend()`, `IsWeekday()`
- 時間帯判定: `IsEarlyMorning()`, `IsMorning()`, `IsAfternoon()`, `IsEvening()`, `IsNight()`
- 日付計算: `AddDays()`, `FirstDayOfMonth()`, `LastDayOfMonth()`

---

### Hour
**パッケージ**: `internal/domain/valueobject/statistics/hour.go`

**制約・バリデーション**:
- **範囲**: 0-23

**ビジネスルール**:
- 時間帯判定:
  - 早朝: `IsEarlyMorning()` (5-8時)
  - 午前: `IsMorning()` (6-11時)
  - 午後: `IsAfternoon()` (12-17時)
  - 夕方: `IsEvening()` (18-21時)
  - 夜間: `IsNight()` (22-5時)
  - 営業時間: `IsBusinessHour()` (9-17時)
  - 深夜: `IsLateNight()` (22-4時)

---

### AverageFocusScore
**パッケージ**: `internal/domain/valueobject/statistics/average_focus_score.go`

**制約・バリデーション**:
- **範囲**: 0.0-100.0
- **浮動小数点精度**: 0.1

**ビジネスルール**:
- 品質レベル判定:
  - 優秀: `IsExcellent()` (90以上)
  - 高: `IsHigh()` (80-89.9)
  - 中: `IsMedium()` (60-79.9)
  - 低: `IsLow()` (40-59.9)
  - 要改善: `NeedsImprovement()` (40未満)
- `IsProductiveDay()`: 70以上
- `IsProductiveWeek()`: 75以上
- 重み付き平均計算: `UpdateAverage()`

---

### TotalRounds
**パッケージ**: `internal/domain/valueobject/statistics/total_rounds.go`

**制約・バリデーション**:
- **範囲**: 0-1000

**ビジネスルール**:
- 品質レベル:
  - 未開始: `IsEmpty()` (0)
  - 軽量: `IsLight()` (1-3)
  - 標準: `IsStandard()` (4-10)
  - 集中: `IsIntensive()` (11-20)
  - 超集中: `IsSuperIntensive()` (21以上)
- `IsProductiveDay()`: 10ラウンド以上
- `IsProductiveWeek()`: 50ラウンド以上

---

### SessionCount
**パッケージ**: `internal/domain/valueobject/statistics/session_count.go`

**制約・バリデーション**:
- **最小値**: 0（上限なし）

**ビジネスルール**:
- アクティビティレベル:
  - 非アクティブ: `IsInactive()` (0)
  - 軽微: `IsMinimal()` (1)
  - 普通: `IsNormal()` (2-3)
  - アクティブ: `IsActive()` (4-7)
  - 超アクティブ: `IsSuperActive()` (8以上)
- `IsActiveDay()`: 2セッション以上
- `IsProductiveDay()`: 5セッション以上

---

### DaysActive
**パッケージ**: `internal/domain/valueobject/statistics/days_active.go`

**制約・バリデーション**:
- **範囲**: 0-7日

**ビジネスルール**:
- アクティビティレベル:
  - 非アクティブ: `IsInactive()` (0日)
  - 低活動: `IsLowActivity()` (1-2日)
  - 中活動: `IsMediumActivity()` (3-4日)
  - 高活動: `IsHighActivity()` (5-6日)
  - 毎日活動: `IsVeryActive()` (7日)
- `IsConsistent()`: 5日以上

---

### TotalWorkMinutes / TotalBreakMinutes
**パッケージ**: `internal/domain/valueobject/statistics/total_work_minutes.go`

**制約・バリデーション**:
- **最小値**: 0

**ビジネスルール**:
- `IsProductiveDay()`: 120分（2時間）以上
- `IsProductiveWeek()`: 600分（10時間）以上
- `IsIntensiveSession()`: 240分（4時間）以上
- `ToHours()`: 時間単位変換

---

### WeekPeriod
**パッケージ**: `internal/domain/valueobject/statistics/week_period.go`

**制約・バリデーション**:
- **開始日**: 月曜日必須
- **終了日**: 日曜日必須
- **期間**: 7日以内
- **開始日 ≤ 終了日**

**ビジネスルール**:
- 週境界計算（月曜日開始）
- 平日・週末日付取得: `GetWeekdays()`, `GetWeekends()`
- 前週・次週計算: `PreviousWeek()`, `NextWeek()`

---

## 4. ラウンド関連（Round）

### RoundID
**パッケージ**: `internal/domain/valueobject/round/round_id.go`

**制約・バリデーション**:
- UUID形式必須

---

### WorkTime
**パッケージ**: `internal/domain/valueobject/round/work_time.go`

**制約・バリデーション**:
- **範囲**: 1-120分
- **デフォルト値**: 25分

**ビジネスルール**:
- `IsStandard()`: 25分判定
- `IsShort()`: 15分以下
- `IsLong()`: 45分以上
- `Duration()`: time.Duration変換

---

### BreakTime
**パッケージ**: `internal/domain/valueobject/round/break_time.go`

**制約・バリデーション**:
- **範囲**: 0-60分
- **デフォルト値**: 5分

**ビジネスルール**:
- `IsStandard()`: 5分判定
- `IsShort()`: 3分以下
- `IsLong()`: 10分以上
- `IsNoBreak()`: 0分判定

---

### FocusScore
**パッケージ**: `internal/domain/valueobject/round/focus_score.go`

**制約・バリデーション**:
- **範囲**: 0-100
- **整数値**

**ビジネスルール**:
- `IsHigh()`: 80以上
- `IsMedium()`: 60-79
- `IsLow()`: 59以下

---

### RoundOrder
**パッケージ**: `internal/domain/valueobject/round/round_order.go`

**制約・バリデーション**:
- **範囲**: 1-50
- **最小値**: 1

**ビジネスルール**:
- `IsFirst()`: 最初のラウンド判定
- `CanBeNextOf()`: 連続性検証
- `IsConsecutiveWith()`: 連続順序判定
- `DistanceFrom()`: 順序間距離計算

---

## 5. セッション関連（Session）

### SessionID
**パッケージ**: `internal/domain/valueobject/session/session_id.go`

**制約・バリデーション**:
- UUID形式必須

---

### RoundCount
**パッケージ**: `internal/domain/valueobject/session/round_count.go`

**制約・バリデーション**:
- **範囲**: 0-50

**ビジネスルール**:
- セッション品質:
  - 未開始: `IsEmpty()` (0)
  - 短時間: `IsShort()` (1)
  - 標準: `IsStandard()` (2-4)
  - 集中: `IsIntensive()` (5-8)
  - 超集中: `IsSuperIntensive()` (9以上)
- `IsProductiveSession()`: 3ラウンド以上
- `IsIntensiveSession()`: 6ラウンド以上

---

### SessionRounds（計画値）
**パッケージ**: `internal/domain/valueobject/session/session_rounds.go`

**制約・バリデーション**:
- **範囲**: 1-10
- **デフォルト値**: 3
- **最適範囲**: 3-5

**ビジネスルール**:
- 集中度レベル:
  - 軽度: `IsShort()` (1-2)
  - 標準: `IsOptimal()` (3-5)
  - 集中: `IsLong()` (6以上)
- `EstimatedDurationMinutes()`: 推定時間計算（30分/ラウンド）
- `GetRecommendedBreakMinutes()`: 推奨休憩時間計算（ラウンド数に応じて10-30分）
- `GetIntensityLevel()`: 集中度レベル取得

---

### AverageFocus
**パッケージ**: `internal/domain/valueobject/session/average_focus.go`

**制約・バリデーション**:
- **範囲**: 0.0-100.0
- **浮動小数点**

**ビジネスルール**:
- `IsHigh()`: 80以上
- `IsMedium()`: 60-79.9
- `IsLow()`: 60未満
- `IsProductiveSession()`: 70以上

---

### TotalWorkMin
**パッケージ**: `internal/domain/valueobject/session/total_work_minutes.go`

**制約・バリデーション**:
- **範囲**: 0-1440分（24時間）

**ビジネスルール**:
- `IsProductiveSession()`: 60分以上
- `IsLongSession()`: 120分以上

---

## 6. 最適化関連（Optimization）

### RoundOptimizationRecord
**パッケージ**: `internal/domain/valueobject/optimization/round_optimization_record.go`

**制約・バリデーション**:
- 作業時間: 1分以上
- 休憩時間: 0分以上

**ビジネスルール**:
- `IsEffectiveOptimization()`: 集中度60以上、作業時間15-60分、休憩時間15分以下
- 品質評価: 集中度に基づく3段階評価

---

### SessionOptimizationRecord
**パッケージ**: `internal/domain/valueobject/optimization/session_optimization_record.go`

**制約・バリデーション**:
- ラウンド数: 1-10
- 平均集中度: 0-100
- 合計作業時間: 1分以上

**ビジネスルール**:
- `IsEffectiveSession()`: 平均集中度60以上、作業時間30分以上、2ラウンド以上
- セッション効率計算: 分あたり集中度

---

## 依存関係・相互制約

### ID系Value Objects
全てUUID形式で統一:
- UserID, TaskID, SessionID, RoundID

### 時間系Value Objects
分単位で統一、Duration()メソッドでtime.Duration変換:
- WorkTime, BreakTime, TotalWorkMinutes, TotalBreakMinutes

### スコア系Value Objects
0-100範囲で統一、品質レベル判定ロジック共通化:
- FocusScore, AverageFocusScore, AverageFocus

### 計画 vs 実績
- **SessionRounds**: セッション開始前の計画値（1-10）
- **RoundCount**: セッション完了後の実績値（0-50）

### 最適化記録
他のValue Objectsを組み合わせて複合的な制約を実装

---

## まとめ

本プロジェクトでは、型安全性とビジネスルールの封じ込めを徹底するため、プリミティブ型を避けて専用のValue Objectsを実装しています。各Value Objectは以下の原則に従って設計されています：

1. **不変性**: 全てのフィールドがprivate
2. **バリデーション**: コンストラクタでの厳密な検証
3. **ビジネスロジック**: ドメイン固有の判定・計算メソッド
4. **型安全性**: プリミティブ型の混同を完全防止
5. **一貫性**: 同様の概念に対する統一された設計パターン

これにより、コンパイル時の型チェックとランタイムでのバリデーションを組み合わせた堅牢なドメインモデルが実現されています。

---

**作成日**: 2025-01-10  
**対象**: SelfPomodoro Backend Value Objects  
**バージョン**: 全Value Objects実装完了版