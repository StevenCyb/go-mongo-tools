package rule

import (
	"reflect"
	"testing"

	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
	"github.com/stretchr/testify/require"
)

func TestRuleDisallow(t *testing.T) {
	t.Parallel()

	t.Run("AllowedCase", func(t *testing.T) {
		t.Parallel()

		var rule Rule = &DisallowRule{}
		rule, err := rule.NewInstance("a", reflect.Bool, nil, "false")

		require.NoError(t, err)
		require.NoError(t, rule.Validate(operation.Spec{
			Path: "a.c",
		}))
	})

	t.Run("ForbiddenCase", func(t *testing.T) {
		t.Parallel()

		var rule Rule = &DisallowRule{}
		rule, err := rule.NewInstance("a", reflect.Bool, nil, "true")

		require.NoError(t, err)
		require.Error(t, rule.Validate(operation.Spec{
			Path: "a.b",
		}))
	})
}
