package documentloader

import (
	"context"
	"io"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/gen2brain/go-fitz"
	"github.com/gptscript-ai/knowledge/pkg/datastore/types"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
)

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
}

// WithConfig sets the PDF loader configuration.
func WithConfig(config PDFOptions) func(o *PDFOptions) {
	return func(o *PDFOptions) {
		*o = config
	}
}

// PDF represents a PDF document loader that implements the DocumentLoader interface.
type PDF struct {
	opts      PDFOptions
	document  *fitz.Document
	converter *md.Converter
}

// NewPDFFromFile creates a new PDF loader with the given options.
func NewPDF(r io.Reader, optFns ...func(o *PDFOptions)) (*PDF, error) {
	doc, err := fitz.NewFromReader(r)
	if err != nil {
		return nil, err
	}
	opts := PDFOptions{
		StartPage: 1,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if opts.StartPage == 0 {
		opts.StartPage = 1
	}

	converter := md.NewConverter("", true, nil)

	return &PDF{
		opts:      opts,
		document:  doc,
		converter: converter,
	}, nil
}

// Load loads the PDF document and returns a slice of vs.Document containing the page contents and metadata.
func (l *PDF) Load(ctx context.Context) ([]vs.Document, error) {
	docs := make([]vs.Document, 0, l.document.NumPage())

	for pageNum := 0; pageNum < l.document.NumPage(); pageNum++ {
		html, err := l.document.HTML(pageNum, true)
		if err != nil {
			return nil, err
		}

		htmlDoc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			return nil, err
		}
		htmlDoc.Find("img").Remove()

		ret, err := htmlDoc.First().Html()
		if err != nil {
			return nil, err
		}

		markdown, err := l.converter.ConvertString(ret)
		if err != nil {
			return nil, err
		}

		doc := vs.Document{
			Content: strings.TrimSpace(markdown),
			Metadata: map[string]any{
				"page":       pageNum + 1,
				"totalPages": l.document.NumPage(),
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
