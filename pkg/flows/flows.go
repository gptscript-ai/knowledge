package flows

import (
	"context"
	"fmt"
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

func (f *RetrievalFlow) Run(ctx context.Context, store vs.VectorStore, query string, datasetID string) ([]vs.Document, error) {
	var err error
	for _, m := range f.QueryModifiers {
		query, err = m.ModifyQuery(query)
		if err != nil {
			return nil, err
		}
	}
	docs, err := f.Retriever.Retrieve(ctx, store, query, datasetID)
	if err != nil {
		return nil, err
	}

	for _, pp := range f.Postprocessors {
		docs, err = pp.Transform(ctx, docs)
		if err != nil {
			return nil, err
		}
	}

	slog.Debug("Retrieved documents", "num_documents", len(docs), "query", query, "dataset", datasetID)
	return docs, nil
}
