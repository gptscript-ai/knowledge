package datastore

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings"
	etypes "github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/types"
	"log/slog"

	"github.com/acorn-io/z"
	"github.com/google/uuid"
	"github.com/gptscript-ai/knowledge/pkg/datastore/filetypes"
	"github.com/gptscript-ai/knowledge/pkg/datastore/textsplitter"
	"github.com/gptscript-ai/knowledge/pkg/datastore/transformers"
	"github.com/gptscript-ai/knowledge/pkg/flows"
	"github.com/gptscript-ai/knowledge/pkg/index"
)

type IngestOpts struct {
	Filename            *string
	FileMetadata        *index.FileMetadata
	IsDuplicateFuncName string
	IsDuplicateFunc     IsDuplicateFunc
	TextSplitterOpts    *textsplitter.TextSplitterOpts
	IngestionFlows      []flows.IngestionFlow
	CreateDataset       bool
}

// Ingest loads a document from a reader and adds it to the dataset.
func (s *Datastore) Ingest(ctx context.Context, datasetID string, content []byte, opts IngestOpts) ([]string, error) {

	// Get dataset
	ds, err := s.GetDataset(ctx, datasetID)
	if err != nil {
		return nil, err
	}

	// Dataset does not exist - create it if requested, else error out
	if ds == nil {
		// Create dataset if it doesn't exist
		if opts.CreateDataset {
			ds = &index.Dataset{ID: datasetID}
			if err := s.NewDataset(ctx, *ds); err != nil {
				return nil, fmt.Errorf("failed to create dataset %q: %w", datasetID, err)
			}
		} else {
			return nil, fmt.Errorf("dataset %q not found", datasetID)
		}

	}

	// Check if Dataset has an embedding config attached
	if ds.EmbeddingsProviderConfig == nil {
		slog.Info("Embeddingsconfig", "config", s.EmbeddingConfig)
		ncfg, err := etypes.AsEmbeddingModelProviderConfig(s.EmbeddingModelProvider)
		if err != nil {
			return nil, fmt.Errorf("failed to get embedding model provider config: %w", err)
		}
		nds := index.Dataset{
			ID:                       datasetID,
			EmbeddingsProviderConfig: &ncfg,
		}
		ds, err = s.UpdateDataset(ctx, nds, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to update dataset: %w", err)
		}
	}
	if ds.EmbeddingsProviderConfig != nil {
		if s.EmbeddingModelProvider.Name() != ds.EmbeddingsProviderConfig.Type {
			slog.Warn("Embeddings provider mismatch", "dataset", datasetID, "attached", ds.EmbeddingsProviderConfig.Type, "configured", s.EmbeddingModelProvider.Name())
		}

		dsEmbeddingProvider, err := embeddings.ProviderFromConfig(*ds.EmbeddingsProviderConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to get embeddings model provider: %w", err)
		}
		err = embeddings.CompareRequiredFields(s.EmbeddingModelProvider.Config(), dsEmbeddingProvider.Config())
		if err != nil {
			return nil, fmt.Errorf("mismatching embedding provider configs: %w", err)
		}
	}

	// File Deduplication
	isDuplicate := DummyDedupe // default: no deduplication
	if opts.IsDuplicateFuncName != "" {
		df, ok := IsDuplicateFuncs[opts.IsDuplicateFuncName]
		if !ok {
			return nil, fmt.Errorf("unknown deduplication function: %s", opts.IsDuplicateFuncName)
		}
		isDuplicate = df
	} else if opts.IsDuplicateFunc != nil {
		isDuplicate = opts.IsDuplicateFunc
	}

	filename := z.Dereference(opts.Filename)

	// Generate ID
	fUUID, err := uuid.NewUUID()
	if err != nil {
		slog.Error("Failed to generate UUID", "error", err)
		return nil, err
	}
	fileID := fUUID.String()

	/*
	 * Detect filetype
	 */

	filetype, err := filetypes.GetFiletype(filename, content)
	if err != nil {
		return nil, err
	}

	/*
	 * Set filename if not provided
	 */
	if filename == "" {
		filename = "<unnamed_document>"
		*opts.Filename = filename
	}

	slog.Debug("Loading data", "type", filetype, "filename", filename, "size", len(content))

	/*
	 * Exit early if the document is a duplicate
	 */
	isDupe, err := isDuplicate(ctx, s, datasetID, nil, opts)
	if err != nil {
		slog.Error("Failed to check for duplicates", "error", err)
		return nil, fmt.Errorf("failed to check for duplicates: %w", err)
	}
	if isDupe {
		slog.Info("Ignoring duplicate document", "filename", filename, "absolute_path", opts.FileMetadata.AbsolutePath)
		return nil, nil
	}

	/*
	 * Load the ingestion flow - custom or default config or mixture of both
	 */
	ingestionFlow := flows.IngestionFlow{}
	for _, flow := range opts.IngestionFlows {
		if flow.SupportsFiletype(filetype) {
			ingestionFlow = flow
			break
		}
	}
	ingestionFlow.FillDefaults(filetype, opts.TextSplitterOpts)

	// Mandatory Transformation: Add filename to metadata
	em := &transformers.ExtraMetadata{Metadata: map[string]any{"filename": filename, "absPath": opts.FileMetadata.AbsolutePath}}
	ingestionFlow.Transformations = append(ingestionFlow.Transformations, em)

	docs, err := ingestionFlow.Run(ctx, bytes.NewReader(content))
	if err != nil {
		slog.Error("Failed to load documents", "error", err)
		return nil, fmt.Errorf("failed to load documents: %w", err)
	}

	if len(docs) == 0 {
		return nil, nil
	}

	// Before adding doc, we need to remove the existing documents for duplicates or old contents
	where := map[string]string{
		"absPath": opts.FileMetadata.AbsolutePath,
	}
	if err := s.Vectorstore.RemoveDocument(ctx, "", datasetID, where, nil); err != nil {
		return nil, err
	}

	// Add documents to VectorStore -> This generates the embeddings
	slog.Debug("Ingesting documents", "count", len(docs))
	docIDs, err := s.Vectorstore.AddDocuments(ctx, docs, datasetID)
	if err != nil {
		slog.Error("Failed to add documents", "error", err)
		return nil, fmt.Errorf("failed to add documents: %w", err)
	}

	// Record file and documents in database
	dbDocs := make([]index.Document, len(docIDs))
	for idx, docID := range docIDs {
		dbDocs[idx] = index.Document{
			ID:      docID,
			FileID:  fileID,
			Dataset: datasetID,
		}
	}

	dbFile := index.File{
		ID:        fileID,
		Dataset:   datasetID,
		Documents: dbDocs,
		FileMetadata: index.FileMetadata{
			Name: filename,
		},
	}

	if opts.FileMetadata != nil {
		dbFile.FileMetadata.AbsolutePath = opts.FileMetadata.AbsolutePath
		dbFile.FileMetadata.Size = opts.FileMetadata.Size
		dbFile.FileMetadata.ModifiedAt = opts.FileMetadata.ModifiedAt
	}

	tx := s.Index.WithContext(ctx).Create(&dbFile)
	if tx.Error != nil {
		slog.Error("Failed to create file", "error", tx.Error)
		return nil, fmt.Errorf("failed to create file: %w", tx.Error)
	}

	slog.Info("Ingested document", "filename", filename, "count", len(docIDs), "absolute_path", dbFile.FileMetadata.AbsolutePath)

	return docIDs, nil
}
