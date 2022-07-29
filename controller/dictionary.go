package controller

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/hblab-ngocnd/get-started/models"
	"github.com/hblab-ngocnd/get-started/services"
	"github.com/labstack/echo/v4"
)

type dictHandler struct {
	translateService  services.TranslateService
	dictionaryService services.DictionaryService
	cacheData         []models.Word
	mu                sync.Mutex
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
	notCache := c.QueryParam("not_cache")
	if notCache != "true" && f.cacheData != nil {
		return c.JSON(http.StatusOK, f.cacheData)
	}
	ctx := context.Background()
	data, err := f.dictionaryService.GetDictionary(ctx, "https://japanesetest4you.com/jlpt-n1-vocabulary-list/")
	if err != nil {
		log.Fatal(err)
	}
	data = f.translateService.TranslateData(ctx, data)
	f.mu.Lock()
	f.cacheData = data
	f.mu.Unlock()
	return c.JSON(http.StatusOK, data)
}
