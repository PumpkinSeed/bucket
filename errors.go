package bucket

import "errors"

var (
	// ErrDocumentTypeAlreadyExists document type already exist
	ErrDocumentTypeAlreadyExists = errors.New("document type already exists")

	// ErrDocumentTypeDoesntExists document type doesn't exist
	ErrDocumentTypeDoesntExists = errors.New("document type doesn't exist")

	// ErrEmptyField field must be filled
	ErrEmptyField = errors.New("field must be filled")

	// ErrEmptyIndex index must be filled
	ErrEmptyIndex = errors.New("index must be filled")

	// ErrEmptyType source type must be filled
	ErrEmptyType = errors.New("source type must set")

	// ErrEmptySource source name must be filled
	ErrEmptySource = errors.New("source name must set")

	// ErrEmptyRefTag referenced tag must be filled
	ErrEmptyRefTag = errors.New("referenced tag must set")

	// ErrConjunctionAndDisjunctionIsNil conjunction and disjunction are nil
	ErrConjunctionAndDisjunctionIsNil = errors.New("conjunction and disjunction are nil")

	// ErrEndAsTimeZero end as Time is zero instant
	ErrEndAsTimeZero = errors.New("endAsTime is zero instant")

	// ErrFirstParameterNotStruct first parameter is not a struct
	ErrFirstParameterNotStruct = errors.New("first parameter is not a struct")

	// ErrInputStructPointer input struct must be a pointer
	ErrInputStructPointer = errors.New("input struct must be pointer")

	// ErrInvalidBulkContainer bulk container type definition error
	ErrInvalidBulkContainer = errors.New("container must be *[]T, with length of ids array")

	// ErrInvalidGetDocumentTypesParam represents value for get document types should be pointer
	ErrInvalidGetDocumentTypesParam = errors.New("internal error: value should be pointer for getDocumentTypes")
)
