//nolint:funlen
package rsql

import (
	"context"
	"regexp"
	"runtime"
	"testing"

	"github.com/StevenCyb/go-mongo-tools/errs"
	testutil "github.com/StevenCyb/go-mongo-tools/mongo/test_util"
	"github.com/StevenCyb/go-mongo-tools/tokenizer"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestQueryParsingWithEmptyQuery_Success(t *testing.T) {
	t.Parallel()

	testutil.ExecuteSuccessTest(t,
		NewParser(nil),
		"",
		bson.D{},
	)
}

func TestQueryParsingWithInvalidQuery_Fail(t *testing.T) {
	t.Parallel()

	testutil.ExecuteFailedTest(t,
		NewParser(nil),
		"not_gonna_work",
		errs.NewErrUnexpectedInputEnd(FieldNameType.String()),
	)
}

func TestQueryParsingWithDifferentLiterals(t *testing.T) {
	t.Parallel()

	t.Run("==STRING_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`firstName=="steven"`,
			bson.D{bson.E{Key: "firstName", Value: "steven"}},
		)
	})

	t.Run("==OID_Success", func(t *testing.T) {
		t.Parallel()

		oid, err := primitive.ObjectIDFromHex("01234567890abcdef1234567")
		require.NoError(t, err)
		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`_id==$oid(01234567890abcdef1234567)`,
			bson.D{bson.E{Key: "_id", Value: oid}},
		)
	})

	t.Run("==INT_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`year==2022`,
			bson.D{bson.E{Key: "year", Value: int64(2022)}},
		)
	})

	t.Run("==FLOAT_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`pi==3.14159265`,
			bson.D{bson.E{Key: "pi", Value: float64(3.14159265)}},
		)
	})

	t.Run("==BOOL_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`is==TRUE`,
			bson.D{bson.E{Key: "is", Value: true}},
		)
		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`is==false`,
			bson.D{bson.E{Key: "is", Value: false}},
		)
	})

	t.Run("==ARRAY_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`roles==("dev","maintainer")`,
			bson.D{bson.E{Key: "roles", Value: bson.A{"dev", "maintainer"}}},
		)
	})

	t.Run("=in=SLICE_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`coll=in=(1, "a", "b",2)`,
			bson.D{bson.E{Key: "coll", Value: bson.E{
				Key:   "$in",
				Value: bson.A{int64(1), "a", "b", int64(2)},
			}}},
		)
	})
}

func TestQueryParsingWithSingleComparisonOperation(t *testing.T) {
	t.Parallel()

	t.Run("==_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`firstName=="steven"`,
			bson.D{bson.E{Key: "firstName", Value: "steven"}},
		)
	})

	t.Run("!=_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`x!=10`,
			bson.D{
				bson.E{Key: "x", Value: bson.D{
					bson.E{Key: "$ne", Value: int64(10)},
				}},
			},
		)
	})

	t.Run("=sw=_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`msg=sw="LOG_"`,
			bson.D{bson.E{Key: "msg", Value: *regexp.MustCompile("^LOG_")}},
		)
	})

	t.Run("=ew=_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`word=ew="ed"`,
			bson.D{bson.E{Key: "word", Value: *regexp.MustCompile("ed$")}},
		)
	})

	t.Run("=gt=_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`x=gt=10`,
			bson.D{
				bson.E{Key: "x", Value: bson.D{
					bson.E{Key: "$gt", Value: int64(10)},
				}},
			},
		)
	})

	t.Run("=ge=_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`x=ge=10`,
			bson.D{
				bson.E{Key: "x", Value: bson.D{
					bson.E{Key: "$gte", Value: int64(10)},
				}},
			},
		)
	})

	t.Run("=lt=_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`x=lt=10`,
			bson.D{
				bson.E{Key: "x", Value: bson.D{
					bson.E{Key: "$lt", Value: int64(10)},
				}},
			},
		)
	})

	t.Run("=le=_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`x=le=10`,
			bson.D{
				bson.E{Key: "x", Value: bson.D{
					bson.E{Key: "$lte", Value: int64(10)},
				}},
			},
		)
	})

	t.Run("=in=_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`coll=in=(1, "a", "b",2)`,
			bson.D{bson.E{Key: "coll", Value: bson.E{
				Key:   "$in",
				Value: bson.A{int64(1), "a", "b", int64(2)},
			}}},
		)
	})

	t.Run("=out=_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`coll=out=(1, "a", "b",2)`,
			bson.D{bson.E{Key: "coll", Value: bson.E{
				Key:   "$nin",
				Value: bson.A{int64(1), "a", "b", int64(2)},
			}}},
		)
	})
}

