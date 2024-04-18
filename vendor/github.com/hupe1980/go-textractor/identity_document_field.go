package textractor

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

// IdentityDocumentField represents a field extracted from an identity document by Textract.
type IdentityDocumentField struct {
	fieldType       IdentityDocumentFieldType
	value           string
	confidence      float64
	normalizedValue *NormalizedIdentityDocumentFieldValue
	raw             types.IdentityDocumentField
}

// FieldType returns the type of the identity document field.
func (idf *IdentityDocumentField) FieldType() IdentityDocumentFieldType {
	return idf.fieldType
}

// Value returns the value of the identity document field.
func (idf *IdentityDocumentField) Value() string {
	return idf.value
}

// Confidence returns the confidence score associated with the field extraction.
func (idf *IdentityDocumentField) Confidence() float64 {
	return idf.confidence
}

// IsNormalized checks if the field value is normalized.
func (idf *IdentityDocumentField) IsNormalized() bool {
	return idf.normalizedValue != nil
}

// NormalizedValue returns the normalized value of the identity document field.
func (idf *IdentityDocumentField) NormalizedValue() *NormalizedIdentityDocumentFieldValue {
	return idf.normalizedValue
}

// NormalizedIdentityDocumentFieldValue represents a normalized value of an identity document field.
type NormalizedIdentityDocumentFieldValue struct {
	valueType types.ValueType
	value     string
}

// ValueType returns the type of the normalized value.
func (nidfv NormalizedIdentityDocumentFieldValue) ValueType() types.ValueType {
	return nidfv.valueType
}

// Value returns the string representation of the normalized value.
func (nidfv NormalizedIdentityDocumentFieldValue) Value() string {
	return nidfv.value
}

// DateValue returns the time representation of the normalized date value.
func (nidfv NormalizedIdentityDocumentFieldValue) DateValue() (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05", nidfv.value)
}
