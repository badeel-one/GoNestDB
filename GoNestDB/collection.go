package GoNestDB

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"go.mongodb.org/mongo-driver/bson"
)

type Collection struct {
	db   *badger.DB
	name string
}

func (c *Connection) NewCollection(name string) (*Collection, error) {
	collectionPath := fmt.Sprintf("%s/%s", c.basePath, name)
	opts := badger.DefaultOptions(collectionPath)
	opts = opts.WithLogger(nil) // Disable logging
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &Collection{
		db:   db,
		name: name,
	}, nil
}

func (c *Collection) generateKey(documentID string) []byte {
	return []byte(documentID)
}

func (c *Collection) Insert(document *Document) error {
	key := c.generateKey(document.ID)

	data, err := bson.Marshal(document)
	if err != nil {
		return err
	}

	err = c.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, data)
	})

	return err
}

func (c *Collection) FindByID(documentID string) (*Document, error) {
	key := c.generateKey(documentID)
	var document *Document

	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		data, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		document = &Document{}
		err = bson.Unmarshal(data, document)
		return err
	})

	if err != nil {
		return nil, err
	}
	return document, nil
}

func (c *Collection) Update(documentID string, update bson.M) error {
	document, err := c.FindByID(documentID)
	if err != nil {
		return err
	}

	for k, v := range update {
		document.Content[k] = v
	}

	return c.Insert(document)
}

func (c *Collection) Delete(documentID string) error {
	err := c.db.Update(func(txn *badger.Txn) error {
		key := c.generateKey(documentID)
		return txn.Delete(key)
	})

	return err
}

func (c *Collection) Find(query bson.M) ([]*Document, error) {
	matchingDocument := make(map[string]*Document)
	matchingDocumentIDs := make(map[string]int)
	var documents []*Document
	queryLen := len(query)

	// No index found, perform a full collection scan
	err := c.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			data, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			var document Document
			err = bson.Unmarshal(data, &document)
			if err != nil {
				return err
			}

			if queryLen == 0 {
				documents = append(documents, &document)
				continue
			}

			for field, value := range query {
				if Compare(document.Content[field], value) {
					matchingDocument[document.ID] = &document
					matchingDocumentIDs[document.ID]++
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if queryLen > 0 {
		for docID, count := range matchingDocumentIDs {
			if count == queryLen {
				documents = append(documents, matchingDocument[docID])
			}
		}
	}

	return documents, nil
}
