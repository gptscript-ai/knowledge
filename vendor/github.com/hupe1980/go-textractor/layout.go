package textractor

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

type LayoutChild interface {
	ID() string
	Text(optFns ...func(*TextLinearizationOptions)) string
	BoundingBox() *BoundingBox
}

type Layout struct {
	base
	children   []LayoutChild
	noNewLines bool
}

func (l *Layout) AddChildren(children ...LayoutChild) {
	l.children = append(l.children, children...)
}

func (l *Layout) Text(optFns ...func(*TextLinearizationOptions)) string {
	opts := DefaultLinerizationOptions

	for _, fn := range optFns {
		fn(&opts)
	}

	if (l.BlockType() == types.BlockTypeLayoutHeader && opts.HideHeaderLayout) ||
		(l.BlockType() == types.BlockTypeLayoutFooter && opts.HideFooterLayout) ||
		(l.BlockType() == types.BlockTypeLayoutFigure && opts.HideFigureLayout) ||
		(l.BlockType() == types.BlockTypeLayoutPageNumber && opts.HidePageNumberLayout) {
		return ""
	}

	var text string

	switch l.BlockType() { // nolint exhaustive
	case types.BlockTypeLayoutList:
		items := make([]string, 0, len(l.children))

		for _, c := range l.children {
			itemText := c.Text(func(tlo *TextLinearizationOptions) {
				*tlo = opts
			})

			if opts.RemoveNewLinesInListElements {
				itemText = strings.ReplaceAll(itemText, "\n", " ")
			}

			items = append(items, fmt.Sprintf("%s%s%s", opts.ListElementPrefix, itemText, opts.ListElementSuffix))
		}

		text = strings.Join(items, opts.ListElementSeparator)
	case types.BlockTypeLayoutPageNumber:
		text = l.linearizeChildren(l.children, opts)
		text = fmt.Sprintf("%s%s%s", opts.PageNumberPrefix, text, opts.PageNumberSuffix)
		text = opts.OnLinerizedPageNumber(text)
	case types.BlockTypeLayoutTitle:
		text = l.linearizeChildren(l.children, opts)
		text = fmt.Sprintf("%s%s%s", opts.TitlePrefix, text, opts.TitleSuffix)
		text = opts.OnLinerizedTitle(text)
	case types.BlockTypeLayoutSectionHeader:
		text = l.linearizeChildren(l.children, opts)
		text = fmt.Sprintf("%s%s%s", opts.SectionHeaderPrefix, text, opts.SectionHeaderSuffix)
		text = opts.OnLinerizedSectionHeader(text)
	default:
		text = l.linearizeChildren(l.children, opts)
	}

	invalidSeparator := strings.Repeat("\n", opts.MaxNumberOfConsecutiveNewLines+1)
	validSeperator := strings.Repeat("\n", opts.MaxNumberOfConsecutiveNewLines)

	for strings.Contains(text, invalidSeparator) {
		text = strings.ReplaceAll(text, invalidSeparator, validSeperator)
	}

	return text
}

