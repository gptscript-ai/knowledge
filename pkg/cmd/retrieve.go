package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore"
	flowconfig "github.com/gptscript-ai/knowledge/pkg/flows/config"
	"github.com/spf13/cobra"
	"log/slog"
)

type ClientRetrieve struct {
	Client
	Dataset string `usage:"Target Dataset ID" short:"d" default:"default" env:"KNOW_TARGET_DATASET"`
	ClientRetrieveOpts
	FlowsFile string `usage:"Path to a YAML/JSON file containing retrieval flows" env:"KNOW_FLOWS_FILE"`
}

type ClientRetrieveOpts struct {
	TopK int `usage:"Number of sources to retrieve" short:"k" default:"5"`
}

func (s *ClientRetrieve) Customize(cmd *cobra.Command) {
	cmd.Use = "retrieve [--dataset <dataset-id>] <query>"
	cmd.Short = "Retrieve sources for a query from a dataset"
	cmd.Args = cobra.ExactArgs(1)
}

func (s *ClientRetrieve) Run(cmd *cobra.Command, args []string) error {
	c, err := s.getClient()
	if err != nil {
		return err
	}

	datasetID := s.Dataset
	query := args[0]

	retrieveOpts := datastore.RetrieveOpts{
		TopK: s.TopK,
	}

	if s.FlowsFile != "" {
		slog.Debug("Loading retrieval flows from config", "flows_file", s.FlowsFile, "dataset", datasetID)
		flowCfg, err := flowconfig.FromFile(s.FlowsFile)
		if err != nil {
			return err
		}
		flow, err := flowCfg.ForDataset(datasetID) // get flow for the dataset
		if err != nil {
			return err
		}

		if flow.Retrieval == nil {
			slog.Info("No retrieval config in assigned flow", "flows_file", s.FlowsFile, "dataset", datasetID)
		} else {
			rf, err := flow.Retrieval.AsRetrievalFlow()
			if err != nil {
				return err
			}
			retrieveOpts.RetrievalFlow = rf
			slog.Debug("Loaded retrieval flow from config", "flows_file", s.FlowsFile, "dataset", datasetID)
		}
	}

	sources, err := c.Retrieve(cmd.Context(), datasetID, query, retrieveOpts)
	if err != nil {
		return err
	}

	jsonSources, err := json.Marshal(sources)
	if err != nil {
		return err
	}

	fmt.Printf("Retrieved the following %d sources for the query %q from dataset %q: %s\n", len(sources), query, datasetID, jsonSources)

	return nil
}
