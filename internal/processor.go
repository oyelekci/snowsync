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

// Invoker invokes another lambda
type Invoker interface {
	Invoke(*lambda.InvokeInput) (*lambda.InvokeOutput, error)
}

// SQSProcessor processes messages from queue
type SQSProcessor struct {
	inv Invoker
}

// NewSQSProcessor returns a new SQSProcessor
func NewSQSProcessor(i Invoker) *SQSProcessor {
	return &SQSProcessor{inv: i}
}

// Process processes individual messages
func (s *SQSProcessor) Process(ctx context.Context, event *events.SQSEvent) error {

	for _, message := range event.Records {
		log.Printf("Processing message %s | %s", message.MessageId, message.Body)

		log.Printf("debug - message.Body %v", message.Body)

		payload, err := json.Marshal(message.Body)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %v", err)
		}

		fmt.Printf("debug - payload %s", string(payload))

		input := &lambda.InvokeInput{
			FunctionName: aws.String(os.Getenv("SAVER_LAMBDA")),
			Payload:      payload,
		}

		err = input.Validate()
		if err != nil {
			return fmt.Errorf("failed to validate invocation input: %v", err)
		}

		_, err = s.inv.Invoke(input)
		if err != nil {
			return fmt.Errorf("failed to invoke saver function: %v", err)
		}

	}
	return nil
}
