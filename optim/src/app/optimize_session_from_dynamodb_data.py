import uuid
from handler.dynamodb.dynamodb_handler import DynamoDBHandler
from optimizer.bayesian_optimizer import BayesianOptimizer
from datetime import datetime

def optimize_session_from_dynamodb_data(
    user_id: uuid.UUID,
    avg_focus_score: float
):
    """DynamoDBのセッション最適化

    Parameters:
        user_id (uuid.UUID): ユーザーID
        avg_focus_score (float): 集中度スコアの平均

    Returns:
        total_work_time (float): 1セッションの作業時間の合計
        break_time (float): 最適なセッション間休憩時間
        round_count (int): 最適なラウンド繰り返し回数
    """
    dynamodb_handler = DynamoDBHandler(table_name="session_optimization_logs", region_name="ap-northeast-1")
    
    # --- 最新のセッションデータを取得 ---
    latest_data = dynamodb_handler.get_session_data(user_id=str(user_id))
    print("最新のセッションデータ取得完了: ", latest_data)
    latest_time = datetime.now().isoformat()
    
    if latest_data and isinstance(latest_data, list) and len(latest_data) > 0:
        # --- 最新のデータを取得 ---
        converted_data = dynamodb_handler._convert_to_list(latest_data)
        latest_item = converted_data[-1]
        latest_time = latest_item['time']
        total_work_time = latest_item.get('total_work_time')
        break_time = latest_item.get('break_time')
        round_count = latest_item.get('round_count')
    else:
        total_work_time = None
        break_time = None
        round_count = None
    
    # --- 最新のデータに集中度スコアを追加して更新 ---
    dynamodb_handler.put_session_data(
        user_id=str(user_id),
        time=latest_time,
        total_work_time=total_work_time,
        break_time=break_time,
        round_count=round_count,
        avg_focus_score=avg_focus_score
    )
    
    # --- 説明変数と目的変数を取得 ---
    explanatory_variable = dynamodb_handler.get_session_data_list(user_id=str(user_id), columns=["total_work_time", "break_time", "round_count"])
    objective_variable = dynamodb_handler.get_session_data_list(user_id=str(user_id), columns=["avg_focus_score"])
    print("説明変数リスト: ", explanatory_variable)
    print("目的変数リスト: ", objective_variable)
    
    # --- セッション最適化 ---
    opt = BayesianOptimizer("session")
    total_work_time, break_time, round_count = opt.optimize_session(explanatory_variable, objective_variable)
    print("提案された1セッションの作業時間の合計: ", total_work_time)
    print("提案されたセッション間休憩時間: ", break_time)
    print("提案されたラウンド繰り返し回数: ", round_count)
    
    # --- DynamoDBにデータを更新（新しいタイムスタンプで別のエントリとして保存） ---
    update_time = datetime.now().isoformat()
    dynamodb_handler.put_session_data(user_id=str(user_id), time=update_time, total_work_time=total_work_time, break_time=break_time, round_count=round_count)

    return {
        "total_work_time": total_work_time,
        "break_time": break_time,
        "round_count": int(round_count)
    }
