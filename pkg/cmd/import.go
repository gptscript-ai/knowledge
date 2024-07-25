package cmd

import (
	"github.com/spf13/cobra"
)

type ClientImportDatasets struct {
	Client
}

func (s *ClientImportDatasets) Customize(cmd *cobra.Command) {
	cmd.Use = "import <path> [<dataset-id>...]"
	cmd.Short = "Import one or more datasets from an archive (zip) (default: all datasets)"
	cmd.Long = `Import one or more datasets from an archive (zip) (default: all datasets).
## IMPORTANT: Embedding functions
   Embedding functions are not part of exported knowledge base archives, so you'll have to know the embedding function used to import the archive.
   This primarily concerns the choice of the embeddings provider (model).
   Note: This is only relevant if you plan to add more documents to the dataset after importing it.
`
	cmd.Args = cobra.MinimumNArgs(1)
}

func (s *ClientImportDatasets) Run(cmd *cobra.Command, args []string) error {
	c, err := s.getClient()
	if err != nil {
		return err
	}

	return c.ImportDatasets(cmd.Context(), args[0], args[1:]...)
}
