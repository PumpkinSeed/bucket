package odatas

import "errors"

var (
	ErrDocumentTypeAlredyExists = errors.New("document type alredy exists")
	ErrDocumentTypeDoesntExists = errors.New("document type doesn't exist")

	ErrEmptyField  = errors.New("field must be filled")
	ErrEmptyIndex  = errors.New("index must be filled")
	ErrEmptyType   = errors.New("source type must set")
	ErrEmptySource = errors.New("source name must set")

	ErrConjunctionAndDisjunktionIsNil = errors.New("conjunction and disjunction are nil")

	ErrEndAsTimeZero = errors.New("endAsTime is zero instant")

	ErrFirstParameterNotStruct = errors.New("first parameter is not a struct")
)
