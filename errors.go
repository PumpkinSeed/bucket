package odatas

import "errors"

var(
	ErrDocumentTypeAlredyExists = errors.New("document type alredy exists")
	ErrDocumentTypeDoesntExists = errors.New("document type doesn't exist")
)
