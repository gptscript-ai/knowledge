package cmd

import (
	"github.com/gptscript-ai/knowledge/pkg/client"
	"github.com/gptscript-ai/knowledge/pkg/datastore"
	"github.com/gptscript-ai/knowledge/pkg/types"
)

type Client struct {
	Server string `usage:"URL of the Knowledge API Server" default:"" env:"KNOW_SERVER_URL"`
	types.OpenAIConfig
	types.DatabaseConfig
}

func (s *Client) getClient() (client.Client, error) {
	if s.Server == "" || s.Server == "standalone" {
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
