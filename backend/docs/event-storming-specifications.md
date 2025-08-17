# イベントストーミング詳細仕様書 - Self Pomodoro

## 📋 イベント詳細仕様

### 0. ユーザーライフサイクル関連イベント

#### 🔥 UserCreated (ユーザーが作成された)
```json
{
  "eventType": "UserCreated",
  "aggregateId": "user-uuid",
  "userId": "user-uuid",
  "timestamp": "2025-01-07T09:00:00Z",
  "data": {
    "userId": "user-uuid",
    "name": "田中太郎",
    "email": "tanaka@example.com",
    "provider": "Cognito_UserPool",
    "registrationSource": "MOBILE_APP",
    "isFirstTimeUser": true
  }
}
```

**発生条件**: 新規ユーザー登録完了
**集約**: User集約
**後続処理**: オンボーディング開始、デフォルト設定作成

#### 🔥 OnboardingCompleted (オンボーディングが完了された)
```json
{
  "eventType": "OnboardingCompleted",
  "aggregateId": "user-uuid",
  "userId": "user-uuid",
  "timestamp": "2025-01-07T09:05:00Z",
  "data": {
    "userId": "user-uuid",
    "completedSteps": [
      "PROFILE_SETUP",
      "PREFERENCES_CONFIG",
      "TUTORIAL_VIEWED"
    ],
    "preferredWorkTime": 25,
    "preferredBreakTime": 5,
    "preferredSessionRounds": 3,
    "skipedSteps": []
  }
}
```

**発生条件**: ユーザーがオンボーディングフローを完了
**集約**: User集約
**後続処理**: サンプルデータ作成、初期設定適用

#### 🔥 SampleDataCreated (サンプルデータが作成された)
```json
{
  "eventType": "SampleDataCreated",
  "aggregateId": "user-uuid",
  "userId": "user-uuid",
  "timestamp": "2025-01-07T09:06:00Z",
  "data": {
    "userId": "user-uuid",
    "createdData": {
      "sampleTasks": [
        {
          "taskId": "sample-task-1",
          "detail": "プレゼンテーション資料作成",
          "isCompleted": false
        },
        {
          "taskId": "sample-task-2", 
          "detail": "メール返信",
          "isCompleted": true
        }
      ],
      "sampleStatistics": {
        "dailyStats": 7,
        "weeklyStats": 1,
        "totalSessions": 12
      },
      "sampleOptimizationRecords": 3
    },
    "dataType": "DEMO_DATA",
    "retentionDays": 30
  }
}
```

**発生条件**: オンボーディング完了後の自動生成
**集約**: 複数集約（Task、Statistics）
**後続処理**: ダッシュボード表示、ユーザー体験向上

### 1. セッション関連イベント

#### 🔥 SessionStarted (セッションが開始された)
```json
{
  "eventType": "SessionStarted",
  "aggregateId": "session-uuid",
  "userId": "user-uuid",
  "timestamp": "2025-01-07T10:00:00Z",
  "data": {
    "sessionId": "session-uuid",
    "plannedRounds": 3,
    "workTime": 25,
    "breakTime": 5,
    "sessionBreakTime": 15
  }
}
```

**発生条件**: ユーザーが「セッション開始」コマンドを実行
**集約**: Session集約
**後続処理**: セッション統計の初期化

#### 🔥 RoundStarted (ラウンドが開始された)
```json
{
  "eventType": "RoundStarted",
  "aggregateId": "session-uuid",
  "userId": "user-uuid",
  "timestamp": "2025-01-07T10:00:30Z",
  "data": {
    "sessionId": "session-uuid",
    "roundId": "round-uuid",
    "roundNumber": 1,
    "plannedWorkTime": 25,
    "plannedBreakTime": 5
  }
}
```

**発生条件**: セッション内で新しいラウンドが開始
**集約**: Session集約 (Round子エンティティ)
**後続処理**: タイマー開始、統計カウンター更新

#### 🔥 RoundCompleted (ラウンドが完了された)
```json
{
  "eventType": "RoundCompleted",
  "aggregateId": "session-uuid",
  "userId": "user-uuid",
  "timestamp": "2025-01-07T10:25:30Z",
  "data": {
    "sessionId": "session-uuid",
    "roundId": "round-uuid",
    "roundNumber": 1,
    "actualWorkTime": 25,
    "actualBreakTime": 5,
    "focusScore": 85,
    "isOptimizationEligible": true
  }
}
```

**発生条件**: ラウンド完了 + 集中度スコア入力
**集約**: Session集約
**後続処理**: 統計更新、最適化判定

