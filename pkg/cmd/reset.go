package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gptscript-ai/knowledge/pkg/datastore"
	"github.com/spf13/cobra"
)

type ClientResetDatastore struct {
	Client
}

func (s *ClientResetDatastore) Customize(cmd *cobra.Command) {
	cmd.Use = "reset-datastore"
	cmd.Short = "Reset the knowledge datastore (WARNING: this deletes all datasets and ingested data)"
	cmd.Args = cobra.ExactArgs(0)
	cmd.Hidden = true
}

func (s *ClientResetDatastore) Run(cmd *cobra.Command, args []string) error {
	dsn, vectordbPath, _, err := datastore.GetDatastorePaths(s.DSN, s.VectorDBConfig.VectorDBPath)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(strings.TrimPrefix(dsn, "sqlite://")); err != nil {
		return fmt.Errorf("failed to remove database file: %w", err)
	}

	if err := os.RemoveAll(vectordbPath); err != nil {
		return fmt.Errorf("failed to remove vector database directory: %w", err)
	}

	fmt.Printf("Successfully reset datastore (DSN: %q, VectorDBPath: %q)\n", dsn, vectordbPath)
	return nil
}
