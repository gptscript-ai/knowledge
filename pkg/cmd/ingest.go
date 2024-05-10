package cmd

import (
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/client"
	"github.com/gptscript-ai/knowledge/pkg/datastore"
	"github.com/spf13/cobra"
	"strings"
)

type ClientIngest struct {
	Client
	Dataset string `usage:"Target Dataset ID" short:"d" default:"default" env:"KNOW_TARGET_DATASET"`
	ClientIngestOpts
	datastore.TextSplitterOpts
}

type ClientIngestOpts struct {
	IgnoreExtensions string `usage:"Comma-separated list of file extensions to ignore" env:"KNOW_INGEST_IGNORE_EXTENSIONS"`
	Concurrency      int    `usage:"Number of concurrent ingestion processes" short:"c" default:"10" env:"KNOW_INGEST_CONCURRENCY"`
	Recursive        bool   `usage:"Recursively ingest directories" short:"r" default:"false" env:"KNOW_INGEST_RECURSIVE"`
}

func (s *ClientIngest) Customize(cmd *cobra.Command) {
	cmd.Use = "ingest [--dataset <dataset-id>] <path>"
	cmd.Short = "Ingest a file/directory into a dataset (non-recursive)"
	cmd.Args = cobra.ExactArgs(1)
}

func (s *ClientIngest) Run(cmd *cobra.Command, args []string) error {
	c, err := s.getClient()
	if err != nil {
		return err
	}

	datasetID := s.Dataset
	filePath := args[0]

	ingestOpts := &client.IngestPathsOpts{
		IgnoreExtensions: strings.Split(s.IgnoreExtensions, ","),
		Concurrency:      s.Concurrency,
		Recursive:        s.Recursive,
		TextSplitterOpts: &s.TextSplitterOpts,
	}

	filesIngested, err := c.IngestPaths(cmd.Context(), datasetID, ingestOpts, filePath)
	if err != nil {
		return err
	}

	fmt.Printf("Ingested %d files from %q into dataset %q\n", filesIngested, filePath, datasetID)
	return nil
}
