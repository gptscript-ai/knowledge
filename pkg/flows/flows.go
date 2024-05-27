package flows

import (
	"context"
	"github.com/gptscript-ai/knowledge/pkg/datastore/documentloader"
	"github.com/gptscript-ai/knowledge/pkg/datastore/textsplitter"
	dstypes "github.com/gptscript-ai/knowledge/pkg/datastore/types"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
)

type IngestionFlow struct {
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

type RetrievalFlow struct {
	// TODO:
}
