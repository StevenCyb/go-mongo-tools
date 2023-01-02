package rule

import (
	"reflect"
	"testing"

	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
	"github.com/stretchr/testify/require"
)

func TestRuleExpression(t *testing.T) {
	t.Parallel()

	var rule Rule = &ExpressionRule{}
	rule, err := rule.NewInstance("a", reflect.Array, nil, "^[a-z]+$")

	require.NoError(t, err)

	t.Run("AllowedCase", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, rule.Validate(operation.Spec{
			Value: "hello",
		}))
	})

	t.Run("ForbiddenCase", func(t *testing.T) {
		t.Parallel()
		require.Error(t, rule.Validate(operation.Spec{
			Value: "123",
		}))
	})
}
