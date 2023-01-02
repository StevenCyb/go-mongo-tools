//nolint:ireturn
package rule

import (
	"errors"
	"reflect"

	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
)

var ErrAddOperationTypeError = errors.New("add operation only applicable to array")

// MatchingOperationToKindRule is a default rule that is applied to all fields.
// This rules if operation is applyable to kind.
type MatchingOperationToKindRule struct {
	Kind reflect.Kind
}

// UseValue instantiate new rule instance for field.
func (m *MatchingOperationToKindRule) NewInstance(_ string, kind reflect.Kind, _ interface{}, _ string) (Rule, error) {
	return &MatchingOperationToKindRule{
		Kind: kind,
	}, nil
}

// NewInheritInstance instantiate new rule instance based on given rule.
func (m *MatchingOperationToKindRule) NewInheritInstance(_ string, kind reflect.Kind, _ interface{}) (Rule, error) {
	return &MatchingOperationToKindRule{
		Kind: kind,
	}, nil
}

// Validate applies rule on given patch operation specification.
func (m MatchingOperationToKindRule) Validate(operationSpec operation.Spec) error {
	if operationSpec.Operation == operation.AddOperation {
		if m.Kind != reflect.Array && m.Kind != reflect.Slice {
			return ErrAddOperationTypeError
		}
	}

	return nil
}
