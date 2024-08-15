//go:build mupdf

package openai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/acorn-io/z"
	"github.com/gen2brain/go-fitz"
	"github.com/gptscript-ai/knowledge/pkg/datastore/embeddings/openai"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"image"
	"image/png"
	"io"
	"log"
	"log/slog"
	"net/http"
)

type OpenAIOCR struct {
	OpenAI    openai.EmbeddingModelProviderOpenAI
	Prompt    string
	MaxTokens *int
}

type ImagePayload struct {
	URL string `json:"url"`
}

type MessageContent struct {
	Type     string       `json:"type"`
	Text     string       `json:"text,omitempty"`
	ImageURL ImagePayload `json:"image_url,omitempty"`
}

type Message struct {
	Role    string           `json:"role"`
	Content []MessageContent `json:"content"`
}

type Payload struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

func (o *OpenAIOCR) Load(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
	if o.Prompt == "" {
		o.Prompt = "What is in this image?"
	}

	images, err := convertPdfToImages(reader)
	if err != nil {
		log.Fatalf("Error converting PDF to images: %v", err)
	}

	docs := make([]vs.Document, len(images))
	for i, img := range images {
		slog.Debug("Processing PDF image", "page", i+1, "totalPages", len(images))
		base64Image, err := encodeImageToBase64(img)
		if err != nil {
			log.Fatalf("Error encoding image to base64: %v", err)
		}

		result, err := o.sendImageToOpenAI(base64Image)
		if err != nil {
			log.Fatalf("Error sending image to OpenAI: %v", err)
		}

		fmt.Printf("Result for page %d: %v\n", i+1, result)
		docs = append(docs, vs.Document{
			Metadata: map[string]interface{}{
				"page":       i + 1,
				"totalPages": len(images),
			},
			Content: fmt.Sprintf("%v", result),
		})
	}
	return docs, nil
}

func convertPdfToImages(reader io.Reader) ([]image.Image, error) {
	doc, err := fitz.NewFromReader(reader)
	if err != nil {
		return nil, err
	}
	defer doc.Close()

	var images []image.Image
	for i := 0; i < doc.NumPage(); i++ {
		img, err := doc.Image(i)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}

func encodeImageToBase64(img image.Image) (string, error) {
	var buffer bytes.Buffer
	err := png.Encode(&buffer, img)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

func (o *OpenAIOCR) sendImageToOpenAI(base64Image string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/chat/completions", o.OpenAI.BaseURL)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + o.OpenAI.APIKey,
	}

	if o.MaxTokens == nil {
		o.MaxTokens = z.Pointer(300)
	}

	if o.OpenAI.Model == "" {
		o.OpenAI.Model = "gpt-4o"
	}

	payload := Payload{
		Model: o.OpenAI.Model,
		Messages: []Message{
			{
				Role: "user",
				Content: []MessageContent{
					{Type: "text", Text: o.Prompt},
					{Type: "image_url", ImageURL: ImagePayload{URL: "data:image/png;base64," + base64Image}},
				},
			},
		},
		MaxTokens: *o.MaxTokens,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}
