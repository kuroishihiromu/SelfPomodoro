import uuid
from fastapi import FastAPI
from typing import Dict, List, Union
from app.optimize_round_with_focus_score_prediction import optimize_round_from_dynamodb_data_with_focus_score_prediction
from app.optimize_round_from_requestbody_data import optimize_round_from_requestbody_data
from app.optimize_session_from_dynamodb_data import optimize_session_from_dynamodb_data
from app.optimize_session_from_requestbody_data import optimize_session_from_requestbody_data


app = FastAPI()

@app.get("/")
def hello():
    return {"message": "PomodoroOptimizer: ポモドーロ最適化サーバー"}


# --- ラウンド最適化 ---
# CSVデータから最適化
@app.get("/round/v1/{user_id}")
def round_v1(
    user_id: uuid.UUID,
    focus_score: float
):
    # return _optimize_round_from_csv_data(user_id, focus_score)
    return {"message": "このエンドポイントは非推奨です。/round/v4/{user_id} を使用してください。"}


# DynamoDBデータから最適化
@app.get("/round/v2/{user_id}")
def round_v2(
    user_id: uuid.UUID,
    focus_score: float
):
    # return _optimize_round_from_dynamodb_data(user_id, focus_score)
    return {"message": "このエンドポイントは非推奨です。/round/v4/{user_id} を使用してください。"}


# 疲労度とモチベーション予測モデルを使用した最適化
@app.get("/round/v3/{user_id}")
def round_v3(
    user_id: uuid.UUID,
    user_input_focus_score: float
):
    return optimize_round_from_dynamodb_data_with_focus_score_prediction(user_id, user_input_focus_score)


# RequestBodyデータから最適化
@app.post("/round/v4/{user_id}")
def round_v4(
    user_id: uuid.UUID,
    round_data: List[Dict[str, Union[str, float, int]]]
):
    return optimize_round_from_requestbody_data(user_id, round_data)



# --- セッション最適化 ---
# CSVデータから最適化
@app.get("/session/v1/{user_id}")
def session_v1(
    user_id: uuid.UUID,
    avg_focus_score: float
):
    # return _optimize_session_from_csv_data(user_id, avg_focus_score)
    return {"message": "このエンドポイントは非推奨です。/session/v3/{user_id} を使用してください。"}


# DynamoDBデータから最適化
@app.get("/session/v2/{user_id}")
def session_v2(
    user_id: uuid.UUID,
    avg_focus_score: float
):
    return optimize_session_from_dynamodb_data(user_id, avg_focus_score)


# RequestBodyデータから最適化
@app.post("/session/v3/{user_id}")
def session_v3(
    user_id: uuid.UUID,
    session_data: List[Dict[str, Union[str, float, int]]]
):
    return optimize_session_from_requestbody_data(user_id, session_data)
