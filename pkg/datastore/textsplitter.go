package datastore

import lcgosplitter "github.com/tmc/langchaingo/textsplitter"

type TextSplitterOpts struct {
	ChunkSize    int    `usage:"Textsplitter Chunk Size" default:"1024" env:"KNOW_TEXTSPLITTER_CHUNK_SIZE"`
	ChunkOverlap int    `usage:"Textsplitter Chunk Overlap" default:"256" env:"KNOW_TEXTSPLITTER_CHUNK_OVERLAP"`
	ModelName    string `usage:"Textsplitter Model Name" default:"gpt-4" env:"KNOW_TEXTSPLITTER_MODEL_NAME"`
	EncodingName string `usage:"Textsplitter Encoding Name" default:"cl100k_base" env:"KNOW_TEXTSPLITTER_ENCODING_NAME"`
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
