//nolint:dupl,tagliatelle
package validator

import (
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/rule"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type demoRule struct{}

func (d demoRule) NewInstance(_ string, _ reflect.Kind, _ interface{}, value string) (rule.Rule, error) { //nolint:ireturn,lll
	return nil, nil
}

func (d demoRule) NewInheritInstance(_ string, _ reflect.Kind, _ interface{}) (rule.Rule, error) { //nolint:ireturn,lll
	return nil, nil
}

func (d demoRule) Validate(_ operation.Spec) error {
	return nil
}

func TestInstantiation(t *testing.T) {
	t.Parallel()

	t.Run("ValidReference_Success", func(t *testing.T) {
		t.Parallel()

		validator, err := NewValidator(reflect.TypeOf(struct{}{}))
		require.NoError(t, err)
		require.NotNil(t, validator)
	})

	t.Run("InvalidReference_Fail", func(t *testing.T) {
		t.Parallel()

		validator, err := NewValidator(reflect.TypeOf(nil))
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrReferenceIsNil))
		require.Nil(t, validator)
	})
}

func TestRegisterRule(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct{}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	t.Run("ValidRule_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.RegisterRule("jp_name", demoRule{})
		require.NoError(t, err)
	})

	t.Run("DuplicateRuleName_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.RegisterRule("jp_disallow", demoRule{})
		require.Error(t, err)
		require.Equal(t, ErrDuplicateRuleTags, err)
	})

	t.Run("MissingPrefix_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.RegisterRule("name", demoRule{})
		require.Error(t, err)
		require.Equal(t, ErrMissingPrefix, err)
	})
}

func TestHandlingObjectID(t *testing.T) {
	t.Parallel()

	newID := primitive.NewObjectID()
	newIDHex := newID.Hex()
	newIDRaw := []int{99, 137, 14, 13, 185, 204, 14, 42, 43, 253, 38, 103}

	validator, err := NewValidator(reflect.TypeOf(struct { //nolint:govet
		ID    primitive.ObjectID  `bson:"_id"`
		IDPtr *primitive.ObjectID `bson:"_id_ptr"`
	}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	err = validator.Validate(operation.Spec{
		Operation: operation.ReplaceOperation,
		Path:      "_id",
		Value:     newIDRaw,
	})
	require.NoError(t, err)

	err = validator.Validate(operation.Spec{
		Operation: operation.ReplaceOperation,
		Path:      "_id",
		Value:     newID,
	})
	require.NoError(t, err)

	err = validator.Validate(operation.Spec{
		Operation: operation.ReplaceOperation,
		Path:      "_id",
		Value:     newIDHex,
	})
	require.NoError(t, err)

	err = validator.Validate(operation.Spec{
		Operation: operation.ReplaceOperation,
		Path:      "_id_ptr",
		Value:     newID,
	})
	require.NoError(t, err)

	err = validator.Validate(operation.Spec{
		Operation: operation.ReplaceOperation,
		Path:      "_id_ptr",
		Value:     newIDHex,
	})
	require.NoError(t, err)

	err = validator.Validate(operation.Spec{
		Operation: operation.ReplaceOperation,
		Path:      "_id_ptr",
		Value:     newIDRaw,
	})
	require.NoError(t, err)
}

func TestHandlingObjectIDArray(t *testing.T) {
	t.Parallel()

	newID := primitive.NewObjectID()
	newIDHex := newID.Hex()
	newIDRaw := []int{99, 137, 14, 13, 185, 204, 14, 42, 43, 253, 38, 103}

	validator, err := NewValidator(reflect.TypeOf(struct { //nolint:govet
		Arr    []primitive.ObjectID  `bson:"arr"`
		ArrPtr *[]primitive.ObjectID `bson:"arr_ptr"`
	}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	err = validator.Validate(operation.Spec{
		Operation: operation.ReplaceOperation,
		Path:      "arr",
		Value:     []interface{}{newID},
	})
	require.NoError(t, err)

	err = validator.Validate(operation.Spec{
		Operation: operation.ReplaceOperation,
		Path:      "arr",
		Value:     []interface{}{newIDHex},
	})
	require.NoError(t, err)

	err = validator.Validate(operation.Spec{
		Operation: operation.ReplaceOperation,
		Path:      "arr",
		Value:     []interface{}{newIDRaw},
	})
	require.NoError(t, err)

	err = validator.Validate(operation.Spec{
		Operation: operation.ReplaceOperation,
		Path:      "arr_ptr",
		Value:     []interface{}{newID},
	})
	require.NoError(t, err)

	err = validator.Validate(operation.Spec{
		Operation: operation.ReplaceOperation,
		Path:      "arr_ptr",
		Value:     []interface{}{newIDHex},
	})
	require.NoError(t, err)

	err = validator.Validate(operation.Spec{
		Operation: operation.ReplaceOperation,
		Path:      "arr_ptr",
		Value:     []interface{}{newIDRaw},
	})
	require.NoError(t, err)
}

func TestUseReferenceWithSimpleTypes(t *testing.T) { //nolint:funlen
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct{}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	err = validator.UseReference(reflect.TypeOf(struct { //nolint:govet
		A string
		B string            `bson:"b"`
		C int               `bson:"c"`
		D []string          `bson:"d"`
		E []int             `bson:"e"`
		F map[string]string `bson:"f"`
		G map[string]int    `bson:"g"`
		H struct {
			A string
			B string `bson:"b"`
		} `bson:"h"`
	}{}))
	expectedRule := map[operation.Path]map[string]rule.Rule{
		"b": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "b"},
		},
		"c": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Int},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: 0, Path: "c"},
		},
		"d": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Slice},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: []string{}, Path: "d"},
		},
		"e": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Slice},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: []int{}, Path: "e"},
		},
		"f": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Map},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: map[string]string{}, Path: "f"},
		},
		"g": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Map},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: map[string]int{}, Path: "g"},
		},
		"h": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Struct},
			"jp_general_matching_kind": &rule.MatchingKindRule{Reference: struct {
				A string
				B string `bson:"b"`
			}{}, Path: "h"},
		},
		"h.b": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "h.b"},
		},
	}
	expectedWildcardRules := map[operation.Path]map[string]rule.Rule{
		"d.*": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "d.*"},
		},
		"e.*": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Int},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: 0, Path: "e.*"},
		},
		"f.*": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "f.*"},
		},
		"g.*": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Int},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: 0, Path: "g.*"},
		},
	}

	require.NoError(t, err)
	require.Equal(t, expectedRule, validator.rules)
	require.Equal(t, expectedWildcardRules, validator.wildcardRules)
}

