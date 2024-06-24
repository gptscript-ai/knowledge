// Package postprocessors is basically the same as package transformers, but used at a different stage of the RAG pipeline
package postprocessors

import (
	"context"
	"fmt"

	"github.com/gptscript-ai/knowledge/pkg/datastore/transformers"
	"github.com/gptscript-ai/knowledge/pkg/datastore/types"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
)

// Postprocessor is similar to types.DocumentTransformer, but can take into account the retrieval query
type Postprocessor interface {
	Transform(ctx context.Context, query string, docs []vs.Document) ([]vs.Document, error)
}

type TransformerWrapper struct {
	types.DocumentTransformer
}

func NewTransformerWrapper(transformer types.DocumentTransformer) *TransformerWrapper {
	return &TransformerWrapper{DocumentTransformer: transformer}
}

func (t *TransformerWrapper) Transform(ctx context.Context, query string, docs []vs.Document) ([]vs.Document, error) {
	return t.DocumentTransformer.Transform(ctx, docs)
}

var PostprocessorMap = map[string]Postprocessor{
	"extra_metadata":                  NewTransformerWrapper(&transformers.ExtraMetadata{}),
	"keywords":                        NewTransformerWrapper(&transformers.KeywordExtractor{}),
	"filter_markdown_docs_no_content": NewTransformerWrapper(&transformers.FilterMarkdownDocsNoContent{}),
	"similarity":                      &SimilarityPostprocessor{},
	"content_substring_filter":        &ContentSubstringFilterPostprocessor{},
	"content_filter":                  &ContentFilterPostprocessor{},
	"cohere_rerank":                   &CohereRerankPostprocessor{},
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
