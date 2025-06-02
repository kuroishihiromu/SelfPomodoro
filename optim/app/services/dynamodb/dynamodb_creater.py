import boto3

class DynamoDBCreater:
    def __init__(self):
        """DynamoDBの初期化"""
        self.dynamodb_client = boto3.client('dynamodb')

    def create_table(self):
        """DynamoDBのテーブルを作成"""
        try:
            """ラウンドデータ"""
            self.client.create_table(
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
        

        try:
            """セッションデータ"""
            self.client.create_table(
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
        

        # try:
        #     """ユーザ関連データ"""
        #     self.client.create_table(
        #         TableName='user_configs',
        #         KeySchema=[
        #             {'AttributeName': 'user_id', 'KeyType': 'HASH'},
        #         ],
        #         AttributeDefinitions=[
        #             {'AttributeName': 'user_id', 'AttributeType': 'S'},  # uuidも文字列扱い
        #         ],
        #         ProvisionedThroughput={
        #             'ReadCapacityUnits': 1,
        #             'WriteCapacityUnits': 1
        #         }
        #     )

        # except Exception as e:
        #     print(f"DynamoDBのテーブル作成に失敗しました->\n {e}")
        #     raise Exception(f"DynamoDBのテーブル作成に失敗しました->\n {e}")
