package documentloader

import (
	"bytes"
	"context"
	"strings"

	"github.com/gptscript-ai/knowledge/pkg/datastore/types"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"github.com/ledongthuc/pdf"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

/*
 * Credits to https://github.com/hupe1980/golc/blob/main/documentloader/pdf.go
 */

// Compile time check to ensure PDF satisfies the DocumentLoader interface.
var _ types.DocumentLoader = (*PDF)(nil)

type PDFOptions struct {
	// Password for encrypted PDF files.
	Password string

	// Page number to start loading from (default is 1).
	StartPage uint

	// Maximum number of pages to load (0 for all pages).
	MaxPages uint

	// Source is the name of the pdf document
	Source string

	// InterpreterConfig is the configuration for the PDF interpreter.
	InterpreterConfig *pdf.InterpreterConfig
}

// WithConfig sets the PDF loader configuration.
func WithConfig(config PDFOptions) func(o *PDFOptions) {
	return func(o *PDFOptions) {
		*o = config
	}
}

// WithInterpreterConfig sets the interpreter config for the PDF loader.
func WithInterpreterConfig(cfg pdf.InterpreterConfig) func(o *PDFOptions) {
	return func(o *PDFOptions) {
		o.InterpreterConfig = &cfg
	}
}

// WithInterpreterOpts sets the interpreter options for the PDF loader.
func WithInterpreterOpts(opts ...pdf.InterpreterOption) func(o *PDFOptions) {
	return func(o *PDFOptions) {
		if o.InterpreterConfig == nil {
			o.InterpreterConfig = &pdf.InterpreterConfig{}
		}
		for _, opt := range opts {
			opt(o.InterpreterConfig)
		}
	}
}

// PDF represents a PDF document loader that implements the DocumentLoader interface.
type PDF struct {
	data []byte
	opts PDFOptions
}

// NewPDFFromFile creates a new PDF loader with the given options.
func NewPDF(data []byte, optFns ...func(o *PDFOptions)) (*PDF, error) {
	opts := PDFOptions{
		StartPage: 1,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.StartPage == 0 {
		opts.StartPage = 1
	}

	return &PDF{
		data: data,
		opts: opts,
	}, nil
}

// Load loads the PDF document and returns a slice of vs.Document containing the page contents and metadata.
func (l *PDF) Load(ctx context.Context) ([]vs.Document, error) {
	var (
		err error
	)

	rs := bytes.NewReader(l.data)
	pdfReader, err := model.NewPdfReader(rs)
	if err != nil {
		return nil, err
	}

	encrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return nil, err
	}

	if encrypted {
		_, err := pdfReader.Decrypt([]byte(l.opts.Password))
		if err != nil {
			return nil, err
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return nil, err
	}

	docs := make([]vs.Document, 0, numPages)

	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			return nil, err
		}

		ex, err := extractor.New(page)
		if err != nil {
			return nil, err
		}

		pageText, err := ex.ExtractText()
		if err != nil {
			return nil, err
		}

		doc := vs.Document{
			Content: strings.TrimSpace(pageText),
			Metadata: map[string]any{
				"page":       i,
				"totalPages": numPages,
			},
		}

		docs = append(docs, doc)
	}

	return docs, nil
}

// LoadAndSplit loads PDF documents from the provided reader and splits them using the specified text splitter.
func (l *PDF) LoadAndSplit(ctx context.Context, splitter types.TextSplitter) ([]vs.Document, error) {
	docs, err := l.Load(ctx)
	if err != nil {
		return nil, err
	}

	return splitter.SplitDocuments(docs)
}
