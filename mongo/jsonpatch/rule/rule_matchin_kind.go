//nolint:ireturn
package rule

import (
	"reflect"

	"github.com/StevenCyb/go-mongo-tools/mongo/jsonpatch/operation"
)

// MatchingKindRule is a default rule that is applied to all fields.
// This rules checks for type and name matches to prevent input for
// unknown fields or to violate types.
type MatchingKindRule struct {
	Reference interface{}
	Path      string
}

// UseValue instantiate new rule instance for field.
func (m *MatchingKindRule) NewInstance(path string, _ reflect.Kind, reference interface{}, _ string) (Rule, error) {
	return &MatchingKindRule{
		Reference: reference,
		Path:      path,
	}, nil
}

// NewInheritInstance instantiate new rule instance based on given rule.
func (m *MatchingKindRule) NewInheritInstance(
	path string, _ reflect.Kind, reference interface{},
) (Rule, error) {
	return &MatchingKindRule{
		Reference: reference,
		Path:      path,
	}, nil
}

// Validate applies rule on given patch operation specification.
func (m MatchingKindRule) Validate(operationSpec operation.Spec) error {
	if operationSpec.Value == nil {
		return nil
	}

	referenceValue := reflect.Zero(reflect.TypeOf(m.Reference))

	return m.deepCompareType(m.Path, referenceValue, reflect.ValueOf(operationSpec.Value), operationSpec.Operation)
}

// deepCompareType checks recursively one interface against a reference.
func (m MatchingKindRule) deepCompareType(
	path string, referenceValue, objectValue reflect.Value, definedOperation operation.Operation,
) error {
	var (
		err           error
		referenceType = referenceValue.Type()
		objectType    = objectValue.Type()
		referenceKind = referenceType.Kind()
		objectKind    = objectType.Kind()
	)

	if definedOperation == operation.AddOperation {
		if objectKind != reflect.Array && objectKind != reflect.Slice &&
			(referenceKind == reflect.Array || referenceKind == reflect.Slice) {
			referenceValueElem := reflect.Zero(referenceType.Elem())

			return m.deepCompareType(m.Path, referenceValueElem, objectValue, definedOperation)
		}
	}

	if referenceKind == reflect.Struct && objectKind == reflect.Map {
		return m.deepCompareMapWithStruct(path, referenceValue, objectValue, definedOperation)
	}

	if referenceKind != objectKind {
		return TypeMismatchError{name: path, actual: objectKind, expected: referenceKind}
	}

	switch objectType.Kind() { //nolint:exhaustive
	case reflect.Ptr:
		err = m.deepCompareType(path, reflect.Zero(referenceType.Elem()),
			reflect.Zero(objectType.Elem()), definedOperation)
	case reflect.Array, reflect.Map, reflect.Slice:
		referenceValueElem := reflect.Zero(referenceType.Elem())
		objectValueElem := reflect.Zero(objectType.Elem())
		err = m.deepCompareIterable(path, referenceValueElem, objectValueElem, definedOperation)
	case reflect.Struct:
		err = m.deepCompareStruct(referenceValue, objectValue, definedOperation)
	}

	return err
}

func (m MatchingKindRule) deepCompareIterable(
	path string, referenceValue, objectValue reflect.Value, definedOperation operation.Operation,
) error {
	var (
		referenceType = referenceValue.Type()
		objectType    = objectValue.Type()
	)

	if objectType.Kind() == reflect.Map && referenceType.Kind() == reflect.Map {
		if referenceType.Key().Kind() != objectType.Key().Kind() {
			return TypeMismatchError{
				name: path, actual: objectType.Key().Kind(), expected: referenceType.Key().Kind(), forKey: true,
			}
		}
	}

	return m.deepCompareType(path+"(item)", referenceValue, objectValue, definedOperation)
}

func (m MatchingKindRule) deepCompareStruct(
	referenceValue, objectValue reflect.Value, definedOperation operation.Operation,
) error {
	var (
		err           error
		referenceType = referenceValue.Type()
		objectType    = objectValue.Type()
	)

	for i := 0; i < objectType.NumField(); i++ {
		var (
			objectField = objectType.Field(i)
			objectName  = objectField.Name
			found       = false
		)

		for i := 0; i < referenceType.NumField(); i++ {
			var (
				referenceField = referenceType.Field(i)
				referenceName  = referenceField.Tag.Get("bson")
				zeroValue      = reflect.Zero(referenceField.Type)
			)

			if referenceField.Type.Kind() == reflect.Ptr {
				zeroValue = reflect.Zero(referenceField.Type.Elem())
			}

			if objectName == referenceName {
				err = m.deepCompareType(objectName, zeroValue,
					reflect.Zero(objectField.Type), definedOperation)

				found = true

				break
			}
		}

		if !found {
			err = UnknownFieldError{name: objectName}

			break
		}
	}

	return err
}

func (m MatchingKindRule) deepCompareMapWithStruct(
	path string, referenceValue, objectValue reflect.Value, definedOperation operation.Operation,
) error {
	var (
		err           error
		referenceType = referenceValue.Type()
		objectType    = objectValue.Type()
	)

	if objectType.Key().Kind() != reflect.String {
		return TypeMismatchError{
			name: path, actual: objectType.Key().Kind(), expected: referenceType.Key().Kind(), forKey: true,
		}
	}

	for _, key := range objectValue.MapKeys() {
		var (
			objectField = objectValue.MapIndex(key)
			found       = false
		)

		for i := 0; i < referenceType.NumField(); i++ {
			var (
				referenceField = referenceType.Field(i)
				referenceName  = referenceField.Tag.Get("bson")
				zeroValue      = reflect.Zero(referenceField.Type)
			)

			if referenceField.Type.Kind() == reflect.Ptr {
				zeroValue = reflect.Zero(referenceField.Type.Elem())
			}

			if key.String() == referenceName {
				err = m.deepCompareType(key.String(), zeroValue,
					reflect.Zero(objectField.Type()), definedOperation)

				found = true

				break
			}
		}

		if !found {
			err = UnknownFieldError{name: key.String()}

			break
		}
	}

	return err
}
