package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	snowsync "github.com/UKHomeOffice/snowsync/internal"
)

func handler(e snowsync.Envelope) (string, error) {
	return snowsync.Create(e)
}

func main() {
	lambda.Start(handler)
}
