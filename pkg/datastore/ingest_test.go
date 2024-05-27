package datastore

import (
	"context"
	"github.com/gptscript-ai/knowledge/pkg/datastore/transformers"
	"github.com/gptscript-ai/knowledge/pkg/flows"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractPDF(t *testing.T) {
	ctx := context.Background()
	textSplitterOpts := NewTextSplitterOpts()
	err := filepath.WalkDir("testdata/pdf", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			t.Fatalf("filepath.WalkDir() error = %v", err)
		}
		if d.IsDir() {
			return nil
		}
		t.Logf("Processing %s", path)
		f, err := os.Open(path)
		require.NoError(t, err, "os.Open() error = %v", err)

		filetype := ".pdf"

		ingestionFlow := flows.IngestionFlow{
			Load:            DefaultDocLoaderFunc(filetype),
			Split:           DefaultTextSplitter(filetype, &textSplitterOpts).SplitDocuments,
			Transformations: DefaultDocumentTransformers(filetype),
		}

		// Mandatory Transformation: Add filename to metadata
		em := &transformers.ExtraMetadata{Metadata: map[string]any{"filename": d.Name()}}
		ingestionFlow.Transformations = append(ingestionFlow.Transformations, em)

		docs, err := GetDocuments(ctx, f, ingestionFlow)
		require.NoError(t, err, "GetDocuments() error = %v", err)
		require.NotEmpty(t, docs, "GetDocuments() returned no documents")
		return nil
	})
	require.NoError(t, err, "filepath.WalkDir() error = %v", err)
}
