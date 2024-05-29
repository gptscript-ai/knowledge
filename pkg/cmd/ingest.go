package cmd

import (
	"fmt"
	"github.com/acorn-io/z"
	"github.com/gptscript-ai/knowledge/pkg/client"
	"github.com/gptscript-ai/knowledge/pkg/datastore/textsplitter"
	flowconfig "github.com/gptscript-ai/knowledge/pkg/flows/config"
	"github.com/spf13/cobra"
	"log/slog"
	"strings"
)

type ClientIngest struct {
	Client
	Dataset string `usage:"Target Dataset ID" short:"d" default:"default" env:"KNOW_TARGET_DATASET"`
	ClientIngestOpts
	textsplitter.TextSplitterOpts
	ClientFlowsConfig
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

	if s.FlowsFile != "" {
		slog.Debug("Loading ingestion flows from config", "flows_file", s.FlowsFile, "dataset", datasetID)
		flowCfg, err := flowconfig.FromFile(s.FlowsFile)
		if err != nil {
			return err
		}

		var flow *flowconfig.FlowConfigEntry
		if s.Flow != "" {
			flow, err = flowCfg.GetFlow(s.Flow)
			if err != nil {
				return err
			}
		} else {
			flow, err = flowCfg.ForDataset(datasetID) // get flow for the dataset
			if err != nil {
				return err
			}
		}

		for _, ingestionFlowConfig := range flow.Ingestion {
			ingestionFlow, err := ingestionFlowConfig.AsIngestionFlow()
			if err != nil {
				return err
			}
			ingestOpts.IngestionFlows = append(ingestOpts.IngestionFlows, z.Dereference(ingestionFlow))
		}

		slog.Debug("Loaded ingestion flows from config", "flows_file", s.FlowsFile, "dataset", datasetID, "flows", len(ingestOpts.IngestionFlows))
	}

	filesIngested, err := c.IngestPaths(cmd.Context(), datasetID, ingestOpts, filePath)
	if err != nil {
		return err
	}

	fmt.Printf("Ingested %d files from %q into dataset %q\n", filesIngested, filePath, datasetID)
	return nil
}
