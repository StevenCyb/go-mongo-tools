//nolint:funlen
package tokenizer

import (
	"fmt"
	"testing"

	"github.com/StevenCyb/go-mongo-tools/errs"
	"github.com/stretchr/testify/require"
)

func TestTokenizer(t *testing.T) {
	t.Parallel()

	var (
		NoneType  Type = "NONE"
		EqualType Type = "EQUAL"
		WordType  Type = "WORD"
		key            = "hello"
		separator      = "="
		value          = "world"
	)

	t.Run("SpecParsing", func(t *testing.T) {
		t.Parallel()

		tokenizer := NewTokenizer(
			fmt.Sprintf("%s%s%s", key, separator, value),
			NoneType, NoneType,
			[]*Spec{
				NewSpec("^=", EqualType),
				NewSpec("^[a-z]+", WordType),
			},
			nil)

		token, err := tokenizer.GetNextToken()
		require.NoError(t, err)
		require.Equal(t, WordType, token.Type)
		require.Equal(t, key, token.Value)

		token, err = tokenizer.GetNextToken()
		require.NoError(t, err)
		require.Equal(t, EqualType, token.Type)
		require.Equal(t, separator, token.Value)

		token, err = tokenizer.GetNextToken()
		require.NoError(t, err)
		require.Equal(t, WordType, token.Type)
		require.Equal(t, value, token.Value)
	})

	t.Run("SpecParsingWithSkip", func(t *testing.T) {
		t.Parallel()

		var SkipType Type = "SKIP"
		tokenizer := NewTokenizer(
			`  hello  = world `,
			SkipType, NoneType,
			[]*Spec{
				NewSpec(`^\s+`, SkipType),
				NewSpec("^=", EqualType),
				NewSpec("^[a-z]+", WordType),
			},
			nil)

		token, err := tokenizer.GetNextToken()
		require.NoError(t, err)
		require.Equal(t, WordType, token.Type)
		require.Equal(t, key, token.Value)

		token, err = tokenizer.GetNextToken()
		require.NoError(t, err)
		require.Equal(t, EqualType, token.Type)
		require.Equal(t, separator, token.Value)

		token, err = tokenizer.GetNextToken()
		require.NoError(t, err)
		require.Equal(t, WordType, token.Type)
		require.Equal(t, value, token.Value)
	})

	t.Run("PolicyWithoutViolation", func(t *testing.T) {
		t.Parallel()

		t.Run("SpecParsing", func(t *testing.T) {
			t.Parallel()

			tokenizer := NewTokenizer(
				fmt.Sprintf("%s%s%s", key, separator, value),
				NoneType, WordType,
				[]*Spec{
					NewSpec("^=", EqualType),
					NewSpec("^[a-z]+", WordType),
				},
				NewPolicy(WhitelistPolicy, key, value))

			token, err := tokenizer.GetNextToken()
			require.NoError(t, err)
			require.Equal(t, WordType, token.Type)
			require.Equal(t, key, token.Value)

			token, err = tokenizer.GetNextToken()
			require.NoError(t, err)
			require.Equal(t, EqualType, token.Type)
			require.Equal(t, separator, token.Value)

			token, err = tokenizer.GetNextToken()
			require.NoError(t, err)
			require.Equal(t, WordType, token.Type)
			require.Equal(t, value, token.Value)
		})
	})

	t.Run("PolicyWithViolation", func(t *testing.T) {
		t.Parallel()

		t.Run("SpecParsing", func(t *testing.T) {
			t.Parallel()

			tokenizer := NewTokenizer(
				fmt.Sprintf("%s%s%s", key, separator, value),
				NoneType, WordType,
				[]*Spec{
					NewSpec("^=", EqualType),
					NewSpec("^[a-z]+", WordType),
				},
				NewPolicy(WhitelistPolicy, key))

			token, err := tokenizer.GetNextToken()
			require.NoError(t, err)
			require.Equal(t, WordType, token.Type)
			require.Equal(t, key, token.Value)

			token, err = tokenizer.GetNextToken()
			require.NoError(t, err)
			require.Equal(t, EqualType, token.Type)
			require.Equal(t, separator, token.Value)

			_, err = tokenizer.GetNextToken()
			require.Equal(t, errs.NewErrPolicyViolation(value), err)
		})
	})
}
