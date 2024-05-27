package datastore

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/acorn-io/z"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/gptscript-ai/knowledge/pkg/datastore/documentloader"
	"github.com/gptscript-ai/knowledge/pkg/datastore/textsplitter"
	"github.com/gptscript-ai/knowledge/pkg/datastore/transformers"
	"github.com/gptscript-ai/knowledge/pkg/datastore/types"
	"github.com/gptscript-ai/knowledge/pkg/flows"
	"github.com/gptscript-ai/knowledge/pkg/index"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	golcdocloaders "github.com/hupe1980/golc/documentloader"
	"github.com/ledongthuc/pdf"
	"github.com/lu4p/cat"
	lcgodocloaders "github.com/tmc/langchaingo/documentloaders"
	"io"
	"log/slog"
	"path"
	"strings"
)

const (
	defaultTokenModel    = "gpt-4"
	defaultChunkSize     = 1024
	defaultChunkOverlap  = 256
	defaultTokenEncoding = "cl100k_base"
)

var firstclassFileExtensions = map[string]struct{}{
	".pdf":   {},
	".html":  {},
	".md":    {},
	".txt":   {},
	".docx":  {},
	".odt":   {},
	".rtf":   {},
	".csv":   {},
	".ipynb": {},
	".json":  {},
}

type IngestOpts struct {
	Filename            *string
	FileMetadata        *index.FileMetadata
	IsDuplicateFuncName string
	IsDuplicateFunc     IsDuplicateFunc
	TextSplitterOpts    *TextSplitterOpts
}

// Ingest loads a document from a reader and adds it to the dataset.
func (s *Datastore) Ingest(ctx context.Context, datasetID string, content []byte, opts IngestOpts) ([]string, error) {
	isDuplicate := DummyDedupe // default: no deduplication
	if opts.IsDuplicateFuncName != "" {
		df, ok := IsDuplicateFuncs[opts.IsDuplicateFuncName]
		if !ok {
			return nil, fmt.Errorf("unknown deduplication function: %s", opts.IsDuplicateFuncName)
		}
		isDuplicate = df
	} else if opts.IsDuplicateFunc != nil {
		isDuplicate = opts.IsDuplicateFunc
	}

	// Generate ID
	fUUID, err := uuid.NewUUID()
	if err != nil {
		slog.Error("Failed to generate UUID", "error", err)
		return nil, err
	}
	fileID := fUUID.String()

	/*
	 * Detect filetype
	 */
	reader := bytes.NewReader(content)
	var filetype string
	if opts.Filename != nil {
		filetype = path.Ext(*opts.Filename)
		if _, ok := firstclassFileExtensions[filetype]; !ok {
			filetype = ""
		}
	}
	if filetype == "" {
		filetype, _, err = mimetypeFromReader(bytes.NewReader(content))
		if err != nil {
			slog.Error("Failed to detect filetype", "error", err)
			return nil, fmt.Errorf("failed to detect filetype: %w", err)
		}
	}
	if filetype == "" {
		slog.Error("Failed to detect filetype", "filename", *opts.Filename)
		return nil, fmt.Errorf("failed to detect filetype")
	}

	filetype = strings.Split(filetype, ";")[0] // remove charset (mimetype), e.g. from "text/plain; charset=utf-8"

	/*
	 * Set filename if not provided
	 */
	if opts.Filename == nil {
		opts.Filename = z.Pointer("<unnamed_document>")
	}

	slog.Debug("Loading data", "type", filetype, "filename", *opts.Filename)

	/*
	 * Exit early if the document is a duplicate
	 */
	isDupe, err := isDuplicate(ctx, s, datasetID, nil, opts)
	if err != nil {
		slog.Error("Failed to check for duplicates", "error", err)
		return nil, fmt.Errorf("failed to check for duplicates: %w", err)
	}
	if isDupe {
		slog.Info("Ignoring duplicate document", "filename", *opts.Filename, "absolute_path", opts.FileMetadata.AbsolutePath)
		return nil, nil
	}

	ingestionFlow := flows.IngestionFlow{
		Load:            DefaultDocLoaderFunc(filetype),
		Split:           DefaultTextSplitter(filetype, opts.TextSplitterOpts).SplitDocuments,
		Transformations: DefaultDocumentTransformers(filetype),
	}

	// Mandatory Transformation: Add filename to metadata
	em := &transformers.ExtraMetadata{Metadata: map[string]any{"filename": *opts.Filename}}
	ingestionFlow.Transformations = append(ingestionFlow.Transformations, em)

	docs, err := GetDocuments(ctx, reader, ingestionFlow)
	if err != nil {
		slog.Error("Failed to load documents", "error", err)
		return nil, fmt.Errorf("failed to load documents: %w", err)
	}

	if len(docs) == 0 {
		slog.Error("No documents found")
		return nil, fmt.Errorf("no documents found")
	}

	// Add documents to VectorStore -> This generates the embeddings
	slog.Debug("Ingesting documents", "count", len(docs))
	docIDs, err := s.Vectorstore.AddDocuments(ctx, docs, datasetID)
	if err != nil {
		slog.Error("Failed to add documents", "error", err)
		return nil, fmt.Errorf("failed to add documents: %w", err)
	}

	// Record file and documents in database
	dbDocs := make([]index.Document, len(docIDs))
	for idx, docID := range docIDs {
		dbDocs[idx] = index.Document{
			ID:      docID,
			FileID:  fileID,
			Dataset: datasetID,
		}
	}

	dbFile := index.File{
		ID:        fileID,
		Dataset:   datasetID,
		Documents: dbDocs,
		FileMetadata: index.FileMetadata{
			Name: *opts.Filename,
		},
	}

	if opts.FileMetadata != nil {
		dbFile.FileMetadata.AbsolutePath = opts.FileMetadata.AbsolutePath
		dbFile.FileMetadata.Size = opts.FileMetadata.Size
		dbFile.FileMetadata.ModifiedAt = opts.FileMetadata.ModifiedAt
	}

	tx := s.Index.WithContext(ctx).Create(&dbFile)
	if tx.Error != nil {
		slog.Error("Failed to create file", "error", tx.Error)
		return nil, fmt.Errorf("failed to create file: %w", tx.Error)
	}

	slog.Info("Ingested document", "filename", *opts.Filename, "count", len(docIDs), "absolute_path", dbFile.FileMetadata.AbsolutePath)

	return docIDs, nil
}

