package rule

import (
	"regexp"
	"testing"

	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
	"github.com/stretchr/testify/require"
)

func TestGetBoolIfNotEmpty(t *testing.T) {
	t.Parallel()

	t.Run("AllowedCase", func(t *testing.T) {
		t.Parallel()

		value, err := getBoolIfNotEmpty("true", "a", "x")
		require.NoError(t, err)
		require.NotNil(t, value)
		require.True(t, *value)
	})

	t.Run("ForbiddenCase", func(t *testing.T) {
		t.Parallel()

		value, err := getBoolIfNotEmpty("??", "a", "x")
		require.Error(t, err)
		require.Nil(t, value)
	})
}

func TestGetFloat64IfNotEmpty(t *testing.T) {
	t.Parallel()

	t.Run("AllowedCase", func(t *testing.T) {
		t.Parallel()

		value, err := getFloat64IfNotEmpty("1.3", "a", "x")
		require.NoError(t, err)
		require.NotNil(t, value)
		require.Equal(t, float64(1.3), *value)
	})

	t.Run("ForbiddenCase", func(t *testing.T) {
		t.Parallel()

		value, err := getFloat64IfNotEmpty("??", "a", "x")
		require.Error(t, err)
		require.Nil(t, value)
	})
}

func TestGetRegexpIfNotEmpty(t *testing.T) {
	t.Parallel()

	t.Run("AllowedCase", func(t *testing.T) {
		t.Parallel()

		value, err := getRegexpIfNotEmpty("/d+", "a", "x")
		require.NoError(t, err)
		require.NotNil(t, value)
		require.Equal(t, *regexp.MustCompile("/d+"), *value)
	})

	t.Run("ForbiddenCase", func(t *testing.T) {
		t.Parallel()

		value, err := getRegexpIfNotEmpty("?]?", "a", "x")
		require.Error(t, err)
		require.Nil(t, value)
	})
}

func TestGetOperationsIfNotEmpty(t *testing.T) {
	t.Parallel()

	t.Run("AllowedCase", func(t *testing.T) {
		t.Parallel()

		value, err := getOperationsIfNotEmpty("add,remove", "a", "x")
		require.NoError(t, err)
		require.NotNil(t, value)
		require.Equal(t, []operation.Operation{operation.AddOperation, operation.RemoveOperation}, *value)
	})

	t.Run("ForbiddenCase", func(t *testing.T) {
		t.Parallel()

		value, err := getOperationsIfNotEmpty("add,", "a", "x")
		require.Error(t, err)
		require.Nil(t, value)
	})
}
