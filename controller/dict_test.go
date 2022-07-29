package controller

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func InitTest() {
	err := godotenv.Load("../.env-test")
	if err != nil {
		panic(err)
	}
}
func TestDictHandler_ApiDict(t *testing.T) {
	InitTest()
	bucketSize = 2
	data, err := getData("https://japanesetest4you.com/jlpt-n1-vocabulary-list/")
	fmt.Printf("%+v", data)
	assert.Equal(t, nil, err)
}

func TestDictHandler_getDetail(t *testing.T) {
	detail, err := getDetail("https://japanesetest4you.com/flashcard/%e8%b5%a4%e5%ad%97-akaji/", 1)
	fmt.Println(detail)
	assert.Equal(t, nil, err)
}

func TestDictHandler_translateToVN(t *testing.T) {
	InitTest()
	res := translateToVN([]string{"listed in public documentation", "test", "alone"})
	assert.NotNil(t, res)
}

func TestDictHandler_translateData(t *testing.T) {
	InitTest()
	data := []Word{
		{
			Index:   3,
			MeanEng: "cache",
		},
		{
			Index:   4,
			MeanEng: "computer",
		},
	}
	res := translateData(data)
	fmt.Println(res)
}
