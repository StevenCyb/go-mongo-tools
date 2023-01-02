package rule

import (
	"reflect"
	"testing"

	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
	"github.com/stretchr/testify/require"
)

func TestRuleDisallowedOperation(t *testing.T) {
	t.Parallel()

	var rule Rule = &DisallowedOperationsRule{}
	rule, err := rule.NewInstance("a", reflect.Array, nil, "replace,move")

	require.NoError(t, err)

	t.Run("AllowedCase", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, rule.Validate(operation.Spec{
			Operation: operation.AddOperation,
		}))
	})

	t.Run("ForbiddenCase", func(t *testing.T) {
		t.Parallel()
		require.Error(t, rule.Validate(operation.Spec{
			Operation: operation.ReplaceOperation,
		}))
	})
}
