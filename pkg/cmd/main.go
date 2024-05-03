package cmd

import (
	"github.com/acorn-io/cmd"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	return cmd.Command(
		&Knowledge{},
		new(Server),
		new(ClientCreateDataset),
		new(ClientGetDataset),
		new(ClientListDatasets),
		new(ClientIngest),
		new(ClientDeleteDataset),
		new(ClientRetrieve),
		new(ClientResetDatastore),
	)
}

type Knowledge struct{}

func (c *Knowledge) Run(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
