package main

import (
	"context"
	"encoding/json"
	notion "notionboy/internal/pkg/notion"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

type RequestData struct {
	Config  notion.NotionConfig `json:"config"`
	Content notion.Content      `json:"content"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var data RequestData
	err := json.Unmarshal([]byte(request.Body), &data)
	if err != nil {
		resp := Response{
			StatusCode:      400,
			IsBase64Encoded: false,
			Body:            request.Body,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}

		return resp, nil
	}

	msg, _ := notion.CreateNewRecord(ctx, &data.Config, &data.Content)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            msg,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
