package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/gptscript-ai/knowledge/pkg/config"
	"github.com/gptscript-ai/knowledge/pkg/datastore"
	"github.com/gptscript-ai/knowledge/pkg/server"
)

// Server is the Server CLI command
type Server struct {
	ServerURL     string `usage:"Server URL" default:"http://localhost" env:"KNOW_SERVER_URL"`
	ServerPort    string `usage:"Server port" default:"8000" env:"KNOW_SERVER_PORT"`
	ServerAPIBase string `usage:"Server API base" default:"/v1" env:"KNOW_SERVER_API_BASE"`

	EmbeddingModelProvider string `usage:"Embedding model provider" default:"openai" env:"KNOW_EMBEDDING_MODEL_PROVIDER" name:"embedding-model-provider" koanf:"provider"`
	ConfigFile             string `usage:"Path to the configuration file" env:"KNOW_CONFIG_FILE" default:"" short:"c"`

	config.DatabaseConfig
	config.VectorDBConfig
}

func (s *Server) Customize(cmd *cobra.Command) {
	cmd.Use = "server"
	cmd.Short = "Run the Knowledge API Server"
	cmd.Long = `Run the Knowledge API Server.`

	cmd.Hidden = true
}

func (s *Server) Run(cmd *cobra.Command, _ []string) error {

	slog.Warn("The knowledge server is underdeveloped and lacking behind the standalone client right now, use at your own risk!") // FIXME: Bring the server on par with the standalone client and drop this warning

	cfg, err := config.LoadConfig(s.ConfigFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if s.EmbeddingModelProvider != "" {
		cfg.EmbeddingsConfig.Provider = s.EmbeddingModelProvider
	}

	ds, err := datastore.NewDatastore(s.DSN, s.AutoMigrate == "true", s.VectorDBConfig.VectorDBPath, cfg.EmbeddingsConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize datastore: %w", err)
	}

	ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	defer cancel()
	return server.NewServer(ds).Start(ctx, server.Config{
		ServerURL: s.ServerURL,
		Port:      s.ServerPort,
		APIBase:   s.ServerAPIBase,
	})
}
