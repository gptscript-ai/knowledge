package textractor

import "github.com/aws/aws-sdk-go-v2/service/textract/types"

type expenseDocumentParser struct {
	blocks []types.Block
}

func newExpenseDocumentParser(identityDocument types.ExpenseDocument) *expenseDocumentParser {
	return &expenseDocumentParser{
		blocks: identityDocument.Blocks,
	}
}

func (edp *expenseDocumentParser) createExpenseDocument() *ExpenseDocument {
	return &ExpenseDocument{
		document: edp.createDocument(),
	}
}

func (edp *expenseDocumentParser) createDocument() *Document {
	parser := newBlockParser(edp.blocks)
	return parser.createDocument()
}
