package textsplitter

import (
	"dario.cat/mergo"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/defaults"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"github.com/mitchellh/mapstructure"
	lcgosplitter "github.com/tmc/langchaingo/textsplitter"
	"log/slog"
)

type SplitterFunc func([]vs.Document) ([]vs.Document, error)

type TextSplitterOpts struct {
	ChunkSize    int    `json:"chunkSize" mapstructure:"chunkSize" usage:"Textsplitter Chunk Size" default:"1024" env:"KNOW_TEXTSPLITTER_CHUNK_SIZE" name:"textsplitter-chunk-size"`
	ChunkOverlap int    `json:"chunkOverlap" mapstructure:"chunkOverlap" usage:"Textsplitter Chunk Overlap" default:"256" env:"KNOW_TEXTSPLITTER_CHUNK_OVERLAP" name:"textsplitter-chunk-overlap"`
	ModelName    string `json:"modelName" mapstructure:"modelName" usage:"Textsplitter Model Name" default:"gpt-4" env:"KNOW_TEXTSPLITTER_MODEL_NAME" name:"textsplitter-model-name"`
	EncodingName string `json:"encodingName" mapstructure:"encodingName" usage:"Textsplitter Encoding Name" default:"cl100k_base" env:"KNOW_TEXTSPLITTER_ENCODING_NAME" name:"textsplitter-encoding-name"`
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
		cfg := NewTextSplitterOpts()
		if config != nil {
			var customCfg TextSplitterOpts
			if err := mapstructure.Decode(config, &customCfg); err != nil {
				return nil, fmt.Errorf("failed to decode text splitter configuration: %w", err)
			}
			slog.Debug("GetTextSplitterFunc Text (before merge)", "config", customCfg)
			if err := mergo.Merge(&customCfg, cfg); err != nil {
				return nil, fmt.Errorf("failed to merge text splitter configuration: %w", err)
			}
			cfg = customCfg
		}
		slog.Debug("TextSplitterFunc", "config", cfg)
		return FromLangchain(NewLcgoTextSplitter(cfg)).SplitDocuments, nil
	case "markdown":
		cfg := NewTextSplitterOpts()
		if config != nil {
			var customCfg TextSplitterOpts
			if err := mapstructure.Decode(config, &customCfg); err != nil {
				return nil, fmt.Errorf("failed to decode text splitter configuration: %w", err)
			}
			slog.Debug("GetTextSplitterFunc Markdown (before merge)", "config", customCfg)
			if err := mergo.Merge(&customCfg, cfg); err != nil {
				return nil, fmt.Errorf("failed to merge text splitter configuration: %w", err)
			}
			cfg = customCfg
		}
		slog.Debug("MarkdownSplitterFunc", "config", cfg)
		return FromLangchain(NewLcgoMarkdownSplitter(cfg)).SplitDocuments, nil
	default:
		return nil, fmt.Errorf("unknown text splitter %q", name)
	}
}
