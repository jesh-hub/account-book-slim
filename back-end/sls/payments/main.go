package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jesh-hub/account-book-slack/back-end/pkg/abs"
	"log"
	"os"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	slackToken := os.Getenv("slackToken")
	channelId := os.Getenv("channelId")
	botId := os.Getenv("botId")

	slackClient := abs.NewSlackClient(slackToken, channelId, botId)
	messages := slackClient.GetMessages()
	messages = slackClient.FilterMessages(messages)

	body, err := json.Marshal(messages)
	if err != nil {
		return Response{StatusCode: 404}, err
	}

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(body),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "payments-handler",
		},
	}

	return resp, nil
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	lambda.Start(Handler)
}