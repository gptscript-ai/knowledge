package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/client"
	"github.com/spf13/cobra"
	"strings"
)

type ClientAskDir struct {
	Client
	Path string `usage:"Path to the directory to query" short:"p" default:"./knowledge"`
	ClientIngestOpts
	ClientRetrieveOpts
}

func (s *ClientAskDir) Customize(cmd *cobra.Command) {
	cmd.Use = "askdir [--path <path>] <query>"
	cmd.Short = "Retrieve sources for a query from a dataset generated from a directory"
	cmd.Args = cobra.ExactArgs(1)
}

func (s *ClientAskDir) Run(cmd *cobra.Command, args []string) error {
	c, err := s.getClient()
	if err != nil {
		return err
	}

	path := s.Path
	query := args[0]

	ingestOpts := &client.IngestPathsOpts{
		IgnoreExtensions: strings.Split(s.IgnoreExtensions, ","),
		Concurrency:      s.Concurrency,
		Recursive:        s.Recursive,
	}

	retrieveOpts := &client.RetrieveOpts{
		TopK: s.TopK,
	}

	sources, err := c.AskDirectory(cmd.Context(), path, query, ingestOpts, retrieveOpts)
	if err != nil {
		return fmt.Errorf("failed to retrieve sources: %w", err)
	}

	if len(sources) == 0 {
		fmt.Printf("No sources found for the query %q from path %q\n", query, path)
		return nil
	}

	jsonSources, err := json.Marshal(sources)
	if err != nil {
		return err
	}

	fmt.Printf("Retrieved the following %d sources for the query %q from path %q: %s\n", len(sources), query, path, jsonSources)

	return nil
}
