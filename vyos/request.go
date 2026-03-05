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

// Op is a type for the different operational modes of the VyOS API.
type Op string

const (
	// defaultUserAgent is the default user agent used by the VyOS API client.
	defaultUserAgent = "go-vyos"

	// Op constants
	OpShowConfig Op = "showConfig"
	OpShowValues Op = "returnValues"
	OpShow       Op = "show"
	OpSet        Op = "set"
	OpComment    Op = "comment"
	OpGenerate   Op = "generate"
	OpConfigure  Op = "configure"
	OpExists     Op = "exists"
	OpReset      Op = "reset"
	OpPowerOff   Op = "poweroff"
	OpReboot     Op = "reboot"
	OpAdd        Op = "add"
	OpDelete     Op = "delete"
	OpSave       Op = "save"
	OpLoad       Op = "load"
	OpConfirm    Op = "confirm"
)

func (c *Client) doRequest(ctx context.Context, req request, resp *response) error {
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

type request interface {
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

type pathRequest struct {
	URLPath string   `json:"-"`
	Op      Op       `json:"op"`
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
