package errs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorChain(t *testing.T) {
	t.Parallel()

	t.Run("WithError", func(t *testing.T) {
		t.Parallel()

		errChain := Chain{}

		errChain.AddIf(nil)
		errChain.AddIf(errors.New("test")) //nolint:goerr113
		errChain.AddIf(nil)

		require.Error(t, errChain.GetError())
	})

	t.Run("WithoutError", func(t *testing.T) {
		t.Parallel()

		errChain := Chain{}

		errChain.AddIf(nil)
		errChain.AddIf(nil)
		errChain.AddIf(nil)

		require.NoError(t, errChain.GetError())
	})
}
