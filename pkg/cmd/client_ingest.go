package cmd

import (
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/client"
	"github.com/spf13/cobra"
	"strings"
)

type ClientIngest struct {
	Client
	IgnoreExtensions string `usage:"Comma-separated list of file extensions to ignore" env:"KNOW_INGEST_IGNORE_EXTENSIONS"`
}

func (s *ClientIngest) Customize(cmd *cobra.Command) {
	cmd.Use = "ingest <dataset-id> <path>"
	cmd.Short = "Ingest a file/directory into a dataset (non-recursive)"
	cmd.Args = cobra.ExactArgs(2)
}

func (s *ClientIngest) Run(cmd *cobra.Command, args []string) error {
	c, err := s.getClient()
	if err != nil {
		return err
	}

	datasetID := args[0]
	filePath := args[1]

	ingestOpts := &client.IngestPathsOpts{
		IgnoreExtensions: strings.Split(s.IgnoreExtensions, ","),
	}

	err = c.IngestPaths(cmd.Context(), datasetID, ingestOpts, filePath)
	if err != nil {
		return err
	}

	fmt.Printf("Ingested %q into dataset %q\n", filePath, datasetID)
	return nil
}
