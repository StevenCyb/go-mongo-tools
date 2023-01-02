package errs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrUnexpectedInput(t *testing.T) {
	t.Parallel()

	key := "c"
	require.Equal(t,
		fmt.Sprintf(errUnexpectedInputMessage, key),
		NewErrUnexpectedInput(key).Error(),
	)
}
