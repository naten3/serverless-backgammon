package wsclient

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
)

type ApiGatewayManagementApi = apigatewaymanagementapi.ApiGatewayManagementApi
type WsClient struct {
	api *ApiGatewayManagementApi
}

var endpoint = "https://" + os.Getenv("WSAPI") + ".execute-api." + os.Getenv("REGION") + ".amazonaws.com/" + os.Getenv("STAGE")
var client = apigatewaymanagementapi.New(sess, aws.NewConfig().WithEndpoint(endpoint))

func New() *WsClient {
	sess := session.Must(session.NewSession())
	// todo get region and stage from environment
	endpoint := "https://" + os.Getenv("WSAPI") + ".execute-api." + os.Getenv("REGION") + ".amazonaws.com/" + os.Getenv("STAGE")
	fmt.Println("endpoint: " + endpoint)
	client := apigatewaymanagementapi.New(sess, aws.NewConfig().WithEndpoint(endpoint))

	return &WsClient{
		api: client,
	}
}

func (client WsClient) Post(connectionId string, action string, object interface{}) {
	body := map[string]interface{}{
		"action": action,
		"data":   object,
	}

	json, err := json.Marshal(body)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	output, err := client.api.PostToConnection(
		&apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: &connectionId,
			Data:         []byte(json),
		},
	)

	if err != nil {
		fmt.Println(output)
		fmt.Println(err.Error())
	}
}
