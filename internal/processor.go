package snowsync

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// LambdaAPI invokes other lambdas
type LambdaAPI interface {
	Invoke(*lambda.InvokeInput) (*lambda.InvokeOutput, error)
}

// NewSQSProcessor returns a new SQSProcessor
func NewSQSProcessor(l LambdaAPI) *SQSProcessor {
	return &SQSProcessor{lam: l}
}

// SQSProcessor processes messages from queue
type SQSProcessor struct {
	lam LambdaAPI
}

// Process processes individual messages
func (s *SQSProcessor) Process(ctx context.Context, event *events.SQSEvent) error {

	for _, message := range event.Records {
		log.Printf("Processing message %s | %s", message.MessageId, message.Body)

		log.Printf("debug - message.Body %v", message.Body)

		p := struct {
			Body string `json:"body"`
		}{
			Body: message.Body,
		}

		payload, err := json.Marshal(p)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %v", err)
		}

		fmt.Printf("debug - payload %s", string(payload))

		input := &lambda.InvokeInput{
			FunctionName:   aws.String(os.Getenv("SAVER_LAMBDA")),
			InvocationType: aws.String("Event"),
			Payload:        payload,
		}

		err = input.Validate()
		if err != nil {
			return fmt.Errorf("failed to validate invocation input: %v", err)
		}

		_, err = s.lam.Invoke(input)
		if err != nil {
			return fmt.Errorf("failed to invoke saver function: %v", err)
		}

		// var dat map[string]interface{}
		// json.Unmarshal(resp.Payload, &dat)
		// log.Printf("debug - response body: %v", dat["body"])
	}
	return nil
}
