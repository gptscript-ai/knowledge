package documentloader

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/bzip2"
	"context"
	"encoding/csv"
	"errors"
	"io"
	"log/slog"
	"strings"

	"code.sajari.com/docconv/v2"
	"github.com/gptscript-ai/knowledge/pkg/datastore/filetypes"
	vs "github.com/gptscript-ai/knowledge/pkg/vectorstore"
	golcdocloaders "github.com/hupe1980/golc/documentloader"
	"github.com/lu4p/cat/rtftxt"
	lcgodocloaders "github.com/tmc/langchaingo/documentloaders"
)

func DefaultDocLoaderFunc(filetype string) func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
	switch filetype {
	case ".pdf", "application/pdf":
		return func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
			r, nerr := NewPDF(reader)
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
	case ".docx", ".odt", ".rtf", "text/rtf", "application/vnd.oasis.opendocument.text", "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
			var text string
			var metadata map[string]string
			var err error
			switch filetype {
			case ".docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
				text, metadata, err = docconv.ConvertDocx(reader)
			case ".rtf", ".rtfd", "text/rtf":
				buf, err := rtftxt.Text(reader)
				if err != nil {
					return nil, err
				}
				text = buf.String()
			case ".odt", "application/vnd.oasis.opendocument.text":
				text, metadata, err = docconv.ConvertODT(reader)
			}

			if err != nil {
				return nil, err
			}

			docs, err := FromLangchain(lcgodocloaders.NewText(strings.NewReader(text))).Load(ctx)
			if err != nil {
				return nil, err
			}

			for _, doc := range docs {
				m := map[string]any{}
				for k, v := range metadata {
					m[k] = v
				}
				doc.Metadata = m
			}

			return docs, nil
		}
	// todo: OCR support is commented out for now as it relies on external dependencies.
	// We might add it back later.
	//case "image/png", "image/jpeg":
	//	return func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
	//		client := gosseract.NewClient()
	//		defer client.Close()
	//		data, nerr := io.ReadAll(reader)
	//		if nerr != nil {
	//			return nil, fmt.Errorf("failed to read %s data: %w", filetype, nerr)
	//		}
	//		if err := client.SetImageFromBytes(data); err != nil {
	//			return nil, fmt.Errorf("failed to feed data into OCR: %w", nerr)
	//		}
	//		text, err := client.Text()
	//		if err != nil {
	//			return nil, fmt.Errorf("failed to convert image data into OCR")
	//		}
	//		return []vs.Document{
	//			{
	//				Content: text,
	//			},
	//		}, nil
	//	}
	case "application/zip", ".zip":
		var result []vs.Document
		return func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
			data, err := io.ReadAll(reader)
			if err != nil {
				return nil, err
			}
			zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
			if err != nil {
				return nil, err
			}
			for _, f := range zipReader.File {
				if f.FileInfo().IsDir() {
					continue
				}
				rc, err := f.Open()
				if err != nil {
					return nil, err
				}
				content, err := io.ReadAll(rc)
				if err != nil {
					return nil, err
				}
				ft, err := filetypes.GetFiletype(f.Name, content)
				if err != nil {
					return nil, err
				}
				docs, err := DefaultDocLoaderFunc(ft)(ctx, bytes.NewReader(content))
				if err != nil {
					return nil, err
				}
				result = append(result, docs...)
			}
			return result, nil
		}

	case "application/x-bzip2", ".bz2":
		return func(ctx context.Context, reader io.Reader) ([]vs.Document, error) {
			tarReader := tar.NewReader(bzip2.NewReader(reader))
			var result []vs.Document
			for {
				header, err := tarReader.Next()
				if err == io.EOF {
					break
				}
				if err != nil {
					return nil, err
				}

				// ignore any apple metadata files https://en.wikipedia.org/wiki/AppleSingle_and_AppleDouble_formats
				if strings.HasPrefix(header.Name, "._") {
					continue
				}

				var buf bytes.Buffer
				if _, err := io.Copy(&buf, tarReader); err != nil {
					return nil, err
				}
				content := buf.Bytes()
				ft, err := filetypes.GetFiletype(header.Name, content)
				if err != nil {
					return nil, err
				}
				docs, err := DefaultDocLoaderFunc(ft)(ctx, bytes.NewReader(content))
				if err != nil {
					return nil, err
				}
				result = append(result, docs...)
			}
			return result, nil
		}

	default:
		slog.Error("Unsupported file type", "type", filetype)
		return nil
	}
}
