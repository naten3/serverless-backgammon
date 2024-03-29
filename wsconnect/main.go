package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, handler interface{}) {
	fmt.Println("open function ran")
}

func main() {
	lambda.Start(Handler)
}
