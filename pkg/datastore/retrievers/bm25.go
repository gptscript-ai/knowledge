package retrievers

import (
	"context"
	"log/slog"
	"sort"

	"github.com/gptscript-ai/knowledge/pkg/datastore/lib/bm25"
	"github.com/gptscript-ai/knowledge/pkg/datastore/postprocessors"
	"github.com/gptscript-ai/knowledge/pkg/datastore/store"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"github.com/philippgille/chromem-go"
)

const BM25RetrieverName = postprocessors.BM25PostprocessorName

type BM25Retriever struct {
	postprocessors.BM25Postprocessor `json:",inline" mapstructure:",squash" yaml:",squash,inline"`
}

func (r *BM25Retriever) Name() string {
	return BM25RetrieverName
}
func (r *BM25Retriever) Retrieve(ctx context.Context, store store.Store, query string, datasetIDs []string, where map[string]string, whereDocument []chromem.WhereDocument) ([]vs.Document, error) {
	log := slog.With("component", "BM25Retriever")

	var docs []vs.Document
	for _, datasetID := range datasetIDs {
		log.Debug("Retrieving documents from dataset", "dataset", datasetID)
		docsDataset, err := store.GetDocuments(ctx, datasetID, where, whereDocument)
		if err != nil {
			log.Error("Failed to retrieve documents from dataset", "dataset", datasetID, "error", err)
			return nil, err
		}
		docs = append(docs, docsDataset...)
	}

	scores, err := bm25.BM25Run(docs, query, r.K1, r.B, r.CleanStopWords)
	if err != nil {
		log.Error("Failed to run BM25", "error", err)
		return nil, err
	}

	for i, doc := range docs {
		doc.Metadata["bm25Score"] = scores[i]
	}

	sort.Slice(docs, func(i, j int) bool {
		return docs[i].Metadata["bm25Score"].(float64) > docs[i].Metadata["bm25Score"].(float64)
	})

	return docs[:r.TopN-1], nil
}
