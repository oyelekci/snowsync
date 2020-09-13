package processor

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"

	"github.com/UKHomeOffice/snowsync/internal/client"
)

// Invoker invokes another lambda
type Invoker interface {
	Invoke(*lambda.InvokeInput) (*lambda.InvokeOutput, error)
}

// Processor processes messages from queue
type Processor struct {
	inv Invoker
}

// NewProcessor returns a new Processor
func NewProcessor(i Invoker) *Processor {
	return &Processor{inv: i}
}

// Process processes individual SQS messages
func (p *Processor) Process(event *events.SQSEvent) error {

	for _, message := range event.Records {
		fmt.Printf("Processing message %s | %s", message.MessageId, message.Body)

		payload, err := json.Marshal(message.Body)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %v", err)
		}

		// call checker lambda, expect external_identifier if it exists
		input := &lambda.InvokeInput{
			FunctionName: aws.String(os.Getenv("CHECKER_LAMBDA")),
			Payload:      payload,
		}

		res, err := p.inv.Invoke(input)
		if err != nil {
			return fmt.Errorf("failed to invoke checker function: %v", err)
		}

		// note: invocation output requires the backticks
		if string(res.Payload) != `""` {

			// call updater lambda
			fmt.Println("A SNOW id exists, updating the existing record...")
			e := client.Envelope{
				MsgID:   "HO_SIAM_IN_REST_INC_UPDATE_JSON_ACP_Incident_Update",
				ExtID:   string(res.Payload),
				Payload: string(payload),
			}

			ticketUpdate, err := json.Marshal(e)
			if err != nil {
				return fmt.Errorf("failed to marshal updater payload: %v", err)
			}

			input = &lambda.InvokeInput{
				FunctionName: aws.String(os.Getenv("UPDATER_LAMBDA")),
				Payload:      ticketUpdate,
			}

			_, err = p.inv.Invoke(input)
			if err != nil {
				return fmt.Errorf("failed to invoke updater %v", err)
			}

			var dat map[string]interface{}
			err = json.Unmarshal([]byte(message.Body), &dat)
			if err != nil {
				return fmt.Errorf("failed to unmarshal: %v", err)
			}

			itemUpdate, err := json.Marshal(dat)
			if err != nil {
				return fmt.Errorf("failed to marshal dbupdater payload: %v", err)
			}

			// call dbupdater lambda
			input := &lambda.InvokeInput{
				FunctionName: aws.String(os.Getenv("DBUPDATER_LAMBDA")),
				Payload:      itemUpdate,
			}

			_, err = p.inv.Invoke(input)
			if err != nil {
				return fmt.Errorf("failed to invoke dbupdater: %v", err)
			}

			return nil
		}
		// call creator lambda, expect SNOW identifier in return
		fmt.Println("No SNOW id found, creating a new record...")

		e := client.Envelope{
			MsgID:   "HO_SIAM_IN_REST_INC_POST_JSON_ACP_Incident_Create",
			Payload: string(payload),
		}

		newTicket, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("failed to marshal creator payload: %v", err)
		}

		input = &lambda.InvokeInput{
			FunctionName: aws.String(os.Getenv("CREATOR_LAMBDA")),
			Payload:      newTicket,
		}

		output, err := p.inv.Invoke(input)
		if err != nil {
			return fmt.Errorf("failed to invoke creator: %v", err)
		}

		// call dbputter lambda, adding SNOW identifier to payload
		var dat map[string]interface{}
		err = json.Unmarshal([]byte(message.Body), &dat)
		if err != nil {
			return fmt.Errorf("failed to unmarshal: %v", err)
		}
		sid := strings.Trim(string(output.Payload), `", \`)
		dat["internal_identifier"] = sid

		newItem, err := json.Marshal(dat)
		if err != nil {
			return fmt.Errorf("failed to marshal dbputter payload: %v", err)
		}

		input = &lambda.InvokeInput{
			FunctionName: aws.String(os.Getenv("DBPUTTER_LAMBDA")),
			Payload:      newItem,
		}

		_, err = p.inv.Invoke(input)
		if err != nil {
			return fmt.Errorf("failed to invoke dbputter: %v", err)
		}
	}
	return nil
}