#### 🔥 SessionCompleted (セッションが完了された)
```json
{
  "eventType": "SessionCompleted",
  "aggregateId": "session-uuid",
  "userId": "user-uuid",
  "timestamp": "2025-01-07T11:30:00Z",
  "data": {
    "sessionId": "session-uuid",
    "completedRounds": 3,
    "totalWorkTime": 75,
    "totalBreakTime": 10,
    "averageFocusScore": 82.3,
    "sessionQuality": "高品質",
    "optimizationData": {
      "eligibleForOptimization": true,
      "recommendedWorkTime": 27,
      "recommendedBreakTime": 6
    }
  }
}
```

**発生条件**: 全ラウンド完了またはユーザーによる早期終了
**集約**: Session集約
**後続処理**: 最適化実行、統計集計、フィードバック生成

### 2. 最適化関連イベント

#### 🔥 OptimizationPreferencesUpdated (最適化設定が更新された)
```json
{
  "eventType": "OptimizationPreferencesUpdated",
  "aggregateId": "optimization-prefs-uuid",
  "userId": "user-uuid",
  "timestamp": "2025-01-07T12:00:00Z",
  "data": {
    "userId": "user-uuid",
    "previousSettings": {
      "workTime": 25,
      "breakTime": 5,
      "sessionRounds": 3,
      "sessionBreakTime": 15
    },
    "newSettings": {
      "workTime": 27,
      "breakTime": 6,
      "sessionRounds": 3,
      "sessionBreakTime": 15
    },
    "optimizationReason": "集中度向上のため作業時間を延長"
  }
}
```

**発生条件**: 自動最適化またはユーザー手動設定変更
**集約**: OptimizationPreferences集約
**後続処理**: 設定UI更新、次回セッション反映

#### 🔥 OptimizationExecuted (最適化が実行された)
```json
{
  "eventType": "OptimizationExecuted",
  "aggregateId": "optimization-uuid",
  "userId": "user-uuid",
  "timestamp": "2025-01-07T11:30:30Z",
  "data": {
    "userId": "user-uuid",
    "sessionId": "session-uuid",
    "optimizationType": "SESSION_BASED",
    "analysisData": {
      "sessionCount": 15,
      "averageFocusScore": 78.5,
      "recommendationConfidence": 0.85
    },
    "recommendations": {
      "workTime": 27,
      "breakTime": 6,
      "reasoning": "高集中度維持のため2分延長を推奨"
    }
  }
}
```

**発生条件**: SessionCompleted後の自動最適化判定
**集約**: 最適化ドメインサービス
**後続処理**: OptimizationPreferences更新、通知送信

### 3. 統計関連イベント

#### 🔥 StatisticsUpdated (統計データが更新された)
```json
{
  "eventType": "StatisticsUpdated",
  "aggregateId": "statistics-uuid",
  "userId": "user-uuid",
  "timestamp": "2025-01-07T10:25:45Z",
  "data": {
    "userId": "user-uuid",
    "updateType": "ROUND_COMPLETION",
    "updatedStatistics": {
      "daily": {
        "date": "2025-01-07",
        "totalRounds": 8,
        "avgFocusScore": 83.2,
        "totalWorkMin": 200,
        "sessionCount": 3
      },
      "hourly": {
        "date": "2025-01-07",
        "hour": 10,
        "totalRounds": 3,
        "avgFocusScore": 85.0
      }
    }
  }
}
```

**発生条件**: RoundCompleted、SessionCompleted発生時
**集約**: Statistics集約
**後続処理**: ダッシュボード更新、トレンド分析

### 4. タスク関連イベント

#### 🔥 TaskCreated (タスクが作成された)
```json
{
  "eventType": "TaskCreated",
  "aggregateId": "task-uuid",
  "userId": "user-uuid",
  "timestamp": "2025-01-07T09:30:00Z",
  "data": {
    "taskId": "task-uuid",
    "userId": "user-uuid",
    "detail": "プレゼンテーション資料作成",
    "estimatedPomodoros": 3,
    "priority": "HIGH"
  }
}
```

**発生条件**: ユーザーによる新規タスク作成
**集約**: Task集約
**後続処理**: タスクリスト表示更新

#### 🔥 TaskCompleted (タスクが完了された)
```json
{
  "eventType": "TaskCompleted",
  "aggregateId": "task-uuid",
  "userId": "user-uuid",
  "timestamp": "2025-01-07T12:30:00Z",
  "data": {
    "taskId": "task-uuid",
    "userId": "user-uuid",
    "completionTime": 90,
    "actualPomodoros": 2,
    "efficiency": 1.5
  }
}
```

