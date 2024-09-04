package gptscript

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/gptscript-ai/go-gptscript"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
)

func GSRead(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
	gs, err := gptscript.NewGPTScript(gptscript.GlobalOptions{})
	if err != nil {
		return nil, err
	}

	tmpFile, err := os.CreateTemp(os.TempDir(), "gptscript-ai-")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(tmpFile, reader)

	filename := tmpFile.Name()

	run, err := gs.Run(ctx, "sys.read", gptscript.Options{
		Input: fmt.Sprintf("{\"filename\": \"%s\"}", filename),
	})
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
