package textractor

import (
	"strings"
)

// Compile time check to ensure Line satisfies the LayoutChild interface.
var _ LayoutChild = (*Line)(nil)

type Line struct {
	base
	words []*Word
}

func (l *Line) Words() []*Word {
	return l.words
}

func (l *Line) Text(_ ...func(*TextLinearizationOptions)) string {
	texts := make([]string, len(l.words))
	for i, w := range l.words {
		texts[i] = w.Text()
	}

	return strings.Join(texts, " ")
}

func (l *Line) String() string {
	return l.Text()
}
