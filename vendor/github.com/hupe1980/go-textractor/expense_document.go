package textractor

type ExpenseDocument struct {
	document     *Document
	summaryField []*ExpenseField
}

func (ed *ExpenseDocument) SummaryFields() []*ExpenseField {
	return ed.summaryField
}
