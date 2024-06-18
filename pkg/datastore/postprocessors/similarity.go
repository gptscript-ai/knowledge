package postprocessors

import (
	"context"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
)

type SimilarityPostprocessor struct {
	Threshold float32
}

func (s *SimilarityPostprocessor) Transform(_ context.Context, query string, docs []vs.Document) ([]vs.Document, error) {
	var filteredDocs []vs.Document
	for _, doc := range docs {
		if doc.SimilarityScore >= s.Threshold {
			filteredDocs = append(filteredDocs, doc)
		}
	}
	return filteredDocs, nil
}
