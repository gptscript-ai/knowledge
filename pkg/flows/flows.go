package flows

import (
	"context"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/store"
	"github.com/philippgille/chromem-go"
	"io"
	"log/slog"
	"slices"

	"github.com/gptscript-ai/knowledge/pkg/datastore/documentloader"
	"github.com/gptscript-ai/knowledge/pkg/datastore/postprocessors"
	"github.com/gptscript-ai/knowledge/pkg/datastore/querymodifiers"
	"github.com/gptscript-ai/knowledge/pkg/datastore/retrievers"
	"github.com/gptscript-ai/knowledge/pkg/datastore/textsplitter"
	"github.com/gptscript-ai/knowledge/pkg/datastore/transformers"
	dstypes "github.com/gptscript-ai/knowledge/pkg/datastore/types"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
)

type IngestionFlow struct {
	Filetypes       []string
	Load            documentloader.LoaderFunc
	Split           textsplitter.SplitterFunc
	Transformations []dstypes.DocumentTransformer
}

func (f *IngestionFlow) Transform(ctx context.Context, docs []vs.Document) ([]vs.Document, error) {
	var err error
	for _, t := range f.Transformations {
		docs, err = t.Transform(ctx, docs)
		if err != nil {
			return nil, err
		}
	}
	return docs, nil
}

func NewDefaultIngestionFlow(filetype string, textsplitterOpts *textsplitter.TextSplitterOpts) IngestionFlow {
	ingestionFlow := IngestionFlow{
		Filetypes: []string{filetype},
	}
	ingestionFlow.FillDefaults(filetype, textsplitterOpts)
	return ingestionFlow
}

func (f *IngestionFlow) SupportsFiletype(filetype string) bool {
	return slices.Contains(f.Filetypes, filetype)
}

func (f *IngestionFlow) FillDefaults(filetype string, textsplitterOpts *textsplitter.TextSplitterOpts) {
	if f.Load == nil {
		f.Load = documentloader.DefaultDocLoaderFunc(filetype)
	}
	if f.Split == nil {
		f.Split = textsplitter.DefaultTextSplitter(filetype, textsplitterOpts).SplitDocuments
	}
	if len(f.Transformations) == 0 {
		f.Transformations = transformers.DefaultDocumentTransformers(filetype)
	}
}

func (f *IngestionFlow) Run(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
	var err error
	var docs []vs.Document

	/*
	 * Load documents from the content
	 * For now, we're using documentloaders from both langchaingo and golc
	 * and translate them to our document schema.
	 */

	if f.Load == nil {
		return nil, nil
	}

	docs, err = f.Load(ctx, reader)
	if err != nil {
		slog.Error("Failed to load documents", "error", err)
		return nil, fmt.Errorf("failed to load documents: %w", err)
	}

	/*
	 * Split documents - Chunking
	 */
	docs, err = f.Split(docs)
	if err != nil {
		slog.Error("Failed to split documents", "error", err)
		return nil, fmt.Errorf("failed to split documents: %w", err)
	}

	/*
	 * Transform documents
	 */
	docs, err = f.Transform(ctx, docs)
	if err != nil {
		slog.Error("Failed to transform documents", "error", err)
		return nil, fmt.Errorf("failed to transform documents: %w", err)
	}

	return docs, nil
}

type RetrievalFlow struct {
	QueryModifiers []querymodifiers.QueryModifier
	Retriever      retrievers.Retriever
	Postprocessors []postprocessors.Postprocessor
}

func (f *RetrievalFlow) FillDefaults(topK int) {
	if f.Retriever == nil {
		slog.Debug("No retriever specified, using basic retriever")
		f.Retriever = &retrievers.BasicRetriever{TopK: topK}
	}
}

type RetrievalFlowOpts struct {
	Where         map[string]string
	WhereDocument []chromem.WhereDocument
}

func (f *RetrievalFlow) Run(ctx context.Context, store store.Store, query string, datasetIDs []string, opts *RetrievalFlowOpts) (*dstypes.RetrievalResponse, error) {
	if opts == nil {
		opts = &RetrievalFlowOpts{}
	}

	queries := []string{query}
	for _, m := range f.QueryModifiers {
		mq, err := m.ModifyQueries(queries)
		if err != nil {
			return nil, fmt.Errorf("failed to modify queries %v with QueryModifier %q: %w", queries, m.Name(), err)
		}
		slog.Debug("Modified queries", "before", queries, "queryModifier", m.Name(), "after", mq)
		queries = mq
	}
	slog.Debug("Updated query set", "query", query, "modified_query_set", queries)

	response := &dstypes.RetrievalResponse{
		Query:     query,
		Responses: make(map[string][]vs.Document, len(queries)),
	}
	for _, q := range queries {

		docs, err := f.Retriever.Retrieve(ctx, store, q, datasetIDs, opts.Where, opts.WhereDocument)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve documents for query %q using retriever %q: %w", q, f.Retriever.Name(), err)
		}
		slog.Debug("Retrieved documents", "num_documents", len(docs), "query", q, "datasets", datasetIDs, "retriever", f.Retriever.Name())
		response.Responses[q] = docs
	}

	for _, pp := range f.Postprocessors {
		err := pp.Transform(ctx, response)
		if err != nil {
			return nil, fmt.Errorf("failed to postprocess retrieval response with Postprocessor %q: %w", pp.Name(), err)
		}
	}
	slog.Debug("Postprocessed RetrievalResponse", "num_responses", len(response.Responses), "original_query", query)

	return response, nil
}
