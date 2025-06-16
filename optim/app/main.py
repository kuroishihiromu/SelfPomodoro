import uuid
from datetime import datetime
from fastapi import FastAPI
from app.services.csv.csv_handler import CSVHandler
from app.services.dynamodb.dynamodb_handler import DynamoDBHandler
from app.services.optimize.bayesian_optimizer import BayesianOptimizer

app = FastAPI()

@app.get("/")
def hello():
    return "PomodoroOptimizationServer: ポモドーロ最適化サーバー"


@app.get("/csv_round_optimizer/{user_id}")
def csv_round_optimizer(user_id: uuid.UUID, focus_score: float):
    """ラウンド最適化API

    Args:
        user_id (uuid.UUID): ユーザーID
        focus_score (float): 集中度スコア

    Returns:
        work_time (float): 最適な作業時間
        break_time (float): 最適な休憩時間
    """
    csv_handler = CSVHandler(f"app/data/round/{user_id}.csv")
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


@app.get("/csv_session_optimizer/{user_id}")
def csv_session_optimizer(user_id: uuid.UUID, avg_focus_score: float):
    """セッション最適化API

    Args:
        user_id (uuid.UUID): ユーザーID
        avg_focus_score (float): 集中度スコアの平均

    Returns:
        total_work_time (float): 1セッションの作業時間の合計
        break_time (float): 最適なセッション間休憩時間
        round_count (int): 最適なラウンド繰り返し回数
    """
    csv_handler = CSVHandler(f"app/data/session/{user_id}.csv")
    # CSVファイル更新
    new_data = [avg_focus_score]
    columns = ["avg_focus_score"]
    csv_handler.update_data(new_data=new_data, columns=columns)
    # セッション最適化の説明変数と目的変数を取得
    explanatory_variable = csv_handler.make_chosen_data_list(columns=["total_work_time", "break_time", "round_count"])
    objective_variable = csv_handler.make_chosen_data_list(columns=["avg_focus_score"])
    print("説明変数リスト", explanatory_variable)
    print("目的変数リスト", objective_variable)
    # セッション最適化
    opt = BayesianOptimizer("session")
    total_work_time, break_time, round_count = opt.optimize_session(explanatory_variable, objective_variable)
    print("提案されたセッション間休憩時間: ", break_time)
    print("提案されたラウンド繰り返し回数: ", round_count)
    print("提案された1セッションの作業時間の合計: ", total_work_time)
    # CSVファイルの更新
    new_data = [break_time, round_count, total_work_time]
    columns = ["break_time", "round_count", "total_work_time"]
    csv_handler.update_data(new_data=new_data, columns=columns)

    return {
        "total_work_time": total_work_time,
        "break_time": break_time,
        "round_count": int(round_count)
    }


@app.get("/dynamo_round_optimizer/{user_id}")
def dynamo_round_optimizer(user_id: uuid.UUID, focus_score: float):
    """DynamoDBのラウンド最適化API

    Args:
        user_id (uuid.UUID): ユーザーID
        focus_score (float): 集中度スコア

    Returns:
        work_time (float): 最適な作業時間
        break_time (float): 最適な休憩時間
    """
    dynamodb_handler = DynamoDBHandler(table_name="round_optimization_logs", region_name="ap-northeast-1")
    
    # 最新の作業時間・休憩時間のデータを取得
    latest_data = dynamodb_handler.get_round_data(user_id=str(user_id))
    if latest_data:
        # 最新のデータを取得
        latest_data = dynamodb_handler._convert_dynamodb_items_to_list(latest_data)[-1]
        latest_time = latest_data['time']
    else:
        latest_time = datetime.now().isoformat()
    
        # 最新のデータに集中度スコアを追加して更新
        dynamodb_handler.put_round_data(
            user_id=str(user_id),
            time=latest_time,
            work_time=latest_data.get('work_time'),
            break_time=latest_data.get('break_time'),
            focus_score=focus_score
        )

    
    # ラウンド最適化の説明変数と目的変数を取得
    explanatory_variable = dynamodb_handler.make_chosen_data_list(user_id=str(user_id), columns=["work_time", "break_time"])
    objective_variable = dynamodb_handler.make_chosen_data_list(user_id=str(user_id), columns=["focus_score"])
    print("説明変数リスト", explanatory_variable)
    print("目的変数リスト", objective_variable)
    
    # ラウンド最適化
    opt = BayesianOptimizer("round")
    work_time, break_time = opt.optimize_round(explanatory_variable, objective_variable)
    print("提案された作業時間: ", work_time)
    print("提案された休憩時間: ", break_time)
    
    # DynamoDBにデータを更新（新しいタイムスタンプで別のエントリとして保存）
    update_time = datetime.now().isoformat()
    dynamodb_handler.put_round_data(user_id=str(user_id), time=update_time, work_time=work_time, break_time=break_time)

    return {
        "work_time": work_time,
        "break_time": break_time
    }


@app.get("/dynamo_session_optimizer/{user_id}")
def dynamo_session_optimizer(user_id: uuid.UUID, avg_focus_score: float):
    """DynamoDBのセッション最適化API

    Args:
        user_id (uuid.UUID): ユーザーID
        avg_focus_score (float): 集中度スコアの平均

    Returns:
        total_work_time (float): 1セッションの作業時間の合計
        break_time (float): 最適なセッション間休憩時間
        round_count (int): 最適なラウンド繰り返し回数
    """
    dynamodb_handler = DynamoDBHandler(table_name="session_optimization_logs", region_name="ap-northeast-1")
    
    # 最新のセッションデータを取得
    latest_data = dynamodb_handler.get_session_data(user_id=str(user_id))
    if latest_data:
        # 最新のデータを取得
        latest_data = dynamodb_handler._convert_dynamodb_items_to_list(latest_data)[-1]
        latest_time = latest_data['time']
    else:
        latest_time = datetime.now().isoformat()
    
    # 最新のデータに集中度スコアを追加して更新
    dynamodb_handler.put_session_data(
        user_id=str(user_id),
        time=latest_time,
        total_work_time=latest_data.get('total_work_time') if latest_data else None,
        break_time=latest_data.get('break_time') if latest_data else None,
        round_count=latest_data.get('round_count') if latest_data else None,
        avg_focus_score=avg_focus_score
    )
    
    # セッション最適化の説明変数と目的変数を取得
    explanatory_variable = dynamodb_handler.make_chosen_session_data_list(user_id=str(user_id), columns=["total_work_time", "break_time", "round_count"])
    objective_variable = dynamodb_handler.make_chosen_session_data_list(user_id=str(user_id), columns=["avg_focus_score"])
    print("説明変数リスト", explanatory_variable)
    print("目的変数リスト", objective_variable)
    
    # セッション最適化
    opt = BayesianOptimizer("session")
    total_work_time, break_time, round_count = opt.optimize_session(explanatory_variable, objective_variable)
    print("提案された1セッションの作業時間の合計: ", total_work_time)
    print("提案されたセッション間休憩時間: ", break_time)
    print("提案されたラウンド繰り返し回数: ", round_count)
    
    # DynamoDBにデータを更新（新しいタイムスタンプで別のエントリとして保存）
    update_time = datetime.now().isoformat()
    dynamodb_handler.put_session_data(user_id=str(user_id), time=update_time, total_work_time=total_work_time, break_time=break_time, round_count=round_count)

    return {
        "total_work_time": total_work_time,
        "break_time": break_time,
        "round_count": int(round_count)
    }
