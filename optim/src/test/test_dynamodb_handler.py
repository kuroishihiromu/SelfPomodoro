import boto3
import os
from moto import mock_aws
from handler.dynamodb.dynamodb_handler import DynamoDBHandler

def set_env():
    """公式推奨のenv設定
    万が一にも本番環境へ変更が当たらないようにする
    """
    os.environ['AWS_ACCESS_KEY_ID'] = 'test'
    os.environ['AWS_SECRET_ACCESS_KEY'] = 'test'
    os.environ['AWS_DEFAULT_REGION'] = 'ap-northeast-1'
    os.environ['AWS_SECURITY_TOKEN'] = 'test'
    os.environ['AWS_SESSION_TOKEN'] = 'test'

@mock_aws
def test_dynamodb_handler():
    """DynamoDBHandlerのテスト"""
    # --- 環境変数を設定 ---
    set_env()
    
    # --- モックの設定 ---
    with mock_aws():
        dynamodb = boto3.client('dynamodb')
        dynamodb.create_table(
            TableName='round_optimization_logs',
            KeySchema=[
                {'AttributeName': 'user_id', 'KeyType': 'HASH'},
                {'AttributeName': 'time', 'KeyType': 'RANGE'}
            ],
            AttributeDefinitions=[
                {'AttributeName': 'user_id', 'AttributeType': 'S'},
                {'AttributeName': 'time', 'AttributeType': 'S'}
            ],
            ProvisionedThroughput={
                'ReadCapacityUnits': 1,
                'WriteCapacityUnits': 1
            }
        )

    dynamodb_handler = DynamoDBHandler(table_name="round_optimization_logs", region_name="ap-northeast-1")

    # --- 説明変数と目的変数を更新 ---
    dynamodb_handler.put_round_data(
        user_id="123e4567-e89b-12d3-a456-426614174000",
        time="1999-01-01T00:00:00.000000",
        work_time=20,
        break_time=10,
        focus_score=90
    )

    # --- 説明変数と目的変数を取得 ---
    explanatory_variable = dynamodb_handler.get_round_data_list(user_id="123e4567-e89b-12d3-a456-426614174000", columns=["work_time", "break_time"])
    objective_variable = dynamodb_handler.get_round_data_list(user_id="123e4567-e89b-12d3-a456-426614174000", columns=["focus_score"])
    
    print("説明変数リスト: ", explanatory_variable)
    print("目的変数リスト: ", objective_variable)

    # --- 型チェック ---
    assert isinstance(explanatory_variable, list)
    assert isinstance(objective_variable, list)

    # --- 値チェック ---
    assert explanatory_variable[0] == [20, 10]
    assert objective_variable[0] == 90
