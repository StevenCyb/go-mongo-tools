package errs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrUnexpectedInputEnd(t *testing.T) {
	t.Parallel()

	key := "b"
	require.Equal(t,
		fmt.Sprintf(errUnexpectedInputEndMessage, key),
		NewErrUnexpectedInputEnd(key).Error(),
	)
}
