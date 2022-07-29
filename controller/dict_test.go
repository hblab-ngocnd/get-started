package controller

import (
	"context"
	"fmt"
	"testing"

	"github.com/hblab-ngocnd/get-started/services"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func InitTest() {
	err := godotenv.Load("../.env_test")
	if err != nil {
		panic(err)
	}
}
func TestDictHandler_ApiDict(t *testing.T) {
	InitTest()
	services.BucketSize = 2
	ctx := context.Background()
	data, err := getData(ctx, "https://japanesetest4you.com/jlpt-n1-vocabulary-list/")
	fmt.Printf("%+v", data)
	assert.Equal(t, nil, err)
}

func TestDictHandler_getDetail(t *testing.T) {
	detail, err := getDetail("https://japanesetest4you.com/flashcard/%e8%b5%a4%e5%ad%97-akaji/", 1)
	fmt.Println(detail)
	assert.Equal(t, nil, err)
}
