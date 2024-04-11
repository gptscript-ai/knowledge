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
	ServerPort    string `usage:"Server port" default:"8080" env:"KNOW_SERVER_PORT"`
	ServerAPIBase string `usage:"Server API base" default:"/v1" env:"KNOW_SERVER_API_BASE"`

	DSN         string `usage:"Server database connection string" default:"sqlite://knowledge.db" env:"KNOW_DSN"`
	AutoMigrate string `usage:"Auto migrate database" default:"true" env:"KNOW_AUTO_MIGRATE"`
}

func (s *Server) Run(cmd *cobra.Command, _ []string) error {
	gormDB, err := db.New(s.DSN, s.AutoMigrate == "true")
	if err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	defer cancel()
	return server.NewServer(gormDB).Start(ctx, server.Config{
		ServerURL: s.ServerURL,
		Port:      s.ServerPort,
		APIBase:   s.ServerAPIBase,
	})
}
