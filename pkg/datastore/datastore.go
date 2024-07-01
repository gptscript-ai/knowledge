package datastore

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/types"
	"io"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/acorn-io/z"
	"github.com/adrg/xdg"
	"github.com/gptscript-ai/knowledge/pkg/config"
	"github.com/gptscript-ai/knowledge/pkg/index"
	"github.com/gptscript-ai/knowledge/pkg/llm"
	"github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"github.com/gptscript-ai/knowledge/pkg/vectorstore/chromem"
	cg "github.com/philippgille/chromem-go"
)

type Datastore struct {
	LLM         llm.LLM
	Index       *index.DB
	Vectorstore vectorstore.VectorStore
}

// GetDatastorePaths returns the paths for the datastore and vectorstore databases.
// In addition, it returns a boolean indicating whether the datastore is an archive.
func GetDatastorePaths(dsn, vectordbPath string) (string, string, bool, error) {
	var isArchive bool

	if dsn == "" {
		var err error
		dsn, err = xdg.DataFile("gptscript/knowledge/knowledge.db")
		if err != nil {
			return "", "", isArchive, err
		}
		dsn = "sqlite://" + dsn
		slog.Debug("Using default DSN", "dsn", dsn)
	}

	if strings.HasPrefix(dsn, types.ArchivePrefix) {
		dsn = "sqlite://" + strings.TrimPrefix(dsn, types.ArchivePrefix)
		isArchive = true
	}

	if vectordbPath == "" {
		var err error
		vectordbPath, err = xdg.DataFile("gptscript/knowledge/vector.db")
		if err != nil {
			return "", "", isArchive, err
		}
		slog.Debug("Using default VectorDBPath", "vectordbPath", vectordbPath)
	}
	if strings.HasPrefix(vectordbPath, types.ArchivePrefix) {
		vectordbPath = strings.TrimPrefix(vectordbPath, types.ArchivePrefix)
		isArchive = true
	}

	return dsn, vectordbPath, isArchive, nil
}

func NewDatastore(dsn string, automigrate bool, vectorDBPath string, openAIConfig config.OpenAIConfig) (*Datastore, error) {
	dsn, vectorDBPath, isArchive, err := GetDatastorePaths(dsn, vectorDBPath)
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

	var vsdb *cg.DB
	if !isArchive {
		vsdb, err = cg.NewPersistentDB(vectorDBPath, false)
		if err != nil {
			return nil, err
		}
	} else {
		// Import from archive -> in-memory DB, not persisted back to the archive
		vsdb = cg.NewDB()
		if err := vsdb.ImportFromFile(vectorDBPath, ""); err != nil {
			return nil, fmt.Errorf("failed to import vector database: %w", err)
		}
	}

	var embeddingFunc cg.EmbeddingFunc
	if openAIConfig.APIType == "Azure" {
		// TODO: clean this up to support inputting the full deployment URL
		deployment := openAIConfig.AzureOpenAIConfig.Deployment
		if deployment == "" {
			deployment = openAIConfig.EmbeddingModel
		}

		deploymentURL, err := url.Parse(openAIConfig.APIBase)
		if err != nil || deploymentURL == nil {
			return nil, fmt.Errorf("failed to parse OpenAI Base URL %q: %w", openAIConfig.APIBase, err)
		}

		deploymentURL = deploymentURL.JoinPath("openai", "deployments", deployment)

		slog.Debug("Using Azure OpenAI API", "deploymentURL", deploymentURL.String(), "APIVersion", openAIConfig.APIVersion)

		embeddingFunc = cg.NewEmbeddingFuncAzureOpenAI(
			openAIConfig.APIKey,
			deploymentURL.String(),
			openAIConfig.APIVersion,
			"",
		)
	} else {
		embeddingFunc = cg.NewEmbeddingFuncOpenAICompat(
			openAIConfig.APIBase,
			openAIConfig.APIKey,
			openAIConfig.EmbeddingModel,
			z.Pointer(true),
		)
	}

	model, err := llm.NewOpenAI(openAIConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM: %w", err)
	}

	ds := &Datastore{
		LLM:         *model,
		Index:       idx,
		Vectorstore: chromem.New(vsdb, embeddingFunc),
	}

	if isArchive {
		return ds, nil
	}

	// Ensure default dataset exists
	defaultDS, err := ds.GetDataset(context.Background(), "default")
	if err != nil {
		return nil, fmt.Errorf("failed to ensure default dataset: %w", err)
	}

	if defaultDS == nil {
		err = ds.NewDataset(context.Background(), index.Dataset{ID: "default"})
		if err != nil {
			return nil, fmt.Errorf("failed to create default dataset: %w", err)
		}
	}

	return ds, nil
}

func (s *Datastore) ExportDatasetsToFile(ctx context.Context, path string, datasets ...string) error {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "knowledge-export-")
	if err != nil {
		return err
	}

	defer os.RemoveAll(tmpDir)

	if err = s.Index.ExportDatasetsToFile(ctx, tmpDir, datasets...); err != nil {
		return err
	}

	if err = s.Vectorstore.ExportCollectionsToFile(ctx, tmpDir, datasets...); err != nil {
		return err
	}

	finfo, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
	}

	// make sure target path is a file
	if finfo != nil && finfo.IsDir() {
		path = filepath.Join(path, fmt.Sprintf("knowledge-export-%s.zip", time.Now().Format("2006-01-02-15-04-05")))
	}

	// zip it up
	if err = zipDir(tmpDir, path); err != nil {
		return err
	}

	return nil
}

func zipDir(src, dst string) error {
	zipfile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(zipfile)
	defer w.Close()

	// Walk the file tree and add files to the zip archive.
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get the file information
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// Update the header name
			header.Name = filepath.Base(path)

			// Write the header
			writer, err := w.CreateHeader(header)
			if err != nil {
				return err
			}

			// If the file is not a directory, write the file to the archive
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