**発生条件**: ユーザーによるタスク完了マーク
**集約**: Task集約
**後続処理**: 生産性統計更新、達成感フィードバック

## 🔄 ポリシー詳細仕様

### 0. オンボーディング関連ポリシー
```typescript
// Policy: 新規ユーザー作成時にオンボーディングを開始する
when(UserCreated).then(async (event) => {
  const { userId, isFirstTimeUser } = event.data;
  
  if (isFirstTimeUser) {
    // デフォルト最適化設定を作成
    await createDefaultOptimizationPreferences(userId);
    
    // オンボーディングフロー開始
    await startOnboardingFlow(userId);
    
    // ウェルカム通知送信
    await sendWelcomeNotification(userId);
  }
});

// Policy: オンボーディング完了時にサンプルデータを作成する
when(OnboardingCompleted).then(async (event) => {
  const { userId, preferredWorkTime, preferredBreakTime } = event.data;
  
  // 個人設定を適用
  await updateOptimizationPreferences(userId, {
    workTime: preferredWorkTime,
    breakTime: preferredBreakTime
  });
  
  // サンプルデータ作成
  const sampleData = await generateSampleData(userId);
  
  emit(SampleDataCreated, {
    userId,
    createdData: sampleData,
    dataType: 'DEMO_DATA'
  });
});
```

### 1. 統計更新ポリシー
```typescript
// Policy: 集中度が記録されたら統計を更新する
when(RoundCompleted).then(async (event) => {
  const { userId, focusScore, actualWorkTime, timestamp } = event.data;
  
  // 日別統計更新
  await updateDailyStatistics(userId, {
    date: formatDate(timestamp),
    focusScore,
    workTime: actualWorkTime
  });
  
  // 時間別統計更新
  await updateHourlyStatistics(userId, {
    date: formatDate(timestamp),
    hour: getHour(timestamp),
    focusScore,
    workTime: actualWorkTime
  });
  
  emit(StatisticsUpdated);
});
```

### 2. 最適化実行ポリシー
```typescript
// Policy: セッション完了時に最適化を実行する
when(SessionCompleted).then(async (event) => {
  const { userId, sessionId, averageFocusScore, optimizationData } = event.data;
  
  if (optimizationData.eligibleForOptimization) {
    const recommendations = await optimizationService.analyze({
      userId,
      sessionData: event.data,
      historicalData: await getHistoricalSessions(userId, 30)
    });
    
    emit(OptimizationExecuted, {
      userId,
      sessionId,
      recommendations
    });
  }
});
```

### 3. フィードバック生成ポリシー
```typescript
// Policy: 高集中度ならポジティブフィードバック
when(RoundCompleted).then(async (event) => {
  const { focusScore } = event.data;
  
  if (focusScore >= 80) {
    await sendNotification({
      type: 'POSITIVE_FEEDBACK',
      message: `素晴らしい集中力です！スコア: ${focusScore}`,
      icon: '🔥'
    });
  }
});

// Policy: 低集中度なら最適化を提案する
when(RoundCompleted).then(async (event) => {
  const { focusScore, userId } = event.data;
  
  if (focusScore < 60) {
    const suggestions = await generateImprovementSuggestions(userId, focusScore);
    
    await sendNotification({
      type: 'IMPROVEMENT_SUGGESTION',
      message: suggestions.primary,
      actions: suggestions.actions
    });
  }
});
```

## 📊 ビュー更新仕様

### 1. リアルタイム更新対象
- **セッション統計表示**: SessionStarted, RoundStarted, RoundCompleted
- **集中度トレンド**: StatisticsUpdated
- **タスクリスト**: TaskCreated, TaskCompleted

### 2. バッチ更新対象
- **集中度ヒートマップ**: 1時間毎のStatisticsUpdated集計
- **パフォーマンス分析**: 日次バッチ処理
- **週次レポート**: 週次バッチ処理

## 🔍 イベント監査とトレーサビリティ

すべてのドメインイベントは以下の情報を含み、完全な監査証跡を提供：

- **イベントID**: 一意識別子
- **集約ID**: イベント発生元の集約
- **ユーザーID**: 操作実行者
- **タイムスタンプ**: 発生時刻（UTC）
- **バージョン**: 集約バージョン
- **因果関係ID**: 関連イベントチェーン
- **メタデータ**: 追加のコンテキスト情報

この仕様により、Self Pomodoroアプリケーションの全ビジネスロジックがイベント駆動アーキテクチャとして明確に定義され、DDD設計との整合性が保たれています。