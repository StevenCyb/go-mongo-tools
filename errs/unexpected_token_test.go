package errs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrUnexpectedInput2(t *testing.T) {
	t.Parallel()

	pos := 42
	key := "d"
	require.Equal(t,
		fmt.Sprintf(errUnexpectedTokenMessage, key, pos),
		NewErrUnexpectedToken(pos, key).Error(),
	)
}
