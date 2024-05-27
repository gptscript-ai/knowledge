package flows

import (
	"context"
	"github.com/gptscript-ai/knowledge/pkg/datastore/documentloader"
	"github.com/gptscript-ai/knowledge/pkg/datastore/textsplitter"
	"github.com/gptscript-ai/knowledge/pkg/datastore/transformers"
	dstypes "github.com/gptscript-ai/knowledge/pkg/datastore/types"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"slices"
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

type RetrievalFlow struct {
	// TODO:
}
