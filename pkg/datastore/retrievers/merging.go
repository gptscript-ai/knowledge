package retrievers

import (
	"context"
	"log/slog"

	"github.com/gptscript-ai/knowledge/pkg/datastore/store"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"github.com/philippgille/chromem-go"
)

const MergingRetrieverName = "merge"

type MergingRetriever struct {
	TopK       int
	Retrievers []Retriever `json:"retrievers" mapstructure:"retrievers" yaml:"retrievers"`
}

func (r *MergingRetriever) Name() string {
	return MergingRetrieverName
}
func (r *MergingRetriever) Retrieve(ctx context.Context, store store.Store, query string, datasetIDs []string, where map[string]string, whereDocument []chromem.WhereDocument) ([]vs.Document, error) {
	log := slog.With("component", "MergingRetriever")

	var docs []vs.Document
	for _, retriever := range r.Retrievers {
		log.Debug("Retrieving documents from retriever", "retriever", retriever.Name())
		docsRetriever, err := retriever.Retrieve(ctx, store, query, datasetIDs, where, whereDocument)
		if err != nil {
			log.Error("Failed to retrieve documents from retriever", "retriever", retriever.Name(), "error", err)
			return nil, err
		}

	docLoop:
		for _, doc := range docs {
			// check if	doc is already in docs and if so, update similarity score if higher
			for i, r := range docsRetriever {
				if doc.ID == r.ID {
					if doc.SimilarityScore > r.SimilarityScore {
						docs[i].SimilarityScore = doc.SimilarityScore
						continue docLoop
					}
				}
			}
			docs = append(docs, doc)
		}
	}

	return docs[:r.TopK-1], nil
}
