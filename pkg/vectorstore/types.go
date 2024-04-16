package vectorstore

type Document struct {
	Content         string         `json:"content"`
	Metadata        map[string]any `json:"metadata"`
	SimilarityScore float32        `json:"similarity_score"`
}
