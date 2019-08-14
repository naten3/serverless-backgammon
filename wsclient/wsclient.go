package wsclient

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

var endpoint = "https://" + os.Getenv("WSAPI") + ".execute-api." + os.Getenv("REGION") + ".amazonaws.com/" + os.Getenv("STAGE")
var sess = session.Must(session.NewSession())
var client = apigatewaymanagementapi.New(sess, aws.NewConfig().WithEndpoint(endpoint))

// Post post a payload to a specific connection id
func Post(connectionID string, action string, object interface{}) error {
	body := map[string]interface{}{
		"action": action,
		"data":   object,
	}

	json, err := json.Marshal(body)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("posting to connection id " + connectionID)
	output, err := client.PostToConnection(
		&apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: &connectionID,
			Data:         []byte(json),
		},
	)

	if err != nil {
		fmt.Println(output)
		fmt.Println(err.Error())
	}
	return err
}
