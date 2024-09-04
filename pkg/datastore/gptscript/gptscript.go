package gptscript

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/gptscript-ai/go-gptscript"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
)

type GPTScriptTool struct {
	gs     *gptscript.GPTScript
	Ref    string
	Params map[string]any
}

func NewGPTScriptTool(gs *gptscript.GPTScript) *GPTScriptTool {
	return &GPTScriptTool{
		gs: gs,
	}
}

func (g *GPTScriptTool) AsDocumentloader() func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {

	return func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {

		tmpFile, err := os.CreateTemp(os.TempDir(), "gptscript-knowledge-")
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(tmpFile, reader)

		filename := tmpFile.Name()

		opts := gptscript.Options{}

		if g.Params != nil {
			for k, v := range g.Params {
				if s, ok := v.(string); ok {
					if s == "$KNOW_LOAD_FILENAME" || s == "${KNOW_LOAD_FILENAME}" {
						g.Params[k] = filename
					}
				}
			}
			paramJSON, err := json.Marshal(g.Params)
			if err != nil {
				return nil, err
			}
			opts.Input = string(paramJSON)
		}

		run, err := g.gs.Run(ctx, g.Ref, opts)
		if err != nil {
			return nil, err
		}

		output, err := run.Bytes()
		if err != nil {
			return nil, err
		}

		docs := []vs.Document{
			{
				Content: string(output),
			},
		}

		return docs, nil
	}
}