func TestQueryParsingWithMultipleComparisonOperation(t *testing.T) {
	t.Parallel()

	t.Run("WithSingleAnd_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`firstName=="steven";age=ge=18`,
			bson.D{
				bson.E{Key: "$and", Value: bson.A{
					bson.D{
						bson.E{Key: "firstName", Value: "steven"},
					},
					bson.D{
						bson.E{Key: "age", Value: bson.D{
							bson.E{Key: "$gte", Value: int64(18)},
						}},
					},
				}},
			},
		)
	})

	t.Run("WithSingleOr_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`level=="error",level=="warning"`,
			bson.D{
				bson.E{Key: "$or", Value: bson.A{
					bson.D{
						bson.E{Key: "level", Value: "error"},
					},
					bson.D{
						bson.E{Key: "level", Value: "warning"},
					},
				}},
			},
		)
	})

	t.Run("WithAndChain_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`firstName=="steven";age=ge=18;gender=="male"`,
			bson.D{
				bson.E{Key: "$and", Value: bson.A{
					bson.D{
						bson.E{Key: "firstName", Value: "steven"},
					},
					bson.D{
						bson.E{Key: "age", Value: bson.D{
							bson.E{Key: "$gte", Value: int64(18)},
						}},
					},
					bson.D{
						bson.E{Key: "gender", Value: "male"},
					},
				}},
			},
		)
	})

	t.Run("WithOrChain_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`level=="panic",level=="error",level=="warning"`,
			bson.D{
				bson.E{Key: "$or", Value: bson.A{
					bson.D{
						bson.E{Key: "level", Value: "panic"},
					},
					bson.D{
						bson.E{Key: "level", Value: "error"},
					},
					bson.D{
						bson.E{Key: "level", Value: "warning"},
					},
				}},
			},
		)
	})

	t.Run("WithMixed_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`a==1,a==2,a==3,b==1;c==1`,
			bson.D{
				bson.E{Key: "$or", Value: bson.A{
					bson.D{
						bson.E{Key: "a", Value: int64(1)},
					},
					bson.D{
						bson.E{Key: "a", Value: int64(2)},
					},
					bson.D{
						bson.E{Key: "a", Value: int64(3)},
					},
					bson.D{
						bson.E{Key: "$and", Value: bson.A{
							bson.D{
								bson.E{Key: "b", Value: int64(1)},
							},
							bson.D{
								bson.E{Key: "c", Value: int64(1)},
							},
						}},
					},
				}},
			},
		)

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`a==1;b==1,a==2;b==2`,
			bson.D{
				bson.E{Key: "$and", Value: bson.A{
					bson.D{
						bson.E{Key: "a", Value: int64(1)},
					},
					bson.D{
						bson.E{Key: "$or", Value: bson.A{
							bson.D{
								bson.E{Key: "b", Value: int64(1)},
							},
							bson.D{
								bson.E{Key: "$and", Value: bson.A{
									bson.D{
										bson.E{Key: "a", Value: int64(2)},
									},
									bson.D{
										bson.E{Key: "b", Value: int64(2)},
									},
								}},
							},
						}},
					},
				}},
			},
		)
	})

	t.Run("WithContext_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`(a==1;b==1),(a==2;b==2),(a==3;b==3)`,
			bson.D{
				bson.E{Key: "$or", Value: bson.A{
					bson.D{
						bson.E{Key: "$and", Value: bson.A{
							bson.D{
								bson.E{Key: "a", Value: int64(1)},
							},
							bson.D{
								bson.E{Key: "b", Value: int64(1)},
							},
						}},
					},
					bson.D{
						bson.E{Key: "$and", Value: bson.A{
							bson.D{
								bson.E{Key: "a", Value: int64(2)},
							},
							bson.D{
								bson.E{Key: "b", Value: int64(2)},
							},
						}},
					},
					bson.D{
						bson.E{Key: "$and", Value: bson.A{
							bson.D{
								bson.E{Key: "a", Value: int64(3)},
							},
							bson.D{
								bson.E{Key: "b", Value: int64(3)},
							},
						}},
					},
				}},
			},
		)
	})

	t.Run("WithNestedContext_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(nil),
			`(a==1;b==1),((a==2,b==2);(a==3,b==3))`,
			bson.D{
				bson.E{Key: "$or", Value: bson.A{
					bson.D{
						bson.E{Key: "$and", Value: bson.A{
							bson.D{
								bson.E{Key: "a", Value: int64(1)},
							},
							bson.D{
								bson.E{Key: "b", Value: int64(1)},
							},
						}},
					},
					bson.D{
						bson.E{Key: "$and", Value: bson.A{
							bson.D{
								bson.E{Key: "$or", Value: bson.A{
									bson.D{
										bson.E{Key: "a", Value: int64(2)},
									},
									bson.D{
										bson.E{Key: "b", Value: int64(2)},
									},
								}},
							},
							bson.D{
								bson.E{Key: "$or", Value: bson.A{
									bson.D{
										bson.E{Key: "a", Value: int64(3)},
									},
									bson.D{
										bson.E{Key: "b", Value: int64(3)},
									},
								}},
							},
						}},
					},
				}},
			},
		)
	})
}

