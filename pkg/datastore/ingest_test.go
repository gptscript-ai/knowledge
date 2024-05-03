package datastore

import (
	"context"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractPDF(t *testing.T) {
	ctx := context.Background()
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
		docs, err := GetDocuments(ctx, d.Name(), ".pdf", f)
		require.NoError(t, err, "GetDocuments() error = %v", err)
		require.NotEmpty(t, docs, "GetDocuments() returned no documents")
		return nil
	})
	require.NoError(t, err, "filepath.WalkDir() error = %v", err)
}
