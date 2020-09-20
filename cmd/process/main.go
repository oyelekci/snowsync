package main

import (
	"context"
	"os"

	snowsync "github.com/UKHomeOffice/snowsync/internal"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	service "github.com/aws/aws-sdk-go/service/lambda"
)

var sess *session.Session
var svc *service.Lambda

func init() {
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc = service.New(sess, &aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
}

func handler(ctx context.Context, sqsEvent *events.SQSEvent) error {
	return snowsync.NewSQSProcessor(svc).Process(ctx, sqsEvent)
}

func main() {
	lambda.Start(handler)
}
