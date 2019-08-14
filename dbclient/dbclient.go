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
	fmt.Printf("Saving websocket id %v and userId %v in region %v", wsID, userID, os.Getenv("REGION"))
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
