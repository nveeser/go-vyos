package vyos

import (
	"context"
)

// ConfigMode provides the set of methods which map to the configuration mode in VyOS
type ConfigMode Client

func (c *ConfigMode) do(ctx context.Context, req request, resp *response) error {
	return (*Client)(c).do(ctx, req, resp)
}

// ConfigRequest wraps request and provides an exported marker interface
// to distinguish requests that can be passed to Configure. The payloads
// of each of them can be combined into a slice and passed as the
// data of the POST payload.
type ConfigRequest interface {
	isConfigCommand()
	request
}

type SetRequest struct {
	Path string
}

func (s *SetRequest) isConfigCommand() {}
func (s *SetRequest) requestPayload() (path string, payload any) {
	return "/configure", &pathRequest{Op: OpSet, Path: parsePath(s.Path, true)}
}

type DeleteRequest struct {
	Path string
}

func (r *DeleteRequest) isConfigCommand() {}
func (r *DeleteRequest) requestPayload() (path string, payload any) {
	return "/configure", &pathRequest{Op: OpDelete, Path: parsePath(r.Path, true)}
}

type CommentRequest struct {
	Op      op       `json:"op"`
	Path    []string `json:"path"`
	Comment string   `json:"comment"`
}

func (r *CommentRequest) isConfigCommand() {}
func (r *CommentRequest) requestPayload() (path string, payload any) {
	return "/configure", r
}

func (c *ConfigMode) Configure(ctx context.Context, reqs ...ConfigRequest) error {
	req := &configureRequest{Requests: reqs}
	resp := &response{}
	return c.do(ctx, req, resp)
}

type configureRequest struct {
	Requests []ConfigRequest
}

func (r *configureRequest) requestPayload() (path string, payload any) {
	var payloads []any
	for _, req := range r.Requests {
		_, p := req.requestPayload()
		payloads = append(payloads, p)
	}
	return "/configure", payloads
}

func (c *ConfigMode) Show(ctx context.Context, path string) (map[string]any, error) {
	req := &pathRequest{
		URLPath: "/retrieve",
		Op:      OpShowConfig,
		Path:    parsePath(path, true),
	}
	resp := &response{Data: map[string]any{}}
	err := c.do(ctx, req, resp)
	if err != nil {
		return nil, err
	}
	return resp.Data.(map[string]any), nil
}

func (c *ConfigMode) ShowValues(ctx context.Context, path string) ([]any, error) {
	req := &pathRequest{
		URLPath: "/retrieve",
		Op:      OpShowValues,
		Path:    parsePath(path, true),
	}
	resp := &response{}
	err := c.do(ctx, req, resp)
	if err != nil {
		return nil, err
	}
	return resp.Data.([]any), nil
}

func (c *ConfigMode) Exists(ctx context.Context, path string) (bool, error) {
	req := &pathRequest{
		URLPath: "/retrieve",
		Op:      OpExists,
		Path:    parsePath(path, true),
	}
	resp := &response{}
	err := c.do(ctx, req, resp)
	if err != nil {
		return false, err
	}
	return resp.Data.(bool), nil
}

func (c *ConfigMode) Set(ctx context.Context, path string) error {
	req := &pathRequest{
		URLPath: "/configure",
		Op:      OpSet,
		Path:    parsePath(path, true),
	}
	resp := &response{}
	err := c.do(ctx, req, resp)
	if err != nil {
		return err
	}
	return nil
}

func (c *ConfigMode) Delete(ctx context.Context, path string) error {
	req := &pathRequest{
		URLPath: "/configure",
		Op:      OpDelete,
		Path:    parsePath(path, true),
	}
	resp := &response{}
	err := c.do(ctx, req, resp)
	if err != nil {
		return err
	}
	return nil
}

func (c *ConfigMode) Comment(ctx context.Context, comment string, path string) error {
	req := &CommentRequest{
		Op:      OpComment,
		Path:    parsePath(path, true),
		Comment: comment,
	}
	resp := &response{}
	return c.do(ctx, req, resp)
}

func (c *ConfigMode) Save(ctx context.Context, file string) (string, error) {
	req := &cfgFileRequest{
		Op:   OpSave,
		File: file,
	}
	resp := &response{}
	err := c.do(ctx, req, resp)
	if err != nil {
		return "", err
	}
	return resp.Data.(string), nil
}

func (c *ConfigMode) Load(ctx context.Context, file string) (string, error) {
	req := &cfgFileRequest{Op: OpLoad, File: file}
	resp := &response{}
	err := c.do(ctx, req, resp)
	if err != nil {
		return "", err
	}
	return resp.Data.(string), nil
}

func (c *ConfigMode) CommitConfirm(ctx context.Context) (string, error) {
	req := &cfgFileRequest{Op: OpConfirm}
	resp := &response{}
	err := c.do(ctx, req, resp)
	if err != nil {
		return "", err
	}
	return resp.Data.(string), nil
}

type cfgFileRequest struct {
	Op   op     `json:"op"`
	File string `json:"fila,omitempty"`
}

func (s *cfgFileRequest) requestPayload() (path string, payload any) {
	return "/config-file", s
}
