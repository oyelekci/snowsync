package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/UKHomeOffice/snowsync/internal/client"
	"github.com/UKHomeOffice/snowsync/internal/updater"
)

func handler(e client.Envelope) (string, error) {
	return updater.Update(e)
}

func main() {
	lambda.Start(handler)
}
