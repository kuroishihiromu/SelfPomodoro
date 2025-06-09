from fastapi.testclient import TestClient
from app.main import app

client = TestClient(app)

def test_round_optimizer():
    response = client.get("round_optimizer/123e4567-e89b-12d3-a456-426614174000?focus_score=100")
    assert response.status_code == 200
    data = response.json()
    assert isinstance(data["work_time"], (int, float))
    assert isinstance(data["break_time"], (int, float))

def test_session_optimizer():
    response = client.get("session_optimizer/123e4567-e89b-12d3-a456-426614174000?avg_focus_score=100")
    assert response.status_code == 200
    data = response.json()
    assert isinstance(data["break_time"], (int, float))
    assert isinstance(data["round_count"], int)
    assert isinstance(data["total_work_time"], (int, float))
    # 値の範囲チェック
    assert 10 <= data["break_time"] <= 60
    assert 1 <= data["round_count"] <= 8
    assert 60 <= data["total_work_time"] <= 480
