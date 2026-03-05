package vyos

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// OpMode provides the set of methods which map to the operation mode in VyOS
type OpMode Client

func (c *OpMode) do(ctx context.Context, req request, resp *response) error {
	return (*Client)(c).do(ctx, req, resp)
}

type InfoRequest struct {
	Version  bool
	Hostname bool
}

var _ request = (*InfoRequest)(nil)
var _ customRequest = (*InfoRequest)(nil)

func (r InfoRequest) requestPayload() (path string, payload any) {
	panic("InfoRequest uses httpRequest()")
}
func (r InfoRequest) httpRequest(u url.URL, token string) (*http.Request, error) {
	v := url.Values{}
	if r.Version {
		v.Set("version", "1")
	}
	if r.Hostname {
		v.Set("hostname", "1")
	}
	u.Path = "/info"
	u.RawQuery = v.Encode()
	return http.NewRequest("GET", u.String(), nil)
}

// InfoResponse is the response data from calling Info
type InfoResponse struct {
	Version  string `json:"version,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	Banner   string `json:"banner,omitempty"`
}

func (c *OpMode) Info(ctx context.Context, req InfoRequest) (*InfoResponse, error) {
	resp := &response{Data: &InfoResponse{}}
	err := c.do(ctx, req, resp)
	if err != nil {
		return nil, fmt.Errorf("error HTTP client: %w", err)
	}
	if !resp.Success {
		return nil, fmt.Errorf(resp.Error)
	}
	return resp.Data.(*InfoResponse), nil
}

func (c *OpMode) Show(ctx context.Context, path string) (string, error) {
	req := &pathRequest{
		URLPath: "/show",
		Op:      OpShow,
		Path:    parsePath(path, false),
	}
	resp := &response{}
	err := c.do(ctx, req, resp)
	if err != nil {
		return "", err
	}
	return resp.Data.(string), nil
}

func (c *OpMode) Generate(ctx context.Context, path string) (string, error) {
	req := &pathRequest{
		URLPath: "/generate",
		Op:      OpGenerate,
		Path:    parsePath(path, false),
	}
	resp := &response{}
	err := c.do(ctx, req, resp)
	if err != nil {
		return "", err
	}
	return resp.Data.(string), nil
}

type AddImageRequest struct {
	URL string
}

func (s *AddImageRequest) requestPayload() (path string, payload any) {
	return "/image", &imagePayload{Op: OpAdd, URL: s.URL}
}

type imagePayload struct {
	Op   op     `json:"op"`
	URL  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}

func (c *OpMode) AddImage(ctx context.Context, url string) (string, error) {
	req := &AddImageRequest{url}
	resp := &response{}
	err := c.do(ctx, req, resp)
	if err != nil {
		return "", err
	}
	return resp.Data.(string), nil
}

type DeleteImageRequest struct {
	Name string
}

func (s *DeleteImageRequest) requestPayload() (path string, payload any) {
	return "/image", &imagePayload{Op: OpDelete, Name: s.Name}
}

func (c *OpMode) DeleteImage(ctx context.Context, name string) (string, error) {
	req := &DeleteImageRequest{name}
	resp := &response{}
	err := c.do(ctx, req, resp)
	if err != nil {
		return "", err
	}
	return resp.Data.(string), nil
}

// Reset resets the object in the specified path
func (c *OpMode) Reset(ctx context.Context, path string) error {
	req := &pathRequest{
		URLPath: "/reset",
		Op:      OpReset,
		Path:    parsePath(path, true),
	}
	resp := &response{}
	return c.do(ctx, req, resp)
}

// PowerOff is a helper function to power off the VyOS instance from the client struct.
func (c *OpMode) PowerOff(ctx context.Context, path string) error {
	req := &pathRequest{
		URLPath: "/poweroff",
		Op:      OpPowerOff,
		Path:    parsePath(path, true),
	}
	resp := &response{}
	return c.do(ctx, req, resp)
}

// Reboot is a helper function to reboot the VyOS instance from the client struct.
func (c *OpMode) Reboot(ctx context.Context, path string) error {
	req := &pathRequest{
		URLPath: "/reboot",
		Op:      OpReboot,
		Path:    parsePath(path, true),
	}
	resp := &response{}
	return c.do(ctx, req, resp)
}
