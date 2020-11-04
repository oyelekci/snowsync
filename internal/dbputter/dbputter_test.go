package dbputter

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mockDynamoDB struct {
	dynamodbiface.DynamoDBAPI
	err error
}

func (md *mockDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	output := new(dynamodb.PutItemOutput)
	return output, md.err
}

func TestDBPut(t *testing.T) {

	tt := []struct {
		name        string
		cluster     string
		component   string
		description string
		issueID     string
		priority    string
		status      string
		summary     string
		err         string
	}{
		{name: "happy", issueID: "abc-123", status: "investigating",
			summary: "keycloak down", description: "lorem ipsum", cluster: "dev",
			priority: "P2", component: "keycloak"},
		{name: "unhappy", issueID: "", err: "failed to put to db:"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			putter := NewPutter(&mockDynamoDB{})

			// todo: fix unhappy path
			if tc.err != "" {
				in := map[string]interface{}{"issueID": tc.issueID}

				err := putter.DBPut(in)
				if err != nil {
					if msg := err.Error(); !strings.Contains(msg, tc.err) {
						t.Errorf("expected error %q, got: %q", tc.err, msg)
					}
				}
			}

			in := map[string]interface{}{"issueID": tc.issueID, "status": tc.status,
				"summary": tc.summary, "description": tc.description, "cluster": tc.cluster,
				"priority": tc.priority, "component": tc.component}

			err := putter.DBPut(in)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
