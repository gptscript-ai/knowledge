package promptlayer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type GetPromptTemplateInput struct {
	PromptName string `json:"prompt_name,omitempty"`
	Version    int    `json:"version,omitempty"`
}

type GetPromptTemplateOutput struct {
	ID             string         `json:"id"`
	Deleted        bool           `json:"deleted"`
	PromptTemplate PromptTemplate `json:"prompt_template"`
	Version        int            `json:"version"`
	Tags           []string       `json:"tags"`
}

func (o *GetPromptTemplateOutput) UnmarshalJSON(data []byte) error {
	temp := struct {
		ID             uint64         `json:"id"`
		Deleted        bool           `json:"deleted"`
		PromptTemplate PromptTemplate `json:"prompt_template"`
		Version        int            `json:"version"`
		Tags           []string       `json:"tags"`
	}{}

	// Unmarshal JSON data into the temporary struct
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	o.ID = fmt.Sprint(temp.ID)
	o.Deleted = temp.Deleted
	o.PromptTemplate = temp.PromptTemplate
	o.Version = temp.Version
	o.Tags = temp.Tags

	return nil
}

func (c *Client) GetPromptTemplate(ctx context.Context, input *GetPromptTemplateInput) (*GetPromptTemplateOutput, error) {
	if input.PromptName == "" {
		return nil, errors.New("promptName is required")
	}

	params := make(url.Values)
	params.Add("prompt_name", input.PromptName)

	if input.Version != 0 {
		params.Add("version", strconv.Itoa(input.Version))
	}

	url := fmt.Sprintf("%s/rest/get-prompt-template?%s", c.baseURL, params.Encode())

	body, err := c.doRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	output := &GetPromptTemplateOutput{}
	if err := json.Unmarshal(body, output); err != nil {
		return nil, err
	}

	return output, nil
}
