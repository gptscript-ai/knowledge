package transformers

import (
	"context"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"strings"
)

// FilterMarkdownDocsNoContent filters out Markdown documents with no content or only headings
//
// TODO: this may be moved into the MarkdownTextSplitter
type FilterMarkdownDocsNoContent struct{}

func (f *FilterMarkdownDocsNoContent) Transform(_ context.Context, docs []vs.Document) ([]vs.Document, error) {
	var filteredDocs []vs.Document
	for _, doc := range docs {
		if doc.Content != "" {
			for _, line := range strings.Split(doc.Content, "\n") {
				if !strings.HasPrefix(line, "#") {
					filteredDocs = append(filteredDocs, doc)
				}
			}
		}
	}
	return filteredDocs, nil
}
