package dbupdater

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

func (md *mockDynamoDB) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	output := new(dynamodb.UpdateItemOutput)

	return output, md.err
}

func TestDBUpdate(t *testing.T) {

	tt := []struct {
		name          string
		issueID       string
		commentID     string
		commentAuthor string
		commentBody   string
		err           string
	}{
		{name: "happy", issueID: "abc-123", commentID: "1", commentAuthor: "bob", commentBody: "first comment"},
		{name: "unhappy", err: "failed to update on db:"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			updater := NewUpdater(&mockDynamoDB{})
			dataSlice := make([]string, 1)
			dataSlice[0] = tc.commentID + tc.commentAuthor + tc.commentBody
			var interfaceSlice []interface{} = make([]interface{}, len(dataSlice))
			for i, d := range dataSlice {
				interfaceSlice[i] = d
			}

			// todo: fix unhappy path
			if tc.err != "" {
				in := map[string]interface{}{"issue_id": "", "comments": interfaceSlice}
				err := updater.DBUpdate(in)
				if err != nil {
					if msg := err.Error(); !strings.Contains(msg, tc.err) {
						t.Errorf("expected error %q, got: %q", tc.err, msg)
					}
					return
				}
			}
			in := map[string]interface{}{"issue_id": tc.issueID, "comments": interfaceSlice}
			err := updater.DBUpdate(in)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
