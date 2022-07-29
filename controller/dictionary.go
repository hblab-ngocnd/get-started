package controller

import (
	"context"
	"log"
	"net/http"

	"github.com/hblab-ngocnd/get-started/services"
	"github.com/labstack/echo/v4"
)

type dictHandler struct {
	translateService  services.TranslateService
	dictionaryService services.DictionaryService
}

func NewDictHandler() *dictHandler {
	return &dictHandler{
		translateService:  services.NewTranslate(),
		dictionaryService: services.NewDictionary(),
	}
}

func (f *dictHandler) Dict(c echo.Context) error {
	return c.Render(http.StatusOK, "dictionary.html", map[string]interface{}{"router": "dictionary"})
}

func (f *dictHandler) ApiDict(c echo.Context) error {
	ctx := context.Background()
	data, err := f.dictionaryService.GetDictionary(ctx, "https://japanesetest4you.com/jlpt-n1-vocabulary-list/")
	if err != nil {
		log.Fatal(err)
	}
	return c.JSON(http.StatusOK, f.translateService.TranslateData(ctx, data))
}
