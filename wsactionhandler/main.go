package main

import (
	"context"
	"encoding/json"
	"fmt"
	"serverless-backgammon/dbclient"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	// todo better way to import this
)

type Response events.APIGatewayProxyResponse

type wsPayload struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

func Handler(context context.Context, request events.APIGatewayWebsocketProxyRequest) (Response, error) {
	// client := wsclient.New()
	connectionID := request.RequestContext.ConnectionID

	var payload wsPayload
	unmarshalErr := json.Unmarshal([]byte(request.Body), &payload)
	if unmarshalErr == nil {
		var err error
		if payload.Action == "watchGame" {
			err = dbclient.WatchGame(connectionID, payload.Data)
		}
		if payload.Action == "joinGame" {
			joinGame(payload.Data)
		}

		if err == nil {
			return Response{
				StatusCode: 200,
				Body:       "success",
			}, nil
		}
		fmt.Println("Error " + err.Error())
		return Response{
			StatusCode: 400,
			Body:       "error",
		}, nil
	}
	fmt.Println("Error unmarshalling payload " + unmarshalErr.Error())
	return Response{
		StatusCode: 400,
		Body:       "invalid",
	}, nil
}

// create game if not present, add player, notify all watchers of game
func joinGame(gameID string) {

}

func main() {
	lambda.Start(Handler)
}
