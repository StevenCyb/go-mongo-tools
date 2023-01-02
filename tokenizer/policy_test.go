package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWhitelist(t *testing.T) {
	t.Parallel()

	policy := NewPolicy(WhitelistPolicy, "a", "b")
	require.True(t, policy.Allow("a"))
	require.True(t, policy.Allow("b"))
	require.False(t, policy.Allow("c"))
}

func TestBlacklist(t *testing.T) {
	t.Parallel()

	policy := NewPolicy(BlacklistPolicy, "a", "b")
	require.False(t, policy.Allow("a"))
	require.False(t, policy.Allow("b"))
	require.True(t, policy.Allow("c"))
}
