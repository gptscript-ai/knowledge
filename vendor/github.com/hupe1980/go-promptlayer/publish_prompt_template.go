package promptlayer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// PublishPromptTemplateInput represents the input data for publishing a prompt template.
type PublishPromptTemplateInput struct {
	PromptName     string         `json:"prompt_name,omitempty"`
	PromptTemplate PromptTemplate `json:"prompt_template,omitempty"`
	Tags           []string       `json:"tags,omitempty"`
	APIKey         string         `json:"api_key,omitempty"`
}

// PublishPromptTemplateOutput represents the output data for publishing a prompt template.
type PublishPromptTemplateOutput struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
}

func (o *PublishPromptTemplateOutput) UnmarshalJSON(data []byte) error {
	temp := struct {
		ID      uint64 `json:"id"`
		Success bool   `json:"success"`
	}{}

	// Unmarshal JSON data into the temporary struct
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	o.ID = fmt.Sprint(temp.ID)
	o.Success = temp.Success

	return nil
}

// PublishPromptTemplate publishes a prompt template using the PromptLayer API.
func (c *Client) PublishPromptTemplate(ctx context.Context, input *PublishPromptTemplateInput) (*PublishPromptTemplateOutput, error) {
	url := fmt.Sprintf("%s/rest/publish-prompt-template", c.baseURL)

	input.Tags = append(input.Tags, c.tags...)
	input.APIKey = c.apiKey

	body, err := c.doRequest(ctx, http.MethodPost, url, input)
	if err != nil {
		return nil, err
	}

	output := &PublishPromptTemplateOutput{}
	if err := json.Unmarshal(body, output); err != nil {
		return nil, err
	}

	return output, nil
}
