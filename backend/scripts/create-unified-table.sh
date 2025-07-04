#!/bin/bash

# 統合DynamoDBテーブル作成スクリプト
# Single Table Design with GSI support

TABLE_NAME="selfpomodoro_unified_table_dev"
REGION="ap-northeast-1"

echo "Creating unified DynamoDB table: $TABLE_NAME"

# JSON定義ファイルを作成
cat > /tmp/table-definition.json << 'EOF'
{
    "TableName": "selfpomodoro_unified_table_dev",
    "AttributeDefinitions": [
        {
            "AttributeName": "PK",
            "AttributeType": "S"
        },
        {
            "AttributeName": "SK",
            "AttributeType": "S"
        },
        {
            "AttributeName": "GSI1PK",
            "AttributeType": "S"
        },
        {
            "AttributeName": "GSI1SK",
            "AttributeType": "S"
        },
        {
            "AttributeName": "GSI2PK",
            "AttributeType": "S"
        },
        {
            "AttributeName": "GSI2SK",
            "AttributeType": "S"
        }
    ],
    "KeySchema": [
        {
            "AttributeName": "PK",
            "KeyType": "HASH"
        },
        {
            "AttributeName": "SK",
            "KeyType": "RANGE"
        }
    ],
    "GlobalSecondaryIndexes": [
        {
            "IndexName": "GSI1",
            "KeySchema": [
                {
                    "AttributeName": "GSI1PK",
                    "KeyType": "HASH"
                },
                {
                    "AttributeName": "GSI1SK",
                    "KeyType": "RANGE"
                }
            ],
            "Projection": {
                "ProjectionType": "ALL"
            },
            "ProvisionedThroughput": {
                "ReadCapacityUnits": 5,
                "WriteCapacityUnits": 5
            }
        },
        {
            "IndexName": "GSI2",
            "KeySchema": [
                {
                    "AttributeName": "GSI2PK",
                    "KeyType": "HASH"
                },
                {
                    "AttributeName": "GSI2SK",
                    "KeyType": "RANGE"
                }
            ],
            "Projection": {
                "ProjectionType": "ALL"
            },
            "ProvisionedThroughput": {
                "ReadCapacityUnits": 5,
                "WriteCapacityUnits": 5
            }
        }
    ],
    "ProvisionedThroughput": {
        "ReadCapacityUnits": 10,
        "WriteCapacityUnits": 10
    }
}
EOF

# テーブル作成実行
aws dynamodb create-table \
    --cli-input-json file:///tmp/table-definition.json \
    --region $REGION

if [ $? -eq 0 ]; then
    echo "Table creation initiated successfully"
    echo "Waiting for table to become active..."
    
    # テーブルがアクティブになるまで待機
    aws dynamodb wait table-exists --table-name $TABLE_NAME --region $REGION
    
    echo "Table $TABLE_NAME is now active!"
    
    # テーブル情報表示
    aws dynamodb describe-table --table-name $TABLE_NAME --region $REGION --query 'Table.{TableName:TableName,Status:TableStatus,ItemCount:ItemCount}'
    
    # 一時ファイル削除
    rm -f /tmp/table-definition.json
else
    echo "Failed to create table"
    rm -f /tmp/table-definition.json
    exit 1
fi