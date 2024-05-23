package flows

import (
	"context"
	dstypes "github.com/gptscript-ai/knowledge/pkg/datastore/types"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
)

type IngestionFlow struct {
	Load            dstypes.DocumentLoaderFunc
	Split           dstypes.TextSplitterFunc
	Transformations []dstypes.DocumentTransformerFunc
}

func (f *IngestionFlow) Transform(ctx context.Context, docs []vs.Document) ([]vs.Document, error) {
	var err error
	for _, t := range f.Transformations {
		docs, err = t(ctx, docs)
		if err != nil {
			return nil, err
		}
	}
	return docs, nil
}

type RetrievalFlow struct {
	// TODO:
}
