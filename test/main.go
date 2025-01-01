package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Path+ ?" + request.Path,
	}, nil
}

func main() {
	log.Println("127.0.0.1:7954")
	os.Setenv("_LAMBDA_SERVER_PORT", "7954")
	os.Setenv("AWS_LAMBDA_RUNTIME_API", "127.0.0.1:7954")
	lambda.Start(handler)
}
