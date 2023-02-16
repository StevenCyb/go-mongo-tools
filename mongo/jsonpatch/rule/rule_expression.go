//nolint:ireturn
package rule

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
)

// ExpressionRule defines expression that needs to be matched by value.
// Anything is printed to string before e.g. int will be a number as string.
type ExpressionRule struct {
	Expression string
	Regex      regexp.Regexp
}

// NewInstance instantiate new rule instance for field.
func (e *ExpressionRule) NewInstance(path string, _ reflect.Kind, _ interface{}, value string) (Rule, error) {
	regex, err := getRegexpIfNotEmpty(value, path, "ExpressionRule")
	if err != nil {
		return nil, err
	}

	return &ExpressionRule{Expression: value, Regex: *regex}, nil
}

// NewInheritInstance instantiate new rule instance based on given rule.
func (e *ExpressionRule) NewInheritInstance(path string, _ reflect.Kind, _ interface{}) (Rule, error) {
	return &ExpressionRule{Expression: e.Expression, Regex: e.Regex}, nil
}

// Validate applies rule on given patch operation specification.
func (e ExpressionRule) Validate(operationSpec operation.Spec) error {
	value := fmt.Sprintf("%+v", operationSpec.Value)

	if !e.Regex.Match([]byte(value)) {
		return ExpressionNotMatchError{expression: e.Expression, value: value}
	}

	return nil
}
