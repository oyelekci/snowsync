package snowsync

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DBPuter is the minimal interface needed to store a ticket
type DBPuter interface {
	PutItemWithContext(aws.Context, *dynamodb.PutItemInput, ...request.Option) (*dynamodb.PutItemOutput, error)
}

// NewSaver returns a new saver
func NewSaver(c DBPuter) *Saver {
	return &Saver{ddb: c}
}

// Saver is a ticket saver
type Saver struct {
	ddb DBPuter
}

// Record is a db record
type Record struct {
	Body []byte `json:"body,omitempty"`
}

// Save saves a ticket
func (s *Saver) Save(ctx context.Context, rec Record) error {

	//log.Printf("debug - request: %v", req)

	var dat map[string]interface{}
	err := json.Unmarshal(rec.Body, &dat)
	if err != nil {
		return fmt.Errorf("failed to unmarshal invocation input: %s", err)
	}
	log.Printf("debug - input body: %v", dat["body"])

	item, err := dynamodbattribute.MarshalMap(dat["body"])
	if err != nil {
		return fmt.Errorf("failed to marshal db record: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	_, err = s.ddb.PutItemWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put to db: %v", err)
	}
	return nil
}
