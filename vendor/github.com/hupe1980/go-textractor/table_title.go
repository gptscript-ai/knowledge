package textractor

import "strings"

// TableTitle represents the title of a table, containing a collection of words.
type TableTitle struct {
	base
	words []*Word
}

// Words returns the words constituting the table title.
func (tt *TableTitle) Words() []*Word {
	return tt.words
}

// Text returns the concatenated text of the table title, using default or provided linearization options.
func (tt *TableTitle) Text(_ ...func(*TextLinearizationOptions)) string {
	texts := make([]string, len(tt.words))
	for i, w := range tt.words {
		texts[i] = w.Text()
	}

	return strings.Join(texts, " ")
}
