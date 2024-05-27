package documentloader

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	golcdocloaders "github.com/hupe1980/golc/documentloader"
	"github.com/ledongthuc/pdf"
	"github.com/lu4p/cat"
	lcgodocloaders "github.com/tmc/langchaingo/documentloaders"
	"io"
	"log/slog"
	"strings"
)

func DefaultDocLoaderFunc(filetype string) func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
	switch filetype {
	case ".pdf", "application/pdf":
		return func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
			data, nerr := io.ReadAll(reader)
			if nerr != nil {
				return nil, fmt.Errorf("failed to read PDF data: %w", nerr)
			}
			r, nerr := NewPDF(bytes.NewReader(data), int64(len(data)), WithInterpreterOpts(pdf.WithIgnoreDefOfNonNameVals([]string{"CMapName"})))
			if nerr != nil {
				slog.Error("Failed to create PDF loader", "error", nerr)
				return nil, nerr
			}
			return r.Load(ctx)
		}
	case ".html", "text/html":
		return func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
			return FromLangchain(lcgodocloaders.NewHTML(reader)).Load(ctx)
		}
	case ".md", "text/markdown":
		return func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
			return FromLangchain(lcgodocloaders.NewText(reader)).Load(ctx)
		}
	case ".txt", "text/plain":
		return func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
			return FromLangchain(lcgodocloaders.NewText(reader)).Load(ctx)
		}
	case ".csv", "text/csv":
		return func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
			docs, err := FromGolc(golcdocloaders.NewCSV(reader)).Load(ctx)
			if err != nil && errors.Is(err, csv.ErrBareQuote) {
				oerr := err
				err = nil
				var nerr error
				docs, nerr = FromGolc(golcdocloaders.NewCSV(reader, func(o *golcdocloaders.CSVOptions) {
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
			return FromLangchain(lcgodocloaders.NewText(reader)).Load(ctx)
		}
	case ".ipynb":
		return func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
			return FromGolc(golcdocloaders.NewNotebook(reader)).Load(ctx)
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
			return FromLangchain(lcgodocloaders.NewText(strings.NewReader(text))).Load(ctx)
		}
	default:
		slog.Error("Unsupported file type", "type", filetype)
		return nil
	}
}
