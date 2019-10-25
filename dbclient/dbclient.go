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
	WatchGame    string `json:"watchedGame"`
}

// GetAuthenticatedUserID get an authenticated user id associated with a websocket
func GetAuthenticatedUserID(wsID string) (string, error) {

	fmt.Println("getting user for wsID " + wsID)
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

type userInfo struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
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

// GetUserName get user's name
func GetUserName(userID string) (string, error) {
	fmt.Printf("getting name for userId %v\n", userID)
	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("UserInfo"),
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {
				S: aws.String(userID),
			},
		},
	})
	if err != nil {
		fmt.Println("error getting user " + err.Error())
		return "", err
	}

	if result != nil && result.Item != nil {
		ui := userInfo{}

		err = dynamodbattribute.UnmarshalMap(result.Item, &ui)
		return ui.Name, err
	}
	return "", nil
}

// WatchGame add watched game attribute to websocket table
func WatchGame(wsID string, gameID string) error {
	fmt.Printf("Watching game %v for websocket %v\n", gameID, wsID)
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

// DeleteWsConnection remove a websocket connection
func DeleteWsConnection(wsID string) error {
	fmt.Printf("Deleting item for websocket id %v\n", wsID)
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"connectionId": {
				S: aws.String(wsID),
			},
		},
		TableName: aws.String("WsUserTable"),
	}

	_, err := db.DeleteItem(input)
	return err
}

// GetGameWatchers get a list of ws connection ids watching this game
func GetGameWatchers(gameID string) ([]string, error) {
	fmt.Printf("getting watchers for  %v\n", gameID)
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String("WsUserTable"),
		IndexName: aws.String("watchedGame-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"watchedGame": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(gameID),
					},
				},
			},
		},
	}
	var resp1, err1 = db.Query(queryInput)
	if err1 != nil {
		fmt.Println("Error fetching watched games " + err1.Error())
		return nil, err1
	}
	wsUsers := []wsUser{}
	err := dynamodbattribute.UnmarshalListOfMaps(resp1.Items, &wsUsers)
	if err != nil {
		fmt.Println("Error unmarshalling watched games " + err1.Error())
		return nil, err1
	}

	result := []string{}
	for i := range wsUsers {
		result = append(result, wsUsers[i].ConnectionID)
	}
	return result, nil
}

// SaveGame save a game state
func SaveGame(game game.Game) error {
	fmt.Printf("Saving game %v\n", game.Id)
	item, err := dynamodbattribute.MarshalMap(game)
	if err != nil {
		fmt.Println("Error marshalling game " + err.Error())
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("Game"),
		Item:      item,
	}
	_, saveError := db.PutItem(input)
	if saveError != nil {
		fmt.Println("error saveing game " + saveError.Error())
		return err
	}
	return saveError
}

// GetGame get a game by id
func GetGame(gameID string) (*game.Game, error) {
	fmt.Printf("getting game for gameID %v\n", gameID)
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
