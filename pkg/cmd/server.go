package cmd

import (
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/config"
	"github.com/gptscript-ai/knowledge/pkg/datastore"
	"github.com/gptscript-ai/knowledge/pkg/server"
	"github.com/spf13/cobra"
	"os/signal"
	"syscall"
)

// Server is the Server CLI command
type Server struct {
	ServerURL     string `usage:"Server URL" default:"http://localhost" env:"KNOW_SERVER_URL"`
	ServerPort    string `usage:"Server port" default:"8000" env:"KNOW_SERVER_PORT"`
	ServerAPIBase string `usage:"Server API base" default:"/v1" env:"KNOW_SERVER_API_BASE"`

	config.OpenAIConfig
	config.DatabaseConfig
	config.VectorDBConfig
}

func (s *Server) Run(cmd *cobra.Command, _ []string) error {
	ds, err := datastore.NewDatastore(s.DSN, s.AutoMigrate == "true", s.VectorDBConfig.VectorDBPath, s.OpenAIConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize datastore: %w", err)
	}

	ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	defer cancel()
	return server.NewServer(ds, s.OpenAIConfig).Start(ctx, server.Config{
		ServerURL: s.ServerURL,
		Port:      s.ServerPort,
		APIBase:   s.ServerAPIBase,
	})
}
