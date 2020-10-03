package snowsync

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Incident is a type of case record
type Incident struct {
	Cluster     string `json:"cluster,omitempty"`
	Component   string `json:"component,omitempty"`
	Description string `json:"description,omitempty"`
	Identifier  string `json:"external_identifier,omitempty"`
	IssueID     string `json:"issue_id,omitempty"`
	Priority    string `json:"priority,omitempty"`
	Status      string `json:"status,omitempty"`
	Summary     string `json:"summary,omitempty"`
}

// incidentUpdate is an event
type incidentUpdate struct {
	incident *Incident
	sqs      Messenger
}

// New initialises a new Incident
func New() *Incident {
	return &Incident{}
}

// execute publishes an incident update to SQS
func (i *incidentUpdate) publish(ctx context.Context) (*Response, error) {

	sm, err := json.Marshal(i.incident)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal SQS payload: %s", err)
	}

	in := sqs.SendMessageInput{
		MessageBody: aws.String(string(sm)),
		QueueUrl:    aws.String(os.Getenv("QUEUE_URL")),
	}

	var res *Response

	if _, err := i.sqs.SendMessageWithContext(ctx, &in); err != nil {
		res = &Response{
			ResponseType: "failure",
			Text:         "could not publish to SQS",
		}
		return res, fmt.Errorf("failed to publish incident: %s", err)
	}

	res = &Response{
		ResponseType: "success",
		Text:         "inbound payload parsed and published to SQS",
	}

	return res, nil

}
