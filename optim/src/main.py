#! /usr/bin/env python3

import uuid
from fastapi import FastAPI
from app.optimize_round_from_csv_data import optimize_round_from_csv_data
from app.optimize_session_from_csv_data import optimize_session_from_csv_data
from app.optimize_round_from_dynamodb_data import optimize_round_from_dynamodb_data
from app.optimize_session_from_dynamodb_data import optimize_session_from_dynamodb_data

app = FastAPI()

@app.get("/")
def hello():
    return {"message": "PomodoroOptimizationServer: ポモドーロ最適化サーバー"}


@app.get("/round/v1/{user_id}")
def round_v1(
    user_id: uuid.UUID,
    focus_score: float
):
    return optimize_round_from_csv_data(user_id, focus_score)


@app.get("/session/v1/{user_id}")
def session_v1(
    user_id: uuid.UUID,
    avg_focus_score: float
):
    return optimize_session_from_csv_data(user_id, avg_focus_score)


@app.get("/round/v2/{user_id}")
def round_v2(
    user_id: uuid.UUID,
    focus_score: float
):
    return optimize_round_from_dynamodb_data(user_id, focus_score)


@app.get("/session/v2/{user_id}")
def session_v2(
    user_id: uuid.UUID,
    avg_focus_score: float
):
    return optimize_session_from_dynamodb_data(user_id, avg_focus_score)