func (l *Layout) linearizeChildren(children []LayoutChild, opts TextLinearizationOptions) string {
	var (
		text string
		prev LayoutChild
	)

	for _, group := range groupElementsHorizontally(children, opts.HeuristicOverlapRatio) {
		sort.Slice(group, func(i, j int) bool {
			return group[i].BoundingBox().Left() < group[j].BoundingBox().Left()
		})

		addRowSeparatorIfTableLayout := true

		for i, child := range group {
			childText := child.Text(func(tlo *TextLinearizationOptions) {
				*tlo = opts
			})

			switch child.(type) {
			case *Table:
				text += fmt.Sprintf("%s%s%s", opts.TableLayoutPrefix, childText, opts.TableLayoutSuffix)
				addRowSeparatorIfTableLayout = false
			case *KeyValue:
				text += fmt.Sprintf("%s%s%s", opts.KeyValueLayoutPrefix, childText, opts.KeyValueLayoutSuffix)
				addRowSeparatorIfTableLayout = false
			default:
				if l.BlockType() == types.BlockTypeLayoutTable {
					sep := opts.TableColumnSeparator
					if i == 0 {
						sep = ""
					}

					text += sep + childText
				} else if partOfSameParagraph(prev, child, opts) {
					text += opts.SameParagraphSeparator + childText
				} else {
					sep := ""
					if prev != nil {
						sep = opts.LayoutElementSeparator
					}

					text += sep + childText
				}

				prev = child
			}
		}

		if l.BlockType() == types.BlockTypeLayoutTable && addRowSeparatorIfTableLayout {
			text = text + opts.TableRowSeparator
		}

		prev = &Line{
			base: base{
				boundingBox: NewEnclosingBoundingBox(group...),
			},
		}
	}

	if l.noNewLines {
		// Replace all occurrences of \n with a space
		text = strings.ReplaceAll(text, "\n", " ")

		// Replace consecutive spaces with a single space
		for strings.Contains(text, "  ") {
			text = strings.ReplaceAll(text, "  ", " ")
		}
	}

	return text
}

// groupElementsHorizontally groups elements horizontally based on their vertical positions.
// It takes a slice of elements and an overlap ratio as parameters, and returns a 2D slice of grouped elements.
func groupElementsHorizontally(elements []LayoutChild, overlapRatio float64) [][]LayoutChild {
	// Create a copy of the elements to avoid modifying the original slice
	sortedElements := make([]LayoutChild, len(elements))
	copy(sortedElements, elements)

	// Sort elements based on the top position of their bounding boxes
	sort.Slice(sortedElements, func(i, j int) bool {
		return sortedElements[i].BoundingBox().Top() < sortedElements[j].BoundingBox().Top()
	})

	var groupedElements [][]LayoutChild

	// Check if the sorted elements slice is empty
	if len(sortedElements) == 0 {
		return groupedElements
	}

	// verticalOverlap calculates the vertical overlap between two children
	verticalOverlap := func(child1, child2 LayoutChild) float64 {
		t1 := child1.BoundingBox().Top()
		h1 := child1.BoundingBox().Height()
		t2 := child2.BoundingBox().Top()
		h2 := child2.BoundingBox().Height()

		top := math.Max(t1, t2)
		bottom := math.Min(t1+h1, t2+h2)

		return math.Max(bottom-top, 0)
	}

	// shouldGroup determines whether a line should be grouped with an existing group of lines
	shouldGroup := func(child LayoutChild, group []LayoutChild) bool {
		if len(group) == 0 {
			return false
		}

		maxHeight := 0.0
		for _, l := range group {
			maxHeight = math.Max(maxHeight, l.BoundingBox().Height())
		}

		totalOverlap := 0.0
		for _, l := range group {
			totalOverlap += verticalOverlap(child, l)
		}

		return totalOverlap/maxHeight >= overlapRatio
	}

	// Initialize the first group with the first element
	currentGroup := []LayoutChild{sortedElements[0]}

	// Iterate through the sorted elements and group them horizontally
	for _, element := range sortedElements[1:] {
		if shouldGroup(element, currentGroup) {
			currentGroup = append(currentGroup, element)
		} else {
			groupedElements = append(groupedElements, currentGroup)
			currentGroup = []LayoutChild{element}
		}
	}

	// Add the last group to the result
	groupedElements = append(groupedElements, currentGroup)

	return groupedElements
}

func partOfSameParagraph(child1, child2 LayoutChild, options TextLinearizationOptions) bool {
	if child1 != nil && child2 != nil {
		return math.Abs(child1.BoundingBox().Left()-child2.BoundingBox().Left()) <= options.HeuristicHTolerance*child1.BoundingBox().Width() &&
			math.Abs(child1.BoundingBox().Top()+child1.BoundingBox().Height()-child2.BoundingBox().Top()) <= options.HeuristicOverlapRatio*math.Min(child1.BoundingBox().Height(), child2.BoundingBox().Height())
	}

	return false
}
