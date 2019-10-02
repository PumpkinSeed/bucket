package bucket

import "errors"

var (
	// Document type alredy exist
	ErrDocumentTypeAlredyExists = errors.New("document type alredy exists")

	// Document type doesn't exist
	ErrDocumentTypeDoesntExists = errors.New("document type doesn't exist")

	// Field must be filled
	ErrEmptyField = errors.New("field must be filled")

	// Index must be filled
	ErrEmptyIndex = errors.New("index must be filled")

	// Source type must be filled
	ErrEmptyType = errors.New("source type must set")

	// Source name must be filled
	ErrEmptySource = errors.New("source name must set")
	ErrEmptyRefTag = errors.New("referenced tag must set")

	// Conjunktion and disjunction are nil
	ErrConjunctionAndDisjunktionIsNil = errors.New("conjunction and disjunction are nil")

	// End as Time is zero instant
	ErrEndAsTimeZero = errors.New("endAsTime is zero instant")

	// First parameter is not a struct
	ErrFirstParameterNotStruct = errors.New("first parameter is not a struct")

	// Input struct must be a pointer
	ErrInputStructPointer = errors.New("input struct must be pointer")
)
