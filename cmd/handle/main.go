package main

import (
	"context"
	"os"

	snowsync "github.com/UKHomeOffice/snowsync/internal"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var sess *session.Session
var esqs *sqs.SQS

func init() {
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	esqs = sqs.New(sess, &aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
}
func handler(req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx := context.Background()
	return snowsync.NewHandler(esqs).Handle(ctx, req)
}

func main() {
	lambda.Start(handler)
}
