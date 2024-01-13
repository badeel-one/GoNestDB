package GoNestDB

import (
	"crypto/rand"

	"github.com/mr-tron/base58"
	"go.mongodb.org/mongo-driver/bson"
)

type Document struct {
	ID      string `bson:"id"`
	Content bson.M `bson:",inline"`
}

func NewDocument(content bson.M) *Document {
	id, ok := content["_id"].(string)
	if !ok {
		idBytes := make([]byte, 12)
		rand.Read(idBytes)
		id = base58.Encode(idBytes)
		content["_id"] = id
	}

	return &Document{
		ID:      id,
		Content: content,
	}
}
