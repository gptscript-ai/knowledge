package cmd

import (
	"github.com/gptscript-ai/knowledge/pkg/db"
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

	OpenAIAPIBase        string `usage:"OpenAI API base" default:"http://localhost:8080/v1" env:"KNOW_OPENAI_API_BASE"` // clicky-chats
	OpenAIAPIKey         string `usage:"OpenAI API key (not required if used with clicky-chats)" default:"sk-foo" env:"KNOW_OPENAI_API_KEY"`
	OpenAIEmbeddingModel string `usage:"OpenAI Embedding model" default:"text-embedding-ada-002" env:"KNOW_OPENAI_EMBEDDING_MODEL"`

	DSN         string `usage:"Server database connection string" default:"sqlite://knowledge.db" env:"KNOW_DSN"`
	AutoMigrate string `usage:"Auto migrate database" default:"true" env:"KNOW_AUTO_MIGRATE"`
}

func (s *Server) Run(cmd *cobra.Command, _ []string) error {
	gormDB, err := db.New(s.DSN, s.AutoMigrate == "true")
	if err != nil {
		return err
	}

	oaiCfg := server.OpenAIConfig{
		APIBase:        s.OpenAIAPIBase,
		APIKey:         s.OpenAIAPIKey,
		EmbeddingModel: s.OpenAIEmbeddingModel,
	}

	ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	defer cancel()
	return server.NewServer(gormDB, oaiCfg).Start(ctx, server.Config{
		ServerURL: s.ServerURL,
		Port:      s.ServerPort,
		APIBase:   s.ServerAPIBase,
	})
}
