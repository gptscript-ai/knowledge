package textractor

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

type identityDocumentParser struct {
	blocks []types.Block
	fields []types.IdentityDocumentField
}

func newIdentityDocumentParser(identityDocument types.IdentityDocument) *identityDocumentParser {
	return &identityDocumentParser{
		blocks: identityDocument.Blocks,
		fields: identityDocument.IdentityDocumentFields,
	}
}

func (idp *identityDocumentParser) createIdentityDocument() *IdentityDocument {
	fields, fieldsMap := idp.createFields()

	return &IdentityDocument{
		document:  idp.createDocument(),
		fields:    fields,
		fieldsMap: fieldsMap,
	}
}

func (idp *identityDocumentParser) createDocument() *Document {
	parser := newBlockParser(idp.blocks)
	return parser.createDocument()
}

func (idp *identityDocumentParser) createFields() ([]*IdentityDocumentField, map[IdentityDocumentFieldType]*IdentityDocumentField) {
	fields := make([]*IdentityDocumentField, len(idp.fields))
	fieldsMap := make(map[IdentityDocumentFieldType]*IdentityDocumentField, len(idp.fields))

	for i, f := range idp.fields {
		t := aws.ToString(f.Type.Text)

		fieldType := IdentityDocumentFieldTypeOther
		if t != "" {
			fieldType = IdentityDocumentFieldType(t)
		}

		field := &IdentityDocumentField{
			fieldType:  fieldType,
			value:      aws.ToString(f.ValueDetection.Text),
			confidence: float64(aws.ToFloat32(f.ValueDetection.Confidence)),
			raw:        f,
		}

		if f.ValueDetection.NormalizedValue != nil {
			field.normalizedValue = &NormalizedIdentityDocumentFieldValue{
				valueType: f.ValueDetection.NormalizedValue.ValueType,
				value:     aws.ToString(f.ValueDetection.NormalizedValue.Value),
			}
		}

		fields[i] = field
		fieldsMap[fieldType] = field
	}

	return fields, fieldsMap
}
