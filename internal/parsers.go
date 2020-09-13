package snowsync

import (
	"fmt"
	"os"

	"github.com/tidwall/gjson"
)

// parseIncident gets some values from inbound request
func parseIncident(input string) (*incidentUpdate, error) {

	err := checkVars(input)
	if err != nil {
		return nil, err
	}

	i := New()
	i.Cluster = gjson.Get(input, os.Getenv("CLUSTER_FIELD")).Str
	i.Component = gjson.Get(input, os.Getenv("COMPONENT_FIELD")).Str
	i.Description = gjson.Get(input, os.Getenv("DESCRIPTION_FIELD")).Str
	i.IssueID = gjson.Get(input, os.Getenv("ISSUE_ID_FIELD")).Str
	i.Priority = gjson.Get(input, os.Getenv("PRIORITY_FIELD")).Str
	i.Status = gjson.Get(input, os.Getenv("STATUS_FIELD")).Str
	i.Summary = gjson.Get(input, os.Getenv("SUMMARY_FIELD")).Str
	u := incidentUpdate{incident: i}

	fmt.Printf("parsed incident: %v, status: %v\n", i.IssueID, i.Status)
	return &u, nil
}

// checkVars checks incoming payload has necessary fields
func checkVars(input string) error {

	vars := []string{
		"COMPONENT_FIELD",
		"DESCRIPTION_FIELD",
		"ISSUE_ID_FIELD",
		"PRIORITY_FIELD",
		"STATUS_FIELD",
		"SUMMARY_FIELD",
	}

	for _, v := range vars {
		field, ok := os.LookupEnv(v)
		if !ok {
			return fmt.Errorf("missing environment variable")
		}
		value := gjson.Get(input, field)
		if !value.Exists() {
			return fmt.Errorf("missing value in payload")
		}
	}
	return nil
}
