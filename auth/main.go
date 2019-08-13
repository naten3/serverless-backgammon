package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dgrijalva/jwt-go"
)

// Help function to generate an IAM policy
func generatePolicy(principalId, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalId}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}

	// Optional output with custom properties of the String, Number or Boolean type.
	authResponse.Context = map[string]interface{}{
		"stringKey":  "stringval",
		"numberKey":  123,
		"booleanKey": true,
	}
	return authResponse
}

func handleRequest(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	// TODO put in env variable
	secret := "Dba98iE002lTOA8YdQtYvdf2U52Eai7WT1sIoVTO-Q0r5KDdHNbIfBZS8P8Y-yFf6NunyqqFcB3HuvOivsEs-Zi4oka_FK4TbW52G9dSsxGoppciGEtUsTFgpKQYpQ7qyZE7ncvf39bWR0Y1RkP-yf2X2Ffeq7bv75vXE2TWhvZU6oSjSTb1Wno04FlRCtJZ1vD1vJqfS1HI_tDKFwH8avwDM8Qu-voJzJIWEGMv2vF-9KBAsFuengcJNrMxKoOeNrQHq5ELxpgemodcCi5xNkKuoL_Rz8c8-LwsUclLqPk2zb-Yed7rlhMOeQLkgqEdLWIVrA0jhzATYmsTeZEl1A"

	token := event.AuthorizationToken

	var jwtKeyfunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) { return []byte(secret), nil }
	parsedToken, err := jwt.Parse(token, jwtKeyfunc)

	if err == nil {
		claims := parsedToken.Claims.(jwt.MapClaims)
		id := claims["id"].(string)
		authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: id}
		authResponse.Context = map[string]interface{}{
			"id": id,
		}
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   "Allow",
					Resource: []string{event.MethodArn},
				},
			},
		}
		return authResponse, nil
	}
	fmt.Println(err.Error())
	return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")

}

func main() {
	lambda.Start(handleRequest)
}
