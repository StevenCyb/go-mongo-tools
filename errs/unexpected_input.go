package errs

import (
	"fmt"
)

const errUnexpectedInputMessage = "Unexpected input \"%v\""

// UnexpectedInputError is an error
// type for unexpected input.
type UnexpectedInputError struct {
	data interface{}
}

// Error returns the error message text.
func (err UnexpectedInputError) Error() string {
	return fmt.Sprintf(errUnexpectedInputMessage, err.data)
}

// NewErrUnexpectedInput cerate a new error.
func NewErrUnexpectedInput(data interface{}) UnexpectedInputError {
	return UnexpectedInputError{data: data}
}
