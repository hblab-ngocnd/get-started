package main

import (
	"os"

	"github.com/hblab-ngocnd/get-started/handlers"
	"github.com/hblab-ngocnd/get-started/infrastructure"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	err := infrastructure.InitDB()
	if err != nil {
		panic(err)
	}
	e.Static("/static", "./static")
	e.GET("/", func(c echo.Context) error {
		return c.File("static/index.html")
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
