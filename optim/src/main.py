from csv_handler.csv_handler import CSVHandler
# from optimize.optimize_round import optimize_round
from optimize.bayesian_optimizer import BayesianOptimizer
from fastapi import FastAPI


app = FastAPI()

@app.get("/")
def hello():
    return "PomodoroOptimizationServer: ポモドーロ最適化サーバー"

@app.get("/round_optimizer/{user_id}")
def round_optimizer(user_id: str, focus_score: float):
    """ラウンド最適化API

    Args:
        focus_score (float): 集中度スコア
    Returns:
        work_time (float): 最適な作業時間
        break_time (float): 最適な休憩時間
    """
    csv_handler = CSVHandler(f"../data/round_csv/{user_id}.csv")
    # CSVファイル更新
    new_data = [focus_score]
    columns = ["focus_score"]
    csv_handler.update_data(new_data=new_data, columns=columns)
    # ラウンド最適化の説明変数と目的変数を取得
    explanatory_variable = csv_handler.make_chosen_data_list(columns=["work_time", "break_time"])
    objective_variable = csv_handler.make_chosen_data_list(columns=["focus_score"])
    print("説明変数リスト", explanatory_variable)
    print("目的変数リスト", objective_variable)
    # ラウンド最適化
    opt = BayesianOptimizer("round")
    work_time, break_time = opt.optimize_round(explanatory_variable, objective_variable)
    # work_time, break_time = optimize_round(explanatory_variable, objective_variable)
    print("提案された作業時間: ", work_time)
    print("提案された休憩時間: ", break_time)
    # CSVファイルの更新
    new_data = [work_time, break_time]
    columns = ["work_time", "break_time"]
    csv_handler.update_data(new_data=new_data, columns=columns)

    return {
        "work_time": work_time,
        "break_time": break_time
    }

@app.get("/session_optimizer/{user_id}")
def session_optimizer(user_id: str, average_focus_score: float):
    """セッション最適化API

    Args:
        average_focus_score (float): 集中度スコアの平均
    Returns:
        total_work_time (float): 1セッションの作業時間の合計
        session_break_time (float): 最適なセッション間休憩時間
        number_of_round (int): 最適なラウンド繰り返し回数
    """
    csv_handler = CSVHandler(f"../data/session_csv/{user_id}.csv")
    # CSVファイル更新
    new_data = [average_focus_score]
    columns = ["average_focus_score"]
    csv_handler.update_data(new_data=new_data, columns=columns)
    # セッション最適化の説明変数と目的変数を取得
    explanatory_variable = csv_handler.make_chosen_data_list(columns=["total_work_time", "session_break_time", "number_of_round"])
    objective_variable = csv_handler.make_chosen_data_list(columns=["average_focus_score"])
    print("説明変数リスト", explanatory_variable)
    print("目的変数リスト", objective_variable)
    # セッション最適化
    opt = BayesianOptimizer("session")
    total_work_time, session_break_time, number_of_round = opt.optimize_session(explanatory_variable, objective_variable)
    print("提案されたセッション間休憩時間: ", session_break_time)
    print("提案されたラウンド繰り返し回数: ", number_of_round)
    # CSVファイルの更新
    new_data = [session_break_time, number_of_round]
    columns = ["session_break_time", "number_of_round"]
    csv_handler.update_data(new_data=new_data, columns=columns)

    return {
        "session_break_time": session_break_time,
        "number_of_round": int(number_of_round)
    }