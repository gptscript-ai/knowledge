package cmd

import (
	"fmt"
	"github.com/acorn-io/z"
	"github.com/spf13/cobra"
	"log/slog"
	"strings"

	"github.com/gptscript-ai/knowledge/pkg/client"
	"github.com/gptscript-ai/knowledge/pkg/datastore/textsplitter"
	flowconfig "github.com/gptscript-ai/knowledge/pkg/flows/config"
)

type ClientIngest struct {
	Client
	Dataset string `usage:"Target Dataset ID" short:"d" default:"default" env:"KNOW_TARGET_DATASET"`
	Prune   bool   `usage:"Prune deleted files" env:"KNOW_INGEST_PRUNE"`
	ClientIngestOpts
	textsplitter.TextSplitterOpts
	ClientFlowsConfig
}

type ClientIngestOpts struct {
	IgnoreExtensions      string `usage:"Comma-separated list of file extensions to ignore" env:"KNOW_INGEST_IGNORE_EXTENSIONS"`
	IgnoreFile            string `usage:"Path to a .gitignore style file" env:"KNOW_INGEST_IGNORE_FILE"`
	IncludeHidden         bool   `usage:"Include hidden files and directories" default:"false" env:"KNOW_INGEST_INCLUDE_HIDDEN"`
	Concurrency           int    `usage:"Number of concurrent ingestion processes" default:"10" env:"KNOW_INGEST_CONCURRENCY"`
	NoRecursive           bool   `usage:"Don't recursively ingest directories" default:"false" env:"KNOW_NO_INGEST_RECURSIVE"`
	NoCreateDataset       bool   `usage:"Do NOT create the dataset if it doesn't exist" default:"true" env:"KNOW_INGEST_NO_CREATE_DATASET"`
	DeduplicationFuncName string `usage:"Name of the deduplication function to use" name:"dedupe-func" env:"KNOW_INGEST_DEDUPE_FUNC"`
}

func (s *ClientIngest) Customize(cmd *cobra.Command) {
	cmd.Use = "ingest [--dataset <dataset-id>] <path>"
	cmd.Short = "Ingest a file/directory into a dataset"
	cmd.Long = `Ingest a file or directory into a dataset.

## Important Note

The first time you ingest something into a dataset, the embedding function (model provider) you chose will be attached to that dataset.
After that, the client must always use that same embedding function to ingest into this dataset.
Usually, this only concerns the choice of the model, as that commonly defines the embedding dimensionality.
This is a constraint of the Vector Database and Similarity Search, as different models yield differently sized embedding vectors and also represent the semantics differently.
`
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
		IgnoreExtensions:    strings.Split(s.IgnoreExtensions, ","),
		Concurrency:         s.Concurrency,
		Recursive:           !s.NoRecursive,
		TextSplitterOpts:    &s.TextSplitterOpts,
		IgnoreFile:          s.IgnoreFile,
		IncludeHidden:       s.IncludeHidden,
		IsDuplicateFuncName: s.DeduplicationFuncName,
		Prune:               s.Prune,
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
