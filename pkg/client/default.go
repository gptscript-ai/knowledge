package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/acorn-io/z"
	"github.com/gptscript-ai/knowledge/pkg/datastore"
	"github.com/gptscript-ai/knowledge/pkg/index"
	"github.com/gptscript-ai/knowledge/pkg/types"
	"github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type DefaultClient struct {
	ServerURL string
}

func NewDefaultClient(serverURL string) *DefaultClient {
	return &DefaultClient{
		ServerURL: strings.TrimSuffix(serverURL, "/"),
	}
}

func (c *DefaultClient) CreateDataset(_ context.Context, datasetID string) (types.Dataset, error) {
	ds := types.Dataset{
		ID:             datasetID,
		EmbedDimension: nil,
	}

	payload, err := json.Marshal(ds)
	if err != nil {
		return types.Dataset{}, err
	}

	resp, err := c.request(http.MethodPost, "/datasets/create", bytes.NewReader(payload))
	if err != nil {
		return types.Dataset{}, err
	}

	var dataset types.Dataset
	err = json.Unmarshal(resp, &dataset)
	if err != nil {
		return types.Dataset{}, err
	}

	return dataset, nil
}

func (c *DefaultClient) DeleteDataset(_ context.Context, datasetID string) error {
	_, err := c.request(http.MethodDelete, fmt.Sprintf("/datasets/%s", datasetID), nil)
	return err
}

func (c *DefaultClient) GetDataset(_ context.Context, datasetID string) (*index.Dataset, error) {
	resp, err := c.request(http.MethodGet, fmt.Sprintf("/datasets/%s", datasetID), nil)
	if err != nil {
		return nil, err
	}

	var dataset *index.Dataset
	err = json.Unmarshal(resp, dataset)
	if err != nil {
		return nil, err
	}

	return dataset, nil
}

func (c *DefaultClient) ListDatasets(_ context.Context) ([]types.Dataset, error) {
	resp, err := c.request(http.MethodGet, "/datasets", nil)
	if err != nil {
		return []types.Dataset{}, err
	}

	var datasets []types.Dataset
	err = json.Unmarshal(resp, &datasets)
	if err != nil {
		return []types.Dataset{}, err
	}

	return datasets, nil
}

func (c *DefaultClient) Ingest(_ context.Context, datasetID string, data []byte, opts datastore.IngestOpts) ([]string, error) {
	payload := types.Ingest{
		Filename: opts.Filename,
		Content:  base64.StdEncoding.EncodeToString(data),
	}
	if opts.FileMetadata != nil {
		payload.FileMetadata = opts.FileMetadata
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := c.request(http.MethodPost, fmt.Sprintf("/datasets/%s/ingest", datasetID), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	var res types.IngestResponse
	err = json.Unmarshal(resp, &res)
	if err != nil {
		return nil, err
	}

	return res.Documents, nil
}

func (c *DefaultClient) IngestPaths(ctx context.Context, datasetID string, opts *IngestPathsOpts, paths ...string) error {
	ingestFile := func(path string) error {
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}
		// encode to []byte
		b64 := make([]byte, base64.StdEncoding.EncodedLen(len(content)))
		base64.StdEncoding.Encode(b64, content)

		// Gather metadata
		finfo, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("failed to stat file %s: %w", path, err)
		}

		abspath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for %s: %w", path, err)
		}

		payload := datastore.IngestOpts{
			Filename: z.Pointer(filepath.Base(path)),
			FileMetadata: &index.FileMetadata{
				Name:         filepath.Base(path),
				AbsolutePath: abspath,
				Size:         finfo.Size(),
				ModifiedAt:   finfo.ModTime(),
			},
			IsDuplicateFuncName: "file_metadata",
		}
		_, err = c.Ingest(ctx, datasetID, b64, payload)
		return err
	}

	return ingestPaths(ctx, opts, ingestFile, paths...)
}

func (c *DefaultClient) DeleteDocuments(_ context.Context, datasetID string, documentIDs ...string) error {
	for _, documentID := range documentIDs {
		_, err := c.request(http.MethodDelete, fmt.Sprintf("/datasets/%s/documents/%s", datasetID, documentID), nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *DefaultClient) Retrieve(_ context.Context, datasetID string, query string) ([]vectorstore.Document, error) {
	data, err := json.Marshal(types.Query{Prompt: query})
	if err != nil {
		return nil, err
	}

	resp, err := c.request(http.MethodPost, fmt.Sprintf("/datasets/%s/retrieve", datasetID), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	var docs []vectorstore.Document
	err = json.Unmarshal(resp, &docs)
	if err != nil {
		return nil, err
	}

	return docs, nil
}

func (c *DefaultClient) request(method, path string, body io.Reader) ([]byte, error) {
	url := c.ServerURL + path

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("API request failed: %s", res.Status)
	}

	if res.Body != nil {
		defer res.Body.Close()
		return io.ReadAll(res.Body)
	}

	return nil, nil
}