func TestUseReferenceComplexStruct(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct{}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	err = validator.UseReference(reflect.TypeOf(struct {
		A struct {
			A struct {
				B string `bson:"b"`
			} `bson:"a"`
		} `bson:"a"`
	}{}))
	expectedRule := map[operation.Path]map[string]rule.Rule{
		"a": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Struct},
			"jp_general_matching_kind": &rule.MatchingKindRule{Reference: struct {
				A struct {
					B string `bson:"b"`
				} `bson:"a"`
			}{}, Path: "a"},
		},
		"a.a": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Struct},
			"jp_general_matching_kind": &rule.MatchingKindRule{Reference: struct {
				B string `bson:"b"`
			}{}, Path: "a.a"},
		},
		"a.a.b": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "a.a.b"},
		},
	}
	expectedWildcardRules := map[operation.Path]map[string]rule.Rule{}

	require.NoError(t, err)

	require.NoError(t, err)
	assertDeepEqual(t, expectedRule, validator.rules)
	assertDeepEqual(t, expectedWildcardRules, validator.wildcardRules)
}

func assertDeepEqual(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Objects not equal. Expected %+v, but got %+v", expected, actual)
	}
}

func TestUseReferenceComplexArray(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct{}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	err = validator.UseReference(reflect.TypeOf(struct {
		B []struct {
			A string `bson:"a"`
		} `bson:"b"`
	}{}))
	expectedRule := map[operation.Path]map[string]rule.Rule{
		"b": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Slice},
			"jp_general_matching_kind": &rule.MatchingKindRule{Reference: []struct {
				A string `bson:"a"`
			}{}, Path: "b"},
		},
	}
	expectedWildcardRules := map[operation.Path]map[string]rule.Rule{
		"b.*": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Struct},
			"jp_general_matching_kind": &rule.MatchingKindRule{Reference: struct {
				A string `bson:"a"`
			}{}, Path: "b.*"},
		},
		"b.*.a": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "b.*.a"},
		},
	}

	require.NoError(t, err)
	require.Equal(t, expectedRule, validator.rules)
	require.Equal(t, expectedWildcardRules, validator.wildcardRules)
}

