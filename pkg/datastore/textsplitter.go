package datastore

import (
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	lcgosplitter "github.com/tmc/langchaingo/textsplitter"
	"strings"
)

type TextSplitterOpts struct {
	ChunkSize    int    `usage:"Textsplitter Chunk Size" default:"1024" env:"KNOW_TEXTSPLITTER_CHUNK_SIZE" name:"textsplitter-chunk-size"`
	ChunkOverlap int    `usage:"Textsplitter Chunk Overlap" default:"256" env:"KNOW_TEXTSPLITTER_CHUNK_OVERLAP" name:"textsplitter-chunk-overlap"`
	ModelName    string `usage:"Textsplitter Model Name" default:"gpt-4" env:"KNOW_TEXTSPLITTER_MODEL_NAME" name:"textsplitter-model-name"`
	EncodingName string `usage:"Textsplitter Encoding Name" default:"cl100k_base" env:"KNOW_TEXTSPLITTER_ENCODING_NAME" name:"textsplitter-encoding-name"`
}

// NewTextSplitterOpts returns the default options for a text splitter.
func NewTextSplitterOpts() TextSplitterOpts {
	return TextSplitterOpts{
		ChunkSize:    defaultChunkSize,
		ChunkOverlap: defaultChunkOverlap,
		ModelName:    defaultTokenModel,
		EncodingName: defaultTokenEncoding,
	}
}

// NewLcgoTextSplitter returns a new langchain-go text splitter.
func NewLcgoTextSplitter(opts TextSplitterOpts) lcgosplitter.TokenSplitter {
	return lcgosplitter.NewTokenSplitter(
		lcgosplitter.WithChunkSize(opts.ChunkSize),
		lcgosplitter.WithChunkOverlap(opts.ChunkOverlap),
		lcgosplitter.WithModelName(opts.ModelName),
		lcgosplitter.WithEncodingName(opts.EncodingName),
	)
}

func NewLcgoMarkdownSplitter(opts TextSplitterOpts) *lcgosplitter.MarkdownTextSplitter {
	return lcgosplitter.NewMarkdownTextSplitter(
		lcgosplitter.WithChunkSize(opts.ChunkSize),
		lcgosplitter.WithChunkOverlap(opts.ChunkOverlap),
		lcgosplitter.WithModelName(opts.ModelName),
		lcgosplitter.WithEncodingName(opts.EncodingName),
		lcgosplitter.WithHeadingHierarchy(true),
	)
}

// FilterMarkdownDocsNoContent filters out Markdown documents with no content or only headings
//
// TODO: this may be moved into the MarkdownTextSplitter as well
func FilterMarkdownDocsNoContent(docs []vs.Document) []vs.Document {
	var filteredDocs []vs.Document
	for _, doc := range docs {
		if doc.Content != "" {
			for _, line := range strings.Split(doc.Content, "\n") {
				if !strings.HasPrefix(line, "#") {
					filteredDocs = append(filteredDocs, doc)
				}
			}
		}
	}
	return filteredDocs
}
