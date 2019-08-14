package dbclient

import (
	"fmt"
	"os"

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
			"ConnectionId": {
				S: aws.String(wsID),
			},
			"UserId": {
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
			"ConnectionId": {
				S: aws.String(wsID),
			},
		},
	}
	_, err := db.DeleteItem(input)
	return err
}

// GetAuthenticatedUserID get an authenticated user id associated with a websocket
func GetAuthenticatedUserID(wsID string) (string, error) {
	fmt.Printf("getting user for wsID %v", wsID)
	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("WsUserTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionId": {
				N: aws.String(wsID),
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
	return result.Item["UserId"].String(), nil
}

// SetName set user's name
func SetName(userID string, name string) error {
	fmt.Printf("Saving user info for user id %v", userID)
	input := &dynamodb.PutItemInput{
		TableName: aws.String("UserInfo"),
		Item: map[string]*dynamodb.AttributeValue{
			"UserId": {
				S: aws.String(userID),
			},
			"Name": {
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
			"ConnectionId": {
				S: aws.String(wsID),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				S: aws.String(gameID),
			},
		},
		UpdateExpression: aws.String("set WatchedGame = :t"),
	}
	_, err := db.UpdateItem(input)
	return err
}

// SaveGame save a game state
func SaveGame(game game.Game) error {
	fmt.Printf("Saving game %v for websocket", game.ID)
	input := &dynamodb.PutItemInput{
		TableName: aws.String("UserInfo"),
		Item: map[string]*dynamodb.AttributeValue{
			"UserId": {
				S: aws.String(userID),
			},
			"Name": {
				S: aws.String(name),
			},
		},
	}
	_, err := db.PutItem(input)
	return err
}

// GetGame get a game by id
func GetGame(gameID string) (game.Game, error) {
	fmt.Printf("getting user for wsID %v", wsID)
	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("Game"),
		Key: map[string]*dynamodb.AttributeValue{
			"GameId": {
				N: aws.String(gameID),
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
	return result.Item["UserId"].String(), nil
}
