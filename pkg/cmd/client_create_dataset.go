package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"strings"
)

type ClientCreateDataset struct {
	Client
	EmbedDim int `usage:"Embedding dimension" default:"1536"`
}

func (s *ClientCreateDataset) Customize(cmd *cobra.Command) {
	cmd.Use = "create-dataset <dataset-id>"
	cmd.Short = "Create a new dataset"
	cmd.Args = cobra.ExactArgs(1)
}

func (s *ClientCreateDataset) Run(cmd *cobra.Command, args []string) error {

	go s.Client.runServer(cmd)
	s.Client.waitForServer()

	url := s.Client.baseURL() + "/datasets/create"

	payload := strings.NewReader(fmt.Sprintf("{\"id\": \"%s\", \"embed_dim\": %d}", args[0], s.EmbedDim))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("ERROR: server returned status %d", res.StatusCode)
	}

	if res.Body != nil {
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		fmt.Println(string(body))
	}

	return nil

}