func TestUseReferenceComplexMap(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct{}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	err = validator.UseReference(reflect.TypeOf(struct {
		C map[string]struct {
			A string `bson:"a"`
		} `bson:"c"`
	}{}))
	expectedRule := map[operation.Path]map[string]rule.Rule{
		"c": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Map},
			"jp_general_matching_kind": &rule.MatchingKindRule{Reference: map[string]struct {
				A string `bson:"a"`
			}{}, Path: "c"},
		},
	}
	expectedWildcardRules := map[operation.Path]map[string]rule.Rule{
		"c.*": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Struct},
			"jp_general_matching_kind": &rule.MatchingKindRule{Reference: struct {
				A string `bson:"a"`
			}{}, Path: "c.*"},
		},
		"c.*.a": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "c.*.a"},
		},
	}

	require.NoError(t, err)
	require.Equal(t, expectedRule, validator.rules)
	require.Equal(t, expectedWildcardRules, validator.wildcardRules)
}

func TestUseReferenceComplexNested(t *testing.T) { //nolint:funlen
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct{}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	err = validator.UseReference(reflect.TypeOf(struct {
		D map[string][]struct {
			A string `bson:"a"`
		} `bson:"d"`
		E []map[string]struct {
			A string `bson:"a"`
		} `bson:"e"`
	}{}))
	expectedRule := map[operation.Path]map[string]rule.Rule{
		"d": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Map},
			"jp_general_matching_kind": &rule.MatchingKindRule{Reference: map[string][]struct {
				A string `bson:"a"`
			}{}, Path: "d"},
		},
		"e": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Slice},
			"jp_general_matching_kind": &rule.MatchingKindRule{Reference: []map[string]struct {
				A string `bson:"a"`
			}{}, Path: "e"},
		},
	}
	expectedWildcardRules := map[operation.Path]map[string]rule.Rule{
		"d.*": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Slice},
			"jp_general_matching_kind": &rule.MatchingKindRule{Reference: []struct {
				A string `bson:"a"`
			}{}, Path: "d.*"},
		},
		"d.*.*": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Struct},
			"jp_general_matching_kind": &rule.MatchingKindRule{Reference: struct {
				A string `bson:"a"`
			}{}, Path: "d.*.*"},
		},
		"d.*.*.a": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "d.*.*.a"},
		},
		"e.*": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Map},
			"jp_general_matching_kind": &rule.MatchingKindRule{Reference: map[string]struct {
				A string `bson:"a"`
			}{}, Path: "e.*"},
		},
		"e.*.*": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Struct},
			"jp_general_matching_kind": &rule.MatchingKindRule{Reference: struct {
				A string `bson:"a"`
			}{}, Path: "e.*.*"},
		},
		"e.*.*.a": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "e.*.*.a"},
		},
	}

	require.NoError(t, err)
	require.Equal(t, expectedRule, validator.rules)
	require.Equal(t, expectedWildcardRules, validator.wildcardRules)
}

