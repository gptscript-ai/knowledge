package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/acorn-io/z"
	"github.com/gptscript-ai/knowledge/pkg/client"
	flowconfig "github.com/gptscript-ai/knowledge/pkg/flows/config"
	"github.com/spf13/cobra"
	"log/slog"
	"path/filepath"
	"strings"
)

type ClientAskDir struct {
	Client
	Path string `usage:"Path to the directory to query" short:"p" default:"./knowledge"`
	ClientIngestOpts
	ClientRetrieveOpts
	FlowsFile string `usage:"Path to a YAML/JSON file containing ingestion/retrieval flows" env:"KNOW_FLOWS_FILE"`
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

	if s.FlowsFile != "" {
		abspath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to get absolute path from %q: %w", path, err)
		}

		datasetID := client.HashPath(abspath)

		slog.Debug("Loading ingestion flows from config", "flows_file", s.FlowsFile, "dataset", datasetID)
		flowCfg, err := flowconfig.FromFile(s.FlowsFile)
		if err != nil {
			return err
		}
		flow, err := flowCfg.ForDataset(datasetID) // get flow for the dataset
		if err != nil {
			return err
		}

		for _, ingestionFlowConfig := range flow.Ingestion {
			ingestionFlow, err := ingestionFlowConfig.AsIngestionFlow()
			if err != nil {
				return err
			}
			ingestOpts.IngestionFlows = append(ingestOpts.IngestionFlows, z.Dereference(ingestionFlow))
		}

		// TODO: add retrieval flows here

		slog.Debug("Loaded ingestion flows from config", "flows_file", s.FlowsFile, "dataset", datasetID, "flows", len(ingestOpts.IngestionFlows))

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
