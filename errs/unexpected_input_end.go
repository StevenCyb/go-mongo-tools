package errs

import (
	"fmt"
)

const errUnexpectedInputEndMessage = "Unexpected end of input, expected: \"%s\""

// UnexpectedInputEndError is an error
// type for unexpected input end.
type UnexpectedInputEndError struct {
	tokenType string
}

// Error returns the error message text.
func (err UnexpectedInputEndError) Error() string {
	return fmt.Sprintf(errUnexpectedInputEndMessage, err.tokenType)
}

// NewErrUnexpectedInputEnd cerate a new error.
func NewErrUnexpectedInputEnd(tokenType string) UnexpectedInputEndError {
	return UnexpectedInputEndError{tokenType: tokenType}
}
