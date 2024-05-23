package types

import (
	"context"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"io"
)

type DocumentTransformerFunc func(context.Context, []vs.Document) ([]vs.Document, error)

type DocumentTransformer interface {
	Transform(context.Context, []vs.Document) ([]vs.Document, error)
}

type DocumentLoaderFunc func(context.Context, io.Reader) ([]vs.Document, error)

type TextSplitterFunc func([]vs.Document) ([]vs.Document, error)

type DocumentLoader interface {
	Load(ctx context.Context) ([]vs.Document, error)
	LoadAndSplit(ctx context.Context, splitter TextSplitter) ([]vs.Document, error)
}

type TextSplitter interface {
	SplitDocuments(docs []vs.Document) ([]vs.Document, error)
}
