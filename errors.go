package bucket

import "errors"

var (
	ErrDocumentTypeAlredyExists = errors.New("document type alredy exists")
	ErrDocumentTypeDoesntExists = errors.New("document type doesn't exist")

	ErrEmptyField  = errors.New("field must be filled")
	ErrEmptyIndex  = errors.New("index must be filled")
	ErrEmptyType   = errors.New("source type must set")
	ErrEmptySource = errors.New("source name must set")
	ErrEmptyRefTag = errors.New("referenced tag must set")

	ErrConjunctionAndDisjunktionIsNil = errors.New("conjunction and disjunction are nil")

	ErrEndAsTimeZero = errors.New("endAsTime is zero instant")

	ErrFirstParameterNotStruct = errors.New("first parameter is not a struct")
	ErrInputStructPointer      = errors.New("input struct must be pointer")

	ErrInvalidBulkContainer = errors.New("container must be *[]T, with length of ids array")
	ErrInvalidGetDocumentTypesParam = errors.New("internal error: value should be pointer for getDocumentTypes")
)
