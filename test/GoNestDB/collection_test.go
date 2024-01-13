package GoNestDB_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	GoNestDB "badeel-one/GoNestDB"
)

func TestCollectionInsertAndFindByID(t *testing.T) {
	conn, cleanup := getConnection(t)
	defer cleanup()

	coll, err := conn.NewCollection("test")
	require.NoError(t, err)
	defer conn.CloseCollection(coll)

	doc := GoNestDB.NewDocument(bson.M{
		"field1": "value1",
		"field2": 42,
	})

	err = coll.Insert(doc)
	require.NoError(t, err)

	foundDoc, err := coll.FindByID(doc.ID)
	require.NoError(t, err)

	assert.EqualValues(t, fmt.Sprintf("%v", doc), fmt.Sprintf("%v", foundDoc))
}

func TestCollectionUpdate(t *testing.T) {
	conn, cleanup := getConnection(t)
	defer cleanup()

	coll, err := conn.NewCollection("test")
	require.NoError(t, err)
	defer conn.CloseCollection(coll)

	doc := GoNestDB.NewDocument(bson.M{
		"field1": "value1",
		"field2": 42,
	})

	err = coll.Insert(doc)
	require.NoError(t, err)

	err = coll.Update(doc.ID, bson.M{
		"field1": "updatedValue1",
		"field3": true,
	})
	require.NoError(t, err)

	updatedDoc, err := coll.FindByID(doc.ID)
	require.NoError(t, err)

	assert.EqualValues(t, "updatedValue1", updatedDoc.Content["field1"])
	assert.EqualValues(t, 42, updatedDoc.Content["field2"])
	assert.EqualValues(t, true, updatedDoc.Content["field3"])
}

func TestCollectionDelete(t *testing.T) {
	conn, cleanup := getConnection(t)
	defer cleanup()

	coll, err := conn.NewCollection("test")
	require.NoError(t, err)
	defer conn.CloseCollection(coll)

	doc := GoNestDB.NewDocument(bson.M{
		"field1": "value1",
		"field2": 42,
	})

	err = coll.Insert(doc)
	require.NoError(t, err)

	err = coll.Delete(doc.ID)
	require.NoError(t, err)

	deletedDoc, err := coll.FindByID(doc.ID)
	require.Error(t, err)
	assert.Nil(t, deletedDoc)
}

func getConnection(t *testing.T) (*GoNestDB.Connection, func()) {
	tempDir, err := os.MkdirTemp("", "GoNestDB_test")
	require.NoError(t, err)

	conn := GoNestDB.NewConnection(tempDir)

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return conn, cleanup
}
