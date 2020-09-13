package snowsync

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// Invoker invokes another lambda
type Invoker interface {
	Invoke(*lambda.InvokeInput) (*lambda.InvokeOutput, error)
}

// Processor processes messages from queue
type Processor struct {
	inv Invoker
}

// Envelope is the JSON expected by SNOW
type Envelope struct {
	MsgID   string `json:"messageid,omitempty"`
	ExtID   string `json:"external_identifier,omitempty"`
	Payload string `json:"payload,omitempty"`
}

// // Item is a db record
// type Item struct {
// 	Cluster     string `json:"cluster,omitempty"`
// 	Component   string `json:"component,omitempty"`
// 	Description string `json:"description,omitempty"`
// 	Ends        string `json:"endTime,omitempty"`
// 	Identifier  string `json:"external_identifier,omitempty"`
// 	IssueID     string `json:"issue_id,omitempty"`
// 	Priority    string `json:"priority,omitempty"`
// 	Status      string `json:"status,omitempty"`
// 	Summary     string `json:"summary,omitempty"`
// 	SupplierRef string `json:"supplierRef,omitempty"`
// 	Starts      string `json:"startTime,omitempty"`
// 	Title       string `json:"title,omitempty"`
// }

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
			e := Envelope{
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

			// call dbupdater lambda
			itemUpdate := &lambda.InvokeInput{
				FunctionName: aws.String(os.Getenv("DBUPDATER_LAMBDA")),
				Payload:      payload,
			}

			_, err = p.inv.Invoke(itemUpdate)
			if err != nil {
				return fmt.Errorf("failed to invoke dbupdater: %v", err)
			}

			return nil
		}
		// call creator lambda, expect external_identifier in return
		fmt.Println("No SNOW id found, creating a new record...")

		e := Envelope{
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

		fmt.Printf("debug - creator payload: %v", string(newTicket))

		eid, err := p.inv.Invoke(input)
		if err != nil {
			return fmt.Errorf("failed to invoke creator: %v", err)
		}

		// call dbputter lambda, adding external_identifier to payload
		var dat map[string]interface{}
		err = json.Unmarshal([]byte(message.Body), &dat)
		if err != nil {
			return fmt.Errorf("failed to unmarshal: %v", err)
		}
		dat["external_identifier"] = string(eid.Payload)

		newItem, err := json.Marshal(dat)
		if err != nil {
			return fmt.Errorf("failed to marshal dbputter payload: %v", err)
		}

		input = &lambda.InvokeInput{
			FunctionName: aws.String(os.Getenv("DBPUTTER_LAMBDA")),
			Payload:      newItem,
		}

		fmt.Printf("debug - dbputter payload: %v", string(newItem))

		_, err = p.inv.Invoke(input)
		if err != nil {
			return fmt.Errorf("failed to invoke dbputter: %v", err)
		}
	}
	return nil
}
