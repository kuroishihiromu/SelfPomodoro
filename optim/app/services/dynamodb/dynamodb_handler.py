# from datetime import datetime, timedelta
# import boto3

# class DynamoDBHandler:
#     def __init__(self, table_name):
#         self.client = boto3.client('dynamodb')
#         self.table_name = table_name

#     def update_data(self):
#         """DynamoDBのテーブルを更新"""
#         self.client.update_table(
#             TableName = self.table_name,
#             Item={
#              "area": {"S": area},
#              "name": {"S": name},
#              "ttl": {"N": ttl},
#             }
#         )
    
#     def choose_data(self):
#         """DynamoDBの特定のデータを取得"""
#         res = self.client.query(
#             TableName = self.table_name,
#             KeyConditionExpression="area = :area",
#             ExpressionAttributeValues={":area": {"S": "エリア1"}},
#     )

# print(response["Items"])
#     )

# print(response['Item'])



# area = "エリア１"
# name = "名前１"
# ttl = str(int((datetime.now() + timedelta(days=30)).timestamp()))

# client.put_item(
#     TableName=tablename,
#     Item={
#              "area": {"S": area},
#              "name": {"S": name},
#              "ttl": {"N": ttl},
#          }
# )

