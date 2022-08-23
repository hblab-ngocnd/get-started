package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/hblab-ngocnd/get-started/models"
	"github.com/hblab-ngocnd/get-started/services"
)

var PermissionDeniedErr = errors.New("usecase: permission denied")
var InvalidErr = errors.New("usecase: invalid")

type dictUseCase struct {
	translateService  services.TranslateService
	dictionaryService services.DictionaryService
	cacheData         map[string][]models.Word
	mu                sync.Mutex
}

func NewDictUseCase() *dictUseCase {
	return &dictUseCase{
		translateService:  services.NewTranslate(),
		dictionaryService: services.NewDictionary(),
	}
}

type DictUseCase interface {
	GetDict(context.Context, int, int, string, string, string) ([]models.Word, error)
	GetDetail(context.Context, string, int) (*string, error)
}

func (u *dictUseCase) GetDict(ctx context.Context, start, pageSize int, notCache, level, pwd string) ([]models.Word, error) {
	if notCache != "true" && u.cacheData != nil && u.cacheData[level] != nil {
		u.mu.Lock()
		defer u.mu.Unlock()
		if start > len(u.cacheData[level]) {
			start = len(u.cacheData[level])
		}
		end := start + pageSize
		if end > len(u.cacheData[level]) {
			end = len(u.cacheData[level])
		}
		return u.cacheData[level][start:end], nil
	}
	if notCache == "true" {
		if len(strings.TrimSpace(pwd)) == 0 || pwd != os.Getenv("SYNC_PASS") {
			return nil, PermissionDeniedErr
		}
	}
	u.mu.Lock()
	defer u.mu.Unlock()
	url := fmt.Sprintf("https://japanesetest4you.com/jlpt-%s-vocabulary-list/", level)
	data, err := u.dictionaryService.GetDictionary(ctx, url)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	data = u.translateService.TranslateData(ctx, data)
	if u.cacheData == nil {
		u.cacheData = make(map[string][]models.Word)
	}
	u.cacheData[level] = data
	if start > len(data) {
		start = len(data)
	}
	end := start + pageSize
	if end > len(data) {
		end = len(data)
	}
	return data[start:end], nil
}

func (u *dictUseCase) GetDetail(ctx context.Context, level string, index int) (*string, error) {
	if u.cacheData[level] == nil || index >= len(u.cacheData[level]) {
		return nil, InvalidErr
	}
	detailURL := u.cacheData[level][index].Link
	if strings.TrimSpace(detailURL) == "" {
		return nil, nil
	}
	data, err := u.dictionaryService.GetDetail(ctx, detailURL, index)
	if err != nil {
		return nil, err
	}
	u.cacheData[level][index].Detail = data
	return &data, nil
}
