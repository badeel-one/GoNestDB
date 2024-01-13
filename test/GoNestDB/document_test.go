package GoNestDB_test

import (
	"os"
	"testing"

	GoNestDB "badeel-one/GoNestDB"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestDocument(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "GoNestDB_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	conn := GoNestDB.NewConnection(tempDir)
	err = conn.RegisterCollection("test")
	assert.NoError(t, err)

	collection, err := conn.GetCollection("test")
	assert.NoError(t, err)
	defer conn.CloseCollection(collection)

	document := GoNestDB.NewDocument(map[string]interface{}{"key": "value"})
	assert.NotNil(t, document)

	err = collection.Insert(document)
	assert.NoError(t, err)

	foundDocument, err := collection.FindByID(document.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundDocument)
	assert.Equal(t, document.ID, foundDocument.ID)
	assert.Equal(t, document.Content, foundDocument.Content)
}

func TestNewDocument(t *testing.T) {
	content := bson.M{
		"field1": "value1",
		"field2": 42,
	}

	doc := GoNestDB.NewDocument(content)

	assert.NotNil(t, doc.ID)
	assert.EqualValues(t, content, doc.Content)
}
