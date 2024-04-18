package textractor

import "github.com/aws/aws-sdk-go-v2/service/textract/types"

// Word represents a word extracted by Textract.
type Word struct {
	base                     // Embedding the base type for common attributes
	text      string         // The text content of the word
	textType  types.TextType // The text type of the word (e.g., printed or handwriting)
	line      *Line          // The line to which the word belongs
	tableCell *TableCell     // The table cell to which the word belongs
}

// Text returns the text content of the word.
func (w *Word) Text() string {
	return w.text
}

// TextType returns the text type of the word.
func (w *Word) TextType() types.TextType {
	return w.textType
}

// IsPrinted checks if the word is printed text.
func (w *Word) IsPrinted() bool {
	return w.TextType() == types.TextTypePrinted
}

// IsHandwriting checks if the word is handwriting.
func (w *Word) IsHandwriting() bool {
	return w.TextType() == types.TextTypeHandwriting
}
