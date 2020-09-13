package snowsync

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DBUpdater is an abstraction (helpful for testing)
type DBUpdater interface {
	UpdateItem(*dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
}

// Updater is a ticket updater
type Updater struct {
	ddb DBUpdater
}

// NewUpdater returns a new updater
func NewUpdater(u DBUpdater) *Updater {
	return &Updater{ddb: u}
}

// DBUpdate updates a db record
func (u *Updater) DBUpdate(payload string) error {

	// dynamically decode payload
	var dat map[string]interface{}
	err := json.Unmarshal([]byte(payload), &dat)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	// item, err := dynamodbattribute.MarshalMap(i)
	// if err != nil {
	// 	return fmt.Errorf("failed to marshal db record: %s", err)
	// }

	input := &dynamodb.UpdateItemInput{
		TableName:        aws.String(os.Getenv("TABLE_NAME")),
		UpdateExpression: aws.String("SET external_identifier = :cid"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":cid": {
				//			S: aws.String(IntIdent),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"issue_id": {
				//		S: aws.String(SupplierRef),
			},
		},
	}

	_, err = u.ddb.UpdateItem(input)
	if err != nil {
		return fmt.Errorf("failed to put to db: %v", err)
	}
	return nil
}
