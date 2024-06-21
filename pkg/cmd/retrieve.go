package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	vserr "github.com/gptscript-ai/knowledge/pkg/vectorstore/errors"
	"log/slog"

	"github.com/gptscript-ai/knowledge/pkg/datastore"
	flowconfig "github.com/gptscript-ai/knowledge/pkg/flows/config"
	"github.com/spf13/cobra"
)

type ClientRetrieve struct {
	Client
	Dataset string `usage:"Target Dataset ID" short:"d" default:"default" env:"KNOW_TARGET_DATASET"`
	ClientRetrieveOpts
	ClientFlowsConfig
}

type ClientRetrieveOpts struct {
	TopK int `usage:"Number of sources to retrieve" short:"k" default:"10"`
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
		// An empty collection is not a hard error - the LLM session can "recover" from it
		if errors.Is(err, vserr.ErrCollectionEmpty) {
			fmt.Printf("Dataset %q does not contain any documents\n", datasetID)
			return nil
		}
		return err
	}

	jsonSources, err := json.Marshal(sources)
	if err != nil {
		return err
	}

	fmt.Printf("Retrieved the following %d sources for the query %q from dataset %q: %s\n", len(sources), query, datasetID, jsonSources)

	return nil
}
