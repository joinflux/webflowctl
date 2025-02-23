package internal

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client represent an abstraction of the Webflow http calls
type Client struct {
	baseUrl    string
	token      string
	httpClient *http.Client
}

// NewClient creates a webflow Client
func NewClient(token string) *Client {
	return &Client{
		baseUrl:    "api.webflow.com/v2",
		token:      token,
		httpClient: &http.Client{},
	}
}

func (c Client) buildUrl(endpoint []string) string {
	return fmt.Sprintf("%s://%s/%s", "https", c.baseUrl, strings.Join(endpoint, "/"))
}

func (c Client) do(method string, endpoint []string, payload io.Reader) ([]byte, error) {
	request, err := http.NewRequest(method, c.buildUrl(endpoint), payload)
	if err != nil {
		return nil, err
	}

	request.Header.Add("authorization", "Bearer "+c.token)
	request.Header.Add("accept", "application/json")
	request.Header.Add("content-type", "application/json")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request Failed: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Dump body to terminal for troubleshooting
	// fmt.Println(string(body))

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusNoContent {
		return nil, fmt.Errorf("request failed: %s\n%s", resp.Status, string(body))
	}

	return body, nil
}

// Post sends a post request to Webflow given an endpoint and a payload. It will error in the following cases:
// 1. There is a problem creating the request
// 2. There is a problem sending the request
// 3. There is a problem reading the response body
// 4. The response is not a 200 status
func (c Client) Post(endpoint []string, payload io.Reader) ([]byte, error) {
	return c.do("POST", endpoint, payload)
}

// Get sends a get request to Webflow given an endpoint. It will error in the following cases:
// 1. There is a problem creating the request
// 2. There is a problem sending the request
// 3. There is a problem reading the response body
// 4. The response is not a 200 status
func (c Client) Get(endpoint []string) ([]byte, error) {
	return c.do("GET", endpoint, nil)
}

// Delete sends a delete request to Webflow given an endpoint. It will error in the following cases:
// 1. There is a problem creating the request
// 2. There is a problem sending the request
// 3. There is a problem reading the response body
// 4. The response is not a 200 status
func (c Client) Delete(endpoint []string) ([]byte, error) {
	return c.do("DELETE", endpoint, nil)
}
