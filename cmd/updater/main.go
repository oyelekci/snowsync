package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	snowsync "github.com/UKHomeOffice/snowsync/internal"
)

func handler(payload string) ([]byte, error) {
	return snowsync.Update(payload)
}

func main() {
	lambda.Start(handler)
}
