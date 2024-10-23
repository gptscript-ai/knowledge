package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/gptscript-ai/knowledge/pkg/datastore/documentloader"
	"github.com/gptscript-ai/knowledge/pkg/datastore/documentloader/structured"
	"github.com/gptscript-ai/knowledge/pkg/datastore/filetypes"
	"github.com/knadh/koanf/maps"
	"github.com/spf13/cobra"
)

type ClientLoad struct {
	Loader       string `usage:"Choose a document loader to use"`
	OutputFormat string `name:"format" usage:"Choose an output format" default:"structured"`
}

func (s *ClientLoad) Customize(cmd *cobra.Command) {
	cmd.Use = "load <input> <output>"
	cmd.Short = "Load a file and transform it to markdown"
	cmd.Args = cobra.ExactArgs(2)
}

func (s *ClientLoad) Run(cmd *cobra.Command, args []string) error {
	input := args[0]
	output := args[1]

	if !slices.Contains([]string{"structured", "markdown"}, s.OutputFormat) {
		return fmt.Errorf("unsupported output format %q", s.OutputFormat)
	}

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
		return fmt.Errorf("failed to load documents from file %q using loader %q: %w", input, s.Loader, err)
	}

	var text string

	switch s.OutputFormat {
	case "markdown":
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

		text = strings.Join(texts, "\n---docbreak---\n")

	case "structured":
		var structuredInput structured.StructuredInput
		structuredInput.Metadata = map[string]any{}
		structuredInput.Documents = make([]structured.StructuredInputDocument, 0, len(docs))

		commonMetadata := maps.Copy(docs[0].Metadata)
		for _, doc := range docs {
			commonMetadata = extractCommon(commonMetadata, doc.Metadata)
			structuredInput.Documents = append(structuredInput.Documents, structured.StructuredInputDocument{
				Metadata: doc.Metadata,
				Content:  doc.Content,
			})
		}

		commonMetadata["source"] = input
		structuredInput.Metadata = commonMetadata

		for i, doc := range structuredInput.Documents {
			structuredInput.Documents[i].Metadata = dropCommon(doc.Metadata, commonMetadata)
		}

		jsonBytes := bytes.NewBuffer(nil)
		encoder := json.NewEncoder(jsonBytes)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(structuredInput); err != nil {
			return fmt.Errorf("failed to encode structured input: %w", err)
		}
		text = jsonBytes.String()
	default:
		return fmt.Errorf("unsupported output format %q", s.OutputFormat)
	}

	if output == "-" {
		fmt.Println(text)
		return nil
	}

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

func dropCommon(target, common map[string]any) map[string]any {
	for key, _ := range target {
		if _, exists := common[key]; exists {
			delete(target, key)
		}
	}

	return target
}

func extractCommon(target, other map[string]any) map[string]any {
	for key, value := range target {
		if v, exists := other[key]; exists && v == value {
			target[key] = value
		} else {
			delete(target, key)
		}
	}

	return target
}
