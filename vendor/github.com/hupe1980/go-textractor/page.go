package textractor

import (
	"slices"
	"sort"
	"strings"
)

type Page struct {
	id         string
	number     int
	width      float64
	height     float64
	childIDs   []string
	words      []*Word
	lines      []*Line
	keyValues  []*KeyValue
	tables     []*Table
	layouts    []*Layout
	queries    []*Query
	signatures []*Signature
}

func (p *Page) ID() string {
	return p.id
}

func (p *Page) Number() int {
	return p.number
}

func (p *Page) Width() float64 {
	return p.width
}

func (p *Page) Height() float64 {
	return p.height
}

func (p *Page) Words() []*Word {
	return p.words
}

func (p *Page) Lines() []*Line {
	return p.lines
}

func (p *Page) Tables() []*Table {
	return p.tables
}

func (p *Page) KeyValues() []*KeyValue {
	return p.keyValues
}

func (p *Page) Layouts() []*Layout {
	return p.layouts
}

func (p *Page) Queries() []*Query {
	return p.queries
}

func (p *Page) Signatures() []*Signature {
	return p.signatures
}

func (p *Page) AddLayouts(layouts ...*Layout) {
	p.layouts = append(p.layouts, layouts...)
}

func (p *Page) Text(optFns ...func(*TextLinearizationOptions)) string {
	// Create a copy of the layouts to avoid modifying the original slice
	sortedLayouts := make([]*Layout, len(p.layouts))
	copy(sortedLayouts, p.layouts)

	// Sort layouts based on the reading order
	sort.Slice(sortedLayouts, func(i, j int) bool {
		return sortedLayouts[i].BoundingBox().Top() < sortedLayouts[j].BoundingBox().Top()
	})

	pageTexts := make([]string, len(sortedLayouts))

	for i, l := range sortedLayouts {
		text := l.Text(optFns...)

		pageTexts[i] = text
	}

	return strings.Join(pageTexts, "\n")
}

func (p *Page) SearchValueByKey(key string) []*KeyValue {
	searchKey := strings.ToLower(key)

	var result []*KeyValue

	for _, kv := range p.keyValues {
		if key := kv.Key(); key != nil {
			if strings.Contains(strings.ToLower(key.Text()), searchKey) {
				result = append(result, kv)
			}
		}
	}

	return result
}

func (p *Page) isChild(id string) bool {
	return slices.Contains(p.childIDs, id)
}
