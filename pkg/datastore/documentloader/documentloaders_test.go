package documentloader

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDocumentLoaderConfig_ValidLoader(t *testing.T) {
	cfg, err := GetDocumentLoaderConfig("pdf")
	assert.NoError(t, err)
	assert.IsTypef(t, PDFOptions{}, cfg, "cfg is not of type PDFOptions")
}

func TestGetDocumentLoaderConfig_InvalidLoader(t *testing.T) {
	_, err := GetDocumentLoaderConfig("invalid")
	assert.Error(t, err)
}

func TestGetDocumentLoaderFunc_ValidLoaderWithoutConfig(t *testing.T) {
	_, err := GetDocumentLoaderFunc("plaintext", nil)
	assert.NoError(t, err)
}

func TestGetDocumentLoaderFunc_ValidLoaderWithInvalidConfig(t *testing.T) {
	_, err := GetDocumentLoaderFunc("pdf", "invalid")
	assert.Error(t, err)
}

func TestGetDocumentLoaderFunc_ValidLoaderWithValidConfig(t *testing.T) {
	_, err := GetDocumentLoaderFunc("pdf", PDFOptions{})
	assert.NoError(t, err)
}

func TestGetDocumentLoaderFunc_InvalidLoader(t *testing.T) {
	_, err := GetDocumentLoaderFunc("invalid", nil)
	assert.Error(t, err)
}

func TestGetDocumentLoaderFunc_LoadPlainText(t *testing.T) {
	loaderFunc, _ := GetDocumentLoaderFunc("plaintext", nil)
	docs, err := loaderFunc(context.Background(), strings.NewReader("test"))
	assert.NoError(t, err)
	assert.Len(t, docs, 1)
	assert.Equal(t, "test", docs[0].Content)
}

func TestGetDocumentLoaderFunc_LoadPDF(t *testing.T) {
	loaderFunc, _ := GetDocumentLoaderFunc("pdf", PDFOptions{})
	_, err := loaderFunc(context.Background(), strings.NewReader("test"))
	assert.Error(t, err)
}
