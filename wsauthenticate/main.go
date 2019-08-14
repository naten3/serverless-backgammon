package main

import (
	"context"
	"encoding/json"
	"fmt"
	"serverless-backgammon/dbclient"
	"serverless-backgammon/wsclient"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dgrijalva/jwt-go"
	// todo better way to import this
)

type Response events.APIGatewayProxyResponse

type wsPayload struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

func Handler(context context.Context, request events.APIGatewayWebsocketProxyRequest) (Response, error) {
	websocketClient := wsclient.New()
	connectionID := request.RequestContext.ConnectionID

	var payload wsPayload
	unmarshalErr := json.Unmarshal([]byte(request.Body), &payload)
	if unmarshalErr != nil {
		fmt.Println(unmarshalErr.Error())
	}
	token := payload.Data
	fmt.Println("token is " + token)

	// TODO put in env variable
	secret := "Dba98iE002lTOA8YdQtYvdf2U52Eai7WT1sIoVTO-Q0r5KDdHNbIfBZS8P8Y-yFf6NunyqqFcB3HuvOivsEs-Zi4oka_FK4TbW52G9dSsxGoppciGEtUsTFgpKQYpQ7qyZE7ncvf39bWR0Y1RkP-yf2X2Ffeq7bv75vXE2TWhvZU6oSjSTb1Wno04FlRCtJZ1vD1vJqfS1HI_tDKFwH8avwDM8Qu-voJzJIWEGMv2vF-9KBAsFuengcJNrMxKoOeNrQHq5ELxpgemodcCi5xNkKuoL_Rz8c8-LwsUclLqPk2zb-Yed7rlhMOeQLkgqEdLWIVrA0jhzATYmsTeZEl1A"
	var jwtKeyfunc jwt.Keyfunc = func(t *jwt.Token) (interface{}, error) { return []byte(secret), nil }
	fmt.Println("parsing token")
	parsedToken, err := jwt.Parse(token, jwtKeyfunc)

	if err == nil {
		claims := parsedToken.Claims.(jwt.MapClaims)
		fmt.Println("getting id off claim")
		id := claims["id"].(string)

		fmt.Println("saving verified user with id " + id)
		saveError := dbclient.SaveVerifiedWsUser(connectionID, id)

		if saveError != nil {
			fmt.Println("save error " + saveError.Error())
			return Response{
				StatusCode: 400,
				Body:       "save error ",
			}, nil
		}

		websocketClient.Post(connectionID, "authenticated", nil)
		return Response{
			StatusCode: 200,
			Body:       "success",
		}, nil
	}
	fmt.Println("Error parsing token " + err.Error())
	return Response{
		StatusCode: 400,
		Body:       "invalid",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
