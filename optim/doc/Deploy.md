# FastAPI Lambda デプロイ手順

このドキュメントでは、FastAPIアプリケーションをAWS Lambdaにデプロイする手順を説明します。

## 前提条件

- AWS CLIがインストールされていること
- Dockerがインストールされていること
- AWS認証情報が設定されていること (`aws configure`)

## 1. ECRリポジトリの作成

```bash
aws ecr create-repository --repository-name pomodoro-optimizer --region ap-northeast-1
```

## 2. ECRへのログイン

```bash
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin $(aws sts get-caller-identity --query Account --output text).dkr.ecr.ap-northeast-1.amazonaws.com
```

## 3. Dockerイメージのビルド

```bash
cd optim
docker build -t pomodoro-optimizer .
```

## 4. Dockerイメージのタグ付け

```bash
# ECRリポジトリURIの取得
ECR_URI=$(aws ecr describe-repositories --repository-names pomodoro-optimizer --region ap-northeast-1 --query 'repositories[0].repositoryUri' --output text)

# タグ付け
docker tag pomodoro-optimizer:latest $ECR_URI:latest
```

## 5. ECRへのプッシュ

```bash
docker push $ECR_URI:latest
```

## 6. IAMロールの作成

```bash
aws iam create-role --role-name lambda-pomodoro-optimizer-role --assume-role-policy-document '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"lambda.amazonaws.com"},"Action":"sts:AssumeRole"}]}'
```

## 7. IAMポリシーのアタッチ

```bash
aws iam attach-role-policy --role-name lambda-pomodoro-optimizer-role --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
```

## 8. Lambda関数の作成

```bash
# アカウントIDの取得
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)

# Lambda関数の作成
aws lambda create-function \
  --function-name pomodoro-optimizer \
  --role arn:aws:iam::$ACCOUNT_ID:role/lambda-pomodoro-optimizer-role \
  --code ImageUri=$ECR_URI:latest \
  --package-type Image \
  --timeout 60 \
  --memory-size 512 \
  --region ap-northeast-1
```

## 9. Function URLの作成

```bash
aws lambda create-function-url-config \
  --function-name pomodoro-optimizer \
  --auth-type NONE \
  --cors AllowMethods=["*"],AllowHeaders=["*"],AllowOrigins=["*"] \
  --region ap-northeast-1
```

## 10. パブリックアクセス許可の追加

```bash
aws lambda add-permission \
  --function-name pomodoro-optimizer \
  --statement-id public-access \
  --action lambda:invokeFunctionUrl \
  --principal "*" \
  --function-url-auth-type NONE \
  --region ap-northeast-1
```

## 11. Lambda関数の更新（必要に応じて）

新しいDockerイメージでLambda関数を更新する場合：

```bash
# 新しいバージョンでタグ付け
docker tag pomodoro-optimizer:latest $ECR_URI:v2

# プッシュ
docker push $ECR_URI:v2

# Lambda関数の更新
aws lambda update-function-code \
  --function-name pomodoro-optimizer \
  --image-uri $ECR_URI:v2 \
  --region ap-northeast-1
```

## 12. エンドポイントの確認

Function URLの取得：

```bash
aws lambda get-function-url-config \
  --function-name pomodoro-optimizer \
  --region ap-northeast-1 \
  --query 'FunctionUrl' \
  --output text
```

## テスト

```bash
# Function URLでのテスト
curl https://YOUR_FUNCTION_URL/

# 特定のエンドポイントのテスト
curl https://YOUR_FUNCTION_URL/round_optimizer/{user_id}?focus_score=0.8
```

## トラブルシューティング

### ログの確認

```bash
aws logs tail /aws/lambda/pomodoro-optimizer --region ap-northeast-1 --follow
```

### Lambda関数の詳細確認

```bash
aws lambda get-function --function-name pomodoro-optimizer --region ap-northeast-1
```

## クリーンアップ

リソースを削除する場合：

```bash
# Lambda関数の削除
aws lambda delete-function --function-name pomodoro-optimizer --region ap-northeast-1

# IAMロールの削除
aws iam detach-role-policy --role-name lambda-pomodoro-optimizer-role --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
aws iam delete-role --role-name lambda-pomodoro-optimizer-role

# ECRリポジトリの削除
aws ecr delete-repository --repository-name pomodoro-optimizer --force --region ap-northeast-1
```

## 注意事項

- リージョンは `ap-northeast-1` (東京) を使用
- Lambda関数のタイムアウトは60秒に設定
- メモリサイズは512MBに設定
- Function URLは認証なし（NONE）で設定
- CORS設定でオリジン、メソッド、ヘッダーをすべて許可
