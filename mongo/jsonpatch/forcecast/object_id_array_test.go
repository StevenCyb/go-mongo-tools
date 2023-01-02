package forcecast

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestObjectIDArray(t *testing.T) {
	t.Parallel()

	var objectIDArrInterface interface{}

	forceCast := ObjectIDArrayCast{}
	emptyObjectID := []primitive.ObjectID{{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}
	objectID := primitive.NewObjectID()
	objectIDArrInterface = []primitive.ObjectID{objectID}
	objectIDSlice := []interface{}{objectID}
	objectIDSliceWithHex := []interface{}{objectID.Hex()}

	require.Equal(t, emptyObjectID, forceCast.ZeroValue())

	castedObjectID, err := forceCast.Cast(objectIDSlice)
	require.NoError(t, err)
	require.Equal(t, objectIDArrInterface, castedObjectID)

	castedObjectID, err = forceCast.Cast(objectIDSliceWithHex)
	require.NoError(t, err)
	require.Equal(t, objectIDArrInterface, castedObjectID)
}
