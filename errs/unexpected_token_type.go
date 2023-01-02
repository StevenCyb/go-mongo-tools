package errs

import (
	"fmt"
)

const errUnexpectedTokenTypeMessage = "Unexpected token: \"%s\" at position \"%d\", expected: \"%s\""

// UnexpectedTokenTypeError.TokenType is an error
// type for unexpected token type.
type UnexpectedTokenTypeError struct {
	actual   string
	expected string
	position int
}

// Error returns the error message text.
func (err UnexpectedTokenTypeError) Error() string {
	return fmt.Sprintf(errUnexpectedTokenTypeMessage,
		err.actual, err.position, err.expected)
}

// NewErrUnexpectedTokenType cerate a new error.
func NewErrUnexpectedTokenType(position int, actual, expected string) UnexpectedTokenTypeError {
	return UnexpectedTokenTypeError{
		position: position,
		actual:   actual,
		expected: expected,
	}
}
