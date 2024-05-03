package datastore

import (
	"context"
	"fmt"
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

func GetDatastorePaths(dsn, vectordbPath string) (string, string, error) {
	if dsn == "" {
		var err error
		dsn, err = xdg.DataFile("gptscript/knowledge/knowledge.db")
		if err != nil {
			return "", "", err
		}
		dsn = "sqlite://" + dsn
		slog.Debug("Using default DSN", "dsn", dsn)
	}

	if vectordbPath == "" {
		var err error
		vectordbPath, err = xdg.DataFile("gptscript/knowledge/vector.db")
		if err != nil {
			return "", "", err
		}
		slog.Debug("Using default VectorDBPath", "vectordbPath", vectordbPath)
	}

	return dsn, vectordbPath, nil
}

func NewDatastore(dsn string, automigrate bool, vectorDBPath string, openAIConfig types.OpenAIConfig) (*Datastore, error) {
	dsn, vectorDBPath, err := GetDatastorePaths(dsn, vectorDBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to determine datastore paths: %w", err)
	}

	idx, err := index.New(dsn, automigrate)
	if err != nil {
		return nil, err
	}

	if err := idx.AutoMigrate(); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate index: %w", err)
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

	ds := &Datastore{
		Index:       idx,
		Vectorstore: chromem.New(vsdb, embeddingFunc),
	}

	// Ensure default dataset exists
	defaultDS, err := ds.GetDataset(context.Background(), "default")
	if err != nil {
		return nil, fmt.Errorf("failed to ensure default dataset: %w", err)
	}

	if defaultDS == nil {
		err = ds.NewDataset(context.Background(), types.Dataset{ID: "default", EmbedDimension: nil})
		if err != nil {
			return nil, fmt.Errorf("failed to create default dataset: %w", err)
		}
	}

	return ds, nil
}
