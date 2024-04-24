package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"strings"
)

type ClientRetrieve struct {
	Client
	TopK int `usage:"Number of sources to retrieve" default:"3"`
}

func (s *ClientRetrieve) Customize(cmd *cobra.Command) {
	cmd.Use = "retrieve <dataset-id> <query>"
	cmd.Short = "Retrieve sources for a query from a dataset"
	cmd.Args = cobra.ExactArgs(2)
}

func (s *ClientRetrieve) Run(cmd *cobra.Command, args []string) error {

	go s.Client.runServer(cmd)
	s.Client.waitForServer()

	url := s.Client.baseURL() + "/datasets/" + args[0] + "/retrieve"

	req, err := http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf("{\"prompt\": \"%s\", \"topk\": %d}", args[1], s.TopK)))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("ERROR: server returned status %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil

}
