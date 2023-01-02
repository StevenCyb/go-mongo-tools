package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvalidTypeError(t *testing.T) {
	t.Parallel()

	require.Equal(t, "type at 'a.b' is invalid",
		InvalidTypeError{path: "a.b"}.Error())
}

func TestUnknownRuleError(t *testing.T) {
	t.Parallel()

	require.Equal(t, "unknown rule 'test'",
		UnknownRuleError{name: "test"}.Error())
}

func TestInheritNonExistingTagError(t *testing.T) {
	t.Parallel()

	require.Equal(t, "defined tag 'test' for heredity does not exist",
		InheritNonExistingTagError{name: "test"}.Error())
}

func TestUnknownPathError(t *testing.T) {
	t.Parallel()

	require.Equal(t, "defined path 'test' is unknown",
		UnknownPathError{path: "test"}.Error())
}
