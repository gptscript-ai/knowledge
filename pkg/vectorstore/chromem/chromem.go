package chromem

import (
	"context"
	"github.com/google/uuid"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"github.com/philippgille/chromem-go"
	"log/slog"
	"maps"
	"runtime"
	"strconv"
)

type Store struct {
	db            *chromem.DB
	embeddingFunc chromem.EmbeddingFunc
}

// New creates a new Chromem vector store.
func New(db *chromem.DB, embeddingFunc chromem.EmbeddingFunc) Store {
	return Store{
		db:            db,
		embeddingFunc: embeddingFunc,
	}
}

func (s Store) CreateCollection(_ context.Context, name string) error {
	_, err := s.db.CreateCollection(name, nil, s.embeddingFunc)
	if err != nil {
		return err
	}

	return nil
}

func (s Store) AddDocuments(ctx context.Context, docs []vs.Document, collection string) ([]string, error) {
	ids := make([]string, len(docs))
	chromemDocs := make([]chromem.Document, len(docs))
	for docIdx, doc := range docs {
		ids[docIdx] = uuid.NewString()
		mc := make(map[string]any)
		maps.Copy(mc, doc.Metadata)
		if len(doc.Content) == 0 {
			slog.Debug("Document has no content", "id", ids[docIdx], "index", docIdx)
			doc.Content = "<no content>"
		}
		chromemDocs[docIdx] = chromem.Document{
			ID:        ids[docIdx],
			Metadata:  anyMapToStringMap(mc),
			Embedding: nil, // Embeddings will be computed downstream
			Content:   doc.Content,
		}
	}

	col := s.db.GetCollection(collection, s.embeddingFunc)
	if col == nil {
		return nil, vs.ErrCollectionNotFound{Collection: collection}
	}

	err := col.AddDocuments(ctx, chromemDocs, runtime.NumCPU()/2)
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

func (s Store) SimilaritySearch(ctx context.Context, query string, numDocuments int, collection string) ([]vs.Document, error) {
	col := s.db.GetCollection(collection, s.embeddingFunc)
	if col == nil {
		return nil, vs.ErrCollectionNotFound{Collection: collection}
	}

	qr, err := col.Query(ctx, query, numDocuments, nil, nil)
	if err != nil {
		return nil, err
	}

	if len(qr) == 0 {
		return nil, nil
	}

	var sDocs []vs.Document

	for _, qrd := range qr {
		sDocs = append(sDocs, vs.Document{
			Metadata:        convertStringMapToAnyMap(qrd.Metadata),
			SimilarityScore: qrd.Similarity,
			Content:         qrd.Content,
		})
	}

	return sDocs, nil
}

func (s Store) RemoveCollection(_ context.Context, collection string) error {
	return s.db.DeleteCollection(collection)
}

func (s Store) RemoveDocument(ctx context.Context, documentID string, collection string) error {
	col := s.db.GetCollection(collection, s.embeddingFunc)
	if col == nil {
		return vs.ErrCollectionNotFound{Collection: collection}
	}
	return col.Delete(ctx, nil, nil, documentID)
}
