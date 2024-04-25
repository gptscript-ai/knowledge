package datastore

import (
	"github.com/acorn-io/z"
	"github.com/gptscript-ai/knowledge/pkg/index"
	"github.com/gptscript-ai/knowledge/pkg/types"
	"github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"github.com/gptscript-ai/knowledge/pkg/vectorstore/chromem"
	cg "github.com/philippgille/chromem-go"
)

type Datastore struct {
	Index       *index.DB
	Vectorstore vectorstore.VectorStore
}

func NewDatastore(dsn string, automigrate bool, openAIConfig types.OpenAIConfig) (*Datastore, error) {
	idx, err := index.New(dsn, automigrate)
	if err != nil {
		return nil, err
	}

	vsdb, err := cg.NewPersistentDB("vector.db", false)
	if err != nil {
		return nil, err
	}

	embeddingFunc := cg.NewEmbeddingFuncOpenAICompat(
		openAIConfig.APIBase,
		openAIConfig.APIKey,
		openAIConfig.EmbeddingModel,
		z.Pointer(true),
	)

	return &Datastore{
		Index:       idx,
		Vectorstore: chromem.New(vsdb, embeddingFunc),
	}, nil
}
