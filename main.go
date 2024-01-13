package main

import (
	"badeel-one/GoNestDB"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Employee struct {
	ID        string `json:"id" bson:"_id"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Email     string `json:"email" bson:"email"`
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// connection := GoNestDB.NewConnection("data")
	connection := GoNestDB.NewConnection()
	_ = connection.RegisterCollection("employees")

	employeeGroup := e.Group("/employees")
	generateCRUDRoutes(employeeGroup, connection, "employees", Employee{})

	e.Start("localhost:8080")
}

func generateCRUDRoutes(group *echo.Group, connection *GoNestDB.Connection, collectionName string, obj interface{}) {
	collection, _ := connection.GetCollection(collectionName)
	objType := reflect.TypeOf(obj)

	group.GET("", func(c echo.Context) error {
		queryFilter := bson.M{}
		if err := c.Bind(&queryFilter); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
		}
		// This part requires iterating over all the documents in the collection and applying the filter.
		// The current package doesn't support this functionality. You can add this feature in the future.
		return echo.NewHTTPError(http.StatusNotImplemented, "Get list not implemented")
	})

	group.GET("", func(c echo.Context) error {
		query := make(primitive.M)
		for key, values := range c.QueryParams() {
			query[key] = values[0]
		}

		documents, err := collection.Find(query)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
		}
		list := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(objType)), 0, 0)
		for _, document := range documents {
			obj := reflect.New(objType).Interface()
			data, _ := bson.Marshal(document.Content)
			bson.Unmarshal(data, obj)
			list = reflect.Append(list, reflect.ValueOf(obj))
		}

		return c.JSON(http.StatusOK, list.Interface())
	})

	group.GET("/:id", func(c echo.Context) error {
		id := c.Param("id")
		document, err := collection.FindByID(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Document not found")
		}
		obj := reflect.New(objType).Interface()
		data, _ := bson.Marshal(document.Content)
		bson.Unmarshal(data, obj)
		return c.JSON(http.StatusOK, obj)
	})

	group.POST("", func(c echo.Context) error {
		obj := reflect.New(objType).Interface()
		if err := c.Bind(obj); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
		}
		content, _ := bson.Marshal(obj)
		var doc primitive.M
		bson.Unmarshal(content, &doc)
		document := GoNestDB.NewDocument(doc)
		err := collection.Insert(document)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error creating document")
		}
		return c.JSON(http.StatusCreated, obj)
	})

	group.PUT("/:id", func(c echo.Context) error {
		id := c.Param("id")
		update := bson.M{}
		if err := c.Bind(&update); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
		}
		err := collection.Update(id, update)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error updating document")
		}
		document, err := collection.FindByID(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error fetching updated document")
		}
		obj := reflect.New(objType).Interface()
		data, _ := bson.Marshal(document.Content)
		bson.Unmarshal(data, obj)
		return c.JSON(http.StatusOK, obj)
	})

	group.DELETE("/:id", func(c echo.Context) error {
		id := c.Param("id")
		err := collection.Delete(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error deleting document")
		}
		return c.NoContent(http.StatusNoContent)
	})
}
