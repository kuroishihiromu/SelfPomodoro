#! /usr/bin/env python3

import uuid
from handler._csv.csv_handler import CSVHandler
from optimizer.bayesian_optimizer import BayesianOptimizer

def _optimize_round_from_csv_data(
    user_id: uuid.UUID,
    focus_score: float
):
    """ラウンド最適化

    Parameters:
        user_id (uuid.UUID): ユーザーID
        focus_score (float): 集中度スコア

    Returns:
        work_time (float): 最適な作業時間
        break_time (float): 最適な休憩時間
    """
    csv_handler = CSVHandler(f"app/data/round/{user_id}.csv")

    # --- CSVファイル更新 ---
    new_data = [focus_score]
    columns = ["focus_score"]
    csv_handler.update_data(new_data=new_data, columns=columns)

    # --- 説明変数と目的変数を取得 ---
    explanatory_variable = csv_handler.make_chosen_data_list(columns=["work_time", "break_time"])
    objective_variable = csv_handler.make_chosen_data_list(columns=["focus_score"])
    print("説明変数リスト: ", explanatory_variable)
    print("目的変数リスト: ", objective_variable)

    # --- ラウンド最適化 ---
    opt = BayesianOptimizer("round")
    work_time, break_time = opt.optimize_round(explanatory_variable, objective_variable)
    print("提案された作業時間: ", work_time)
    print("提案された休憩時間: ", break_time)

    # --- CSVファイルの更新 ---
    new_data = [work_time, break_time]
    columns = ["work_time", "break_time"]
    csv_handler.update_data(new_data=new_data, columns=columns)

    return {
        "work_time": work_time,
        "break_time": break_time
    }
