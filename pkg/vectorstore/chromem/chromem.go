package chromem

import (
	"context"
	"github.com/google/uuid"
	"github.com/philippgille/chromem-go"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"log/slog"
	"maps"
	"runtime"
	"strconv"
)

type Store struct {
	db         *chromem.DB
	collection *chromem.Collection

	embeddingFunc chromem.EmbeddingFunc
}

// New creates a new Chromem vector store.
func New(opts ...Option) (Store, error) {
	// Apply options
	store, err := applyClientOptions(opts...)
	if err != nil {
		return Store{}, err
	}

	return store, nil
}

func (s Store) AddDocuments(ctx context.Context, docs []schema.Document, options ...vectorstores.Option) ([]string, error) {

	ids := make([]string, len(docs))
	chromemDocs := make([]chromem.Document, len(docs))
	for docIdx, doc := range docs {
		ids[docIdx] = uuid.New().String()
		mc := make(map[string]any)
		maps.Copy(mc, doc.Metadata)
		if len(doc.PageContent) == 0 {
			slog.Debug("Document has no content", "id", ids[docIdx], "index", docIdx)
			doc.PageContent = "<no content>"
		}
		chromemDocs[docIdx] = chromem.Document{
			ID:        ids[docIdx],
			Metadata:  anyMapToStringMap(mc),
			Embedding: nil, // Embeddings will be computed downstream
			Content:   doc.PageContent,
		}
	}

	err := s.collection.AddDocuments(ctx, chromemDocs, runtime.NumCPU()/2)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func anyMapToStringMap(m map[string]any) map[string]string {
	convertedMap := make(map[string]string)

	for key, value := range m {
		switch v := value.(type) {
		case string:
			convertedMap[key] = v
		case int:
			convertedMap[key] = strconv.Itoa(v)
		case bool:
			convertedMap[key] = strconv.FormatBool(v)
		case float64:
			convertedMap[key] = strconv.FormatFloat(v, 'f', -1, 64)
		default:
			x, ok := value.(string)
			if ok {
				convertedMap[key] = x
			}
			// skip unsupported types for now
		}
	}
	return convertedMap
}

func convertStringMapToAnyMap(m map[string]string) map[string]any {
	convertedMap := make(map[string]any)

	for key, value := range m {
		convertedMap[key] = value
	}
	return convertedMap
}

func (s Store) SimilaritySearch(ctx context.Context, query string, numDocuments int, options ...vectorstores.Option) ([]schema.Document, error) {

	qr, err := s.collection.Query(ctx, query, numDocuments, nil, nil)
	if err != nil {
		return nil, err
	}

	if len(qr) == 0 {
		return nil, nil
	}

	var sDocs []schema.Document

	for _, qrd := range qr {
		sDocs = append(sDocs, schema.Document{
			Metadata:    convertStringMapToAnyMap(qrd.Metadata),
			Score:       qrd.Similarity,
			PageContent: qrd.Content,
		})
	}

	return sDocs, nil
}

func (s Store) RemoveCollection(ctx context.Context) error {
	return s.db.DeleteCollection(s.collection.Name)
}