func TestQueryParsingFailCases(t *testing.T) {
	t.Parallel()

	t.Run("WithUnknownCompareOperation_Fail", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteFailedTest(t,
			NewParser(nil),
			`x=7`,
			errs.NewErrUnexpectedToken(1, "="),
		)
	})

	t.Run("WithIncompleteExpression_Fail", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteFailedTest(t,
			NewParser(nil),
			`x==7;`,
			errs.NewErrUnexpectedInputEnd("FIELD_NAME"),
		)
	})

	t.Run("WithWrongCompareOperation_Fail", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteFailedTest(t,
			NewParser(nil),
			`x=in=3`,
			errs.NewErrUnexpectedTokenType(6, "NUMERIC_LITERAL", "("),
		)
	})

	t.Run("WithNotClosedContext_Fail", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteFailedTest(t,
			NewParser(nil),
			`(x==7`,
			errs.NewErrUnexpectedInputEnd(")"),
		)
	})
}

func TestQueryParsingWithPolicy(t *testing.T) {
	t.Parallel()

	t.Run("WithAllowedFieldName_Success", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteSuccessTest(t,
			NewParser(
				tokenizer.NewPolicy(tokenizer.WhitelistPolicy, "name", "age")),
			`name=="steven",age=ge=18`,
			bson.D{
				bson.E{Key: "$or", Value: bson.A{
					bson.D{
						bson.E{Key: "name", Value: "steven"},
					},
					bson.D{
						bson.E{Key: "age", Value: bson.D{
							bson.E{Key: "$gte", Value: int64(18)},
						}},
					},
				}},
			},
		)
	})

	t.Run("WithDisallowedFieldNames_Fail", func(t *testing.T) {
		t.Parallel()

		testutil.ExecuteFailedTest(t,
			NewParser(
				tokenizer.NewPolicy(tokenizer.WhitelistPolicy, "name", "age")),
			`name=="steven",age=ge=18,gender="male"`,
			errs.NewErrPolicyViolation("gender"),
		)
	})
}

func TestInterpretation(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "darwin" {
		t.Skip("Not running on darwin")
	}

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

	t.Run("FilterByGender_Success", func(t *testing.T) {
		t.Parallel()

		parser := NewParser(nil)
		filter, err := parser.Parse(`gender=="female"`)
		require.NoError(t, err)

		testutil.FindCompare(t, collection, filter, nil, items[1], items[2])
	})

	t.Run("FilterByGenderAndAgeGreater_Success", func(t *testing.T) {
		t.Parallel()

		parser := NewParser(nil)
		filter, err := parser.Parse(`gender=="female";age=ge=30`)
		require.NoError(t, err)

		testutil.FindCompare(t, collection, filter, nil, items[2])
	})

	t.Run("FilterAgeBetween20And50_Success", func(t *testing.T) {
		t.Parallel()

		parser := NewParser(nil)
		filter, err := parser.Parse(`age=ge=20;age=le=50`)
		require.NoError(t, err)

		testutil.FindCompare(t, collection, filter, nil, items[1], items[2], items[3])
	})

	t.Run("FilterAgeBetween25And50OrName_Success", func(t *testing.T) {
		t.Parallel()

		parser := NewParser(nil)
		filter, err := parser.Parse(`first_name=="Alexa",(age=ge=25;age=le=50)`)
		require.NoError(t, err)

		testutil.FindCompare(t, collection, filter, nil, items[1], items[2], items[3])
	})
}
