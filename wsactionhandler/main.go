package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"serverless-backgammon/dbclient"
	"serverless-backgammon/game"
	"serverless-backgammon/wsclient"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	// todo better way to import this
)

type Response events.APIGatewayProxyResponse

type wsPayload struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func Handler(context context.Context, request events.APIGatewayWebsocketProxyRequest) (Response, error) {
	connectionID := request.RequestContext.ConnectionID

	var payload wsPayload
	unmarshalErr := json.Unmarshal([]byte(request.Body), &payload)
	if unmarshalErr == nil {

		userID, userFetchError := dbclient.GetAuthenticatedUserID(connectionID)
		if userFetchError != nil || userID == "" {
			fmt.Println("Error " + userFetchError.Error())
			return Response{
				StatusCode: 400,
				Body:       "error",
			}, userFetchError
		}
		println("user id is " + userID)
		println("payload type " + payload.Type)

		var err error
		if payload.Type == "wsWatchGame" {
			err = watchGame(connectionID, payload.Payload)
		} else if payload.Type == "wsJoinGame" {
			err = joinGame(payload.Payload, userID)
		} else {
			err = errors.New("unrecognized action")
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

// watch the game, post to websockets if in progress
func watchGame(connectionID string, gameID string) error {
	dbclient.WatchGame(connectionID, gameID)
	fetchedGame, err := dbclient.GetGame(gameID)
	if err == nil {
		fmt.Println("Posting game status to" + connectionID)
		err = wsclient.Post(connectionID, "wsWatchedGame", fetchedGame)
		return err
	}
	return nil
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
	if (*fetchedGame.White == userID) || (*fetchedGame.Black == userID) {
		fmt.Println(userID + " attempted to join game again")
		return nil
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

	saveErr := dbclient.SaveGame(*fetchedGame)
	if saveErr != nil {
		fmt.Println("save error " + saveErr.Error())
		return saveErr
	}

	//notify watchers that a user joined
	fmt.Println("getting game watchers")
	watchers, err := dbclient.GetGameWatchers(gameID)
	if err != nil {
		return err
	}

	fmt.Println("notifying game watchers")
	name, err := dbclient.GetUserName(userID)
	if err != nil {
		return err
	}
	err = wsclient.PostToMultiple(watchers, "wsUserJoined", map[string]interface{}{
		"userId": userID,
		"color":  color,
		"name":   name,
	})
	return err
}

func main() {
	lambda.Start(Handler)
}
