package GoNestDB

import (
	"os"
	"path/filepath"
)

type Connection struct {
	basePath string
}

func NewConnection(args ...string) *Connection {
	basePath := ""
	if len(args) > 0 {
		basePath = args[0]
	} else {
		if os.Getenv("ENVIRONMENT") == "production" {
			basePath = filepath.Join(GetRoamingDataPath(), "data")
		} else {
			basePath = "data"
		}
	}
	// Ensure that the folder path exists
	err := os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return &Connection{
		basePath: basePath,
	}
}

func (c *Connection) RegisterCollection(name string) error {
	collection, err := c.NewCollection(name)
	if err != nil {
		return err
	}
	_ = collection.db.Close()
	return nil
}

func (c *Connection) GetCollection(name string) (*Collection, error) {
	return c.NewCollection(name)
}

func (c *Connection) CloseCollection(collection *Collection) error {
	return collection.db.Close()
}
