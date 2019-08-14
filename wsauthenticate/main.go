package main

import (
	"context"
	"encoding/json"
	"fmt"
	"serverless-backgammon/dbclient"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dgrijalva/jwt-go"

	// todo better way to import this

	"github.com/naten3/serverless-backgammon/wsclient"
)

type Response events.APIGatewayProxyResponse

type wsPayload struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

func Handler(context context.Context, request events.APIGatewayWebsocketProxyRequest) (Response, error) {
	client := wsclient.New()
	connectionID := request.RequestContext.ConnectionID

	// TODO put in env variable
	secret := "Dba98iE002lTOA8YdQtYvdf2U52Eai7WT1sIoVTO-Q0r5KDdHNbIfBZS8P8Y-yFf6NunyqqFcB3HuvOivsEs-Zi4oka_FK4TbW52G9dSsxGoppciGEtUsTFgpKQYpQ7qyZE7ncvf39bWR0Y1RkP-yf2X2Ffeq7bv75vXE2TWhvZU6oSjSTb1Wno04FlRCtJZ1vD1vJqfS1HI_tDKFwH8avwDM8Qu-voJzJIWEGMv2vF-9KBAsFuengcJNrMxKoOeNrQHq5ELxpgemodcCi5xNkKuoL_Rz8c8-LwsUclLqPk2zb-Yed7rlhMOeQLkgqEdLWIVrA0jhzATYmsTeZEl1A"

	var payload wsPayload
	unmarshalErr := json.Unmarshal([]byte(request.Body), &payload)
	if unmarshalErr != nil {
		fmt.Println(unmarshalErr.Error())
	}
	token := payload.Data
	fmt.Println("token is " + token)

	client.Post(connectionID, "token is "+token)

	var jwtKeyfunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) { return []byte(secret), nil }
	parsedToken, err := jwt.Parse(token, jwtKeyfunc)

	if err != nil {
		claims := parsedToken.Claims.(jwt.MapClaims)
		id := claims["id"].(string)

		saveError := dbclient.SaveVerifiedWsUser(connectionID, id)

		if saveError != nil {
			return Response{
				StatusCode: 400,
				Body:       "save error " + saveError.Error(),
			}, nil
		}

		return Response{
			StatusCode: 200,
			Body:       "success",
		}, nil
	}

	return Response{
		StatusCode: 400,
		Body:       "invalid",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
