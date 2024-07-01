package cmd

import (
	"archive/zip"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/client"
	"github.com/gptscript-ai/knowledge/pkg/config"
	"github.com/gptscript-ai/knowledge/pkg/datastore"
	"github.com/gptscript-ai/knowledge/pkg/datastore/types"
	"io"
	"os"
	"path/filepath"
)

type Client struct {
	Server string `usage:"URL of the Knowledge API Server" default:"" env:"KNOW_SERVER_URL"`
	config.OpenAIConfig
	config.DatabaseConfig
	config.VectorDBConfig
}

type ClientFlowsConfig struct {
	FlowsFile string `usage:"Path to a YAML/JSON file containing ingestion/retrieval flows" env:"KNOW_FLOWS_FILE"`
	Flow      string `usage:"Flow name" env:"KNOW_FLOW"`
}

func (s *Client) getClientFromArchive(archive string) (client.Client, error) {
	// unpack to tempdir
	tmpDir, err := os.MkdirTemp(os.TempDir(), "knowledge-retrieve-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	r, err := zip.OpenReader(archive)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	if len(r.File) != 2 {
		return nil, fmt.Errorf("knowledge archive must contain exactly two files, found %d", len(r.File))
	}

	dbFile := ""
	vectorStoreFile := ""
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()

		path := filepath.Join(tmpDir, f.Name)
		if f.FileInfo().IsDir() {
			continue
		}

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return nil, err
		}
		defer f.Close()

		if _, err := io.Copy(f, rc); err != nil {
			return nil, err
		}
		_ = f.Close()
		_ = rc.Close()

		// FIXME: this should not be static as we may support multiple (vector) DBs at some point
		if filepath.Ext(f.Name()) == ".db" {
			dbFile = path
		} else if filepath.Ext(f.Name()) == ".gob" {
			vectorStoreFile = path
		}
	}

	if dbFile == "" || vectorStoreFile == "" {
		return nil, fmt.Errorf("knowledge archive must contain exactly one .db and one .gob file")
	}

	s.DSN = types.ArchivePrefix + dbFile
	s.VectorDBPath = types.ArchivePrefix + vectorStoreFile

	return s.getClient()
}

func (s *Client) getClient() (client.Client, error) {

	if s.Server == "" || s.Server == "standalone" {
		ds, err := datastore.NewDatastore(s.DSN, s.AutoMigrate == "true", s.VectorDBConfig.VectorDBPath, s.OpenAIConfig)
		if err != nil {
			return nil, err
		}
		c, err := client.NewStandaloneClient(ds)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
	return client.NewDefaultClient(s.Server), nil
}
