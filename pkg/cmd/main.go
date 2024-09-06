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
		new(ClientImportDatasets),
		new(ClientEditDataset),
		new(Version),
	)
}

type Knowledge struct {
	Debug bool `usage:"Enable debug logging" env:"DEBUG" hidden:"true"`
}

func (c *Knowledge) Run(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

func (c *Knowledge) Customize(cmd *cobra.Command) {
	cmd.PersistentPreRun = func(cmd *cobra.Command, _ []string) {
		if c.Debug {
			_ = slog.SetLogLoggerLevel(slog.LevelDebug)
		}
	}
}

type Version struct{}

func (c *Version) Run(cmd *cobra.Command, _ []string) error {
	fmt.Println(version.Version)
	return nil
}
