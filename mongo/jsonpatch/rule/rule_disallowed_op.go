//nolint:ireturn
package rule

import (
	"reflect"

	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
)

// DisallowedOperationsRule defines witch patch operations are disallowed for the field.
type DisallowedOperationsRule struct {
	Operations []operation.Operation
}

// NewInstance instantiate new rule instance for field.
func (d *DisallowedOperationsRule) NewInstance(
	path string, _ reflect.Kind, _ interface{}, value string,
) (Rule, error) {
	operations, err := getOperationsIfNotEmpty(value, path, "DisallowedOperationsRule")
	if err != nil {
		return nil, err
	}

	return &DisallowedOperationsRule{Operations: *operations}, nil
}

// NewInheritInstance instantiate new rule instance based on given rule.
func (d *DisallowedOperationsRule) NewInheritInstance(_ string, _ reflect.Kind, _ interface{}) (Rule, error) {
	return &DisallowedOperationsRule{Operations: d.Operations}, nil
}

// Validate applies rule on given patch operation specification.
func (d DisallowedOperationsRule) Validate(operationSpec operation.Spec) error {
	for _, operation := range d.Operations {
		if operation == operationSpec.Operation {
			return OperationNotAllowedError{operation: operation}
		}
	}

	return nil
}
