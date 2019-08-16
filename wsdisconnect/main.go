package main

import (
	"context"
	"serverless-backgammon/dbclient"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	// todo better way to import this
)

type Response events.APIGatewayProxyResponse

func Handler(context context.Context, request events.APIGatewayWebsocketProxyRequest) (Response, error) {
	connectionID := request.RequestContext.ConnectionID
	dbclient.DeleteWsUser(connectionID)
	return Response{
		StatusCode: 200,
		Body:       "success",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
