package errs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrPolicyViolation(t *testing.T) {
	t.Parallel()

	key := "a"
	require.Equal(t,
		fmt.Sprintf(errPolicyViolationMessage, key),
		NewErrPolicyViolation(key).Error(),
	)
}
