package vectorstore

type ErrCollectionNotFound struct {
	Collection string
}

func (e ErrCollectionNotFound) Error() string {
	return "collection not found: " + e.Collection
}
