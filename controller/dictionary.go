package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/hblab-ngocnd/get-started/models"
	"github.com/hblab-ngocnd/get-started/services"
	"github.com/labstack/echo/v4"
)

type dictHandler struct {
	translateService  services.TranslateService
	dictionaryService services.DictionaryService
	cacheData         map[string][]models.Word
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
	level := c.QueryParam("level")
	start, err := strconv.Atoi(c.QueryParam("start"))
	var end int
	if err != nil {
		start = 0
	}
	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil {
		pageSize = 20
	}
	if level == "" {
		level = "n1"
	}
	switch level {
	case "n1", "n2", "n3", "n4", "n5":
	default:
		return c.NoContent(http.StatusBadRequest)
	}
	if notCache != "true" && f.cacheData != nil && f.cacheData[level] != nil {
		f.mu.Lock()
		defer f.mu.Unlock()
		if start > len(f.cacheData[level]) {
			start = len(f.cacheData[level])
		}
		end = start + pageSize
		if end > len(f.cacheData[level]) {
			end = len(f.cacheData[level])
		}
		return c.JSON(http.StatusOK, f.cacheData[level][start:end])
	}
	if notCache == "true" {
		pwd := c.QueryParam("password")
		if len(strings.TrimSpace(pwd)) == 0 || pwd != os.Getenv("SYNC_PASS") {
			return c.NoContent(http.StatusBadRequest)
		}
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	ctx := context.Background()
	url := fmt.Sprintf("https://japanesetest4you.com/jlpt-%s-vocabulary-list/", level)
	data, err := f.dictionaryService.GetDictionary(ctx, url)
	if err != nil {
		log.Fatal(err)
	}
	data = f.translateService.TranslateData(ctx, data)
	if f.cacheData == nil {
		f.cacheData = make(map[string][]models.Word)
	}
	f.cacheData[level] = data
	if start > len(data) {
		start = len(data)
	}
	end = start + pageSize
	if end > len(data) {
		end = len(data)
	}
	return c.JSON(http.StatusOK, data[start:end])
}

func (f *dictHandler) ApiGetDetail(c echo.Context) error {
	index, err := strconv.Atoi(c.Param("index"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	level := c.QueryParam("level")
	switch level {
	case "n1", "n2", "n3", "n4", "n5":
	default:
		return c.NoContent(http.StatusBadRequest)
	}
	if f.cacheData[level] == nil && index >= len(f.cacheData[level]) {
		return c.NoContent(http.StatusBadRequest)
	}
	detailURL := f.cacheData[level][index].Link
	if strings.TrimSpace(detailURL) == "" {
		return c.String(http.StatusOK, "")
	}
	ctx := context.Background()
	data, err := f.dictionaryService.GetDetail(ctx, detailURL, index)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	f.cacheData[level][index].Detail = data
	return c.String(http.StatusOK, data)
}
