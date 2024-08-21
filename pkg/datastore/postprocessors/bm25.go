package postprocessors

import (
	"context"
	"log/slog"
	"sort"
	"strings"

	"github.com/gptscript-ai/knowledge/pkg/datastore/types"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"github.com/iwilltry42/bm25-go/bm25"
)

const BM25PostprocessorName = "bm25"

type BM25Postprocessor struct {
	TopN         int
	SparseWeight float64 // How to weight the BM25 scores against the similarity scores from dense vector search

	K1 float64 // K1 should be between 1.2 and 2 - controls term frequency saturation
	B  float64 // B should be around 0.75 - controls the influence of document length normalization
}

func (c *BM25Postprocessor) Transform(ctx context.Context, response *types.RetrievalResponse) error {

	if c.K1 == 0 {
		c.K1 = 1.5
	}
	if c.B == 0 {
		c.B = 0.75
	}

	var err error
	for q, docs := range response.Responses {
		response.Responses[q], err = c.transform(ctx, q, docs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *BM25Postprocessor) transform(ctx context.Context, query string, docs []vs.Document) ([]vs.Document, error) {
	slog.Debug("BM25", "topN", c.TopN, "numDocs", len(docs), "sparseWeight", c.SparseWeight)

	corpus := make([]string, len(docs))
	// TODO: remove stopwords, etc.
	for i, doc := range docs {
		corpus[i] = doc.Content
	}

	whiteSpaceTokenizer := func(s string) []string { return strings.Split(s, " ") }

	okapi, err := bm25.NewBM25Okapi(corpus, whiteSpaceTokenizer, c.K1, c.B, nil)
	if err != nil {
		return nil, err
	}

	scores, err := okapi.GetScores(strings.Split(query, " "))
	if err != nil {
		return nil, err
	}

	for i, doc := range docs {
		docs[i].Metadata["bm25Score"] = scores[i]

		// Combine BM25 score with similarity score
		docs[i].Metadata["combinedScore"] = c.SparseWeight*scores[i] + (1-c.SparseWeight)*float64(doc.SimilarityScore)
	}

	// Sort by combined score
	sort.Slice(docs, func(i, j int) bool {
		return docs[i].Metadata["combinedScore"].(float64) > docs[i].Metadata["bm25Score"].(float64)
	})

	return docs[:c.TopN-1], nil
}

func (c *BM25Postprocessor) Name() string {
	return BM25PostprocessorName
}
