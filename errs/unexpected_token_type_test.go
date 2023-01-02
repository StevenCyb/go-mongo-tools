package errs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrUnexpectedTokenTypeMessage(t *testing.T) {
	t.Parallel()

	pos := 42
	key1 := "abc"
	key2 := "abc"
	require.Equal(t,
		fmt.Sprintf(errUnexpectedTokenTypeMessage, key1, pos, key2),
		NewErrUnexpectedTokenType(pos, key1, key2).Error(),
	)
}
