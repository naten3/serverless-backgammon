package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	secret := "Dba98iE002lTOA8YdQtYvdf2U52Eai7WT1sIoVTO-Q0r5KDdHNbIfBZS8P8Y-yFf6NunyqqFcB3HuvOivsEs-Zi4oka_FK4TbW52G9dSsxGoppciGEtUsTFgpKQYpQ7qyZE7ncvf39bWR0Y1RkP-yf2X2Ffeq7bv75vXE2TWhvZU6oSjSTb1Wno04FlRCtJZ1vD1vJqfS1HI_tDKFwH8avwDM8Qu-voJzJIWEGMv2vF-9KBAsFuengcJNrMxKoOeNrQHq5ELxpgemodcCi5xNkKuoL_Rz8c8-LwsUclLqPk2zb-Yed7rlhMOeQLkgqEdLWIVrA0jhzATYmsTeZEl1A"

	id := uuid.NewV4()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(),
			StatusCode: 400}, nil
	}

	jsonResult, _ := json.Marshal(map[string]string{
		"token": tokenString,
	})

	return events.APIGatewayProxyResponse{Body: string(jsonResult),
		StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
