import uuid
import numpy as np
from datetime import datetime
from handler._dynamodb.dynamodb_handler import DynamoDBHandler
from helper.make_time_series_data import make_time_series_data
from model.focus_score_model import FocusScoreModel
from optimizer.bayesian_optimizer import BayesianOptimizer

def _optimize_round_from_dynamodb_data_with_focus_score_prediction(
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

    # --- 最新の日付と過去全てのデータを取得 ---
    latest_data, all_past_data = round_dynamodb_handler.get_latest_day_and_all_past_days_data(user_id=str(user_id))

    print(f"最新の日付のデータ: {latest_data}")
    print(f"過去全てのデータ: {all_past_data}")
    print(f"最新の日付のデータの長さ: {type(latest_data)}")
    print(f"過去全てのデータの長さ: {type(all_past_data)}")

    # --- time_stepを取得 ---
    time_step = max(1, len(latest_data) - 1)
    print(f"time_step: {time_step}")
    print(f"latest_dataの長さ: {len(latest_data)}")

    # --- 時系列データの作成 ---
    # 最新の日付の時系列データ
    latest_time_series_data = make_time_series_data(latest_data)
    # 過去全てのデータの時系列データ
    all_past_time_series_data = make_time_series_data(all_past_data)
    
    print(f"latest_time_series_dataの長さ: {len(latest_time_series_data)}")
    print(f"all_past_time_series_dataの長さ: {len(all_past_time_series_data)}")
    
    # --- モデルの作成 ---
    model = FocusScoreModel(time_step=time_step)
    
    # --- モデルの訓練 ---
    model.fit(train_data=np.array(all_past_time_series_data))

    print(f"予測データ: {latest_time_series_data}")
    
    # --- 予測の実行 ---
    predicted_focus_score = model.predict(
        test_data=np.array(latest_time_series_data)
    )
    print("予測完了")
    print("予測されたfocus_score:", predicted_focus_score)
    
    # --- 最新の作業時間・休憩時間のデータを取得 ---
    latest_data = round_dynamodb_handler.get_round_data(user_id=str(user_id))
    print("最新のラウンドデータ取得完了: ", latest_data)
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
