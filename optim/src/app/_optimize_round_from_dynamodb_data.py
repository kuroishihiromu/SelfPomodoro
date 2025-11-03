#! /usr/bin/env python3

import uuid
from handler.dynamodb.dynamodb_handler import DynamoDBHandler
from optimizer.bayesian_optimizer import BayesianOptimizer
from datetime import datetime


def _optimize_round_from_dynamodb_data(
    user_id: uuid.UUID,
    focus_score: float
):
    """DynamoDBのラウンド最適化

    Parameters:
        user_id (uuid.UUID): ユーザーID
        focus_score (float): 集中度スコア

    Returns:
        work_time (float): 最適な作業時間
        break_time (float): 最適な休憩時間
    """
    dynamodb_handler = DynamoDBHandler(table_name="round_optimization_logs", region_name="ap-northeast-1")
    
    # --- 最新の作業時間・休憩時間のデータを取得 ---
    latest_data = dynamodb_handler.get_round_data(user_id=str(user_id))
    latest_time = datetime.now().isoformat()
    
    if latest_data and isinstance(latest_data, list) and len(latest_data) > 0:
        # --- 最新のデータを取得 ---
        converted_data = dynamodb_handler._convert_to_list(latest_data)
        latest_item = converted_data[-1]
        latest_time = latest_item['time']
        work_time = latest_item.get('work_time')
        break_time = latest_item.get('break_time')
    else:
        work_time = None
        break_time = None
    
    # --- 最新のデータに集中度スコアを追加して更新 ---
    dynamodb_handler.put_round_data(
        user_id=str(user_id),
        time=latest_time,
        work_time=work_time,
        break_time=break_time,
        focus_score=focus_score
    )

    # --- 説明変数と目的変数を取得 ---
    explanatory_variable = dynamodb_handler.get_round_data_list(user_id=str(user_id), columns=["work_time", "break_time"])
    objective_variable = dynamodb_handler.get_round_data_list(user_id=str(user_id), columns=["focus_score"])
    print("説明変数リスト: ", explanatory_variable)
    print("目的変数リスト: ", objective_variable)
    
    # --- ラウンド最適化 ---
    opt = BayesianOptimizer("round")
    work_time, break_time = opt.optimize_round(explanatory_variable, objective_variable)
    print("提案された作業時間: ", work_time)
    print("提案された休憩時間: ", break_time)
    
    # --- DynamoDBにデータを更新（新しいタイムスタンプで別のエントリとして保存） ---
    update_time = datetime.now().isoformat()
    dynamodb_handler.put_round_data(user_id=str(user_id), time=update_time, work_time=work_time, break_time=break_time)

    return {
        "work_time": work_time,
        "break_time": break_time
    }
