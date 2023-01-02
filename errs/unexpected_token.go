package errs

import "fmt"

const errUnexpectedTokenMessage = "Unexpected token: \"%s\" at position \"%d\""

// UnexpectedTokenError is an error
// type for unexpected token.
type UnexpectedTokenError struct {
	token    string
	position int
}

// Error returns the error message text.
func (err UnexpectedTokenError) Error() string {
	return fmt.Sprintf(errUnexpectedTokenMessage,
		err.token,
		err.position)
}

// NewErrUnexpectedToken cerate a new error.
func NewErrUnexpectedToken(position int, token string) UnexpectedTokenError {
	return UnexpectedTokenError{
		position: position,
		token:    token,
	}
}
