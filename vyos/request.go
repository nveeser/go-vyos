package vyos

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// HTTPError is the error returned from the HTTP request.
type HTTPError struct {
	Code int
	Body string
}

func (h *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %s: %s", http.StatusText(h.Code), h.Body)
}

// op is a type for the different operational modes of the VyOS API.
type op string

const (
	// defaultUserAgent is the default user agent used by the VyOS API client.
	defaultUserAgent = "go-vyos"

	// op constants
	OpShowConfig op = "showConfig"
	OpShowValues op = "returnValues"
	OpShow       op = "show"
	OpSet        op = "set"
	OpComment    op = "comment"
	OpGenerate   op = "generate"
	OpConfigure  op = "configure"
	OpExists     op = "exists"
	OpReset      op = "reset"
	OpPowerOff   op = "poweroff"
	OpReboot     op = "reboot"
	OpAdd        op = "add"
	OpDelete     op = "delete"
	OpSave       op = "save"
	OpLoad       op = "load"
	OpConfirm    op = "confirm"
)

func (c *Client) do(ctx context.Context, req request, resp *response) error {
	httpReq, err := c.httpReq(req)
	if err != nil {
		return fmt.Errorf("error building http.Request %w", err)
	}
	httpResp, err := c.httpClient.Do(httpReq.WithContext(ctx))
	if err != nil {
		return err
	}

	defer httpResp.Body.Close()
	if httpResp.StatusCode != 200 {
		body, err := io.ReadAll(httpResp.Body)
		if err != nil {
			log.Printf("error reading failed response body:%s", err)
		}
		return &HTTPError{httpResp.StatusCode, string(body)}
	}
	// Decode the response body into the provided interface.
	err = json.NewDecoder(httpResp.Body).Decode(resp)
	if err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("vyos error(success=false): %s", resp.Error)
	}
	return nil
}

// httpReq creates the HTTP request from the specified request by creating a POST
// where the data is provided by the requestPayload method. If the instance
// supports the custom
func (c *Client) httpReq(req request) (*http.Request, error) {
	if pr, ok := req.(customRequest); ok {
		return pr.httpRequest(*c.baseURL, c.token)
	}

	// The method must always be POST. The VyOS API only supports POST requests.
	method := http.MethodPost

	urlPath, data := req.requestPayload()
	u := *c.baseURL
	u.Path = urlPath

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marhalling request %w", err)
	}
	body := url.Values{}
	body.Set("data", string(jsonData))
	body.Set("key", c.token)

	httpReq, err := http.NewRequest(method, u.String(), strings.NewReader(body.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error building http request %w", err)
	}
	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return httpReq, nil
}

// request is the interface consumed by Client.httpReq to build
// the http.Request instance for a given command.
type request interface {
	// requestPayload should return the URL path (within the base URL)
	// and the payload to pass to the HTTP POST.
	requestPayload() (path string, payload any)
}

type customRequest interface {
	httpRequest(base url.URL, token string) (*http.Request, error)
}

// response represents a raw response from the VyOS API.
type response struct {
	Success bool   `json:"success,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

var _ request = (*pathRequest)(nil)

// pathRequest is the common implementation of the request interface.
type pathRequest struct {
	URLPath string   `json:"-"`
	Op      op       `json:"op"`
	Path    []string `json:"path"`
}

func (r *pathRequest) requestPayload() (path string, payload any) {
	return r.URLPath, r
}

func parsePath(path string, nonNil bool) []string {
	if path == "" && nonNil {
		return []string{}
	}
	return strings.Split(path, " ")
}
