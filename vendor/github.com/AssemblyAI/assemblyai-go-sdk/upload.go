package assemblyai

import (
	"context"
	"io"
)

// Uploads an audio file to AssemblyAI's servers and returns the new URL.
//
// The uploaded file can only be accessed from AssemblyAI's servers. You can use
// the URL to transcript the file, but you can't download it.
//
// https://www.assemblyai.com/docs/API%20reference/upload
func (c *Client) Upload(ctx context.Context, data io.Reader) (string, error) {
	req, err := c.newRequest("POST", "/v2/upload", data)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	var result struct {
		UploadURL string `json:"upload_url"`
	}

	if _, err := c.do(ctx, req, &result); err != nil {
		return "", err
	}

	return result.UploadURL, nil
}
