package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/UKHomeOffice/snowsync/internal/client"
	"github.com/UKHomeOffice/snowsync/internal/creator"
)

func handler(e client.Envelope) (string, error) {
	return creator.Create(e)
}

func main() {
	lambda.Start(handler)
}
