//nolint:dupl
package operation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOperationFromString(t *testing.T) {
	t.Parallel()

	t.Run("RemoveOperationString", func(t *testing.T) {
		t.Parallel()

		operation, err := FromString("reMOve")
		require.NoError(t, err)
		require.NotEmpty(t, operation)
		require.Equal(t, RemoveOperation, *operation)
	})

	t.Run("InvalidOperationString", func(t *testing.T) {
		t.Parallel()

		operation, err := FromString("???")
		require.Error(t, err)
		require.Empty(t, operation)
	})
}

func TestRemoveOperationOperationValidation(t *testing.T) {
	t.Parallel()

	operation := RemoveOperation
	path := Path("a")
	invalidPath := Path(".")

	t.Run("Valid_Success", func(t *testing.T) {
		t.Parallel()
		require.True(t, Spec{Operation: operation, Path: path}.Valid())
	})

	t.Run("WithMissingPathFail", func(t *testing.T) {
		t.Parallel()
		require.False(t, Spec{Operation: operation}.Valid())
	})

	t.Run("WithInvalidPathFail", func(t *testing.T) {
		t.Parallel()
		require.False(t, Spec{Operation: operation, Path: invalidPath}.Valid())
	})
}

func TestAddOperationOperationValidation(t *testing.T) {
	t.Parallel()

	operation := AddOperation
	path := Path("a")
	invalidPath := Path(".")
	value := 1

	t.Run("Valid_Success", func(t *testing.T) {
		t.Parallel()
		require.True(t, Spec{Operation: operation, Path: path, Value: value}.Valid())
	})

	t.Run("WithMissingValueFail", func(t *testing.T) {
		t.Parallel()
		require.False(t, Spec{Operation: operation, Path: path}.Valid())
	})

	t.Run("WithMissingPathFail", func(t *testing.T) {
		t.Parallel()
		require.False(t, Spec{Operation: operation, Value: value}.Valid())
	})

	t.Run("WithInvalidPathFail", func(t *testing.T) {
		t.Parallel()
		require.False(t, Spec{Operation: operation, Path: invalidPath, Value: value}.Valid())
	})
}

func TestReplaceOperationOperationValidation(t *testing.T) {
	t.Parallel()

	operation := ReplaceOperation
	path := Path("a")
	invalidPath := Path(".")
	value := 1

	t.Run("Valid_Success", func(t *testing.T) {
		t.Parallel()
		require.True(t, Spec{Operation: operation, Path: path, Value: value}.Valid())
	})

	t.Run("WithMissingValueFail", func(t *testing.T) {
		t.Parallel()
		require.False(t, Spec{Operation: operation, Path: path}.Valid())
	})

	t.Run("WithMissingPathFail", func(t *testing.T) {
		t.Parallel()
		require.False(t, Spec{Operation: operation, Value: value}.Valid())
	})

	t.Run("WithInvalidPathFail", func(t *testing.T) {
		t.Parallel()
		require.False(t, Spec{Operation: operation, Path: invalidPath, Value: value}.Valid())
	})
}

func TestMoveOperationOperationValidation(t *testing.T) {
	t.Parallel()

	operation := MoveOperation
	pathA := Path("a")
	pathB := Path("b")
	invalidPath := Path(".")

	t.Run("Valid_Success", func(t *testing.T) {
		t.Parallel()
		require.True(t, Spec{Operation: operation, Path: pathB, From: pathA}.Valid())
	})

	t.Run("WithMissingPathFail", func(t *testing.T) {
		t.Parallel()
		require.False(t, Spec{Operation: operation, From: pathA}.Valid())
	})

	t.Run("WithInvalidPathFail", func(t *testing.T) {
		t.Parallel()
		require.False(t, Spec{Operation: operation, Path: invalidPath, From: pathA}.Valid())
	})

	t.Run("WithMissingFromFail", func(t *testing.T) {
		t.Parallel()
		require.False(t, Spec{Operation: operation, Path: pathB}.Valid())
	})

	t.Run("WithInvalidFromFail", func(t *testing.T) {
		t.Parallel()
		require.False(t, Spec{Operation: operation, Path: pathB, From: invalidPath}.Valid())
	})
}

func TestCopyOperationOperationValidation(t *testing.T) {
	t.Parallel()

	operation := CopyOperation
	pathA := Path("a")
	pathB := Path("b")
	invalidPath := Path(".")

	t.Run("Valid_Success", func(t *testing.T) {
		t.Parallel()
		require.True(t, Spec{Operation: operation, Path: pathB, From: pathA}.Valid())
	})

	t.Run("WithMissingPathFail", func(t *testing.T) {
		t.Parallel()
		require.False(t, Spec{Operation: operation, From: pathA}.Valid())
	})

	t.Run("WithInvalidPathFail", func(t *testing.T) {
		t.Parallel()
		require.False(t, Spec{Operation: operation, Path: invalidPath, From: pathA}.Valid())
	})

	t.Run("WithMissingFromFail", func(t *testing.T) {
		t.Parallel()
		require.False(t, Spec{Operation: operation, Path: pathB}.Valid())
	})

	t.Run("WithInvalidFromFail", func(t *testing.T) {
		t.Parallel()
		require.False(t, Spec{Operation: operation, Path: pathB, From: invalidPath}.Valid())
	})
}
