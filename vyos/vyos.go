package vyos

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// Client represents a VyOS API client.
type Client struct {
	httpClient *http.Client // HTTP client used to communicate with the API.
	baseURL    *url.URL
	token      string // token used for authentication.
	userAgent  string // User agent used when communicating with the API.
}

// NewClient creates a new VyOS API client. Host must be
// a value url string (http://host:port) or a host:port string. Otherwise
// an error is returned.
func NewClient(host string, opts ...Option) (*Client, error) {
	baseURL, ok := buildURL(host)
	if !ok {
		return nil, fmt.Errorf("could not parse as URL or host:port: %s", host)
	}
	c := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout:   5 * time.Second,
			Transport: http.DefaultTransport,
		},
		userAgent: defaultUserAgent,
	}
	for _, o := range opts {
		o(c)
	}
	return c, nil
}

func buildURL(host string) (*url.URL, bool) {
	u, err := url.Parse(host)
	if err == nil && u.Scheme != "" && u.Host != "" {
		return u, true
	}
	return &url.URL{
		Scheme: "https",
		Host:   host,
	}, true
}

// OpMode returns the client which provides "op-mode" commands.
func (c *Client) OpMode() *OpMode { return (*OpMode)(c) }

// ConfigMode provides the client for updating configuration.
func (c *Client) ConfigMode() *ConfigMode { return (*ConfigMode)(c) }

// Option is the type passed to NewClient to configure the client
type Option func(*Client)

// Token sets the token for the VyOS API client.
func Token(token string) Option {
	return func(c *Client) {
		c.token = token
	}
}

// Insecure enables http.Transport.TLSClientConfig.InsecureSkipVerify
// for calling hosts with self-signed certificates.
func Insecure() Option {
	return func(c *Client) {
		transport := c.httpClient.Transport
	loop:
		switch x := transport.(type) {
		case *http.Transport:
			if x.TLSClientConfig == nil {
				x.TLSClientConfig = &tls.Config{}
			}
			x.TLSClientConfig.InsecureSkipVerify = true
		case *loggingTransport:
			transport = x.RoundTripper
			goto loop

		case nil:
			c.httpClient.Transport = &http.Transport{}
			goto loop
		}
	}
}

// UserAgent sets the User-Agent.
func UserAgent(s string) Option {
	return func(c *Client) {
		c.userAgent = s
	}
}

// Timeout sets the default timeout for the http.Client.
func Timeout(d time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = d
	}
}

// WithHTTPClient updates the Client to use the provided HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// DebugLogging enables printing the HTTP Request and HTTP Response
// of each interaction for debugging.
func DebugLogging() Option {
	return func(c *Client) {
		c.httpClient.Transport = &loggingTransport{c.httpClient.Transport}
	}
}

// loggingTransport wraps an http.RoundTripper to log HTTP request and response
// values.
type loggingTransport struct {
	http.RoundTripper
}

func (t *loggingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	// Dump the request data including the body (true)
	reqDump, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		log.Printf("Error dumping request: %v", err)
	} else {
		fmt.Printf("REQUEST:\n%s\n", reqDump)
	}

	// Execute the actual request
	resp, err := t.RoundTripper.RoundTrip(r)
	if err != nil {
		// Log error and return
		log.Printf("Error making request: %v", err)
		return resp, err
	}

	// Dump the response data including the body (true)
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Printf("Error dumping response: %v", err)
	} else {
		fmt.Printf("RESPONSE:\n%s\n", respDump)
	}
	return resp, err
}
