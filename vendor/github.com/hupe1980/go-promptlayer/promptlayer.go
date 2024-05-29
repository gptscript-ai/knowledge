package promptlayer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HTTPClient is an interface representing an HTTP client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientOptions struct {
	BaseURL    string
	Tags       []string
	HTTPClient HTTPClient
}

// Client represents the PromptLayer API client.
type Client struct {
	apiKey     string
	baseURL    string
	tags       []string
	httpClient HTTPClient
}

// NewClient creates a new PromptLayer API client with the provided API key.
func NewClient(apiKey string, optFns ...func(o *ClientOptions)) *Client {
	opts := ClientOptions{
		BaseURL:    "https://api.promptlayer.com",
		HTTPClient: http.DefaultClient,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &Client{
		apiKey:     apiKey,
		baseURL:    opts.BaseURL,
		tags:       opts.Tags,
		httpClient: opts.HTTPClient,
	}
}

// PromptTemplate represents a prompt template.
type PromptTemplate struct {
	Template       string   `json:"template,omitempty"`
	InputVariables []string `json:"input_variables,omitempty"`
}

// MarshalJSON is a custom JSON marshaler for PromptTemplate.
func (pt *PromptTemplate) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type           string   `json:"_type"`
		Template       string   `json:"template,omitempty"`
		InputVariables []string `json:"input_variables,omitempty"`
	}{
		Type:           "prompt",
		Template:       pt.Template,
		InputVariables: pt.InputVariables,
	})
}

// APIError represents an error response from the PromptLayer API.
type APIError struct {
	// Message contains the error message.
	Message string `json:"message"`
}

// doRequest performs an HTTP request with the specified method, URL, and payload.
// It returns the response body as a byte array and any error encountered during the request.
func (c *Client) doRequest(ctx context.Context, method, url string, payload any) ([]byte, error) {
	var body io.Reader

	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		body = bytes.NewReader(b)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("X-API-KEY", c.apiKey)

	res, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		apiError := APIError{}
		if err := json.Unmarshal(resBody, &apiError); err != nil {
			return nil, fmt.Errorf("prompt layer error: %s", string(resBody))
		}

		return nil, fmt.Errorf("prompt layer error: %s", apiError.Message)
	}

	return resBody, nil
}
