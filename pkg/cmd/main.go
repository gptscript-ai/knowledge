package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/acorn-io/cmd"
	"github.com/gptscript-ai/knowledge/version"
	"github.com/spf13/cobra"
)

func init() {
	if os.Getenv("DEBUG") != "" {
		_ = slog.SetLogLoggerLevel(slog.LevelDebug)
	}
}

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
		new(ClientAskDir),
		new(ClientExportDatasets),
		new(Version),
	)
}

type Knowledge struct{}

func (c *Knowledge) Run(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

type Version struct{}

func (c *Version) Run(cmd *cobra.Command, _ []string) error {
	fmt.Println(version.Version)
	return nil
}
