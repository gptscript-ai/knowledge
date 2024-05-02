package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

type ClientRetrieve struct {
	Client
	TopK int `usage:"Number of sources to retrieve" default:"3"`
}

func (s *ClientRetrieve) Customize(cmd *cobra.Command) {
	cmd.Use = "retrieve <dataset-id> <query>"
	cmd.Short = "Retrieve sources for a query from a dataset"
	cmd.Args = cobra.ExactArgs(2)
}

func (s *ClientRetrieve) Run(cmd *cobra.Command, args []string) error {
	c, err := s.getClient()
	if err != nil {
		return err
	}

	datasetID := args[0]
	query := args[1]

	sources, err := c.Retrieve(cmd.Context(), datasetID, query)
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
