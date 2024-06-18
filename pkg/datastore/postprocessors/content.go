package postprocessors

import (
	"context"
	"encoding/json"
	"github.com/gptscript-ai/knowledge/pkg/llm"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
)

type ContentFilterPostprocessor struct {
	Question string // Question about the content, that can be answered with yes or no
	Include  bool   // Whether to include or exclude the documents for which the answer is yes
	LLM      llm.LLM
}

var promptTemplate string = `You're an expert content analyst.
You're given a question about document contents and have to answer it with true or false.
No other option is allowed.

The question is: "{.question}"

The content is:
{.content}

--- End of Content ---
Reply only with {"result": <true-or-false>}.`

type response struct {
	Result bool `json:"result"`
}

func (c *ContentFilterPostprocessor) Transform(ctx context.Context, _ string, docs []vs.Document) ([]vs.Document, error) {
	var filteredDocs []vs.Document
	for _, doc := range docs {
		res, err := c.LLM.Prompt(ctx, promptTemplate, map[string]any{
			"question": c.Question,
			"content":  doc.Content,
		})
		if err != nil {
			return nil, err
		}

		var resp response
		err = json.Unmarshal([]byte(res), &resp)
		if err != nil {
			return nil, err
		}

		if resp.Result == c.Include {
			filteredDocs = append(filteredDocs, doc)
		}
	}
	return filteredDocs, nil
}
