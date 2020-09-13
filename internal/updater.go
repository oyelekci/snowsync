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

// Update makes an outbound update request
func Update(payload string) ([]byte, error) {

	surl, err := url.Parse(os.Getenv("SNOW_URL"))
	if err != nil {
		return nil, err
	}

	c := &Client{
		BaseURL:    surl,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	req, err := c.newRequest("", []byte(payload))
	if err != nil {
		return nil, err
	}

	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("sent request, SNOW replied with: %v", string(body))

	// dynamically decode SNOW response
	var dat map[string]interface{}
	err = json.Unmarshal(body, &dat)
	if err != nil {
		return nil, err
	}

	// check for external_identifier
	rts := dat["result"].(map[string]interface{})
	ini := rts["external_identifier"].(string)
	if ini == "" {
		return nil, fmt.Errorf("request failed, SNOW did not return an identifier")
	}

	rlg := rts["log"].(string)
	if strings.Contains(rlg, "Updating") {
		fmt.Printf("SNOW updated identifier: %v", ini)
		return body, nil
	}
	fmt.Printf("unexpected SNOW response: %v", string(body))
	return nil, fmt.Errorf("unexpected SNOW response")

	//return body, nil
}
