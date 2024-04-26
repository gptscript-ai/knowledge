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
	"github.com/gptscript-ai/knowledge/pkg/index"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	golcdocloaders "github.com/hupe1980/golc/documentloader"
	golcschema "github.com/hupe1980/golc/schema"
	lcgodocloaders "github.com/tmc/langchaingo/documentloaders"
	lcgoschema "github.com/tmc/langchaingo/schema"
	"io"
	"log/slog"
	"path"
	"strings"
)

var firstclassFileExtensions = map[string]struct{}{
	".pdf":   {},
	".html":  {},
	".md":    {},
	".txt":   {},
	".csv":   {},
	".ipynb": {},
}

type IngestOpts struct {
	Filename            *string
	FileMetadata        *index.FileMetadata
	IsDuplicateFuncName string
	IsDuplicateFunc     IsDuplicateFunc
}

// Ingest loads a document from a reader and adds it to the dataset.
func (s *Datastore) Ingest(ctx context.Context, datasetID string, reader io.Reader, opts IngestOpts) ([]string, error) {
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

	slog.Info("IngestOpts", "opts", opts)

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
	var filetype string
	if opts.Filename != nil {
		filetype = path.Ext(*opts.Filename)
		if _, ok := firstclassFileExtensions[filetype]; !ok {
			filetype = ""
		}
	}
	if filetype == "" {
		filetype, reader, err = mimetypeFromReader(reader)
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

	/*
	 * Load documents from the content
	 * For now, we're using documentloaders from both langchaingo and golc
	 * and translate them to our document schema.
	 */

	var lcgodocs []lcgoschema.Document
	var golcdocs []golcschema.Document

	switch filetype {
	case ".pdf", "application/pdf":
		// The PDF loader requires a size argument, so we can either read the whole file into memory
		// or write it to a temporary file and pass load directly from that file.
		// We choose the former for now.
		data, err := io.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to read PDF data: %w", err)
		}
		r, err := golcdocloaders.NewPDF(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			slog.Error("Failed to create PDF loader", "error", err)
			return nil, err
		}
		golcdocs, err = r.Load(ctx)
	case ".html", "text/html":
		lcgodocs, err = lcgodocloaders.NewHTML(reader).Load(ctx)
	case ".md", ".txt", "text/plain", "text/markdown":
		lcgodocs, err = lcgodocloaders.NewText(reader).Load(ctx)
	case ".csv", "text/csv":
		golcdocs, err = golcdocloaders.NewCSV(reader).Load(ctx)
		if err != nil && errors.Is(err, csv.ErrBareQuote) {
			oerr := err
			err = nil
			var nerr error
			golcdocs, nerr = golcdocloaders.NewCSV(reader, func(o *golcdocloaders.CSVOptions) {
				o.LazyQuotes = true
			}).Load(ctx)
			if nerr != nil {
				err = errors.Join(oerr, nerr)
			}
		}
	case ".ipynb":
		golcdocs, err = golcdocloaders.NewNotebook(reader).Load(ctx)
	default:
		// TODO(@iwilltry42): Fallback to plaintext reader? Example: Makefile, Dockerfile, Source Files, etc.
		slog.Error("Unsupported file type", "filename", *opts.Filename, "type", filetype)
		return nil, fmt.Errorf("unsupported file type: %s", filetype)
	}

	if err != nil {
		slog.Error("Failed to load document", "error", err)
		return nil, fmt.Errorf("failed to load document: %w", err)
	}

	docs := make([]vs.Document, len(lcgodocs)+len(golcdocs))
	for idx, doc := range lcgodocs {
		doc.Metadata["filename"] = *opts.Filename
		docs[idx] = vs.Document{
			Metadata: doc.Metadata,
			Content:  doc.PageContent,
		}
	}

	for idx, doc := range golcdocs {
		doc.Metadata["filename"] = *opts.Filename
		docs[idx] = vs.Document{
			Metadata: doc.Metadata,
			Content:  doc.PageContent,
		}
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
