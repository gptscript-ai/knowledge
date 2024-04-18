package textractor

import (
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

// SelectionElement represents an element with selection status.
type SelectionElement struct {
	base
	status types.SelectionStatus
}

// Status returns the selection status of the element.
func (se *SelectionElement) Status() types.SelectionStatus {
	return se.status
}

// IsSelected checks if the element is selected.
func (se *SelectionElement) IsSelected() bool {
	return se.Status() == types.SelectionStatusSelected
}

// Text returns the text representation of the selection element.
// It considers the selection status and applies linearization options.
func (se *SelectionElement) Text(optFns ...func(*TextLinearizationOptions)) string {
	opts := DefaultLinerizationOptions

	for _, fn := range optFns {
		fn(&opts)
	}

	text := opts.SelectionElementNotSelected
	if se.IsSelected() {
		text = opts.SelectionElementSelected
	}

	return text
}

func (se *SelectionElement) String() string {
	return se.Text()
}
