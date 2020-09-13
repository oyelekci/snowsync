package checker

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DBChecker is an abstraction (helpful for testing)
type DBChecker interface {
	GetItem(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
}

// NewChecker returns a new checker
func NewChecker(d DBChecker) *Checker {
	return &Checker{ddb: d}
}

// Checker is a record checker
type Checker struct {
	ddb DBChecker
}

// Check checks if a prior db record has external identifier
func (c *Checker) Check(payload string) (string, error) {

	// dynamically decode payload
	var dat map[string]interface{}
	err := json.Unmarshal([]byte(payload), &dat)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	iid := dat["issue_id"].(string)

	if iid == "" {
		return "", fmt.Errorf("no issue_id in payload")
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"issue_id": {
				S: aws.String(iid),
			},
		},
	}

	res, err := c.ddb.GetItem(input)
	if err != nil {
		return "", fmt.Errorf("failed to get item: %v", err)
	}

	// dynamically decode db item
	var itm map[string]interface{}
	err = dynamodbattribute.UnmarshalMap(res.Item, &itm)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal item: %v", err)
	}

	if itm["internal_identifier"] != nil {
		fmt.Printf("Issue id %v has a SNOW identifier: %v", iid, itm["internal_identifier"].(string))
		return itm["internal_identifier"].(string), nil
	}

	fmt.Printf("Issue id %v has no SNOW identifier", iid)
	return "", nil
}
