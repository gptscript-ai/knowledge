package postprocessors

import (
	"context"
	"github.com/gptscript-ai/knowledge/pkg/datastore/types"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"slices"
)

const ReducePostprocessorName = "reduce"

type ReducePostprocessor struct {
	TopK int
}

func (s *ReducePostprocessor) Transform(ctx context.Context, response *types.RetrievalResponse) error {
	for q, docs := range response.Responses {

		slices.SortFunc(docs, func(i, j vs.Document) int {
			if i.SimilarityScore > j.SimilarityScore {
				return -1
			}
			if i.SimilarityScore < j.SimilarityScore {
				return 1
			}
			return 0
		})

		topK := s.TopK
		if topK > len(docs) {
			topK = len(docs) - 1
		}

		response.Responses[q] = docs[:topK]
	}
	return nil
}

func (s *ReducePostprocessor) Name() string {
	return ReducePostprocessorName
}
