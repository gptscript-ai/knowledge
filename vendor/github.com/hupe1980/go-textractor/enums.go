package textractor

// IdentityDocumentFieldType represents the type of fields in an identity document.
type IdentityDocumentFieldType string

const (
	IdentityDocumentFieldTypeFirstName        IdentityDocumentFieldType = "FIRST_NAME"
	IdentityDocumentFieldTypeLastName         IdentityDocumentFieldType = "LAST_NAME"
	IdentityDocumentFieldTypeMiddleName       IdentityDocumentFieldType = "MIDDLE_NAME"
	IdentityDocumentFieldTypeSuffix           IdentityDocumentFieldType = "Suffix"
	IdentityDocumentFieldTypeCityInAddress    IdentityDocumentFieldType = "CITY_IN_ADDRESS"
	IdentityDocumentFieldTypeZipCodeInAddress IdentityDocumentFieldType = "ZIP_CODE_IN_ADDRESS"
	IdentityDocumentFieldTypeStateInAddress   IdentityDocumentFieldType = "STATE_IN_ADDRESS"
	IdentityDocumentFieldTypeStateName        IdentityDocumentFieldType = "STATE_NAME"
	IdentityDocumentFieldTypeDocumentNumber   IdentityDocumentFieldType = "DOCUMENT_NUMBER"
	IdentityDocumentFieldTypeExpirationDate   IdentityDocumentFieldType = "EXPIRATION_DATE"
	IdentityDocumentFieldTypeDateOfBirth      IdentityDocumentFieldType = "DATE_OF_BIRTH"
	IdentityDocumentFieldTypeDateOfIssue      IdentityDocumentFieldType = "DATE_OF_ISSUE"
	IdentityDocumentFieldTypeIDType           IdentityDocumentFieldType = "ID_TYPE"
	IdentityDocumentFieldTypeEndorsements     IdentityDocumentFieldType = "ENDORSEMENTS"
	IdentityDocumentFieldTypeVeteran          IdentityDocumentFieldType = "VETERAN"
	IdentityDocumentFieldTypeRestrictions     IdentityDocumentFieldType = "RESTRICTIONS"
	IdentityDocumentFieldTypeClass            IdentityDocumentFieldType = "CLASS"
	IdentityDocumentFieldTypeAddress          IdentityDocumentFieldType = "ADDRESS"
	IdentityDocumentFieldTypeCounty           IdentityDocumentFieldType = "COUNTY"
	IdentityDocumentFieldTypePlaceOfBirth     IdentityDocumentFieldType = "PLACE_OF_BIRTH"
	IdentityDocumentFieldTypeMRZCode          IdentityDocumentFieldType = "MRZ_CODE"
	IdentityDocumentFieldTypeOther            IdentityDocumentFieldType = "Other"
)

// IdentityDocumentType represents the type of an identity document.
type IdentityDocumentType string

const (
	IdentityDocumentTypeDriverLicenseFront IdentityDocumentType = "DRIVER LICENSE FRONT"
	IdentityDocumentTypePassport           IdentityDocumentType = "PASSPORT"
	IdentityDocumentTypeOther              IdentityDocumentType = "OTHER"
)
