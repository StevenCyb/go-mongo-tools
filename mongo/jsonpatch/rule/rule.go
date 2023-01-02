package rule

import (
	"reflect"

	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
)

// Rule defines the interface for a patch operation rule.
type Rule interface {
	NewInstance(patch string, kind reflect.Kind, instance interface{}, value string) (Rule, error)
	NewInheritInstance(patch string, kind reflect.Kind, instance interface{}) (Rule, error)
	Validate(operationSpec operation.Spec) error
}
