package dbclient

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
	fmt.Printf("Deleting user %v", wsID)
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

// WatchGame add watched game attribute to websocket table
func WatchGame(wsID string, gameID string) error {
	fmt.Printf("Deleting user %v", wsID)
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("WsUserTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"ConnectionId": {
				S: aws.String(wsID),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#AT": aws.String("WatchedGame"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				S: aws.String(gameID),
			},
		},
	}
	_, err := db.UpdateItem(input)
	return err
}

/*func getItem(isbn string) (*book, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("Books"),
		Key: map[string]*dynamodb.AttributeValue{
			"ISBN": {
				S: aws.String(isbn),
			},
		},
	}

	result, err := db.GetItem(input)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}

	bk := new(book)
	err = dynamodbattribute.UnmarshalMap(result.Item, bk)
	if err != nil {
		return nil, err
	}

	return bk, nil


// Add a book record to DynamoDB.
func putItem(bk *book) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String("Books"),
		Item: map[string]*dynamodb.AttributeValue{
			"ISBN": {
				S: aws.String(bk.ISBN),
			},
			"Title": {
				S: aws.String(bk.Title),
			},
			"Author": {
				S: aws.String(bk.Author),
			},
		},
	}

	_, err := db.PutItem(input)
	return err
}*/
