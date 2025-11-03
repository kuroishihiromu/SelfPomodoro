import uuid
from typing import List, Dict, Union
from helper.make_data_list import make_data_list
from optimizer.bayesian_optimizer import BayesianOptimizer


def optimize_session_from_requestbody_data(
    user_id: uuid.UUID,
    session_data: List[Dict[str, Union[str, float, int]]]
):
    """リクエストボディデータを使用したセッション最適化

    Parameters:
        user_id (uuid.UUID): ユーザーID
        session_data (List[Dict[str, Union[str, float, int]]]): セッションデータのリスト

    Returns:
        total_work_time (float): 1セッションの作業時間の合計
        break_time (float): 最適なセッション間休憩時間
        round_count (int): 最適なラウンド繰り返し回数
    """
    print(f"ユーザーID: {user_id} がセッション最適化をリクエスト")
    print(f"セッションデータ: {session_data}")

    # --- 説明変数と目的変数を取得 ---
    explanatory_variable = make_data_list(session_data, ["total_work_time", "break_time", "round_count"])
    objective_variable = make_data_list(session_data, ["avg_focus_score"])
    print("説明変数リスト: ", explanatory_variable)
    print("目的変数リスト: ", objective_variable)
    
    # --- セッション最適化 ---
    opt = BayesianOptimizer("session")
    total_work_time, break_time, round_count = opt.optimize_session(explanatory_variable, objective_variable)
    print("提案された1セッションの作業時間の合計: ", total_work_time)
    print("提案されたセッション間休憩時間: ", break_time)
    print("提案されたラウンド繰り返し回数: ", round_count)
    

    return {
        "total_work_time": total_work_time,
        "break_time": break_time,
        "round_count": int(round_count)
    }
