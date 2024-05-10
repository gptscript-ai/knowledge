package cmd

import (
	"github.com/gptscript-ai/knowledge/pkg/client"
	"github.com/gptscript-ai/knowledge/pkg/config"
	"github.com/gptscript-ai/knowledge/pkg/datastore"
)

type Client struct {
	Server string `usage:"URL of the Knowledge API Server" default:"" env:"KNOW_SERVER_URL"`
	config.OpenAIConfig
	config.DatabaseConfig
	config.VectorDBConfig
}

func (s *Client) getClient() (client.Client, error) {
	if s.Server == "" || s.Server == "standalone" {
		ds, err := datastore.NewDatastore(s.DSN, s.AutoMigrate == "true", s.VectorDBConfig.VectorDBPath, s.OpenAIConfig)
		if err != nil {
			return nil, err
		}
		c, err := client.NewStandaloneClient(ds)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
	return client.NewDefaultClient(s.Server), nil
}
