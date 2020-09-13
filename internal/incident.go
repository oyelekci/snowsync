package snowsync

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Incident is a type of ticket
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

// publish writes an incident to SQS
func (i *incidentUpdate) publish() error {

	sm, err := json.Marshal(i.incident)
	if err != nil {
		return fmt.Errorf("failed to marshal SQS payload: %s", err)
	}

	in := sqs.SendMessageInput{
		MessageBody: aws.String(string(sm)),
		QueueUrl:    aws.String(os.Getenv("QUEUE_URL")),
	}

	if _, err := i.sqs.SendMessage(&in); err != nil {
		return fmt.Errorf("failed to publish incident: %s", err)
	}

	return nil

}
