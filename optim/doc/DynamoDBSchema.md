# DynamoDB Schema

## round_optimization_logs
| key | type | note |
| --- | --- | --- |
| user_id + time | PK/SK | タイムスタンプ or session_id |
| work_time | int | 次ラウンドの作業時間 |
| break_time | int | 次ラウンドの休憩時間 |
| focus_score | int | 入力スコア |
| timestamp | string | ISO8601 |

## session_optimization_logs
| key | type | note |
| --- | --- | --- |
| user_id + time | PK/SK | タイムスタンプ or session_id |
| round_count | int | 次セッションのラウンド数 |
| break_time | int | 次セッションの長休憩時間 |
| avg_focus_score | float | 平均集中度スコア |
| total_work_time | int | 合計作業時間（分） |
| timestamp | string | ISO8601 |

## user_configs
| column | type | note |
| --- | --- | --- |
| user_id | uuid PK |  |
| round_work_time | int | 現在のデフォルト作業時間 |
| round_break_time | int | 現在のデフォルト休憩時間 |
| session_rounds | int | セッション内のラウンド数 |
| session_break | int | セッション後の長休憩時間 |


# DynamoDB Data Structure
## RoundData
### Structure
```
type UserRoundData = RoundRecord[];

interface RoundRecord {
  user_id: string;
  time: string;
  work_time: number;
  break_time: number;
  focus_score?: number; 
}
```
### Sample
```
[
  {
    'break_time': {'N': '15'},
    'user_id': {'S': '123e4567-e89b-12d3-a456-426614174000'},
    'work_time': {'N': '25'},
    'focus_score': {'N': '85'},
    'time': {'S': '2025-01-01T12:07:55'}
  },
  {
    'break_time': {'N': '8.58383252694626'},
    'user_id': {'S': '123e4567-e89b-12d3-a456-426614174000'},
    'work_time': {'N': '35.94517823265169'},
    'focus_score': {'N': '88'},
    'time': {'S': '2025-06-23T01:33:54.991892'}
  },
  {
    'break_time': {'N': '10.169277459808319'},
    'user_id': {'S': '123e4567-e89b-12d3-a456-426614174000'},
    'work_time': {'N': '35.811200506982274'},
    'time': {'S': '2025-06-23T01:41:08.042795'}
  },
]
```

## SessionData
### Structure
```
type UserSessionSummaryData = SessionSummaryRecord[];

interface SessionSummaryRecord {
  user_id: string;
  time: string;
  round_count: number;
  total_work_time: number;
  break_time: number;
  avg_focus_score?: number; 
}
### Sample
[
  {
    'break_time': {'N': '17.424599264740213'},
    'round_count': {'N': '1'},
    'user_id': {'S': '123e4567-e89b-12d3-a456-426614174000'},
    'total_work_time': {'N': '191.4315558046318'},
    'avg_focus_score': {'N': '77'},
    'time': {'S': '2025-06-09T20:59:43.042898'}
  },
  {
    'break_time': {'N': '15.06805918680749'},
    'round_count': {'N': '6'},
    'user_id': {'S': '123e4567-e89b-12d3-a456-426614174000'},
    'total_work_time': {'N': '412.0773336698032'},
    'avg_focus_score': {'N': '44'},
    'time': {'S': '2025-06-23T01:42:57.842767'}
  },
  {
    'break_time': {'N': '31.743478787386643'},
    'round_count': {'N': '6'},
    'user_id': {'S': '123e4567-e89b-12d3-a456-426614174000'},
    'total_work_time': {'N': '138.39673756892154'},
    'time': {'S': '2025-06-23T01:43:16.915882'}
  }
]
```
