FROM python:3.12-slim

# AWS Lambdaへのデプロイに必要
COPY --from=public.ecr.aws/awsguru/aws-lambda-adapter:0.9.0 /lambda-adapter /opt/extensions/lambda-adapter 

WORKDIR /code

COPY ./requirements.txt /code/requirements.txt

RUN apt-get update && apt-get install -y curl && \
  curl -LsSf https://astral.sh/uv/install.sh | sh && \
  /root/.local/bin/uv pip install --system --no-cache-dir -r /code/requirements.txt

COPY . .

ENV AWS_LAMBDA_FUNCTION_HANDLER=app.main:app
ENV AWS_LAMBDA_FUNCTION_MEMORY_SIZE=128
ENV AWS_LAMBDA_FUNCTION_NAME=pomodoro-optimizer
ENV AWS_LAMBDA_FUNCTION_VERSION=1
ENV AWS_LAMBDA_RUNTIME_API=127.0.0.1

# FastAPIを起動
CMD [ "uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8080" ]
