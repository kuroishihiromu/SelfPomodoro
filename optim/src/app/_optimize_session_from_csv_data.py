#! /usr/bin/env python3

import uuid
from optim.src.handler._csv.csv_handler import CSVHandler
from optimizer.bayesian_optimizer import BayesianOptimizer

def _optimize_session_from_csv_data(
    user_id: uuid.UUID,
    avg_focus_score: float
):
    """セッション最適化

    Parameters:
        user_id (uuid.UUID): ユーザーID
        avg_focus_score (float): 集中度スコアの平均

    Returns:
        total_work_time (float): 1セッションの作業時間の合計
        break_time (float): 最適なセッション間休憩時間
        round_count (int): 最適なラウンド繰り返し回数
    """
    csv_handler = CSVHandler(f"app/data/session/{user_id}.csv")

    # --- CSVファイル更新 ---
    new_data = [avg_focus_score]
    columns = ["avg_focus_score"]
    csv_handler.update_data(new_data=new_data, columns=columns)

    # --- 説明変数と目的変数を取得 ---
    explanatory_variable = csv_handler.make_chosen_data_list(columns=["total_work_time", "break_time", "round_count"])
    objective_variable = csv_handler.make_chosen_data_list(columns=["avg_focus_score"])
    print("説明変数リスト", explanatory_variable)
    print("目的変数リスト", objective_variable)

    # --- セッション最適化 ---
    opt = BayesianOptimizer("session")
    total_work_time, break_time, round_count = opt.optimize_session(explanatory_variable, objective_variable)
    print("提案されたセッション間休憩時間: ", break_time)
    print("提案されたラウンド繰り返し回数: ", round_count)
    print("提案された1セッションの作業時間の合計: ", total_work_time)

    # --- CSVファイルの更新 ---
    new_data = [break_time, round_count, total_work_time]
    columns = ["break_time", "round_count", "total_work_time"]
    csv_handler.update_data(new_data=new_data, columns=columns)

    return {
        "total_work_time": total_work_time,
        "break_time": break_time,
        "round_count": int(round_count)
    }
