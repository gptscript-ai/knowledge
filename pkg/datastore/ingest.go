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
	"github.com/lu4p/cat"
	lcgodocloaders "github.com/tmc/langchaingo/documentloaders"
	lcgoschema "github.com/tmc/langchaingo/schema"
	lcgosplitter "github.com/tmc/langchaingo/textsplitter"
	"io"
	"log/slog"
	"path"
	"strings"
)

const defaultTokenModel = "gpt-4"

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
}

type IngestOpts struct {
	Filename            *string
	FileMetadata        *index.FileMetadata
	IsDuplicateFuncName string
	IsDuplicateFunc     IsDuplicateFunc
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

	docs, err := GetDocuments(ctx, *opts.Filename, filetype, reader)
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

func GetDocuments(ctx context.Context, filename, filetype string, reader io.Reader) ([]vs.Document, error) {
	/*
	 * Load documents from the content
	 * For now, we're using documentloaders from both langchaingo and golc
	 * and translate them to our document schema.
	 */

	var lcgodocs []lcgoschema.Document
	var golcdocs []golcschema.Document

	var err error

	switch filetype {
	case ".pdf", "application/pdf":
		// The PDF loader requires a size argument, so we can either read the whole file into memory
		// or write it to a temporary file and pass load directly from that file.
		// We choose the former for now.
		data, nerr := io.ReadAll(reader)
		if nerr != nil {
			return nil, fmt.Errorf("failed to read PDF data: %w", nerr)
		}
		r, nerr := golcdocloaders.NewPDF(bytes.NewReader(data), int64(len(data)))
		if nerr != nil {
			slog.Error("Failed to create PDF loader", "error", nerr)
			return nil, nerr
		}
		rdocs, nerr := r.Load(ctx)
		if nerr != nil {
			slog.Error("Failed to load PDF", "filename", filename, "error", nerr)
			return nil, fmt.Errorf("failed to load PDF %q: %w", filename, nerr)
		}

		// TODO: consolidate splitters in this repo, so we don't have to convert back and forth
		splitter := lcgosplitter.NewTokenSplitter(lcgosplitter.WithModelName(defaultTokenModel))
		lcgodocs = make([]lcgoschema.Document, len(rdocs))
		for idx, rdoc := range rdocs {
			lcgodocs[idx] = lcgoschema.Document{
				PageContent: rdoc.PageContent,
				Metadata:    rdoc.Metadata,
			}
		}
		lcgodocs, err = lcgosplitter.SplitDocuments(splitter, lcgodocs)
	case ".html", "text/html":
		splitter := lcgosplitter.NewTokenSplitter(lcgosplitter.WithModelName(defaultTokenModel))
		lcgodocs, err = lcgodocloaders.NewHTML(reader).LoadAndSplit(ctx, splitter)
	case ".md", "text/markdown":
		splitter := lcgosplitter.NewMarkdownTextSplitter(lcgosplitter.WithModelName(defaultTokenModel))
		lcgodocs, err = lcgodocloaders.NewText(reader).LoadAndSplit(ctx, splitter)
	case ".txt", "text/plain":
		splitter := lcgosplitter.NewTokenSplitter(lcgosplitter.WithModelName(defaultTokenModel))
		lcgodocs, err = lcgodocloaders.NewText(reader).LoadAndSplit(ctx, splitter)
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
	case ".json", "application/json":
		splitter := lcgosplitter.NewTokenSplitter(lcgosplitter.WithModelName(defaultTokenModel))
		lcgodocs, err = lcgodocloaders.NewText(reader).LoadAndSplit(ctx, splitter)
	case ".ipynb":
		golcdocs, err = golcdocloaders.NewNotebook(reader).Load(ctx)
	case ".docx", ".odt", ".rtf", "application/vnd.oasis.opendocument.text", "text/rtf", "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		data, nerr := io.ReadAll(reader)
		if nerr != nil {
			return nil, fmt.Errorf("failed to read %s data: %w", filetype, nerr)
		}
		text, nerr := cat.FromBytes(data)
		if nerr != nil {
			return nil, fmt.Errorf("failed to extract text from %s: %w", filetype, nerr)
		}
		splitter := lcgosplitter.NewTokenSplitter(lcgosplitter.WithModelName(defaultTokenModel))
		lcgodocs, err = lcgodocloaders.NewText(strings.NewReader(text)).LoadAndSplit(ctx, splitter)
	default:
		// TODO(@iwilltry42): Fallback to plaintext reader? Example: Makefile, Dockerfile, Source Files, etc.
		slog.Error("Unsupported file type", "filename", filename, "type", filetype)
		return nil, fmt.Errorf("file %q has unsupported file type %q", filename, filetype)
	}

	if err != nil {
		slog.Error("Failed to load document", "error", err)
		return nil, fmt.Errorf("failed to load document: %w", err)
	}

	docs := make([]vs.Document, len(lcgodocs)+len(golcdocs))
	for idx, doc := range lcgodocs {
		doc.Metadata["filename"] = filename
		docs[idx] = vs.Document{
			Metadata: doc.Metadata,
			Content:  doc.PageContent,
		}
	}

	for idx, doc := range golcdocs {
		doc.Metadata["filename"] = filename
		docs[idx] = vs.Document{
			Metadata: doc.Metadata,
			Content:  doc.PageContent,
		}
	}

	return docs, nil
}
