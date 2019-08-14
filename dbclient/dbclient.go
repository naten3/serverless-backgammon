package dbclient

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"serverless-backgammon/game"
)

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion(os.Getenv("REGION")))

// SaveVerifiedWsUser add a user and session id combination
func SaveVerifiedWsUser(wsID string, userID string) error {
	fmt.Printf("Saving websocket id %v and userId %v", wsID, userID)
	input := &dynamodb.PutItemInput{
		TableName: aws.String("WsUserTable"),
		Item: map[string]*dynamodb.AttributeValue{
			"connectionId": {
				S: aws.String(wsID),
			},
			"userId": {
				S: aws.String(userID),
			},
		},
	}
	_, err := db.PutItem(input)
	return err
}

// DeleteWsUser delete a websocket connection id
func DeleteWsUser(wsID string) error {
	fmt.Printf("Deleting ws connection %v", wsID)
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("WsUserTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"connectionId": {
				S: aws.String(wsID),
			},
		},
	}
	_, err := db.DeleteItem(input)
	return err
}

type wsUser struct {
	ConnectionID string `json:"connectionId"`
	UserID       string `json:"userId"`
}

// GetAuthenticatedUserID get an authenticated user id associated with a websocket
func GetAuthenticatedUserID(wsID string) (string, error) {

	fmt.Printf("getting user for wsID %v", wsID)
	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("WsUserTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"connectionId": {
				S: aws.String(wsID),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	if result == nil {
		return "", nil
	}

	ws := wsUser{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &ws)

	if err != nil {
		return "", err
	}

	return ws.UserID, nil
}

// SetName set user's name
func SetName(userID string, name string) error {
	fmt.Printf("Saving user info for user id %v", userID)
	input := &dynamodb.PutItemInput{
		TableName: aws.String("UserInfo"),
		Item: map[string]*dynamodb.AttributeValue{
			"userId": {
				S: aws.String(userID),
			},
			"name": {
				S: aws.String(name),
			},
		},
	}
	_, err := db.PutItem(input)
	return err
}

// WatchGame add watched game attribute to websocket table
func WatchGame(wsID string, gameID string) error {
	fmt.Printf("Watching game %v for websocket %v", gameID, wsID)
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("WsUserTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"connectionId": {
				S: aws.String(wsID),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				S: aws.String(gameID),
			},
		},
		UpdateExpression: aws.String("set watchedGame = :t"),
	}
	_, err := db.UpdateItem(input)
	return err
}

// SaveGame save a game state
func SaveGame(game game.Game) error {
	fmt.Printf("Saving game %v", game.Id)
	item, err := dynamodbattribute.MarshalMap(game)
	for k := range item {
		fmt.Println("Key " + k)
	}
	if err != nil {
		fmt.Println("Error saving game " + err.Error())
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Game"),
		Item:      item,
	}
	_, saveError := db.PutItem(input)
	return saveError
}

// GetGame get a game by id
func GetGame(gameID string) (*game.Game, error) {
	fmt.Printf("getting game for gameID %v", gameID)
	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("Game"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(gameID),
			},
		},
	})
	if err != nil {
		fmt.Println("error getting game " + err.Error())
		return nil, err
	}

	if result != nil && result.Item != nil {
		g := game.Game{}

		err = dynamodbattribute.UnmarshalMap(result.Item, &g)
		return &g, err
	}
	return nil, nil
}
