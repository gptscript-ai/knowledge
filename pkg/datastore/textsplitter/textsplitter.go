package textsplitter

import (
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/defaults"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	lcgosplitter "github.com/tmc/langchaingo/textsplitter"
)

type SplitterFunc func([]vs.Document) ([]vs.Document, error)

type TextSplitterOpts struct {
	ChunkSize    int    `usage:"Textsplitter Chunk Size" default:"1024" env:"KNOW_TEXTSPLITTER_CHUNK_SIZE" name:"textsplitter-chunk-size"`
	ChunkOverlap int    `usage:"Textsplitter Chunk Overlap" default:"256" env:"KNOW_TEXTSPLITTER_CHUNK_OVERLAP" name:"textsplitter-chunk-overlap"`
	ModelName    string `usage:"Textsplitter Model Name" default:"gpt-4" env:"KNOW_TEXTSPLITTER_MODEL_NAME" name:"textsplitter-model-name"`
	EncodingName string `usage:"Textsplitter Encoding Name" default:"cl100k_base" env:"KNOW_TEXTSPLITTER_ENCODING_NAME" name:"textsplitter-encoding-name"`
}

// NewTextSplitterOpts returns the default options for a text splitter.
func NewTextSplitterOpts() TextSplitterOpts {
	return TextSplitterOpts{
		ChunkSize:    defaults.TextSplitterChunkSize,
		ChunkOverlap: defaults.TextSplitterChunkOverlap,
		ModelName:    defaults.TextSplitterTokenModel,
		EncodingName: defaults.TextSplitterTokenEncoding,
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

func GetTextSplitterConfig(name string) (any, error) {
	// TODO: expose splitter-specific config, not only our top-level options
	switch name {
	case "text", "markdown":
		return TextSplitterOpts{}, nil
	default:
		return nil, fmt.Errorf("unknown text splitter %q", name)
	}
}

func GetTextSplitterFunc(name string, config any) (SplitterFunc, error) {
	switch name {
	case "text":
		if config == nil {
			config = NewTextSplitterOpts()
		}
		config, ok := config.(TextSplitterOpts)
		if !ok {
			return nil, fmt.Errorf("invalid text splitter configuration")
		}
		return FromLangchain(NewLcgoTextSplitter(config)).SplitDocuments, nil
	case "markdown":
		if config == nil {
			config = NewTextSplitterOpts()
		}
		config, ok := config.(TextSplitterOpts)
		if !ok {
			return nil, fmt.Errorf("invalid markdown text splitter configuration")
		}
		return FromLangchain(NewLcgoMarkdownSplitter(config)).SplitDocuments, nil
	default:
		return nil, fmt.Errorf("unknown text splitter %q", name)
	}
}
