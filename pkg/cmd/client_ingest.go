package cmd

import (
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

type ClientIngest struct {
	Client
}

func (s *ClientIngest) Customize(cmd *cobra.Command) {
	cmd.Use = "ingest <dataset-id> <path>"
	cmd.Short = "Ingest a file into a dataset"
	cmd.Args = cobra.ExactArgs(2)
}

func (s *ClientIngest) Run(cmd *cobra.Command, args []string) error {

	go s.Client.runServer(cmd)
	s.Client.waitForServer()

	url := s.Client.baseURL() + "/datasets/" + args[0] + "/ingest"

	filecontent, err := os.ReadFile(args[1])
	if err != nil {
		return err
	}

	b64data := base64.StdEncoding.EncodeToString(filecontent)

	req, err := http.NewRequest("POST", url, strings.NewReader(fmt.Sprintf("{\"content\": \"%s\", \"filename\": \"%s\"}", b64data, path.Base(args[1]))))
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
