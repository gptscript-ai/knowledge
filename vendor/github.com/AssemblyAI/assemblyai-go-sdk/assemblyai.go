package assemblyai

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	version              = "1.3.0"
	defaultBaseURLScheme = "https"
	defaultBaseURLHost   = "api.assemblyai.com"
	defaultUserAgent     = "assemblyai-go/" + version
)

// Client manages the communication with the AssemblyAI API.
type Client struct {
	baseURL   *url.URL
	userAgent string
	apiKey    string

	httpClient *http.Client

	Transcripts *TranscriptService
	LeMUR       *LeMURService
}

// NewClientWithOptions returns a new configurable AssemblyAI client. If you provide client
// options, they override the default values. Most users will want to use
// [NewClientWithAPIKey].
func NewClientWithOptions(opts ...ClientOption) *Client {
	defaultAPIKey := os.Getenv("ASSEMBLYAI_API_KEY")

	c := &Client{
		baseURL: &url.URL{
			Scheme: defaultBaseURLScheme,
			Host:   defaultBaseURLHost,
		},
		userAgent:  defaultUserAgent,
		httpClient: &http.Client{},
		apiKey:     defaultAPIKey,
	}

	for _, f := range opts {
		f(c)
	}

	c.Transcripts = &TranscriptService{client: c}
	c.LeMUR = &LeMURService{client: c}

	return c
}

// NewClient returns a new authenticated AssemblyAI client.
func NewClient(apiKey string) *Client {
	return NewClientWithOptions(WithAPIKey(apiKey))
}

// ClientOption lets you configure the AssemblyAI client.
type ClientOption func(*Client)

// WithHTTPClient sets the http.Client used for making requests to the API.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL sets the API endpoint used by the client. Mainly used for testing.
func WithBaseURL(rawurl string) ClientOption {
	return func(c *Client) {
		if u, err := url.Parse(rawurl); err == nil {
			c.baseURL = u
		}
	}
}

// WithUserAgent sets the user agent used by the client.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// WithAPIKey sets the API key used for authentication.
func WithAPIKey(key string) ClientOption {
	return func(c *Client) {
		c.apiKey = key
	}
}

func (c *Client) newJSONRequest(method, path string, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter

	if body != nil {
		buf = new(bytes.Buffer)

		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := c.newRequest(method, path, buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	rawurl := c.baseURL.ResolveReference(rel).String()

	req, err := http.NewRequest(method, rawurl, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Authorization", c.apiKey)

	return req, err
}

func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var apierr APIError

		if err := json.NewDecoder(resp.Body).Decode(&apierr); err != nil {
			return nil, err
		}

		apierr.Status = resp.StatusCode

		return nil, apierr
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, err
}