// mimetypeFromReader returns the MIME type of input and a new reader which still has the whole input
func mimetypeFromReader(reader io.Reader) (string, io.Reader, error) {
	header := bytes.NewBuffer(nil)
	mtype, err := mimetype.DetectReader(io.TeeReader(reader, header))
	if err != nil {
		return "", nil, err
	}

	// Get back complete input reader
	newReader := io.MultiReader(header, reader)

	return mtype.String(), newReader, err
}

func DefaultDocLoaderFunc(filetype string) func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
	switch filetype {
	case ".pdf", "application/pdf":
		return func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
			data, nerr := io.ReadAll(reader)
			if nerr != nil {
				return nil, fmt.Errorf("failed to read PDF data: %w", nerr)
			}
			r, nerr := documentloader.NewPDF(bytes.NewReader(data), int64(len(data)), documentloader.WithInterpreterOpts(pdf.WithIgnoreDefOfNonNameVals([]string{"CMapName"})))
			if nerr != nil {
				slog.Error("Failed to create PDF loader", "error", nerr)
				return nil, nerr
			}
			return r.Load(ctx)
		}
	case ".html", "text/html":
		return func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
			return documentloader.FromLangchain(lcgodocloaders.NewHTML(reader)).Load(ctx)
		}
	case ".md", "text/markdown":
		return func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
			return documentloader.FromLangchain(lcgodocloaders.NewText(reader)).Load(ctx)
		}
	case ".txt", "text/plain":
		return func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
			return documentloader.FromLangchain(lcgodocloaders.NewText(reader)).Load(ctx)
		}
	case ".csv", "text/csv":
		return func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
			docs, err := documentloader.FromGolc(golcdocloaders.NewCSV(reader)).Load(ctx)
			if err != nil && errors.Is(err, csv.ErrBareQuote) {
				oerr := err
				err = nil
				var nerr error
				docs, nerr = documentloader.FromGolc(golcdocloaders.NewCSV(reader, func(o *golcdocloaders.CSVOptions) {
					o.LazyQuotes = true
				})).Load(ctx)
				if nerr != nil {
					err = errors.Join(oerr, nerr)
				}
			}
			return docs, err
		}
	case ".json", "application/json":
		return func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
			return documentloader.FromLangchain(lcgodocloaders.NewText(reader)).Load(ctx)
		}
	case ".ipynb":
		return func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
			return documentloader.FromGolc(golcdocloaders.NewNotebook(reader)).Load(ctx)
		}
	case ".docx", ".odt", ".rtf", "application/vnd.oasis.opendocument.text", "text/rtf", "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
			data, nerr := io.ReadAll(reader)
			if nerr != nil {
				return nil, fmt.Errorf("failed to read %s data: %w", filetype, nerr)
			}
			text, nerr := cat.FromBytes(data)
			if nerr != nil {
				return nil, fmt.Errorf("failed to extract text from %s: %w", filetype, nerr)
			}
			return documentloader.FromLangchain(lcgodocloaders.NewText(strings.NewReader(text))).Load(ctx)
		}
	default:
		slog.Error("Unsupported file type", "type", filetype)
		return nil
	}
}

func DefaultTextSplitter(filetype string, textSplitterOpts *TextSplitterOpts) types.TextSplitter {
	if textSplitterOpts == nil {
		textSplitterOpts = z.Pointer(NewTextSplitterOpts())
	}
	genericTextSplitter := textsplitter.FromLangchain(NewLcgoTextSplitter(*textSplitterOpts))
	markdownTextSplitter := textsplitter.FromLangchain(NewLcgoMarkdownSplitter(*textSplitterOpts))

	switch filetype {
	case ".md", "text/markdown":
		return markdownTextSplitter
	default:
		return genericTextSplitter
	}
}

func DefaultDocumentTransformers(filetype string) []types.DocumentTransformer {
	return []types.DocumentTransformer{}
}

func GetDocuments(ctx context.Context, reader io.Reader, ingestionFlow flows.IngestionFlow) ([]vs.Document, error) {
	var err error
	var docs []vs.Document

	/*
	 * Load documents from the content
	 * For now, we're using documentloaders from both langchaingo and golc
	 * and translate them to our document schema.
	 */

	docs, err = ingestionFlow.Load(ctx, reader)
	if err != nil {
		slog.Error("Failed to load documents", "error", err)
		return nil, fmt.Errorf("failed to load documents: %w", err)
	}

	/*
	 * Split documents - Chunking
	 */
	docs, err = ingestionFlow.Split(docs)
	if err != nil {
		slog.Error("Failed to split documents", "error", err)
		return nil, fmt.Errorf("failed to split documents: %w", err)
	}

	/*
	 * Transform documents
	 */
	docs, err = ingestionFlow.Transform(ctx, docs)
	if err != nil {
		slog.Error("Failed to transform documents", "error", err)
		return nil, fmt.Errorf("failed to transform documents: %w", err)
	}

	return docs, nil
}
