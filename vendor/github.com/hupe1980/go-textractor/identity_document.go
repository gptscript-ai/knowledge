package textractor

type IdentityDocument struct {
	document  *Document
	fields    []*IdentityDocumentField
	fieldsMap map[IdentityDocumentFieldType]*IdentityDocumentField
}

func (id *IdentityDocument) Document() *Document {
	return id.document
}

func (id *IdentityDocument) IdentityDocumentType() IdentityDocumentType {
	if f := id.FieldByType(IdentityDocumentFieldTypeIDType); f != nil {
		return IdentityDocumentType(f.Value())
	}

	return IdentityDocumentTypeOther
}

func (id *IdentityDocument) Fields() []*IdentityDocumentField {
	return id.fields
}

func (id *IdentityDocument) FieldByType(ft IdentityDocumentFieldType) *IdentityDocumentField {
	return id.fieldsMap[ft]
}
