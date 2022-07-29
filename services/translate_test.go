package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/hblab-ngocnd/get-started/models"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func InitTest() {
	err := godotenv.Load("../.env_test")
	if err != nil {
		panic(err)
	}
}

func TestTranslateService_translateToVN(t *testing.T) {
	InitTest()
	ctx := context.Background()
	res := translateToVN(ctx, []string{"listed in public documentation", "test", "alone"})
	assert.NotNil(t, res)
}

func TestTranslateService_TranslateData(t *testing.T) {
	InitTest()
	ctx := context.Background()
	data := []models.Word{
		{
			Index:   3,
			MeanEng: "cache",
		},
		{
			Index:   4,
			MeanEng: "computer",
		},
	}
	res := translateData(ctx, data)
	fmt.Println(res)
	assert.NotNil(t, res)
}
