package postprocessors

import (
	"context"
	"github.com/gptscript-ai/knowledge/pkg/datastore/types"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
)

const SimilarityPostprocessorName = "similarity"

type SimilarityPostprocessor struct {
	Threshold float32
}

func (s *SimilarityPostprocessor) Transform(ctx context.Context, response *types.RetrievalResponse) error {
	for q, docs := range response.Responses {
		var filteredDocs []vs.Document
		for _, doc := range docs {
			if doc.SimilarityScore >= s.Threshold {
				filteredDocs = append(filteredDocs, doc)
			}
		}
		response.Responses[q] = filteredDocs
	}
	return nil
}

func (s *SimilarityPostprocessor) Name() string {
	return SimilarityPostprocessorName
}
