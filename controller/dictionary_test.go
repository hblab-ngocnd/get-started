package controller

import (
	"github.com/joho/godotenv"
)

func InitTest() {
	err := godotenv.Load("../.env_test")
	if err != nil {
		panic(err)
	}
}
