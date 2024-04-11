package types

// Dataset represents a new knowledge vector space
type Dataset struct {
	Name     string `json:"name" binding:"required"`
	EmbedDim *int   `json:"embed_dim" example:"1536" swaggertype:"integer"`
}

// Query represents an incoming user query
type Query struct {
	Prompt string `json:"prompt" binding:"required"`
	Topk   *int   `json:"topk" example:"5" swaggertype:"integer"`
}

// Ingest represents incoming content that should be ingested
type Ingest struct {
	Filename *string `json:"filename"`
	FileID   *string `json:"file_id"`
	Content  string  `json:"content" binding:"required"`
}
