package main

import (
	"context"
	"os"

	snowsync "github.com/UKHomeOffice/snowsync/internal"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var sess *session.Session
var ddb *dynamodb.DynamoDB

func init() {
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	ddb = dynamodb.New(sess, &aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
}

func handler(ctx context.Context, rec snowsync.Record) error {
	return snowsync.NewSaver(ddb).Save(ctx, rec)
}

func main() {
	lambda.Start(handler)
}
