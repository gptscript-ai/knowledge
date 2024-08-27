package retrievers

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/acorn-io/z"
	"github.com/gptscript-ai/knowledge/pkg/datastore/lib/scores"
	"github.com/gptscript-ai/knowledge/pkg/datastore/store"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"github.com/mitchellh/mapstructure"
	"github.com/philippgille/chromem-go"
)

const MergingRetrieverName = "merge"

type MergingRetriever struct {
	TopK       int
	Retrievers []RetrieverToMerge `json:"retrievers" mapstructure:"retrievers" yaml:"retrievers"`
	retrievers []Retriever
}

type RetrieverToMerge struct {
	Name    string         `json:"name,omitempty" mapstructure:"name" yaml:"name"`
	Weight  *float32       `json:"weight,omitempty" mapstructure:"weight" yaml:"weight"`
	Options map[string]any `json:"options,omitempty" mapstructure:"options" yaml:"options"`
}

func (r *MergingRetriever) Name() string {
	return MergingRetrieverName
}

func (r *MergingRetriever) DecodeConfig(cfg map[string]any) error {
	if err := mapstructure.Decode(cfg, &r); err != nil {
		return fmt.Errorf("failed to decode merging retriever configuration: %w", err)
	}

	r.retrievers = make([]Retriever, len(r.Retrievers))
	for i, retrieverConfig := range r.Retrievers {
		retriever, err := GetRetriever(retrieverConfig.Name)
		if err != nil {
			return err
		}
		if err := retriever.DecodeConfig(retrieverConfig.Options); err != nil {
			return err
		}
		r.retrievers[i] = retriever
	}

	return nil
}

func (r *MergingRetriever) Retrieve(ctx context.Context, store store.Store, query string, datasetIDs []string, where map[string]string, whereDocument []chromem.WhereDocument) ([]vs.Document, error) {
	log := slog.With("component", "MergingRetriever")

	// Set default weight to 1.0 if not provided
	for i, retriever := range r.Retrievers {
		if retriever.Weight == nil {
			r.Retrievers[i].Weight = z.Pointer(float32(1.0))
		}
	}

	var resultDocs []vs.Document
	for ri, retriever := range r.Retrievers {
		log.Debug("Retrieving documents from retriever", "retriever", retriever.Name)
		retrievedDocs, err := r.retrievers[ri].Retrieve(ctx, store, query, datasetIDs, where, whereDocument)
		if err != nil {
			log.Error("Failed to retrieve documents from retriever", "retriever", retriever.Name, "error", err)
			return nil, err
		}

		min, max := scores.FindMinMaxScores(retrievedDocs)

	docLoop:
		for _, retrievedDoc := range retrievedDocs {
			for i, resultDoc := range resultDocs {
				// check if	resultDoc is already in resultDocs and if so, update similarity score if higher
				if resultDoc.ID == retrievedDoc.ID {
					// Note that this was found by another retriever and note it's similarityScore
					resultDocs[i].Metadata["retriever"] = fmt.Sprintf("%s,%s", resultDocs[i].Metadata["retriever"], retriever.Name)
					resultDocs[i].Metadata["retrieverScore::"+retriever.Name] = retrievedDoc.SimilarityScore
					normalizedScore := scores.NormalizeScore(retrievedDoc.SimilarityScore, min, max)
					resultDocs[i].Metadata["retrieverScoreNormalized::"+retriever.Name] = normalizedScore
					resultDocs[i].SimilarityScore += normalizedScore * z.Dereference(retriever.Weight)
					continue docLoop
				}
			}
			// not in resultDocs yet, add it
			retrievedDoc.Metadata["retriever"] = retriever.Name
			retrievedDoc.Metadata["retrieverScore::"+retriever.Name] = retrievedDoc.SimilarityScore
			normalizedScore := scores.NormalizeScore(retrievedDoc.SimilarityScore, min, max)
			retrievedDoc.Metadata["retrieverScoreNormalized::"+retriever.Name] = normalizedScore
			retrievedDoc.SimilarityScore = normalizedScore * z.Dereference(retriever.Weight)
			resultDocs = append(resultDocs, retrievedDoc)
		}
	}

	// Sort the resultDocs by similarity score
	slices.SortFunc(resultDocs, func(i, j vs.Document) int {
		if i.SimilarityScore > j.SimilarityScore {
			return -1
		}
		if i.SimilarityScore < j.SimilarityScore {
			return 1
		}
		return 0
	})

	return resultDocs[:r.TopK-1], nil
}
