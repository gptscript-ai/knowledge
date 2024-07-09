// Package postprocessors is basically the same as package transformers, but used at a different stage of the RAG pipeline
package postprocessors

import (
	"context"
	"fmt"

	"github.com/gptscript-ai/knowledge/pkg/datastore/transformers"
	"github.com/gptscript-ai/knowledge/pkg/datastore/types"
)

// Postprocessor is similar to types.DocumentTransformer, but can take into account the retrieval query
type Postprocessor interface {
	Transform(ctx context.Context, response *types.RetrievalResponse) error
	Name() string
}

type TransformerWrapper struct {
	types.DocumentTransformer
}

func NewTransformerWrapper(transformer types.DocumentTransformer) *TransformerWrapper {
	return &TransformerWrapper{DocumentTransformer: transformer}
}

func (t *TransformerWrapper) Transform(ctx context.Context, response *types.RetrievalResponse) error {
	for q, docs := range response.Responses {
		newDocs, err := t.DocumentTransformer.Transform(ctx, docs)
		if err != nil {
			return err
		}
		response.Responses[q] = newDocs
	}
	return nil
}

func (t *TransformerWrapper) Name() string {
	return t.DocumentTransformer.Name()
}

var PostprocessorMap = map[string]Postprocessor{
	transformers.ExtraMetadataName:               NewTransformerWrapper(&transformers.ExtraMetadata{}),
	transformers.KeywordExtractorName:            NewTransformerWrapper(&transformers.KeywordExtractor{}),
	transformers.FilterMarkdownDocsNoContentName: NewTransformerWrapper(&transformers.FilterMarkdownDocsNoContent{}),
	SimilarityPostprocessorName:                  &SimilarityPostprocessor{},
	ContentSubstringFilterPostprocessorName:      &ContentSubstringFilterPostprocessor{},
	ContentFilterPostprocessorName:               &ContentFilterPostprocessor{},
	CohereRerankPostprocessorName:                &CohereRerankPostprocessor{},
}

func GetPostprocessor(name string) (Postprocessor, error) {
	var postprocessor Postprocessor
	var ok bool
	postprocessor, ok = PostprocessorMap[name]
	if !ok {
		return nil, fmt.Errorf("unknown postprocessor %q", name)
	}
	return postprocessor, nil
}
