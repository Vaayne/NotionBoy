package main

import (
	"context"
	"encoding/json"

	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	notionapi "github.com/dstotijn/go-notion"
	"github.com/sirupsen/logrus"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

type RequestData struct {
	Config  NotionConfig `json:"config"`
	Content Content      `json:"content"`
}

type Notion interface {
	ParseContent()
	CreateNewRecord()
}

type NotionConfig struct {
	DatabaseID  string `json:"database_id"`
	BearerToken string `json:"bearer_token"`
}

func GetNotionClient(token string) *notionapi.Client {
	return notionapi.NewClient(token, nil)
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

	msg := CreateNewRecord(ctx, &data.Config, &data.Content)

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