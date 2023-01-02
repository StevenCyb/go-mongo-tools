package operation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValidPath(t *testing.T) {
	t.Parallel()

	t.Run("ShortPath_Success", func(t *testing.T) {
		t.Parallel()

		path := Path("root")
		require.True(t, path.Valid())
	})

	t.Run("LongPath_Success", func(t *testing.T) {
		t.Parallel()

		path := Path("parent.child")
		require.True(t, path.Valid())
	})

	t.Run("InvalidPath_Fail", func(t *testing.T) {
		t.Parallel()

		path := Path("now..valid")
		require.False(t, path.Valid())

		path = Path(".fds.")
		require.False(t, path.Valid())

		path = Path("")
		require.False(t, path.Valid())

		path = Path(".##")
		require.False(t, path.Valid())
	})
}

//nolint:gocritic
func TestComparePath(t *testing.T) {
	t.Parallel()

	require.True(t, Path("").Equal(Path("")))
	require.True(t, Path("a.b.c").Equal(Path("a.b.c")))
	require.True(t, Path("*.b.c").Equal(Path("a.b.c")))
	require.True(t, Path("a.*.c").Equal(Path("a.b.c")))
	require.True(t, Path("a.b.*").Equal(Path("a.b.c")))
	require.True(t, Path("*.*").Equal(Path("a.b")))
	require.True(t, Path("a.*").Equal(Path("a.0")))
	require.False(t, Path("a.*.b").Equal(Path("a.b.c")))
	require.False(t, Path("*.*").Equal(Path("x")))
	require.False(t, Path("*.*").Equal(Path("x.y.z")))
}
