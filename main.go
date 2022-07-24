package main

import (
	"os"

	"github.com/hblab-ngocnd/get-started/infrastructure"
)

func main() {
	err := infrastructure.InitDB()
	if err != nil {
		panic(err)
	}
	e := infrastructure.SetupServer()
	//When running on Cloud Foundry, get the PORT from the environment variable.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" //Local
	}
	e.Logger.Fatal(e.Start(":" + port))
}
