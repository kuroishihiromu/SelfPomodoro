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
