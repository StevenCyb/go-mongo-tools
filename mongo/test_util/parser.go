package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

type TextParser interface {
	Parse(query string) (bson.D, error)
}

func ExecuteSuccessTest(t *testing.T, parser TextParser, query string, expect bson.D) {
	t.Helper()

	actual, err := parser.Parse(query)
	require.NoError(t, err)
	require.Equal(t, expect, actual)
}

func ExecuteFailedTest(t *testing.T, parser TextParser, query string, expectedError error) {
	t.Helper()

	_, err := parser.Parse(query)
	require.Equal(t, expectedError, err)
}
