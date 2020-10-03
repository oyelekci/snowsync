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

// Save saves a ticket
func (s *Saver) Save(ctx context.Context, st string) error {

	log.Printf("debug - into saver: %+v", st)

	rec, err := json.Marshal(st)
	if err != nil {
		return fmt.Errorf("failed to marshal db record: %v", err)
	}

	item, err := dynamodbattribute.MarshalMap(rec)
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
