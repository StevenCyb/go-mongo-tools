//nolint:ireturn
package rule

import (
	"reflect"

	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
)

// AllowedOperationsRule defines witch patch operations are allowed for the field.
type AllowedOperationsRule struct {
	Operations []operation.Operation
}

// NewInstance instantiate new rule instance for field.
func (a *AllowedOperationsRule) NewInstance(
	path string, _ reflect.Kind, _ interface{}, value string,
) (Rule, error) {
	operations, err := getOperationsIfNotEmpty(value, path, "AllowedOperationsRule")
	if err != nil {
		return nil, err
	}

	return &AllowedOperationsRule{Operations: *operations}, nil
}

// NewInheritInstance instantiate new rule instance based on given rule.
func (a *AllowedOperationsRule) NewInheritInstance(_ string, _ reflect.Kind, _ interface{}) (Rule, error) {
	return &AllowedOperationsRule{Operations: a.Operations}, nil
}

// Validate applies rule on given patch operation specification.
func (a AllowedOperationsRule) Validate(operationSpec operation.Spec) error {
	for _, operation := range a.Operations {
		if operation == operationSpec.Operation {
			return nil
		}
	}

	return OperationNotAllowedError{operation: operationSpec.Operation}
}
