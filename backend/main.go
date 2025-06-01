package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// HealthResponse はヘルスチェックのレスポンス
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// handler はLambdaのエントリーポイント
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := HealthResponse{
		Status:  "OK",
		Version: "serverless-dev",
	}

	body, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Internal server error"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type",
			"Access-Control-Allow-Methods": "GET,OPTIONS",
		},
		Body: string(body),
	}, nil
}

func main() {
	lambda.Start(handler)
}
