package checker

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mockDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	err error
}

func (md *mockDynamoDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {

	output := new(dynamodb.GetItemOutput)

	existing := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"issue_id": {
				S: aws.String("abc-124"),
			},
		},
		TableName: aws.String(""),
	}

	if input.Key == nil {
		return nil, md.err
	} else if input.String() == existing.String() {
		return output, md.err
	}
	output = &dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"internal_identifier": {
				S: aws.String("inc-123"),
			},
		},
	}
	return output, md.err
}

func TestCheck(t *testing.T) {

	tt := []struct {
		name    string
		issueID string
		err     string
	}{
		{name: "happy_old", issueID: "abc-123"},
		{name: "happy_new", issueID: "abc-124"},
		{name: "unhappy", err: "no issue_id in payload"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			checker := NewChecker(&mockDynamoDB{})
			m := map[string]string{"issue_id": tc.issueID}

			in, err := json.Marshal(m)
			if err != nil {
				t.Fatalf("could not marshal test payload: %v", err)
			}

			_, err = checker.Check(string(in))

			if tc.err == "" {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
			if err != nil {
				if msg := err.Error(); !strings.Contains(msg, tc.err) {
					t.Errorf("expected error %q, got: %q", tc.err, msg)
				}
				return
			}
		})

	}
}
