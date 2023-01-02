//nolint:goconst
package jsonpatch

import (
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
	"github.com/stretchr/testify/require"
)

func TestPolicyInterface(t *testing.T) {
	t.Parallel()
	// just test the assigning

	var (
		policy1 Policy = DisallowPathPolicy{}
		policy2 Policy = DisallowOperationOnPathPolicy{}
		policy3 Policy = ForceTypeOnPathPolicy{}
		policy4 Policy = ForceRegexMatchPolicy{}
		policy5 Policy = StrictPathPolicy{}
	)

	// so the variables are used
	require.NotNil(t, policy1)
	require.NotNil(t, policy2)
	require.NotNil(t, policy3)
	require.NotNil(t, policy4)
	require.NotNil(t, policy5)
}

func TestDisallowPathPolicy(t *testing.T) {
	t.Parallel()

	details := "something"
	forbiddenPath := operation.Path("user.password")
	allowedPath := operation.Path("user.userdetails")
	policy := DisallowPathPolicy{Details: details, Path: forbiddenPath}

	require.Equal(t, details, policy.Details)
	require.False(t, policy.Test(operation.Spec{Path: forbiddenPath}))
	require.True(t, policy.Test(operation.Spec{Path: allowedPath}))
}

func TestDisallowOperationOnPathPolicy(t *testing.T) {
	t.Parallel()

	details := "something"
	path := operation.Path("user.userdetails")
	notEffected := operation.Path("user.disabled")
	policy := DisallowOperationOnPathPolicy{
		Details: details, Path: path, Operations: []operation.Operation{operation.RemoveOperation, operation.AddOperation},
	}

	require.Equal(t, details, policy.Details)
	require.False(t, policy.Test(operation.Spec{Path: path, Operation: operation.RemoveOperation}))
	require.False(t, policy.Test(operation.Spec{Path: path, Operation: operation.AddOperation}))
	require.True(t, policy.Test(operation.Spec{Path: path, Operation: operation.ReplaceOperation}))
	require.True(t, policy.Test(operation.Spec{Path: notEffected, Operation: operation.RemoveOperation}))
}

func TestForceTypeOnPathPolicy(t *testing.T) {
	t.Parallel()

	details := "something"
	path := operation.Path("product.price")
	notEffected := operation.Path("product.tag")
	policy := ForceTypeOnPathPolicy{
		Details: details, Path: path, Kind: reflect.Float64,
	}

	require.Equal(t, details, policy.Details)
	require.False(t, policy.Test(operation.Spec{Path: path, Value: "not a price"}))
	require.True(t, policy.Test(operation.Spec{Path: path, Value: 6.99}))
	require.True(t, policy.Test(operation.Spec{Path: notEffected, Value: "vegetable"}))
}

func TestForceRegexMatchPolicy(t *testing.T) {
	t.Parallel()

	details := "something"
	path := operation.Path("product.version")
	notEffected := operation.Path("product.tag")
	policy := ForceRegexMatchPolicy{
		Details: details, Path: path,
		Expression: *regexp.MustCompile(`^v[0-9]*\.[0-9]*\.[0-9]*$`),
	}

	policy2 := ForceRegexMatchPolicy{
		Details: details, Path: "*.version",
		Expression: *regexp.MustCompile(`^v[0-9]*\.[0-9]*\.[0-9]*$`),
	}

	require.Equal(t, details, policy.Details)
	require.False(t, policy.Test(operation.Spec{Path: path, Value: "va.0.0"}))
	require.True(t, policy.Test(operation.Spec{Path: path, Value: "v0.3.7"}))
	require.True(t, policy.Test(operation.Spec{Path: notEffected, Value: "backend"}))
	require.True(t, policy2.Test(operation.Spec{Path: "api.version", Value: "v0.3.7"}))
	require.True(t, policy2.Test(operation.Spec{Path: "backend.version", Value: "v0.3.7"}))
}

func TestStrictPathPolicy(t *testing.T) {
	t.Parallel()

	details := "something"
	path1 := operation.Path("product.version")
	path2 := operation.Path("product.timestamp")
	path3 := operation.Path("product.tags.*")
	path4 := operation.Path("product.meta.*.name")
	path5 := operation.Path("*.*.*.title")
	invalidPath := operation.Path("product.tag")
	policy := StrictPathPolicy{
		Details: details, Path: []operation.Path{path1, path2, path3, path4, path5},
	}

	require.Equal(t, details, policy.Details)
	require.True(t, policy.Test(operation.Spec{Path: path1, Value: "v0.3.7"}))
	require.True(t, policy.Test(operation.Spec{Path: path2, Value: "24.11.20022"}))
	require.True(t, policy.Test(operation.Spec{
		Path: operation.Path(strings.ReplaceAll(string(path3), "*", "key")), Value: "xyz",
	}))
	require.False(t, policy.Test(operation.Spec{
		Path:  operation.Path(strings.ReplaceAll(string(path3), "*", "key")) + ".key2",
		Value: "xyz",
	}))
	require.True(t, policy.Test(operation.Spec{
		Path: operation.Path(strings.ReplaceAll(string(path4), "*", "key")), Value: "xyz",
	}))
	require.True(t, policy.Test(operation.Spec{Path: operation.Path("product.group.a.title"), Value: "xyz"}))
	require.True(t, policy.Test(operation.Spec{Path: operation.Path("product.group.b.title"), Value: "xyz"}))
	require.False(t, policy.Test(operation.Spec{Path: invalidPath, Value: "something"}))
}

func TestForceOperationOnPathPolicy(t *testing.T) {
	t.Parallel()

	details := "something"
	path := operation.Path("user.userdetails")
	policy := ForceOperationOnPathPolicy{
		Details: details, Path: path, Operation: operation.AddOperation,
	}

	require.Equal(t, details, policy.Details)
	require.True(t, policy.Test(operation.Spec{Path: path, Operation: operation.AddOperation}))
	require.False(t, policy.Test(operation.Spec{Path: path, Operation: operation.RemoveOperation}))
	require.False(t, policy.Test(operation.Spec{Path: path, Operation: operation.ReplaceOperation}))
}
