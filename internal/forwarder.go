package snowsync

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Client is a HTTP client
type Client struct {
	BaseURL    *url.URL
	httpClient *http.Client
}

// newRequest creates a HTTP request
func (c *Client) newRequest(path string, body []byte) (*http.Request, error) {

	p, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	u := c.BaseURL.ResolveReference(p)

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	cerr := errors.New("missing credentials")
	admin, ok := os.LookupEnv("SNOW_USER")
	if !ok {
		return nil, cerr
	}
	password, ok := os.LookupEnv("SNOW_PASS")
	if !ok {
		return nil, cerr
	}
	req.SetBasicAuth(admin, password)

	return req, nil
}

// do makes a HTTP request
func (c *Client) do(req *http.Request) (*http.Response, error) {

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, err
}

// Forward makes the outbound request
func Forward(ctx context.Context, payload string) error {

	surl, err := url.Parse(os.Getenv("SNOW_URL"))
	if err != nil {
		return err
	}

	c := &Client{
		BaseURL:    surl,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	req, err := c.newRequest("", []byte(payload))
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}

	log.Printf("snow response code: %v", resp.StatusCode)

	return nil
}
