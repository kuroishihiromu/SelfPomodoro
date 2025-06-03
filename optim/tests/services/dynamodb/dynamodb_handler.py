from app.services.dynamodb.dynamodb_handler import DynamoDBHandler

def test_round_dynamodb_handler():
    """DynamoDBのラウンドデータのテスト"""
    try:
        print("[テスト] DynamoDBのラウンドデータのテスト")
        print()

        # インスタンス作成
        dynamodb_handler = DynamoDBHandler(table_name='round_optimization_logs', region_name='ap-northeast-1')
        
        # テーブル名の確認
        assert dynamodb_handler.table_name == 'round_optimization_logs'

        # 作業時間と小休憩時間の追加
        dynamodb_handler.put_round_data(
            user_id='testuser-abcd-1234-5678-dynamodbtest',
            time='2023-01-01T10:00:00Z',
            work_time=20,
            break_time=6,
            timestamp='2023-01-01T10:00:00.000Z'
        )
        
        # ラウンドデータの取得
        round_data = dynamodb_handler.get_round_data(
            user_id='testuser-abcd-1234-5678-dynamodbtest',
            time='2023-01-01T10:00:00Z'
        )
        
        # 集中度スコアのリスト化
        focus_score_list = dynamodb_handler.make_chosen_data_list(
            user_id='testuser-abcd-1234-5678-dynamodbtest',
            columns=['focus_score']
        )

        # 作業時間と小休憩時間のリスト化
        work_time_break_time_list = dynamodb_handler.make_chosen_data_list(
            user_id='testuser-abcd-1234-5678-dynamodbtest',
            columns=['work_time', 'break_time']
        )

        print("ラウンドデータの取得:\n", round_data)
        assert isinstance(round_data, dict)
        print()
        print("集中度スコアのリスト化:\n", focus_score_list)
        print()
        print("作業時間と小休憩時間のリスト化:\n", work_time_break_time_list)
        assert isinstance(focus_score_list, list)
        assert isinstance(work_time_break_time_list, list)
        print()
        print()

        print("[成功]: DynamoDBのラウンドデータのテスト完了")
        print()
    
    except ValueError as e:
        print(f"[失敗]: DynamoDBのラウンドデータのテストに失敗しました: {e}")
        print()
        raise
    except Exception as e:
        print(f"[失敗]: DynamoDBのラウンドデータのテストに失敗しました: {e}")
        raise


def test_session_dynamodb_handler():
    """セッションデータのテスト"""
    try:
        print("[テスト] DynamoDBのセッションデータのテスト")
        print()

        # インスタンス作成
        dynamodb_handler = DynamoDBHandler(table_name='session_optimization_logs', region_name='ap-northeast-1')
        
        # テーブル名の確認
        assert dynamodb_handler.table_name == 'session_optimization_logs'

        # ラウンド数と長休憩時間の追加
        dynamodb_handler.put_session_data(
            user_id='testuser-abcd-1234-5678-dynamodbtest',
            time='2023-01-01T10:00:00Z',
            round_count=2,
            break_time=5,
            total_work_time=50,
            timestamp='2023-01-01T10:00:00.000Z'
        )
        
        # セッションデータの取得
        session_data = dynamodb_handler.get_session_data(
            user_id='testuser-abcd-1234-5678-dynamodbtest'
        )
        
        # 平均集中度スコアのリスト化
        avg_focus_score_list = dynamodb_handler.make_chosen_session_data_list(
            user_id='testuser-abcd-1234-5678-dynamodbtest',
            columns=['avg_focus_score']
        )

        # ラウンド数と長休憩時間と作業時間のリスト化
        round_count_break_time_list = dynamodb_handler.make_chosen_session_data_list(
            user_id='testuser-abcd-1234-5678-dynamodbtest',
            columns=['round_count', 'break_time', 'total_work_time']
        )
        
        print("セッションデータの取得:\n", session_data)
        assert isinstance(session_data, list)
        print()
        print("平均集中度スコアのリスト化:\n", avg_focus_score_list)
        print()
        print("ラウンド数と長休憩時間と作業時間のリスト化:\n", round_count_break_time_list)
        assert isinstance(avg_focus_score_list, list)
        assert isinstance(round_count_break_time_list, list)
        print()
        print()
        

        print("[成功]: DynamoDBのセッションデータのテスト完了")
        print()

    except Exception as e:
        print(f"セッションデータのリスト化に失敗しました: {e}")
        raise

    
if __name__ == "__main__":
    test_round_dynamodb_handler()
    test_session_dynamodb_handler()
