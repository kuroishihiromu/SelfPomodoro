#! /usr/bin/env python3

import uuid
import numpy as np
from concurrent.futures import ThreadPoolExecutor
from datetime import datetime
from handler.dynamodb.dynamodb_handler import DynamoDBHandler
from helper.make_time_series_data import make_time_series_data
from model.focus_score_model import FocusScoreModel
from optimizer.bayesian_optimizer import BayesianOptimizer

def optimize_round_with_focus_score_prediction(
    user_id: uuid.UUID,
    user_input_focus_score: float
):
    """疲労度とモチベーション予測モデルを追加したDynamoDBのラウンド最適化

    Parameters:
        user_id (uuid.UUID): ユーザーID

    Returns:
        round_data (list): ラウンドデータ
    """
    # --- DynamoDBインスタンス ---
    round_dynamodb_handler = DynamoDBHandler(
        table_name="round_optimization_logs",
        region_name="ap-northeast-1"
    )

    # --- 過去4日間のデータを取得 ---
    latest_data, second_latest_data, third_latest_data, fourth_latest_data = round_dynamodb_handler.get_four_past_days_data(user_id=str(user_id))

    # --- time_stepを取得 ---
    time_step = min(len(latest_data), len(second_latest_data), len(third_latest_data), len(fourth_latest_data))

    # --- 時系列データの作成 ---
    # 最新の日付の時系列データ + 集中度スコアのないデータ
    latest_time_series_data, latest_no_focus_score_data = make_time_series_data(latest_data)
    # 過去3日間の時系列データ
    past_time_datas = [second_latest_data, third_latest_data, fourth_latest_data]
    with ThreadPoolExecutor(max_workers=len(past_time_datas)) as executor:
        futures = [executor.submit(make_time_series_data, data) for data in past_time_datas]
        results = [future.result() for future in futures]
    past_time_series_data_list = [result[0] for result in results]
    # 時系列データの結合
    past_time_series_data_list.append(latest_time_series_data)
    
    # --- モデルの作成 ---
    model = FocusScoreModel(time_step=time_step)
    
    # --- モデルの訓練 ---
    model.fit(train_data=np.concatenate(past_time_series_data_list, axis=0))

    print(f"予測データ: {latest_no_focus_score_data}")
    
    # --- 予測の実行 ---
    predicted_focus_score = model.predict(
        test_data=np.array(latest_no_focus_score_data)
    )
    print("予測完了")
    print("予測されたfocus_score:", predicted_focus_score)
    
    # --- 最新の作業時間・休憩時間のデータを取得 ---
    latest_data = round_dynamodb_handler.get_round_data(user_id=str(user_id))
    latest_time = datetime.now().isoformat()
    
    if latest_data and isinstance(latest_data, list) and len(latest_data) > 0:
        # --- 最新のデータを取得 ---
        converted_data = round_dynamodb_handler._convert_to_list(latest_data)
        latest_item = converted_data[-1]
        latest_time = latest_item['time']
        work_time = latest_item.get('work_time')
        break_time = latest_item.get('break_time')
    else:
        work_time = None
        break_time = None

    # --- 集中度スコアを導出 ---
    focus_score = 2/3 * user_input_focus_score + 1/3 * predicted_focus_score
    print("導出された集中度スコア: ", focus_score)
    
    # --- 最新のデータに集中度スコアを追加して更新 ---
    round_dynamodb_handler.put_round_data(
        user_id=str(user_id),
        time=latest_time,
        work_time=work_time,
        break_time=break_time,
        focus_score=focus_score
    )

    # --- 説明変数と目的変数を取得 ---
    explanatory_variable = round_dynamodb_handler.get_round_data_list(user_id=str(user_id), columns=["work_time", "break_time"])
    objective_variable = round_dynamodb_handler.get_round_data_list(user_id=str(user_id), columns=["focus_score"])
    print("説明変数リスト: ", explanatory_variable)
    print("目的変数リスト: ", objective_variable)
    
    # --- ラウンド最適化 ---
    opt = BayesianOptimizer("round")
    work_time, break_time = opt.optimize_round(explanatory_variable, objective_variable)
    print("提案された作業時間: ", work_time)
    print("提案された休憩時間: ", break_time)
    
    # --- DynamoDBにデータを更新（新しいタイムスタンプで別のエントリとして保存） ---
    update_time = datetime.now().isoformat()
    round_dynamodb_handler.put_round_data(user_id=str(user_id), time=update_time, work_time=work_time, break_time=break_time)

    return {
        "work_time": work_time,
        "break_time": break_time
    }
