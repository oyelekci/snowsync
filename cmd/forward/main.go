package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	snowsync "github.com/UKHomeOffice/snowsync/internal"
)

func init() {
}

func handler(ctx context.Context, payload string) error {
	return snowsync.Forward(ctx, payload)
}

func main() {
	lambda.Start(handler)
}
