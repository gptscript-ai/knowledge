package transformers

import (
	"context"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
)

type ExtraMetadata struct {
	Metadata map[string]any
}

func (e *ExtraMetadata) Transform(_ context.Context, docs []vs.Document) ([]vs.Document, error) {
	for i := range docs {
		docs[i].Metadata = e.Metadata
	}
	return docs, nil
}
