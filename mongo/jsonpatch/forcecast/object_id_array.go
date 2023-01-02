package forcecast

import (
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrNoArrayType           = errors.New("value is not valid array")
	ErrImpossibleCastToArray = errors.New("impossible to cast input to array")
)

type ObjectIDArrayCast struct{}

func (o ObjectIDArrayCast) ZeroValue() interface{} {
	return []primitive.ObjectID{emptyObjectID}
}

func (o ObjectIDArrayCast) Cast(input interface{}) (interface{}, error) {
	if input == nil {
		return nil, ErrNoArrayType
	}

	if kind := reflect.TypeOf(input).Kind(); kind != reflect.Array && kind != reflect.Slice {
		return nil, ErrNoArrayType
	}

	objectIDCast := ObjectIDCast{}
	arr := []primitive.ObjectID{}

	inputArr, ok := input.([]interface{})
	if !ok {
		return nil, ErrImpossibleCastToArray
	}

	for i, possibleObjectID := range inputArr {
		objectID, err := objectIDCast.Cast(possibleObjectID)
		if err != nil {
			return nil, fmt.Errorf("casting object id on index '%d' failed: %w", i, err)
		}

		arr = append(arr, objectID.(primitive.ObjectID)) //nolint:forcetypeassert
	}

	return arr, nil
}
