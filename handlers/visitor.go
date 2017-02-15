package handlers

import (
	"net/http"

	"github.com/hblab-ngocnd/get-started/infrastructure"
	"github.com/hblab-ngocnd/get-started/models"
	"github.com/labstack/echo/v4"
	"github.com/timjacobi/go-couchdb"
)

type alldocsResult struct {
	TotalRows int `json:"total_rows"`
	Offset    int
	Rows      []map[string]interface{}
}

func CreateVisitor(c echo.Context) error {
	var visitor models.Visitor
	if c.Bind(&visitor) == nil {
		infrastructure.GetDB().Post(visitor)
		c.String(200, "Hello "+visitor.Name)
	}
	return nil
}
func ListVisitor(c echo.Context) error {
	var result alldocsResult
	err := infrastructure.GetDB().AllDocs(&result, couchdb.Options{"include_docs": true})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "unable to fetch docs"})
	}
	return c.JSON(200, result.Rows)
}
