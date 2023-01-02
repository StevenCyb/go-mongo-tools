package validator

import (
	"errors"
	"fmt"
)

var (
	ErrMissingPrefix     = errors.New("missing tag prefix '" + prefix + "'")
	ErrDuplicateRuleTags = errors.New("tag name already registered")
	ErrNilRule           = errors.New("rule is nil")
	ErrReferenceIsNil    = errors.New("reference is nil")
)

// InvalidTypeError indicate that a type is invalid.
type InvalidTypeError struct {
	path string
}

func (i InvalidTypeError) Error() string {
	return fmt.Sprintf("type at '%s' is invalid", i.path)
}

// UnknownRuleError indicate that requested rule is not known.
type UnknownRuleError struct {
	name string
}

func (u UnknownRuleError) Error() string {
	return fmt.Sprintf("unknown rule '%s'", u.name)
}

// InheritNonExistingTagError indicate that defined heredity does not exist.
type InheritNonExistingTagError struct {
	name string
}

func (i InheritNonExistingTagError) Error() string {
	return fmt.Sprintf("defined tag '%s' for heredity does not exist", i.name)
}

// UnknownPathError indicate that defined path does not exist.
type UnknownPathError struct {
	path string
}

func (u UnknownPathError) Error() string {
	return fmt.Sprintf("defined path '%s' is unknown", u.path)
}
