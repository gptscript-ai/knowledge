package textractor

import (
	"strings"

	"github.com/hupe1980/go-textractor/internal"
)

// Document represents a document consisting of multiple pages.
type Document struct {
	pages []*Page
}

// Pages returns the slice of Page objects in the document.
func (d *Document) Pages() []*Page {
	return d.pages
}

// Words returns a slice containing all the words in the document.
func (d *Document) Words() []*Word {
	words := make([][]*Word, 0, len(d.Pages()))

	for _, p := range d.Pages() {
		words = append(words, p.Words())
	}

	return internal.Concatenate(words...)
}

// Lines returns a slice containing all the lines in the document.
func (d *Document) Lines() []*Line {
	lines := make([][]*Line, 0, len(d.Pages()))

	for _, p := range d.Pages() {
		lines = append(lines, p.Lines())
	}

	return internal.Concatenate(lines...)
}

// Tables returns a slice containing all the tables in the document.
func (d *Document) Tables() []*Table {
	tables := make([][]*Table, 0, len(d.Pages()))

	for _, p := range d.Pages() {
		tables = append(tables, p.Tables())
	}

	return internal.Concatenate(tables...)
}

// KeyValues returns a slice containing all the key-value pairs in the document.
func (d *Document) KeyValues() []*KeyValue {
	keyValues := make([][]*KeyValue, 0, len(d.Pages()))

	for _, p := range d.Pages() {
		keyValues = append(keyValues, p.KeyValues())
	}

	return internal.Concatenate(keyValues...)
}

// Signatures returns a slice containing all the signatures in the document.
func (d *Document) Signatures() []*Signature {
	signatures := make([][]*Signature, 0, len(d.Pages()))

	for _, p := range d.Pages() {
		signatures = append(signatures, p.Signatures())
	}

	return internal.Concatenate(signatures...)
}

// Text linearizes the document into a single text string, optionally applying specified options.
func (d *Document) Text(optFns ...func(*TextLinearizationOptions)) string {
	pageTexts := make([]string, len(d.Pages()))

	for i, p := range d.Pages() {
		pageTexts[i] = p.Text(optFns...)
	}

	return strings.Join(pageTexts, "\n")
}