func TestUseReferenceWithSimpleRules(t *testing.T) { //nolint:funlen
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct{}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	err = validator.UseReference(reflect.TypeOf(struct {
		A string `bson:"a"`
		B string `bson:"b" jp_disallow:"true"`
		C string `bson:"c" jp_min:"3"`
		D string `bson:"d" jp_max:"3"`
		E string `bson:"e" jp_expression:"^\\w+$"`
		F string `bson:"f" jp_op_allowed:"add,remove"`
		G string `bson:"g" jp_op_disallowed:"add,remove"`
	}{}))
	expectedRule := map[operation.Path]map[string]rule.Rule{
		"a": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "a"},
		},
		"b": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "b"},
			"jp_disallow":                           &rule.DisallowRule{Disallow: true},
		},
		"c": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "c"},
			"jp_min":                                &rule.MinRule{Min: 3},
		},
		"d": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "d"},
			"jp_max":                                &rule.MaxRule{Max: 3},
		},
		"e": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "e"},
			"jp_expression": &rule.ExpressionRule{
				Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`),
			},
		},
		"f": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "f"},
			"jp_op_allowed": &rule.AllowedOperationsRule{
				Operations: []operation.Operation{operation.AddOperation, operation.RemoveOperation},
			},
		},
		"g": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "g"},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{
				Operations: []operation.Operation{operation.AddOperation, operation.RemoveOperation},
			},
		},
	}
	expectedWildcardRules := map[operation.Path]map[string]rule.Rule{}

	require.NoError(t, err)
	require.Equal(t, expectedRule, validator.rules)
	require.Equal(t, expectedWildcardRules, validator.wildcardRules)
}

func TestUseReferenceWithHeredityStruct(t *testing.T) { //nolint:funlen
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct{}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	type demoStruct struct {
		DoInherit   string `bson:"do_inherit"`
		DoOverwrite string `bson:"do_overwrite" jp_disallow:"false" jp_min:"2" jp_max:"2" jp_expression:"^\\d+$" jp_op_allowed:"move" jp_op_disallowed:"replace"` //nolint:lll
		Nested      struct {
			DoInherit string `bson:"do_inherit"`
			Nested    struct {
				DoInherit string `bson:"do_inherit"`
			} `bson:"nested"`
		} `bson:"nested"`
	}

	err = validator.UseReference(reflect.TypeOf(struct {
		A demoStruct `bson:"a" jp_inherit:"jp_disallow,jp_min,jp_max,jp_expression,jp_op_allowed,jp_op_disallowed" jp_disallow:"true" jp_min:"3" jp_max:"3" jp_expression:"^\\w+$" jp_op_allowed:"add" jp_op_disallowed:"remove"` //nolint:lll
	}{}))
	require.NoError(t, err)

	expectedRule := map[operation.Path]map[string]rule.Rule{
		"a": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Struct},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: demoStruct{}, Path: "a"},
			"jp_disallow":                           &rule.DisallowRule{Disallow: true},
			"jp_min":                                &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.do_inherit": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "a.do_inherit"},
			"jp_disallow":                           &rule.DisallowRule{Disallow: true},
			"jp_min":                                &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.do_overwrite": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "a.do_overwrite"},
			"jp_disallow":                           &rule.DisallowRule{Disallow: false},
			"jp_min":                                &rule.MinRule{Min: 2}, "jp_max": &rule.MaxRule{Max: 2},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\d+$`, Regex: *regexp.MustCompile(`^\d+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.MoveOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.ReplaceOperation}},
		},
		"a.nested": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Struct},
			"jp_general_matching_kind": &rule.MatchingKindRule{Reference: struct {
				DoInherit string `bson:"do_inherit"`
				Nested    struct {
					DoInherit string `bson:"do_inherit"`
				} `bson:"nested"`
			}{}, Path: "a.nested"},
			"jp_disallow": &rule.DisallowRule{Disallow: true},
			"jp_min":      &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.nested.do_inherit": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "a.nested.do_inherit"},
			"jp_disallow":                           &rule.DisallowRule{Disallow: true},
			"jp_min":                                &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.nested.nested": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Struct},
			"jp_general_matching_kind": &rule.MatchingKindRule{Reference: struct {
				DoInherit string `bson:"do_inherit"`
			}{}, Path: "a.nested.nested"},
			"jp_disallow": &rule.DisallowRule{Disallow: true},
			"jp_min":      &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.nested.nested.do_inherit": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "a.nested.nested.do_inherit"},
			"jp_disallow":                           &rule.DisallowRule{Disallow: true},
			"jp_min":                                &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
	}
	expectedWildcardRules := map[operation.Path]map[string]rule.Rule{}

	require.Equal(t, expectedRule, validator.rules)
	require.Equal(t, expectedWildcardRules, validator.wildcardRules)
}

func TestUseReferenceWithHeredityArray(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct{}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	type demoStruct struct {
		DoInherit   string `bson:"do_inherit"`
		DoOverwrite string `bson:"do_overwrite" jp_disallow:"false" jp_min:"2" jp_max:"2" jp_expression:"^\\d+$" jp_op_allowed:"move" jp_op_disallowed:"replace"` //nolint:lll
	}

	err = validator.UseReference(reflect.TypeOf(struct {
		A []demoStruct `bson:"a" jp_inherit:"jp_disallow,jp_min,jp_max,jp_expression,jp_op_allowed,jp_op_disallowed" jp_disallow:"true" jp_min:"3" jp_max:"3" jp_expression:"^\\w+$" jp_op_allowed:"add" jp_op_disallowed:"remove"` //nolint:lll
	}{}))
	require.NoError(t, err)

	expectedRule := map[operation.Path]map[string]rule.Rule{
		"a": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Slice},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: []demoStruct{}, Path: "a"},
			"jp_disallow":                           &rule.DisallowRule{Disallow: true},
			"jp_min":                                &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
	}
	expectedWildcardRules := map[operation.Path]map[string]rule.Rule{
		"a.*": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Struct},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: demoStruct{}, Path: "a.*"},
			"jp_disallow":                           &rule.DisallowRule{Disallow: true},
			"jp_min":                                &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.*.do_inherit": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "a.*.do_inherit"},
			"jp_disallow":                           &rule.DisallowRule{Disallow: true},
			"jp_min":                                &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.*.do_overwrite": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "a.*.do_overwrite"},
			"jp_disallow":                           &rule.DisallowRule{Disallow: false},
			"jp_min":                                &rule.MinRule{Min: 2}, "jp_max": &rule.MaxRule{Max: 2},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\d+$`, Regex: *regexp.MustCompile(`^\d+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.MoveOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.ReplaceOperation}},
		},
	}

	require.Equal(t, expectedRule, validator.rules)
	require.Equal(t, expectedWildcardRules, validator.wildcardRules)
}

func TestUseReferenceWithHeredityMap(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct{}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	type demoStruct struct {
		DoInherit   string `bson:"do_inherit"`
		DoOverwrite string `bson:"do_overwrite" jp_disallow:"false" jp_min:"2" jp_max:"2" jp_expression:"^\\d+$" jp_op_allowed:"move" jp_op_disallowed:"replace"` //nolint:lll
	}

	err = validator.UseReference(reflect.TypeOf(struct {
		A map[string]demoStruct `bson:"a" jp_inherit:"jp_disallow,jp_min,jp_max,jp_expression,jp_op_allowed,jp_op_disallowed" jp_disallow:"true" jp_min:"3" jp_max:"3" jp_expression:"^\\w+$" jp_op_allowed:"add" jp_op_disallowed:"remove"` //nolint:lll
	}{}))
	require.NoError(t, err)

	expectedRule := map[operation.Path]map[string]rule.Rule{
		"a": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Map},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: map[string]demoStruct{}, Path: "a"},
			"jp_disallow":                           &rule.DisallowRule{Disallow: true},
			"jp_min":                                &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
	}
	expectedWildcardRules := map[operation.Path]map[string]rule.Rule{
		"a.*": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.Struct},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: demoStruct{}, Path: "a.*"},
			"jp_disallow":                           &rule.DisallowRule{Disallow: true},
			"jp_min":                                &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.*.do_inherit": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "a.*.do_inherit"},
			"jp_disallow":                           &rule.DisallowRule{Disallow: true},
			"jp_min":                                &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.*.do_overwrite": {
			"jp_general_matching_operation_to_kind": &rule.MatchingOperationToKindRule{Kind: reflect.String},
			"jp_general_matching_kind":              &rule.MatchingKindRule{Reference: "", Path: "a.*.do_overwrite"},
			"jp_disallow":                           &rule.DisallowRule{Disallow: false},
			"jp_min":                                &rule.MinRule{Min: 2}, "jp_max": &rule.MaxRule{Max: 2},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\d+$`, Regex: *regexp.MustCompile(`^\d+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.MoveOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.ReplaceOperation}},
		},
	}

	require.Equal(t, expectedRule, validator.rules)
	require.Equal(t, expectedWildcardRules, validator.wildcardRules)
}

func TestValidateTypecheckRule(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct {
		A string `bson:"a"`
	}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	t.Run("ValidType_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.ReplaceOperation, Path: "a", Value: "new"})
		require.NoError(t, err)
	})

	t.Run("InvalidType_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.ReplaceOperation, Path: "a", Value: 123})
		require.Error(t, err)
		require.Equal(t, "operation no allowed: 'a' has invalid kind 'int', must be 'string'", err.Error())
	})
}

type demoValidateStruct struct { //nolint:govet
	A int `bson:"a"`
	B struct {
		C int `bson:"c"`
	} `bson:"b"`
	D []struct {
		E int `bson:"e"`
	} `bson:"d"`
	F map[string]struct {
		G int `bson:"g"`
	} `bson:"f"`
}

func TestValidatePath(t *testing.T) { //nolint:funlen
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(demoValidateStruct{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	t.Run("ValidRootPath_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.ReplaceOperation, Path: "a", Value: 1})
		require.NoError(t, err)
	})

	t.Run("ValidNestedPath_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.ReplaceOperation, Path: "b.c", Value: 1})
		require.NoError(t, err)
	})

	t.Run("ValidNestedArrayPath_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.ReplaceOperation, Path: "d.0.e", Value: 1})
		require.NoError(t, err)
	})

	t.Run("ValidNestedMapPath_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.ReplaceOperation, Path: "f.0.g", Value: 1})
		require.NoError(t, err)
	})

	t.Run("InvalidRootPath_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.ReplaceOperation, Path: "aa", Value: 1})
		require.Error(t, err)
		require.Equal(t, "defined path 'aa' is unknown", err.Error())
	})

	t.Run("InvalidNestedPath_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.ReplaceOperation, Path: "x.f", Value: 1})
		require.Error(t, err)
		require.Equal(t, "defined path 'x.f' is unknown", err.Error())
	})

	t.Run("InvalidNestedPath_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.ReplaceOperation, Path: "d.0.c", Value: 1})
		require.Error(t, err)
		require.Equal(t, "defined path 'd.0.c' is unknown", err.Error())
	})

	t.Run("InvalidNestedPath_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.ReplaceOperation, Path: "f.0.c", Value: 1})
		require.Error(t, err)
		require.Equal(t, "defined path 'f.0.c' is unknown", err.Error())
	})
}

func TestValidateDisallowRule(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct {
		A string `bson:"a" jp_disallow:"false"`
		B string `bson:"b" jp_disallow:"true"`
	}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	t.Run("ValidType_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.ReplaceOperation, Path: "a", Value: "new"})
		require.NoError(t, err)
	})

	t.Run("InvalidType_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.ReplaceOperation, Path: "b", Value: 123})
		require.Error(t, err)
	})
}

func TestValidateMinRule(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct {
		A int `bson:"a" jp_min:"3"`
	}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	t.Run("ValidType_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.ReplaceOperation, Path: "a", Value: 3})
		require.NoError(t, err)
	})

	t.Run("InvalidType_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.ReplaceOperation, Path: "a", Value: 1})
		require.Error(t, err)
		require.Equal(t, "operation no allowed: value is less then specified: '1.000000' < '3.000000'", err.Error())
	})
}

func TestValidateMaxRule(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct {
		A int `bson:"a" jp_max:"3"`
	}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	t.Run("ValidType_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.ReplaceOperation, Path: "a", Value: 3})
		require.NoError(t, err)
	})

	t.Run("InvalidType_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.ReplaceOperation, Path: "a", Value: 123})
		require.Error(t, err)
		require.Equal(t, "operation no allowed: value is greater then specified: '123.000000' > '3.000000'", err.Error())
	})
}

func TestValidateExpressionRule(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct {
		A string `bson:"a" jp_expression:"^[a-z]+$"`
	}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	t.Run("ValidType_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.ReplaceOperation, Path: "a", Value: "abc"})
		require.NoError(t, err)
	})

	t.Run("InvalidType_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.ReplaceOperation, Path: "a", Value: "123"})
		require.Error(t, err)
		require.Equal(t, "operation no allowed: expression '^[a-z]+$' not match 123", err.Error())
	})
}

func TestValidateAllowedOperationsRule(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct {
		A string `bson:"a" jp_op_allowed:"replace"`
	}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	t.Run("ValidType_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.ReplaceOperation, Path: "a", Value: "abc"})
		require.NoError(t, err)
	})

	t.Run("InvalidType_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.RemoveOperation, Path: "a"})
		require.Error(t, err)
		require.Equal(t, "operation no allowed: operation 'remove' not allowed", err.Error())
	})
}

func TestValidateDisallowedOperationsRule(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct {
		A string `bson:"a" jp_op_disallowed:"remove"`
	}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	t.Run("ValidType_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.ReplaceOperation, Path: "a", Value: "abc"})
		require.NoError(t, err)
	})

	t.Run("InvalidType_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.RemoveOperation, Path: "a"})
		require.Error(t, err)
		require.Equal(t, "operation no allowed: operation 'remove' not allowed", err.Error())
	})
}
