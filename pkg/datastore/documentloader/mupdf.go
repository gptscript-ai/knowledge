//go:build mupdf

package documentloader

import (
	"context"
	"fmt"
	"github.com/gptscript-ai/knowledge/pkg/datastore/documentloader/pdf/defaults"
	"github.com/gptscript-ai/knowledge/pkg/datastore/documentloader/pdf/mupdf"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	"github.com/mitchellh/mapstructure"
	"io"
	"log/slog"
)

func init() {

	defaults.DefaultPDFReaderFunc = func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
		slog.Debug("Default PDF Reader is MuPDF")
		r, err := mupdf.NewPDF(reader)
		if err != nil {
			slog.Error("Failed to create MuPDF loader", "error", err)
			return nil, err
		}
		return r.Load(ctx)
	}

	MuPDFGetter = func(config any) (LoaderFunc, error) {
		var pdfConfig mupdf.PDFOptions
		if config != nil {
			slog.Debug("PDF custom config", "config", config)
			if err := mapstructure.Decode(config, &pdfConfig); err != nil {
				return nil, fmt.Errorf("failed to decode PDF document loader configuration: %w", err)
			}
			slog.Debug("PDF custom config (decoded)", "pdfConfig", pdfConfig)
		}
		return func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
			r, err := mupdf.NewPDF(reader, mupdf.WithConfig(pdfConfig))
			if err != nil {
				slog.Error("Failed to create PDF loader", "error", err)
				return nil, err
			}
			return r.Load(ctx)
		}, nil
	}

	MuPDFConfig = mupdf.PDFOptions{}
}
