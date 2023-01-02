//nolint:funlen
package sort

import (
	"context"
	"testing"

	"github.com/StevenCyb/go-mongo-tools/errs"
	testutil "github.com/StevenCyb/go-mongo-tools/mongo/test_util"
	"github.com/StevenCyb/go-mongo-tools/tokenizer"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func TestParsing(t *testing.T) {
	t.Parallel()

	t.Run("Query", func(t *testing.T) {
		t.Parallel()

		t.Run("WithEmptyQuery_Success", func(t *testing.T) {
			t.Parallel()

			testutil.ExecuteSuccessTest(t,
				NewParser(nil),
				"",
				bson.D{},
			)
		})

		t.Run("WithSingleAscCriteria_Success", func(t *testing.T) {
			t.Parallel()

			testutil.ExecuteSuccessTest(t,
				NewParser(nil),
				"name=asc",
				bson.D{bson.E{Key: "name", Value: 1}},
			)
			testutil.ExecuteSuccessTest(t,
				NewParser(nil),
				"name=1",
				bson.D{bson.E{Key: "name", Value: 1}},
			)
		})

		t.Run("WithSingleDescCriteria_Success", func(t *testing.T) {
			t.Parallel()

			testutil.ExecuteSuccessTest(t,
				NewParser(nil),
				"name=desc",
				bson.D{bson.E{Key: "name", Value: -1}},
			)
			testutil.ExecuteSuccessTest(t,
				NewParser(nil),
				"name=-1",
				bson.D{bson.E{Key: "name", Value: -1}},
			)
		})

		t.Run("WithMultipleSortingCriteria_Success", func(t *testing.T) {
			t.Parallel()

			testutil.ExecuteSuccessTest(t,
				NewParser(nil),
				"firstName=asc,age=desc, something = asc",
				bson.D{
					bson.E{Key: "firstName", Value: 1},
					bson.E{Key: "age", Value: -1},
					bson.E{Key: "something ", Value: 1},
				},
			)
		})

		t.Run("WithUnknownLiteral_Fail", func(t *testing.T) {
			t.Parallel()

			testutil.ExecuteFailedTest(t,
				NewParser(nil),
				"firstName=?",
				errs.NewErrUnexpectedTokenType(11, "FIELD_NAME", "SORT_CRITERIA"),
			)
		})
		t.Run("WithUnexpectedSeparator_Fail", func(t *testing.T) {
			t.Parallel()

			testutil.ExecuteFailedTest(t,
				NewParser(nil),
				"firstName=asc+lastName=asc",
				errs.NewErrUnexpectedTokenType(22, "FIELD_NAME", ","),
			)
		})
	})

	t.Run("WithPolicy", func(t *testing.T) {
		t.Parallel()

		t.Run("WithAllowedFieldName_Success", func(t *testing.T) {
			t.Parallel()

			testutil.ExecuteSuccessTest(t,
				NewParser(
					tokenizer.NewPolicy(tokenizer.WhitelistPolicy, "a", "b")),
				"a=asc,b=asc",
				bson.D{bson.E{Key: "a", Value: 1}, bson.E{Key: "b", Value: 1}},
			)
		})

		t.Run("WithDisallowedFieldName_Fail", func(t *testing.T) {
			t.Parallel()

			testutil.ExecuteFailedTest(t,
				NewParser(
					tokenizer.NewPolicy(tokenizer.WhitelistPolicy, "a", "b")),
				"a=asc,b=desc,c=desc",
				errs.NewErrPolicyViolation("c"),
			)
		})
	})
}

func TestInterpretation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	server := testutil.NewStrikemongoServer(t)
	mongoClient, collection, database := testutil.NewClientWithCollection(t, server)

	//nolint:errcheck
	t.Cleanup(func() {
		server.Stop()
		mongoClient.Disconnect(ctx)
		database.Drop(ctx)
	})

	items := []testutil.DummyDoc{
		{FirstName: "Max", LastName: "Muster", Gender: "male", Age: 52},
		{FirstName: "Alexa", LastName: "Amaizon", Gender: "female", Age: 22},
		{FirstName: "Tina", LastName: "Someone", Gender: "female", Age: 33},
		{FirstName: "Samal", LastName: "Someone", Gender: "male", Age: 26},
	}

	itemsInterface := []interface{}{}
	for _, item := range items {
		itemsInterface = append(itemsInterface, item)
	}

	testutil.Populate(t, collection, itemsInterface)

	t.Run("SortByName_Success", func(t *testing.T) {
		t.Parallel()

		parser := NewParser(nil)
		sort, err := parser.Parse(`first_name=asc`)
		require.NoError(t, err)

		testutil.FindCompare(t, collection, nil, sort, items[1], items[0], items[3], items[2])
	})

	t.Run("SortByGenderAscAndAgeDesc_Success", func(t *testing.T) {
		t.Parallel()

		parser := NewParser(nil)
		sort, err := parser.Parse(`gender=asc,age=desc`)
		require.NoError(t, err)

		testutil.FindCompare(t, collection, nil, sort, items[2], items[1], items[0], items[3])
	})
}
