// Package tiktoken provides functionality for tokenizing and encoding text using the tiktoken algorithm.
// The package includes various functions for text processing and encoding using the tiktoken algorithm.
package tiktoken

import (
	"fmt"
	"strings"
)

// Constants for different encodings.
const (
	CL100kBase string = "cl100k_base"
	P50kBase   string = "p50k_base"
	P50kEdit   string = "p50k_edit"
	R50kBase   string = "r50k_base"
	GPT2       string = "gpt2"
)

// ModelPrefixToEncoding maps model prefixes to encodings.
var ModelPrefixToEncoding = map[string]string{
	// chat
	"gpt-4-":         CL100kBase, // e.g., gpt-4-0314, etc., plus gpt-4-32k
	"gpt-3.5-turbo-": CL100kBase, // e.g, gpt-3.5-turbo-0301, -0401, etc.
	"gpt-35-turbo":   CL100kBase, // Azure deployment name
}

// ModelToEncoding maps models to encodings.
var ModelToEncoding = map[string]string{
	// chat
	"gpt-4":         CL100kBase,
	"gpt-3.5-turbo": CL100kBase,
	"gpt-35-turbo":  CL100kBase, // Azure deployment name
	// text
	"text-davinci-003": P50kBase,
	"text-davinci-002": P50kBase,
	"text-davinci-001": R50kBase,
	"text-curie-001":   R50kBase,
	"text-babbage-001": R50kBase,
	"text-ada-001":     R50kBase,
	"davinci":          R50kBase,
	"curie":            R50kBase,
	"babbage":          R50kBase,
	"ada":              R50kBase,
	// code
	"code-davinci-002": P50kBase,
	"code-davinci-001": P50kBase,
	"code-cushman-002": P50kBase,
	"code-cushman-001": P50kBase,
	"davinci-codex":    P50kBase,
	"cushman-codex":    P50kBase,
	// edit
	"text-davinci-edit-001": P50kEdit,
	"code-davinci-edit-001": P50kEdit,
	// embeddings
	"text-embedding-ada-002": CL100kBase,
	"text-embedding-3-small": CL100kBase,
	"text-embedding-3-large": CL100kBase,
	// old embeddings
	"text-similarity-davinci-001":  R50kBase,
	"text-similarity-curie-001":    R50kBase,
	"text-similarity-babbage-001":  R50kBase,
	"text-similarity-ada-001":      R50kBase,
	"text-search-davinci-doc-001":  R50kBase,
	"text-search-curie-doc-001":    R50kBase,
	"text-search-babbage-doc-001":  R50kBase,
	"text-search-ada-doc-001":      R50kBase,
	"code-search-babbage-code-001": R50kBase,
	"code-search-ada-code-001":     R50kBase,
	// open source
	"gpt2": GPT2,
}

// NewEncodingForModel returns a new Encoding based on the given model.
// It checks the ModelToEncoding map and ModelPrefixToEncoding map to find a matching encoding.
func NewEncodingForModel(model string) (*Encoding, error) {
	if encoding, ok := ModelToEncoding[model]; ok {
		return NewEncodingByName(encoding)
	} else {
		for prefix, encoding := range ModelPrefixToEncoding {
			if strings.HasPrefix(model, prefix) {
				return NewEncodingByName(encoding)
			}
		}
	}

	return nil, fmt.Errorf("no encoding for model %s", model)
}
