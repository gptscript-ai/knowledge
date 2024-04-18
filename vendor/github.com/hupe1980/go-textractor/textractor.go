package textractor

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

// DocumentAPIOutput represents the output of the Textract Document API.
type DocumentAPIOutput struct {
	DocumentMetadata *types.DocumentMetadata `json:"DocumentMetadata"`
	Blocks           []types.Block           `json:"Blocks"`
}

// ParseDocumentAPIOutput parses the Textract Document API output into a Document.
func ParseDocumentAPIOutput(output *DocumentAPIOutput) (*Document, error) {
	parser := newBlockParser(output.Blocks)

	document := parser.createDocument()

	if len(document.pages) != int(aws.ToInt32(output.DocumentMetadata.Pages)) {
		return nil, fmt.Errorf("number of pages %d does not match metadata %d", len(document.pages), aws.ToInt32(output.DocumentMetadata.Pages))
	}

	return document, nil
}

// AnalyzeIDOutput represents the output of the Textract Analyze ID API.
type AnalyzeIDOutput struct {
	DocumentMetadata  *types.DocumentMetadata  `json:"DocumentMetadata"`
	IdentityDocuments []types.IdentityDocument `json:"IdentityDocuments"`
}

// ParseAnalyzeIDOutput parses the Textract Analyze ID API output into a slice of IdentityDocument.
func ParseAnalyzeIDOutput(output *AnalyzeIDOutput) ([]*IdentityDocument, error) {
	parsedIdentityDocuments := make([]*IdentityDocument, len(output.IdentityDocuments))

	for i, d := range output.IdentityDocuments {
		parser := newIdentityDocumentParser(d)
		parsedIdentityDocuments[i] = parser.createIdentityDocument()
	}

	if len(parsedIdentityDocuments) != int(aws.ToInt32(output.DocumentMetadata.Pages)) {
		return nil, fmt.Errorf("number of pages %d does not match metadata %d", len(parsedIdentityDocuments), aws.ToInt32(output.DocumentMetadata.Pages))
	}

	return parsedIdentityDocuments, nil
}

// AnalyzeExpenseOutput represents the output of the Textract Analyze Expense API.
type AnalyzeExpenseOutput struct {
	DocumentMetadata *types.DocumentMetadata `json:"DocumentMetadata"`
	ExpenseDocuments []types.ExpenseDocument `json:"ExpenseDocuments"`
}

// ParseAnalyzeExpenseOutput parses the Textract Analyze Expense API output into a slice of ExpenseDocument.
func ParseAnalyzeExpenseOutput(output *AnalyzeExpenseOutput) ([]*ExpenseDocument, error) {
	parsedExpenseDocuments := make([]*ExpenseDocument, len(output.ExpenseDocuments))

	for i, d := range output.ExpenseDocuments {
		parser := newExpenseDocumentParser(d)
		parsedExpenseDocuments[i] = parser.createExpenseDocument()
	}

	if len(parsedExpenseDocuments) != int(aws.ToInt32(output.DocumentMetadata.Pages)) {
		return nil, fmt.Errorf("number of pages %d does not match metadata %d", len(parsedExpenseDocuments), aws.ToInt32(output.DocumentMetadata.Pages))
	}

	return parsedExpenseDocuments, nil
}
