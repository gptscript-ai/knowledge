package defaults

const (
	EmbeddingDimension int = 1536
	TopK               int = 5

	TextSplitterTokenModel    = "gpt-4"
	TextSplitterChunkSize     = 1024
	TextSplitterChunkOverlap  = 256
	TextSplitterTokenEncoding = "cl100k_base"
)
