package rule

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
)

var ErrOperationsNotAllowed = errors.New("operations are not allowed")

// LessThenError indicate that a value is less then the reference.
type LessThenError struct {
	ref   float64
	value float64
}

func (l LessThenError) Error() string {
	return fmt.Sprintf("value is less then specified: '%f' < '%f'", l.value, l.ref)
}

// GreaterThenError indicate that a value is greater then the reference.
type GreaterThenError struct {
	ref   float64
	value float64
}

func (g GreaterThenError) Error() string {
	return fmt.Sprintf("value is greater then specified: '%f' > '%f'", g.value, g.ref)
}

// OperationNotAllowedError indicate that a given JSON patch operation is not allowed.
type OperationNotAllowedError struct {
	operation operation.Operation
}

func (o OperationNotAllowedError) Error() string {
	return fmt.Sprintf("operation '%s' not allowed", o.operation)
}

// UnknownFieldError indicate that a field is not known.
type UnknownFieldError struct {
	name string
}

func (u UnknownFieldError) Error() string {
	return fmt.Sprintf("unknown field '%s'", u.name)
}

// TypeMismatchError indicate that a given type not match a reference.
type TypeMismatchError struct {
	name     string
	expected reflect.Kind
	actual   reflect.Kind
	forKey   bool
}

func (t TypeMismatchError) Error() string {
	if t.forKey {
		return fmt.Sprintf("'%s' key has invalid kind '%s', must be '%s'", t.name, t.actual.String(), t.expected.String())
	}

	return fmt.Sprintf("'%s' has invalid kind '%s', must be '%s'", t.name, t.actual.String(), t.expected.String())
}

// ExpressionNotMatchError indicate that given value not match expression.
type ExpressionNotMatchError struct {
	expression string
	value      string
}

func (e ExpressionNotMatchError) Error() string {
	return fmt.Sprintf("expression '%s' not match %s", e.expression, e.value)
}
