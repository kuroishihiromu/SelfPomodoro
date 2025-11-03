# API Schema

## Round Optimization API

### Endpoint
```
POST /optimize/round/{user_id}
```

### Path Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| user_id   | UUID | Yes      | ユーザーの一意識別子 |

### Request Body Structure
```typescript
type RoundRequestData = RoundRecord[];

interface RoundRecord {
  time: string;           // ISO 8601 format
  work_time: number;      // 作業時間 （ 15-30分 ）
  break_time: number;     // 休憩時間 （ 3-20分 ）
  focus_score: number;    // 集中度スコア （ 0-100点 ）
}
```

### Sample Request
```
POST /optimize/round/123e4567-e89b-12d3-a456-426614174000

Content-Type: application/json

Request Body:
[
    {
        "break_time": 15,
        "work_time": 25,
        "focus_score": 85,
        "time": "2025-01-01T12:07:55"
    },
    {
        "break_time": 8.58383252694626,
        "work_time": 35.94517823265169,
        "focus_score": 88,
        "time": "2025-06-23T01:33:54.991892"
    },
    {
        "break_time": 10.169277459808319,
        "work_time": 35.811200506982274,
        "focus_score": 92,
        "time": "2025-06-23T01:41:08.042795"
    }
]
```

## Session Optimization API

### Endpoint
```
POST /optimize/session/{user_id}
```

### Path Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| user_id   | UUID | Yes      | ユーザーの一意識別子 |

### Request Body Structure
```typescript
type SessionRequestData = SessionRecord[];

interface SessionRecord {
  time: string;              // ISO 8601 format
  round_count: number;       // ラウンド数 ( 1-6回 )
  total_work_time: number;   // 総作業時間 （ 15-360分） 
  break_time: number;        // 休憩時間 （ 5-60分 ）
  avg_focus_score: number;   // 平均集中度スコア （ 0-100点 ）
}
```

### Sample Request
```
POST /optimize/session/123e4567-e89b-12d3-a456-426614174000

Content-Type: application/json

Request Body:
[
    {
        "break_time": 17.424599264740213,
        "round_count": 1,
        "total_work_time": 191.4315558046318,
        "avg_focus_score": 77,
        "time": "2025-06-09T20:59:43.042898"
    },
    {
        "break_time": 15.06805918680749,
        "round_count": 6,
        "total_work_time": 312.0773336698032,
        "avg_focus_score": 44,
        "time": "2025-06-23T01:42:57.842767"
    },
    {
        "break_time": 31.743478787386643,
        "round_count": 6,
        "total_work_time": 138.39673756892154,
        "avg_focus_score": 33,
        "time": "2025-06-23T01:43:16.915882"
    }
]
```
