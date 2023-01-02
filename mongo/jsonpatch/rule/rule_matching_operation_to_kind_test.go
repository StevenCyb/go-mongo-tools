package rule

import (
	"reflect"
	"testing"

	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
	"github.com/stretchr/testify/require"
)

func TestMatchingOperationToKindRule(t *testing.T) {
	t.Parallel()

	t.Run("ArrayWithAddOperation_Success", func(t *testing.T) {
		t.Parallel()

		rule := MatchingOperationToKindRule{Kind: reflect.Array}
		require.NoError(t, rule.Validate(operation.Spec{Operation: operation.AddOperation}))
	})

	t.Run("StringWithAddOperation_Success", func(t *testing.T) {
		t.Parallel()

		rule := MatchingOperationToKindRule{Kind: reflect.String}
		require.Error(t, rule.Validate(operation.Spec{Operation: operation.AddOperation}))
	})

	t.Run("ArrayWithRemoveOperation_Success", func(t *testing.T) {
		t.Parallel()

		rule := MatchingOperationToKindRule{Kind: reflect.Array}
		require.NoError(t, rule.Validate(operation.Spec{Operation: operation.RemoveOperation}))
	})
}
