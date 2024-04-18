package textractor

import "strings"

// TableFooter represents the footer of a table block.
type TableFooter struct {
	base
	words []*Word
}

// Words returns the words within the table footer.
func (tf *TableFooter) Words() []*Word {
	return tf.words
}

// Text returns the concatenated text of all words in the table footer.
func (tf *TableFooter) Text(_ ...func(*TextLinearizationOptions)) string {
	texts := make([]string, len(tf.words))
	for i, w := range tf.words {
		texts[i] = w.Text()
	}

	return strings.Join(texts, " ")
}
