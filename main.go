package main

import (
	"errors"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/hblab-ngocnd/get-started/handlers"
	"github.com/hblab-ngocnd/get-started/infrastructure"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRegistry struct {
	templates map[string]*template.Template
}

// Render renders a template document
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}
	return tmpl.ExecuteTemplate(w, "base.html", data)
}

func main() {
	e := echo.New()
	templates := make(map[string]*template.Template)
	templates["home.html"] = template.Must(template.ParseFiles("public/views/home.html", "public/views/base.html"))
	templates["upload.html"] = template.Must(template.ParseFiles("public/views/upload.html", "public/views/base.html"))

	e.Renderer = &TemplateRegistry{
		templates: templates,
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	err := infrastructure.InitDB()
	if err != nil {
		panic(err)
	}
	e.Static("/static", "./static")
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "home.html", map[string]interface{}{"router": "home"})
	})
	e.GET("/upload", handlers.UploadFiles)
	e.POST("/api/visitors", handlers.CreateVisitor)
	e.GET("/api/visitors", handlers.ListVisitor)
	e.POST("/api/upload", handlers.ApiUpload)
	//When running on Cloud Foundry, get the PORT from the environment variable.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" //Local
	}
	e.Logger.Fatal(e.Start(":" + port))
}
