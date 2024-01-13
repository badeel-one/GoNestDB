package GoNestDB_test

import (
	"os"
	"testing"

	GoNestDB "badeel-one/GoNestDB"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnection(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "GoNestDB_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	conn := GoNestDB.NewConnection(tempDir)

	// Register a new collection
	err = conn.RegisterCollection("test")
	assert.NoError(t, err)

	// Get the registered collection
	collection, err := conn.GetCollection("test")
	assert.NoError(t, err)
	assert.NotNil(t, collection)

	// Close the collection
	err = conn.CloseCollection(collection)
	assert.NoError(t, err)
}

func TestNewConnection(t *testing.T) {
	conn, cleanup := getConnection(t)
	defer cleanup()

	assert.NotNil(t, conn)
}

func TestRegisterAndCloseCollection(t *testing.T) {
	conn, cleanup := getConnection(t)
	defer cleanup()

	coll, err := conn.NewCollection("test")
	require.NoError(t, err)

	assert.NotNil(t, coll)

	err = conn.CloseCollection(coll)
	require.NoError(t, err)
}
