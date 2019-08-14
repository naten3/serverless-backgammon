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

func (client WsClient) Post(connectionId string, object interface{}) {
	json, error := json.Marshal(object)

	output, error := client.api.PostToConnection(
		&apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: &connectionId,
			Data:         []byte(json),
		},
	)

	if error != nil {
		fmt.Println(output)
		fmt.Println(error.Error())
	}
}
