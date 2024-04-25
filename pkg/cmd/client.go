package cmd

import (
	cmd2 "github.com/acorn-io/cmd"
	"github.com/gptscript-ai/knowledge/pkg/client"
	"github.com/gptscript-ai/knowledge/pkg/datastore"
	"github.com/gptscript-ai/knowledge/pkg/types"
	"github.com/spf13/cobra"
)

type Client struct {
	Server string `usage:"URL of the Knowledge API Server (set to 'standalone' for standalone client)" default:"http://localhost:8000/v1" env:"KNOW_SERVER_URL"`
	types.OpenAIConfig
	types.DatabaseConfig
}

func (s *Client) Customize(cmd *cobra.Command) {
	cmd.Use = "client"
	cmd.Short = "Create and Retrieve Knowledge"
	cmd.Args = cobra.NoArgs
	cmd.AddCommand(cmd2.Command(&ClientCreateDataset{Client: *s}))
	cmd.AddCommand(cmd2.Command(&ClientIngest{Client: *s}))
	cmd.AddCommand(cmd2.Command(&ClientRetrieve{Client: *s}))
	cmd.AddCommand(cmd2.Command(&ClientDeleteDataset{Client: *s}))
}

func (s *Client) Run(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

func (s *Client) getClient() (client.Client, error) {
	if s.Server == "standalone" {
		ds, err := datastore.NewDatastore(s.DSN, s.AutoMigrate == "true", s.OpenAIConfig)
		if err != nil {
			return nil, err
		}
		c, err := client.NewStandaloneClient(ds)
		if err != nil {
			return nil, err
		}
		return c, c.Datastore.Index.AutoMigrate()
	}
	return client.NewDefaultClient(s.Server), nil
}
