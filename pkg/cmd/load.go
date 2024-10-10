package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gptscript-ai/knowledge/pkg/datastore/documentloader"
	"github.com/gptscript-ai/knowledge/pkg/datastore/filetypes"
	"github.com/spf13/cobra"
)

type ClientLoad struct {
	Loader string `usage:"Choose a document loader to use"`
}

func (s *ClientLoad) Customize(cmd *cobra.Command) {
	cmd.Use = "load <input> <output>"
	cmd.Short = "Load a file and transform it to markdown"
	cmd.Args = cobra.ExactArgs(2)
}

func (s *ClientLoad) Run(cmd *cobra.Command, args []string) error {
	input := args[0]
	output := args[1]

	inputBytes, err := os.ReadFile(input)
	if err != nil {
		return fmt.Errorf("failed to read input file %q: %w", input, err)
	}

	filetype, err := filetypes.GetFiletype(input, inputBytes)
	if err != nil {
		return fmt.Errorf("failed to get filetype for input file %q: %w", input, err)
	}

	var loader documentloader.LoaderFunc

	if s.Loader == "" {
		loader = documentloader.DefaultDocLoaderFunc(filetype, documentloader.DefaultDocLoaderFuncOpts{})
	} else {
		var err error
		loader, err = documentloader.GetDocumentLoaderFunc(s.Loader, nil)
		if err != nil {
			return fmt.Errorf("failed to get document loader function %q: %w", s.Loader, err)
		}
	}

	if loader == nil {
		return fmt.Errorf("unsupported file type %q", input)
	}

	docs, err := loader(cmd.Context(), bytes.NewReader(inputBytes))
	if err != nil {
		return fmt.Errorf("failed to load documents: %w", err)
	}

	var texts []string
	for _, doc := range docs {

		if len(doc.Content) == 0 {
			continue
		}

		metadata, err := json.Marshal(doc.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}

		content := fmt.Sprintf("!metadata %s\n%s", metadata, doc.Content)

		texts = append(texts, content)
	}

	text := strings.Join(texts, "\n---docbreak---\n")

	outputFile, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("failed to create output file %q: %w", output, err)
	}

	_, err = outputFile.WriteString(text)
	if err != nil {
		return fmt.Errorf("failed to write to output file %q: %w", output, err)
	}

	return nil
}
