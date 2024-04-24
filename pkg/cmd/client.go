package cmd

import (
	"fmt"
	cmd2 "github.com/acorn-io/cmd"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type Client struct {
	ExternalServer bool `usage:"Run Server in parallel if it's not running elsewhere" default:"false" env:"KNOW_CLIENT_RUN_SERVER"`
	Server
}

func (s *Client) Customize(cmd *cobra.Command) {
	cmd.Use = "client"
	cmd.Short = "Interact with the Knowledge API server (which can be run in parallel)"
	cmd.Args = cobra.NoArgs
	cmd.AddCommand(cmd2.Command(&ClientCreateDataset{Client: *s}))
	cmd.AddCommand(cmd2.Command(&ClientIngest{Client: *s}))
	cmd.AddCommand(cmd2.Command(&ClientRetrieve{Client: *s}))
	cmd.AddCommand(cmd2.Command(&ClientDeleteDataset{Client: *s}))
}

func (s *Client) Run(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

func (s *Client) runServer(cmd *cobra.Command) {
	if s.ExternalServer {
		slog.Info("Using external server")
		return
	}
	slog.Info("Starting server")
	gin.SetMode(gin.ReleaseMode)
	if err := s.Server.Run(cmd, nil); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}

func (s *Client) baseURL() string {
	return fmt.Sprintf("%s:%s/%s", s.ServerURL, s.ServerPort, strings.Trim(s.ServerAPIBase, "/"))
}

func (s *Client) waitForServer() {
	// Wait for /docs enddpoint to be available and timeout after 10s
	url := s.baseURL() + "/docs"
	for i := 0; i < 10; i++ {
		if _, err := http.Get(url); err == nil {
			return
		}
		slog.Debug("Waiting for server to start")
		time.Sleep(1 * time.Second)
	}

}
