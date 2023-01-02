package rule

import (
	"reflect"
	"testing"

	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
	"github.com/stretchr/testify/require"
)

func TestLessThenError(t *testing.T) {
	t.Parallel()

	require.Equal(t, "value is less then specified: '2.000000' < '3.000000'",
		LessThenError{ref: 3, value: 2}.Error())
}

func TestGreaterThenError(t *testing.T) {
	t.Parallel()

	require.Equal(t, "value is greater then specified: '3.000000' > '2.000000'",
		GreaterThenError{ref: 2, value: 3}.Error())
}

func TestOperationNotAllowedError(t *testing.T) {
	t.Parallel()

	require.Equal(t, "operation 'add' not allowed",
		OperationNotAllowedError{operation: operation.AddOperation}.Error())
}

func TestUnknownFieldError(t *testing.T) {
	t.Parallel()

	require.Equal(t, "unknown field 'test'",
		UnknownFieldError{name: "test"}.Error())
}

func TestTypeMismatchError(t *testing.T) {
	t.Parallel()

	require.Equal(t, "'test' has invalid kind 'int', must be 'string'",
		TypeMismatchError{name: "test", expected: reflect.String, actual: reflect.Int}.Error())
	require.Equal(t, "'test' key has invalid kind 'int', must be 'string'",
		TypeMismatchError{name: "test", expected: reflect.String, actual: reflect.Int, forKey: true}.Error())
}

func TestExpressionNotMatchError(t *testing.T) {
	t.Parallel()

	require.Equal(t, `expression '^\d$' not match a`,
		ExpressionNotMatchError{expression: `^\d$`, value: "a"}.Error())
}
