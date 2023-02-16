//nolint:ireturn
package rule

import (
	"reflect"

	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
)

// DisallowRule defines if operations on field are disallowed.
type DisallowRule struct {
	Disallow bool
}

// NewInstance instantiate new rule instance for field.
func (d *DisallowRule) NewInstance(path string, _ reflect.Kind, _ interface{}, value string) (Rule, error) {
	disallow, err := getBoolIfNotEmpty(value, path, "DisallowRule")
	if err != nil {
		return nil, err
	}

	return &DisallowRule{Disallow: *disallow}, nil
}

// NewInheritInstance instantiate new rule instance based on given rule.
func (d *DisallowRule) NewInheritInstance(_ string, _ reflect.Kind, _ interface{}) (Rule, error) {
	return &DisallowRule{Disallow: d.Disallow}, nil
}

// Validate applies rule on given patch operation specification.
func (d DisallowRule) Validate(operationSpec operation.Spec) error {
	if !d.Disallow {
		return nil
	}

	return ErrOperationsNotAllowed
}
