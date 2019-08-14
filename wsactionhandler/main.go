package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"serverless-backgammon/dbclient"
	"serverless-backgammon/game"

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

		userID, userFetchError := dbclient.GetAuthenticatedUserID(connectionID)
		if userFetchError != nil {
			fmt.Println("Error " + userFetchError.Error())
			return Response{
				StatusCode: 400,
				Body:       "error",
			}, userFetchError
		}

		var err error
		if payload.Action == "watchGame" {
			err = dbclient.WatchGame(connectionID, payload.Data)
		}
		if payload.Action == "joinGame" {
			err = joinGame(payload.Data, userID)
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
		}, err
	}
	fmt.Println("Error unmarshalling payload " + unmarshalErr.Error())
	return Response{
		StatusCode: 400,
		Body:       "invalid",
	}, nil
}

// create game if not present, add player, notify all watchers of game
func joinGame(gameID string, userID string) error {
	fetchedGame, err := dbclient.GetGame(gameID)
	if err != nil {
		return err
	}
	if fetchedGame == nil {
		newGame := game.NewGame(gameID)
		fetchedGame = &newGame
	}
	if fetchedGame.White != nil && fetchedGame.Black != nil {
		return errors.New("Game is full")
	}
	if fetchedGame.White != nil && (*fetchedGame.White == userID) || fetchedGame.Black != nil && (*fetchedGame.Black == userID) {
		return errors.New(userID + " attempted to join game again")
	}
	var color game.Color
	if fetchedGame.White == nil && fetchedGame.Black == nil {
		if rand.Intn(2)%2 == 0 {
			color = game.Black
		} else {
			color = game.White
		}
	}
	if fetchedGame.White == nil {
		color = game.White
	} else {
		color = game.Black
	}

	if color == game.White {
		fetchedGame.White = &userID
	} else {
		fetchedGame.Black = &userID
	}

	dbclient.SaveGame(*fetchedGame)
	return nil
}

func main() {
	lambda.Start(Handler)
}
