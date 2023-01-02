package forcecast

import (
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrInvalidTypeForObjectID = errors.New("ObjectID must be string (24 characters) or array (12 bytes)")

const (
	objectIDStringLen = 24
	objectIDArrayLen  = 12
)

var emptyObjectID = primitive.ObjectID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} //nolint:gochecknoglobals

type ObjectIDCast struct{}

func (o ObjectIDCast) ZeroValue() interface{} {
	return emptyObjectID
}

func (o ObjectIDCast) Cast(input interface{}) (interface{}, error) {
	if input == nil {
		return nil, ErrInvalidTypeForObjectID
	}

	value := reflect.ValueOf(input)

	switch reflect.TypeOf(input).Kind() { //nolint:exhaustive
	case reflect.String:
		if value.Len() != objectIDStringLen {
			return nil, ErrInvalidTypeForObjectID
		}

		objectID, err := primitive.ObjectIDFromHex(input.(string))
		if err != nil {
			return nil, ErrInvalidTypeForObjectID
		}

		return objectID, nil
	case reflect.Array, reflect.Slice:
		if value.Len() != objectIDArrayLen {
			return nil, ErrInvalidTypeForObjectID
		}

		if value.Type().Elem().Kind() != reflect.Uint8 || value.Kind() == reflect.Slice {
			byteArr := primitive.ObjectID{}

			for i := 0; i < value.Len(); i++ {
				byteArr[i] = uint8(value.Index(i).Int())
			}

			input = byteArr
		}

		objectID, ok := input.(primitive.ObjectID)
		if !ok {
			return nil, ErrInvalidTypeForObjectID
		}

		return objectID, nil
	}

	return nil, ErrInvalidTypeForObjectID
}
