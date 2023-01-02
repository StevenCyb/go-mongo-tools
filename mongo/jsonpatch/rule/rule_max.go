//nolint:ireturn
package rule

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
)

var ErrMaxRuleViolation = errors.New("value greater then specified")

// MaxRule defines the maximum size/value:
/*
 * Applies on len for Array, Chan, Map, Slice and String.
 * Applies on value for any numeric type.
 */
type MaxRule struct {
	Max float64
}

// NewInstance instantiate new rule instance for field.
func (m *MaxRule) NewInstance(path string, _ reflect.Kind, _ interface{}, value string) (Rule, error) {
	max, err := getFloat64IfNotEmpty(value, path, "MaxRule")
	if err != nil {
		return nil, err
	}

	return &MaxRule{
		Max: *max,
	}, nil
}

// NewInheritInstance instantiate new rule instance based on given rule.
func (m *MaxRule) NewInheritInstance(_ string, _ reflect.Kind, _ interface{}) (Rule, error) {
	return &MaxRule{
		Max: m.Max,
	}, nil
}

// Validate applies rule on given patch operation specification.
func (m MaxRule) Validate(operationSpec operation.Spec) error {
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

	if m.Max < size {
		return GreaterThenError{ref: m.Max, value: size}
	}

	return nil
}
