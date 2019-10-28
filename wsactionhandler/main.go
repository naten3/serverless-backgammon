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
		} else if payload.Type == "wsChangeName" {
			err = changeName(payload.Payload, userID)
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
	if err != nil {
		return err
	}
	if err == nil && fetchedGame != nil {
		err = wsclient.Post(connectionID, "wsWatchedGame", addDisplayNames(fetchedGame))
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
		// TODO if two players join at once there's maybe a race condition
		newGame := game.NewGame(gameID)
		fetchedGame = &newGame
	}
	if len(fetchedGame.White) > 0 && len(fetchedGame.Black) > 0 {
		return errors.New("Game is full")
	}
	if fetchedGame.White == userID || fetchedGame.Black == userID {
		fmt.Println(userID + " attempted to join game again")
		return nil
	}
	var color game.Color
	if len(fetchedGame.White) == 0 && len(fetchedGame.Black) == 0 {
		if rand.Intn(2)%2 == 0 {
			color = game.Black
		} else {
			color = game.White
		}
	}
	if len(fetchedGame.White) == 0 {
		color = game.White
	} else {
		color = game.Black
	}

	if color == game.White {
		fetchedGame.White = userID
	} else {
		fetchedGame.Black = userID
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
	err = wsclient.PostToMultiple(watchers, "wsUserJoined", addDisplayNames(fetchedGame))
	return err
}

func changeName(name string, userID string) error {
	err := dbclient.SetName(userID, name)
	if err != nil {
		return err
	}
	watchers, err2 := getUserGameWatchers(userID)
	if err2 != nil {
		return err2
	}
	err = wsclient.PostToMultiple(watchers, "wsUserNameChanged", map[string]interface{}{
		"userId":      userID,
		"displayName": name,
	})
	return err
}

// get all watchers for all games this user is a player in
func getUserGameWatchers(userID string) ([]string, error) {
	userGames, err := dbclient.GetAllUserGames(userID)
	if err != nil {
		return nil, err
	}
	allWatchers := make([]string, 0)
	// Loop through the games this user is a player in and get all watchers for each game
	// This is gross but I don't see a better way
	for _, game := range userGames {
		watchers, err2 := dbclient.GetGameWatchers(game.Id)
		if err2 != nil {
			return nil, err2
		}
		allWatchers = append(allWatchers, watchers...)
	}
	return allWatchers, nil
}

func addDisplayNames(game *game.Game) interface{} {
	whiteName := ""
	blackName := ""
	whiteName, _ = dbclient.GetUserName(game.White)
	blackName, _ = dbclient.GetUserName(game.Black)
	return map[string]interface{}{
		"game":      game,
		"whiteName": whiteName,
		"blackName": blackName,
	}
}

func main() {
	lambda.Start(Handler)
}
