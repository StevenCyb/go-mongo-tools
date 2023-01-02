//nolint:ireturn
package rule

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
)

// MinRule defines the minimum size/value:
/*
 * Applies on len for Array, Chan, Map, Slice and String.
 * Applies on value for any numeric type.
 */
type MinRule struct {
	Min float64
}

// NewInstance instantiate new rule instance for field.
func (m *MinRule) NewInstance(path string, _ reflect.Kind, _ interface{}, value string) (Rule, error) {
	min, err := getFloat64IfNotEmpty(value, path, "MinRule")
	if err != nil {
		return nil, err
	}

	return &MinRule{
		Min: *min,
	}, nil
}

// NewInheritInstance instantiate new rule instance based on given rule.
func (m *MinRule) NewInheritInstance(_ string, _ reflect.Kind, _ interface{}) (Rule, error) {
	return &MinRule{
		Min: m.Min,
	}, nil
}

// Validate applies rule on given patch operation specification.
func (m MinRule) Validate(operationSpec operation.Spec) error {
	var (
		kind = reflect.TypeOf(operationSpec.Value).Kind()
		size float64
	)

	switch kind { //nolint:exhaustive
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		size = float64(reflect.ValueOf(operationSpec.Value).Len())
	default:
		converted, err := strconv.ParseFloat(fmt.Sprint(operationSpec.Value), 64)
		if err != nil {
			return fmt.Errorf("converting failed: %w", err)
		}

		size = converted
	}

	if m.Min > size {
		return LessThenError{ref: m.Min, value: size}
	}

	return nil
}
