package snowsync

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Create makes an outbound create request and returns a SNOW identifier
func Create(e Envelope) (string, error) {

	fmt.Printf("debug - into creator: %v", e)

	surl, err := url.Parse(os.Getenv("SNOW_URL"))
	if err != nil {
		return "", fmt.Errorf("no SNOW URL provided: %v", err)
	}

	c := &Client{
		BaseURL:    surl,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	out, err := json.Marshal(e)
	if err != nil {
		return "", fmt.Errorf("failed to marshal snow payload: %v", err)
	}

	req, err := c.newRequest("", out)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}

	res, err := c.do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call SNOW: %v", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	fmt.Printf("sent request, SNOW replied with: %v", string(body))

	// dynamically decode SNOW response
	var dat map[string]interface{}
	err = json.Unmarshal(body, &dat)
	if err != nil {
		return "", err
	}

	// check for external_identifier
	rts := dat["result"].(map[string]interface{})
	ini := rts["external_identifier"].(string)
	if ini == "" {
		return "", fmt.Errorf("request failed, SNOW did not return an identifier")
	}

	rlg := rts["log"].(string)
	if strings.Contains(rlg, "Inserting") {
		fmt.Printf("SNOW returned new identifier: %v", ini)
		return ini, nil
	}
	fmt.Printf("unexpected SNOW response: %v", string(body))
	return "", fmt.Errorf("unexpected SNOW response")
}
