package documentloader

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/types"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"io"
	"slices"
	"strings"
)

// Compile time check to ensure CSV satisfies the DocumentLoader interface.
var _ types.DocumentLoader = (*CSV)(nil)

type CSVDocumentFormat string

const (
	// OriginalFormat represents the original format of the CSV document.
	OriginalFormat CSVDocumentFormat = "original"

	// JSONFormat represents the JSON format of the CSV document.
	JSONFormat CSVDocumentFormat = "json"

	// MarkdownFormat represents the Markdown format of the CSV document.
	MarkdownFormat CSVDocumentFormat = "markdown"
)

// CSVOptions contains options for configuring the CSV loader.
type CSVOptions struct {
	// Separator is the rune used to separate fields in the CSV file.
	Separator rune

	// LazyQuotes controls whether the CSV reader should use lazy quotes mode.
	LazyQuotes bool

	// Columns is a list of column names to filter and include in the loaded documents.
	Columns []string

	// ConcatRows controls whether to concatenate rows into a single document.
	ConcatRows bool

	// RowSeparator is the string used to separate rows in the concatenated document. Default is "\n".
	RowSeparator string

	// MaxConcatRows is the maximum number of rows to concatenate into a single document.
	MaxConcatRows int

	// Format is the format in which the documents will be stored. Default is "original".
	Format CSVDocumentFormat
}

// CSV represents a CSV document loader.
type CSV struct {
	r    io.Reader
	opts CSVOptions
}

// NewCSV creates a new CSV loader with an io.Reader and optional configuration options.
// It returns a pointer to the created CSV loader.
func NewCSV(r io.Reader, optFns ...func(o *CSVOptions)) *CSV {
	opts := CSVOptions{
		Separator:     ',',
		LazyQuotes:    false,
		ConcatRows:    false,
		MaxConcatRows: 100,
		RowSeparator:  "\n",
		Format:        OriginalFormat,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	return &CSV{
		r:    r,
		opts: opts,
	}
}

// Load loads CSV documents from the provided reader.
func (l *CSV) Load(ctx context.Context) ([]vs.Document, error) {
	var (
		header []string
		docs   []vs.Document
		rown   uint

		err error

		docContent []string // content of a single document
	)

	reader := csv.NewReader(l.r)
	reader.Comma = l.opts.Separator
	reader.LazyQuotes = l.opts.LazyQuotes

	// Read header
	header, err = reader.Read()
	if err != nil {
		if err == io.EOF {
			return docs, nil
		}
		return nil, err
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		var content []string

		// Transposed Markdown format
		if l.opts.Format == MarkdownFormat {
			for i, value := range row {
				if len(l.opts.Columns) > 0 && !slices.Contains(l.opts.Columns, header[i]) {
					continue
				}

				line := fmt.Sprintf("%s: %s", header[i], value)
				content = append(content, line)
			}
		}

		// Original format
		if l.opts.Format == OriginalFormat {

		}

		rown++

		// We're not concatenating, so just append to the result slice and continue
		if !l.opts.ConcatRows {
			doc := vs.Document{
				Content:  strings.Join(content, "\n"),
				Metadata: map[string]any{"row": rown},
			}
			docs = append(docs, doc)
			continue
		}

		// Concatenating rows

	}

	return docs, nil
}

// LoadAndSplit loads CSV documents from the provided reader and splits them using the specified text splitter.
func (l *CSV) LoadAndSplit(ctx context.Context, splitter types.TextSplitter) ([]vs.Document, error) {
	docs, err := l.Load(ctx)
	if err != nil {
		return nil, err
	}

	return splitter.SplitDocuments(docs)
}
