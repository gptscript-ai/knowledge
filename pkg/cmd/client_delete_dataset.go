package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
)

type ClientDeleteDataset struct {
	Client
}

func (s *ClientDeleteDataset) Customize(cmd *cobra.Command) {
	cmd.Use = "delete-dataset <dataset-id>"
	cmd.Short = "Delete a dataset"
	cmd.Args = cobra.ExactArgs(1)
}

func (s *ClientDeleteDataset) Run(cmd *cobra.Command, args []string) error {

	go s.Client.runServer(cmd)
	s.Client.waitForServer()

	url := s.Client.baseURL() + "/datasets/" + args[0]

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

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
