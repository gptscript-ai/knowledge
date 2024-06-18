package postprocessors

import (
	"context"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"strings"
)

type ContentSubstringFilterPostprocessor struct {
	Contains    []string
	NotContains []string
}

func (c *ContentSubstringFilterPostprocessor) Transform(_ context.Context, _ string, docs []vs.Document) ([]vs.Document, error) {
	var filteredDocs []vs.Document
	for _, doc := range docs {
		containsOK := true
		for _, contains := range c.Contains {
			if !strings.Contains(doc.Content, contains) {
				containsOK = false
				break
			}
		}

		notContainsOK := true
		for _, notContains := range c.NotContains {
			if strings.Contains(doc.Content, notContains) {
				notContainsOK = false
				break
			}
		}

		if containsOK && notContainsOK {
			filteredDocs = append(filteredDocs, doc)
		}
	}
	return filteredDocs, nil
}
