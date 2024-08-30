package postprocessors

import (
	"context"
	"log/slog"
	"slices"

	"github.com/gptscript-ai/knowledge/pkg/datastore/lib/scores"
	"github.com/gptscript-ai/knowledge/pkg/datastore/types"
)

const ReducePostprocessorName = "reduce"

type ReducePostprocessor struct {
	TopK int
}

func (s *ReducePostprocessor) Transform(ctx context.Context, response *types.RetrievalResponse) error {
	for q, docs := range response.Responses {

		topK := s.TopK

		if len(docs) <= topK {
			continue
		}

		slices.SortFunc(docs, scores.SortBySimilarityScore)

		if topK > len(docs) {
			topK = len(docs)
		}
		if topK <= 0 {
			continue
		}

		slog.Info("Reducing topK", "topK", topK, "len(docs)", len(docs))

		response.Responses[q] = docs[:topK]
	}
	return nil
}

func (s *ReducePostprocessor) Name() string {
	return ReducePostprocessorName
}
