#! /usr/bin/env python3

import boto3


class DynamoDBCreator:
    """DynamoDBのテーブルを作成するクラス
    
    Attributes:
        dynamodb_client (boto3.client): DynamoDBのクライアント
    
    Tables:
        - round_optimization_logs: ラウンドデータを保存するテーブル
        - session_optimization_logs: セッションデータを保存するテーブル
        - user_configs: ユーザ関連の設定データを保存するテーブル
    """
    def __init__(
        self,
        region_name: str = 'ap-northeast-1'
    ) -> None:
        """DynamoDBクライアントを作成
        
        Parameters:
            region_name (str): AWSリージョン
        """
        self.dynamodb_client = boto3.client('dynamodb', region_name=region_name)

    def _create_table(
        self
    ) -> None:
        """DynamoDBのテーブルを作成"""
        # --- ラウンドデータ ---
        try:
            self.dynamodb_client.create_table(
                TableName='round_optimization_logs',
                KeySchema=[
                    {'AttributeName': 'user_id', 'KeyType': 'HASH'},  # パーティションキー
                    {'AttributeName': 'time', 'KeyType': 'RANGE'},    # ソートキー
                ],
                AttributeDefinitions=[
                    {'AttributeName': 'user_id', 'AttributeType': 'S'},  # 文字列
                    {'AttributeName': 'time', 'AttributeType': 'S'},
                ],
                ProvisionedThroughput={
                    'ReadCapacityUnits': 1,
                    'WriteCapacityUnits': 1
                }
            )

        except Exception as e:
            print(f"ラウンドデータのDynamoDBテーブル作成に失敗しました->\n {e}")
            raise Exception(f"ラウンドデータのDynamoDBテーブル作成に失敗しました->\n {e}")
        
        # --- セッションデータ ---
        try:
            self.dynamodb_client.create_table(
                TableName = "session_optimization_logs",
                KeySchema=[
                    {"AttributeName": 'user_id', 'KeyType': 'HASH'},
                    {'AttributeName': 'time', 'KeyType': 'RANGE'},
                ],
                AttributeDefinitions=[
                    {'AttributeName': 'user_id', 'AttributeType': 'S'},
                    {'AttributeName': 'time', 'AttributeType': 'S'},
                ],
                ProvisionedThroughput={
                    'ReadCapacityUnits': 1,
                    'WriteCapacityUnits': 1
                }
            )

        except Exception as e:
            print(f"セッションデータのDynamoDBテーブル作成に失敗しました->\n {e}")
            raise Exception(f"セッションデータのDynamoDBテーブル作成に失敗しました->\n {e}")
        
        # --- ユーザ関連データ ---
        try:
            self.dynamodb_client.create_table(
                TableName='user_configs',
                KeySchema=[
                    {'AttributeName': 'user_id', 'KeyType': 'HASH'},
                ],
                AttributeDefinitions=[
                    {'AttributeName': 'user_id', 'AttributeType': 'S'},  # uuidも文字列扱い
                ],
                ProvisionedThroughput={
                    'ReadCapacityUnits': 1,
                    'WriteCapacityUnits': 1
                }
            )

        except Exception as e:
            print(f"DynamoDBのテーブル作成に失敗しました->\n {e}")
            raise Exception(f"DynamoDBのテーブル作成に失敗しました->\n {e}")
