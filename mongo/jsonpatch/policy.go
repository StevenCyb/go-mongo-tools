package jsonpatch

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
)

// Policy specifies the interface for an policy.
type Policy interface {
	GetDetails() string
	Test(operationSpec operation.Spec) bool
}

// DisallowPathPolicy specifies a path that is not allowed.
type DisallowPathPolicy struct {
	Details string
	Path    operation.Path
}

// GetDetails returns the name of this policy.
func (d DisallowPathPolicy) GetDetails() string {
	return d.Details
}

// Test if given operation specification is valid or not.
func (d DisallowPathPolicy) Test(operationSpec operation.Spec) bool {
	return !d.Path.Equal(operationSpec.Path)
}

// DisallowOperationOnPathPolicy disallows specified operations on path.
type DisallowOperationOnPathPolicy struct {
	Details    string
	Path       operation.Path
	Operations []operation.Operation
}

// GetDetails returns the name of this policy.
func (d DisallowOperationOnPathPolicy) GetDetails() string {
	return d.Details
}

// Test if given operation specification is valid or not.
func (d DisallowOperationOnPathPolicy) Test(operationSpec operation.Spec) bool {
	if !d.Path.Equal(operationSpec.Path) {
		return true
	}

	for _, operation := range d.Operations {
		if operation == operationSpec.Operation {
			return false
		}
	}

	return true
}

// ForceTypeOnPathPolicy forces the value of a specif path to be from given type.
type ForceTypeOnPathPolicy struct {
	Details string
	Path    operation.Path
	Kind    reflect.Kind
}

// GetDetails returns the name of this policy.
func (f ForceTypeOnPathPolicy) GetDetails() string {
	return f.Details
}

// Test if given operation specification is valid or not.
func (f ForceTypeOnPathPolicy) Test(operationSpec operation.Spec) bool {
	if !f.Path.Equal(operationSpec.Path) {
		return true
	}

	return reflect.TypeOf(operationSpec.Value).Kind() == f.Kind
}

// ForceRegexMatchPolicy forces the value of a specif path to match expression.
type ForceRegexMatchPolicy struct {
	Details    string
	Path       operation.Path
	Expression regexp.Regexp
}

// GetDetails returns the name of this policy.
func (f ForceRegexMatchPolicy) GetDetails() string {
	return f.Details
}

// Test if given operation specification is valid or not.
func (f ForceRegexMatchPolicy) Test(operationSpec operation.Spec) bool {
	if !f.Path.Equal(operationSpec.Path) {
		return true
	}

	return f.Expression.MatchString(fmt.Sprintf("%+v", operationSpec.Value))
}

// StrictPathPolicy forces path to be strictly one of.
type StrictPathPolicy struct {
	Details string
	Path    []operation.Path
}

// GetDetails returns the name of this policy.
func (s StrictPathPolicy) GetDetails() string {
	return s.Details
}

// Test if given operation specification is valid or not.
func (s StrictPathPolicy) Test(operationSpec operation.Spec) bool {
	for _, path := range s.Path {
		if path.Equal(operationSpec.Path) {
			return true
		}
	}

	return false
}

// ForceOperationOnPathPolicy force specified operation on path.
type ForceOperationOnPathPolicy struct {
	Details   string
	Path      operation.Path
	Operation operation.Operation
}

// GetDetails returns the name of this policy.
func (d ForceOperationOnPathPolicy) GetDetails() string {
	return d.Details
}

// Test if given operation specification is valid or not.
func (d ForceOperationOnPathPolicy) Test(operationSpec operation.Spec) bool {
	if !d.Path.Equal(operationSpec.Path) {
		return true
	}

	return d.Operation == operationSpec.Operation
}
