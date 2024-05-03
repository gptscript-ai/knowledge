package cmd

import (
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/client"
	"github.com/spf13/cobra"
	"strings"
)

type ClientIngest struct {
	Client
	Dataset          string `usage:"Target Dataset ID" default:"default" env:"KNOW_TARGET_DATASET"`
	IgnoreExtensions string `usage:"Comma-separated list of file extensions to ignore" env:"KNOW_INGEST_IGNORE_EXTENSIONS"`
	Concurrency      int    `usage:"Number of concurrent ingestion processes" default:"10" env:"KNOW_INGEST_CONCURRENCY"`
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
	}

	err = c.IngestPaths(cmd.Context(), datasetID, ingestOpts, filePath)
	if err != nil {
		return err
	}

	fmt.Printf("Ingested %q into dataset %q\n", filePath, datasetID)
	return nil
}
