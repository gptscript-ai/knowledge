package flows

import (
	dstypes "github.com/gptscript-ai/knowledge/pkg/datastore/types"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
)

type IngestionFlow struct {
	DocumentLoader       *dstypes.DocumentLoader
	TextSplitter         *dstypes.TextSplitter
	DocumentTransformers []func([]vs.Document) ([]vs.Document, error)
}

type RetrievalFlow struct {
	// TODO:
}
