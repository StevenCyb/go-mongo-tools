package errs

import (
	"fmt"
)

const errPolicyViolationMessage = "Policy violation, policy disallow \"%s\""

// PolicyViolationError is an error
// type for policy violation.
type PolicyViolationError struct {
	key string
}

// Error returns the error message text.
func (err PolicyViolationError) Error() string {
	return fmt.Sprintf(errPolicyViolationMessage, err.key)
}

// NewErrUnexpectedInputEnd cerate a new error.
func NewErrPolicyViolation(key string) PolicyViolationError {
	return PolicyViolationError{key: key}
}
