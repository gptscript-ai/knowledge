package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

type ClientExportDatasets struct {
	Client
	Output string `usage:"Output path" default:"."`
}

func (s *ClientExportDatasets) Customize(cmd *cobra.Command) {
	cmd.Use = "export <dataset-id> [<dataset-id>...]"
	cmd.Short = "Export one or more datasets as an archive (zip)"
	cmd.Args = cobra.MinimumNArgs(1)
}

func (s *ClientExportDatasets) Run(cmd *cobra.Command, args []string) error {
	c, err := s.getClient()
	if err != nil {
		return err
	}

	for _, datasetID := range args {
		ds, err := c.GetDataset(cmd.Context(), datasetID)
		if err != nil {
			return err
		}

		if ds.ID == "" {
			return fmt.Errorf("dataset %q not found", datasetID)
		}
	}

	return c.ExportDatasets(cmd.Context(), s.Output, args...)
}
