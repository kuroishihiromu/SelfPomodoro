import uuid
from typing import List, Dict, Union
from helper.make_data_list import make_data_list
from optimizer.bayesian_optimizer import BayesianOptimizer


def optimize_round_from_requestbody_data(
    user_id: uuid.UUID,
    round_data: List[Dict[str, Union[str, float, int]]]
):
    """リクエストボディデータを使用したラウンド最適化

    Parameters:
        user_id (uuid.UUID): ユーザーID
        round_data (List[Dict[str, Union[str, float, int]]]): ラウンドデータのリスト

    Returns:
        work_time (float): 最適な作業時間
        break_time (float): 最適な休憩時間
    """
    
    print(f"ユーザーID: {user_id} がラウンド最適化をリクエスト")
    print(f"ラウンドデータ: {round_data}")
    
    # --- 説明変数と目的変数を取得 ---
    explanatory_variable = make_data_list(round_data, ["work_time", "break_time"])
    objective_variable = make_data_list(round_data, ["focus_score"])
    print("説明変数リスト: ", explanatory_variable)
    print("目的変数リスト: ", objective_variable)
    
    # --- ラウンド最適化 ---
    opt = BayesianOptimizer("round")
    work_time, break_time = opt.optimize_round(explanatory_variable, objective_variable)
    print("提案された作業時間: ", work_time)
    print("提案された休憩時間: ", break_time)
    

    return {
        "work_time": work_time,
        "break_time": break_time
    }
