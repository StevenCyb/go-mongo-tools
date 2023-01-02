package operation

import (
	"errors"
	"fmt"
	"strings"
)

var ErrUnknownOperation = errors.New("unknown operation")

// Operation represents an JSON patch operation.
type Operation string

const (
	// RemoveOperation is an operation to remove the value at the target location.
	// Requires `path`.
	RemoveOperation Operation = "remove"
	// AddOperation add a value or array to an array at the target location.
	// Requires `path` and `value`. Target of the path must be an array.
	AddOperation Operation = "add"
	// ReplaceOperation replaces the value at the target location
	// with a new value. Requires `path` and `value`.
	ReplaceOperation Operation = "replace"
	// MoveOperation removes the value at a specified location and
	// adds it to the target location. Requires `from` and `path`.
	MoveOperation Operation = "move"
	// CopyOperation copies the value from a specified location to the
	// target location. Requires `from` and `path`.
	CopyOperation Operation = "copy"
)

// FromString return operation that matches string or nil with error if no match.
func FromString(operationString string) (*Operation, error) {
	operation := Operation(strings.ToLower(operationString))

	if operation != RemoveOperation &&
		operation != AddOperation &&
		operation != ReplaceOperation &&
		operation != MoveOperation &&
		operation != CopyOperation {
		return nil, fmt.Errorf("%w: %s", ErrUnknownOperation, operationString)
	}

	return &operation, nil
}

// Spec specify an path operation.
type Spec struct {
	From      Path        `json:"from"`
	Path      Path        `json:"path"`
	Value     interface{} `json:"value"`
	Operation Operation   `json:"op"` //nolint:tagliatelle
}

// Valid check if operation is valid.
func (s Spec) Valid() bool {
	if s.Operation == "" {
		return false
	}

	switch s.Operation {
	case RemoveOperation:
		if !s.Path.Valid() {
			return false
		}
	case AddOperation:
		if !s.Path.Valid() {
			return false
		} else if s.Value == nil {
			return false
		}
	case ReplaceOperation:
		if !s.Path.Valid() {
			return false
		} else if s.Value == nil {
			return false
		}
	case MoveOperation:
		if !s.Path.Valid() {
			return false
		}

		if !s.From.Valid() {
			return false
		}
	case CopyOperation:
		if !s.Path.Valid() {
			return false
		}

		if !s.From.Valid() {
			return false
		}
	default:
		return false
	}

	return true
}
