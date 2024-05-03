package datastore

import (
	"github.com/acorn-io/z"
	"github.com/adrg/xdg"
	"github.com/gptscript-ai/knowledge/pkg/index"
	"github.com/gptscript-ai/knowledge/pkg/types"
	"github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"github.com/gptscript-ai/knowledge/pkg/vectorstore/chromem"
	cg "github.com/philippgille/chromem-go"
	"log/slog"
)

type Datastore struct {
	Index       *index.DB
	Vectorstore vectorstore.VectorStore
}

func NewDatastore(dsn string, automigrate bool, vectorDBPath string, openAIConfig types.OpenAIConfig) (*Datastore, error) {
	if dsn == "" {
		var err error
		dsn, err = xdg.DataFile("gptscript/knowledge/knowledge.db")
		if err != nil {
			return nil, err
		}
		dsn = "sqlite://" + dsn
		slog.Debug("Using default DSN", "dsn", dsn)
	}

	idx, err := index.New(dsn, automigrate)
	if err != nil {
		return nil, err
	}

	if vectorDBPath == "" {
		vectorDBPath, err = xdg.DataFile("gptscript/knowledge/vector.db")
		if err != nil {
			return nil, err
		}
		slog.Debug("Using default VectorDBPath", "vectorDBPath", vectorDBPath)
	}

	vsdb, err := cg.NewPersistentDB(vectorDBPath, false)
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
