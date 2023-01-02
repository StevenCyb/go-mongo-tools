package forcecast

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestObjectID(t *testing.T) {
	t.Parallel()

	forceCast := ObjectIDCast{}
	emptyObjectID := primitive.ObjectID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	objectID := primitive.NewObjectID()
	objectIDHex := objectID.Hex()
	objectIDSlice := []int{}
	objectIDArray := [12]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	for i, value := range objectID {
		objectIDArray[i] = int(value)
		objectIDSlice = append(objectIDSlice, int(value))
	}

	require.Equal(t, emptyObjectID, forceCast.ZeroValue())

	castedObjectID, err := forceCast.Cast(objectID)
	require.NoError(t, err)
	require.Equal(t, objectID, castedObjectID)

	castedObjectID, err = forceCast.Cast(objectIDHex)
	require.NoError(t, err)
	require.Equal(t, objectID, castedObjectID)

	castedObjectID, err = forceCast.Cast(objectIDArray)
	require.NoError(t, err)
	require.Equal(t, objectID, castedObjectID)

	castedObjectID, err = forceCast.Cast(objectIDSlice)
	require.NoError(t, err)
	require.Equal(t, objectID, castedObjectID)
}
