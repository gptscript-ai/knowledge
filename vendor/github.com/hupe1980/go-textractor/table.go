package textractor

import (
	"cmp"
	"encoding/csv"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/hupe1980/go-textractor/internal"
	"github.com/olekukonko/tablewriter"
)

// Compile time check to ensure Table satisfies the LayoutChild interface.
var _ LayoutChild = (*Table)(nil)

type Table struct {
	base
	title       *TableTitle
	footers     []*TableFooter
	mergedCells []*TableMergedCell
	cells       []*TableCell
}

func (t *Table) Words() []*Word {
	words := make([][]*Word, 0, len(t.cells))

	for _, c := range t.cells {
		words = append(words, c.Words())
	}

	return internal.Concatenate(words...)
}

func (t *Table) Text(optFns ...func(*TextLinearizationOptions)) string {
	opts := DefaultLinerizationOptions

	for _, fn := range optFns {
		fn(&opts)
	}

	var tableText string

	switch opts.TableLinearizationFormat {
	case "plaintext":
		texts := []string{}

		for _, r := range t.Rows() {
			cellText := ""

			for i, c := range r.Cells() {
				if i == 0 {
					cellText += c.Text()
				} else {
					cellText += opts.TableColumnSeparator + c.Text()
				}
			}

			texts = append(texts, cellText)
		}

		tableText = strings.Join(texts, opts.TableRowSeparator)
	case "markdown":
		tableString := &strings.Builder{}

		tw := tablewriter.NewWriter(tableString)
		tw.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		tw.SetCenterSeparator("|")

		var (
			header []string
			data   [][]string
		)

		for i, r := range t.Rows() {
			if i == 0 {
				header = make([]string, 0, len(r.Cells()))
				for _, c := range r.Cells() {
					header = append(header, c.Text())
				}
			} else {
				rowData := make([]string, 0, len(r.Cells()))
				for _, c := range r.Cells() {
					rowData = append(rowData, c.Text())
				}

				data = append(data, rowData)
			}
		}

		tw.SetHeader(header)
		tw.AppendBulk(data)
		tw.Render()

		tableText = tableString.String()
	default:
		panic(fmt.Sprintf("unknown table format: %s", opts.TableLinearizationFormat))
	}

	return fmt.Sprintf("%s%s%s", opts.TablePrefix, tableText, opts.TableSuffix)
}

func (t *Table) RowCount() int {
	max := slices.MaxFunc(t.cells, func(a, b *TableCell) int {
		return cmp.Compare(a.rowIndex+a.rowSpan-1, b.rowIndex+b.rowSpan-1)
	})

	return max.rowIndex
}

type CellAtOptions struct {
	IgnoreMergedCells bool
}

func (t *Table) CellAt(rowIndex, columnIndex int, optFns ...func(*CellAtOptions)) Cell {
	opts := CellAtOptions{
		IgnoreMergedCells: true,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	if !opts.IgnoreMergedCells {
		for _, mc := range t.mergedCells {
			if mc.columnIndex <= columnIndex &&
				mc.columnIndex+mc.columnSpan > columnIndex &&
				mc.rowIndex <= rowIndex &&
				mc.rowIndex+mc.rowSpan > rowIndex {
				return mc
			}
		}
	}

	for _, c := range t.cells {
		if c.columnIndex == columnIndex && c.rowIndex == rowIndex {
			return c
		}
	}

	return nil
}

type RowCellsAtOptions struct {
	IgnoreMergedCells bool
}

func (t *Table) RowCellsAt(rowIndex int, optFns ...func(*RowCellsAtOptions)) []Cell {
	opts := RowCellsAtOptions{
		IgnoreMergedCells: true,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	cells := make([]Cell, 0)
	mergedCellIDs := make([]string, 0)

	if opts.IgnoreMergedCells {
		for _, mc := range t.mergedCells {
			if mc.rowIndex <= rowIndex && mc.rowIndex+mc.rowSpan > rowIndex {
				cells = append(cells, mc)
				for _, c := range mc.cells {
					mergedCellIDs = append(mergedCellIDs, c.ID())
				}
			}
		}
	}

	for _, c := range t.cells {
		if c.rowIndex == rowIndex && !slices.Contains(mergedCellIDs, c.ID()) {
			cells = append(cells, c)
		}
	}

	return cells
}

type TableRow struct {
	cells []Cell
}

func (tr *TableRow) Cells() []Cell {
	return tr.cells
}

// OCRConfidence returns the OCR confidence for the table row.
func (tr *TableRow) OCRConfidence() *OCRConfidence {
	meanValues := make([]float64, 0, len(tr.cells))
	maxValues := make([]float64, 0, len(tr.cells))
	minValues := make([]float64, 0, len(tr.cells))

	for _, cell := range tr.cells {
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

type RowsOptions struct {
	IgnoreMergedCells bool
}

func (t *Table) Rows(optFns ...func(*RowsOptions)) []*TableRow {
	opts := RowsOptions{
		IgnoreMergedCells: true,
	}

	for _, fn := range optFns {
		fn(&opts)
	}

	rowCount := t.RowCount()
	rows := make([]*TableRow, 0, rowCount)

	for i := 1; i <= rowCount; i++ {
		rows = append(rows, &TableRow{
			cells: t.RowCellsAt(i, func(rcao *RowCellsAtOptions) {
				rcao.IgnoreMergedCells = opts.IgnoreMergedCells
			}),
		})
	}

	return rows
}

func (t *Table) ToCSV(w io.Writer) error {
	cw := csv.NewWriter(w)

	defer cw.Flush()

	var (
		header []string
		data   [][]string
	)

	for i, r := range t.Rows(func(ro *RowsOptions) {
		ro.IgnoreMergedCells = true
	}) {
		if i == 0 {
			header = make([]string, 0, len(r.Cells()))
			for _, c := range r.Cells() {
				header = append(header, c.Text())
			}
		} else {
			rowData := make([]string, 0, len(r.Cells()))
			for _, c := range r.Cells() {
				rowData = append(rowData, c.Text())
			}

			data = append(data, rowData)
		}
	}

	if err := cw.Write(header); err != nil {
		return err
	}

	if err := cw.WriteAll(data); err != nil {
		return err
	}

	return nil
}
