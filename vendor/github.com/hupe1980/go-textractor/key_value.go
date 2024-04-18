package textractor

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/hupe1980/go-textractor/internal"
)

// Compile time check to ensure KeyValue satisfies the LayoutChild interface.
var _ LayoutChild = (*KeyValue)(nil)

// KeyValue represents a key-value pair in a document.
type KeyValue struct {
	base
	key   *Key
	value *Value
	page  *Page
}

// Key returns the key of the key-value pair.
func (kv *KeyValue) Key() *Key {
	return kv.key
}

// Value returns the value of the key-value pair.
func (kv *KeyValue) Value() *Value {
	return kv.value
}

// Confidence calculates the confidence score for a key value.
func (kv *KeyValue) Confidence() float64 {
	scores := make([]float64, 0)

	if kv.Key() != nil {
		scores = append(scores, kv.Key().Confidence())
	}

	if kv.Value() != nil {
		scores = append(scores, kv.Value().Confidence())
	}

	return internal.Mean(scores)
}

// OCRConfidence returns the OCR confidence for the key-value pair.
func (kv *KeyValue) OCRConfidence() *OCRConfidence {
	keyOCR := &OCRConfidence{}
	if kv.Key() != nil {
		keyOCR = kv.Key().OCRConfidence()
	}

	valueOCR := &OCRConfidence{}
	if kv.Value() != nil {
		valueOCR = kv.Value().OCRConfidence()
	}

	return &OCRConfidence{
		mean: internal.Mean([]float64{keyOCR.Mean(), valueOCR.Mean()}),
		min:  math.Min(keyOCR.Min(), valueOCR.Min()),
		max:  math.Max(keyOCR.Max(), valueOCR.Max()),
	}
}

// BoundingBox returns the bounding box that encloses the key-value pair.
func (kv *KeyValue) BoundingBox() *BoundingBox {
	return NewEnclosingBoundingBox[BoundingBoxAccessor](kv.Key(), kv.Value())
}

// Polygon returns the polygon representing the key-value pair.
func (kv *KeyValue) Polygon() Polygon {
	// TODO
	panic("not implemented")
}

// Words returns the words in the key-value pair.
func (kv *KeyValue) Words() []*Word {
	return internal.Concatenate(kv.Key().Words(), kv.Value().Words())
}

// Text returns the text content of the key-value pair.
func (kv *KeyValue) Text(optFns ...func(*TextLinearizationOptions)) string {
	opts := DefaultLinerizationOptions

	for _, fn := range optFns {
		fn(&opts)
	}

	keyText := kv.Key().Text()
	keyText = fmt.Sprintf("%s%s%s", opts.KeyPrefix, keyText, opts.KeySuffix)

	valueText := kv.Value().Text()
	valueText = fmt.Sprintf("%s%s%s", opts.ValuePrefix, valueText, opts.ValueSuffix)

	if len(keyText) == 0 && len(valueText) == 0 {
		return ""
	}

	text := fmt.Sprintf("%s%s%s", keyText, opts.SameParagraphSeparator, valueText)
	if kv.Value().SelectionElement() != nil {
		text = fmt.Sprintf("%s%s%s", valueText, opts.SameParagraphSeparator, keyText)
	}

	return fmt.Sprintf("%s%s%s", opts.KeyValuePrefix, text, opts.KeyValueSuffix)
}

// String returns the string representation of the key-value pair.
func (kv *KeyValue) String() string {
	if kv.Value().SelectionElement() != nil {
		return fmt.Sprintf("%s %s", kv.Value(), kv.Key())
	}

	return fmt.Sprintf("%s : %s", kv.Key(), kv.Value())
}

// Key represents the key part of a key-value pair.
type Key struct {
	base
	words []*Word
}

// Words returns the words in the key.
func (k *Key) Words() []*Word {
	return k.words
}

// Text returns the text content of the key.
func (k *Key) Text() string {
	texts := make([]string, len(k.words))
	for i, w := range k.words {
		texts[i] = w.Text()
	}

	return strings.Join(texts, " ")
}

// OCRConfidence returns the OCR confidence for the key.
func (k *Key) OCRConfidence() *OCRConfidence {
	scores := make([]float64, len(k.words))
	for i, w := range k.words {
		scores[i] = w.Confidence()
	}

	return &OCRConfidence{
		mean: internal.Mean(scores),
		max:  slices.Max(scores),
		min:  slices.Min(scores),
	}
}

// String returns the string representation of the key.
func (k *Key) String() string {
	return k.Text()
}

// Value represents the value part of a key-value pair.
type Value struct {
	base
	words            []*Word
	selectionElement *SelectionElement
}

// Words returns the words in the value.
func (v *Value) Words() []*Word {
	return v.words
}

// SelectionElement returns the selection element associated with the table cell.
func (v *Value) SelectionElement() *SelectionElement {
	return v.selectionElement
}

// Text returns the text content of the value.
func (v *Value) Text(optFns ...func(*TextLinearizationOptions)) string {
	if v.selectionElement != nil {
		return v.selectionElement.Text(optFns...)
	}

	texts := make([]string, len(v.words))
	for i, w := range v.words {
		texts[i] = w.Text()
	}

	text := strings.Join(texts, " ")

	// Replace all occurrences of \n with a space
	text = strings.ReplaceAll(text, "\n", " ")

	// Replace consecutive spaces with a single space
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	return text
}

// OCRConfidence returns the OCR confidence for the value.
func (v *Value) OCRConfidence() *OCRConfidence {
	if v.selectionElement != nil {
		return &OCRConfidence{
			mean: v.SelectionElement().Confidence(),
			min:  v.SelectionElement().Confidence(),
			max:  v.SelectionElement().Confidence(),
		}
	}

	scores := make([]float64, len(v.words))
	for i, w := range v.words {
		scores[i] = w.Confidence()
	}

	return &OCRConfidence{
		mean: internal.Mean(scores),
		max:  slices.Max(scores),
		min:  slices.Min(scores),
	}
}

// String returns the string representation of the value.
func (v *Value) String() string {
	return v.Text()
}
