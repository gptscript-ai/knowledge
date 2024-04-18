package textractor

import (
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/hupe1980/go-textractor/internal"
)

// Cell defines the interface for a table cell in Textract.
type Cell interface {
	Words() []*Word
	Text(optFns ...func(*TextLinearizationOptions)) string
	Confidence() float64
	OCRConfidence() *OCRConfidence
	IsColumnHeader() bool
	IsTableTitle() bool
	IsTableFooter() bool
	IsTableSummary() bool
	IsTableSectionTitle() bool
	IsMerged() bool
}

// TableMergedCell represents a merged cell in a table.
type TableMergedCell struct {
	cell
	cells []*TableCell
}

// Words returns the words in the merged cell.
func (tmc *TableMergedCell) Words() []*Word {
	words := make([][]*Word, 0, len(tmc.cells))

	for _, c := range tmc.cells {
		words = append(words, c.Words())
	}

	return internal.Concatenate(words...)
}

// Text returns the text content of the merged cell.
func (tmc *TableMergedCell) Text(_ ...func(*TextLinearizationOptions)) string {
	words := tmc.Words()

	texts := make([]string, len(words))
	for i, w := range words {
		texts[i] = w.Text()
	}

	return strings.Join(texts, " ")
}

// OCRConfidence returns the OCR confidence for the merged cell.
func (tmc *TableMergedCell) OCRConfidence() *OCRConfidence {
	meanValues := make([]float64, 0, len(tmc.cells))
	maxValues := make([]float64, 0, len(tmc.cells))
	minValues := make([]float64, 0, len(tmc.cells))

	for _, cell := range tmc.cells {
		c := cell.OCRConfidence()
		meanValues = append(meanValues, c.Mean())
		maxValues = append(maxValues, c.Max())
		minValues = append(minValues, c.Min())
	}

	return &OCRConfidence{
		mean: internal.Mean(meanValues),
		max:  slices.Max(maxValues),
		min:  slices.Min(minValues),
	}
}

// TableCell represents a cell in a table.
type TableCell struct {
	cell
	words            []*Word
	selectionElement *SelectionElement
}

// Words returns the words in the table cell.
func (tc *TableCell) Words() []*Word {
	return tc.words
}

// SelectionElement returns the selection element associated with the table cell.
func (tc *TableCell) SelectionElement() *SelectionElement {
	return tc.selectionElement
}

// Text returns the text content of the table cell.
func (tc *TableCell) Text(optFns ...func(*TextLinearizationOptions)) string {
	if tc.selectionElement != nil {
		return tc.selectionElement.Text(optFns...)
	}

	texts := make([]string, len(tc.words))
	for i, w := range tc.words {
		texts[i] = w.Text()
	}

	return strings.Join(texts, " ")
}

// OCRConfidence returns the OCR confidence for the table cell.
func (tc *TableCell) OCRConfidence() *OCRConfidence {
	if tc.selectionElement != nil {
		return &OCRConfidence{
			mean: tc.SelectionElement().Confidence(),
			min:  tc.SelectionElement().Confidence(),
			max:  tc.SelectionElement().Confidence(),
		}
	}

	scores := make([]float64, len(tc.words))
	for i, w := range tc.words {
		scores[i] = w.Confidence()
	}

	return &OCRConfidence{
		mean: internal.Mean(scores),
		max:  slices.Max(scores),
		min:  slices.Min(scores),
	}
}

// cell represents the base information shared among different types of table cells.
type cell struct {
	base
	rowIndex    int
	columnIndex int
	rowSpan     int
	columnSpan  int
	entityTypes []types.EntityType
}

// newCell creates a new cell instance from the provided Textract block and page information.
func newCell(cb types.Block, p *Page) cell {
	return cell{
		base:        newBase(cb, p),
		rowIndex:    int(aws.ToInt32(cb.RowIndex)),
		columnIndex: int(aws.ToInt32(cb.ColumnIndex)),
		rowSpan:     int(aws.ToInt32(cb.RowSpan)),
		columnSpan:  int(aws.ToInt32(cb.ColumnSpan)),
	}
}

// IsColumnHeader checks if the cell is a column header.
func (c *cell) IsColumnHeader() bool {
	return slices.Contains(c.entityTypes, types.EntityTypeColumnHeader)
}

// IsTableTitle checks if the cell is a table title.
func (c *cell) IsTableTitle() bool {
	return slices.Contains(c.entityTypes, types.EntityTypeTableTitle)
}

// IsTableFooter checks if the cell is a table footer.
func (c *cell) IsTableFooter() bool {
	return slices.Contains(c.entityTypes, types.EntityTypeTableFooter)
}

// IsTableSummary checks if the cell is a table summary.
func (c *cell) IsTableSummary() bool {
	return slices.Contains(c.entityTypes, types.EntityTypeTableSummary)
}

// IsTableSectionTitle checks if the cell is a table section title.
func (c *cell) IsTableSectionTitle() bool {
	return slices.Contains(c.entityTypes, types.EntityTypeTableSectionTitle)
}

// IsMerged checks if the cell is part of a merged group.
func (c *cell) IsMerged() bool {
	return c.rowSpan > 1 || c.columnSpan > 1
}
