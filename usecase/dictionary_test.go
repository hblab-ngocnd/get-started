package usecase

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hblab-ngocnd/get-started/models"
	"github.com/hblab-ngocnd/get-started/services"
	"github.com/hblab-ngocnd/get-started/services/mock_services"
	"github.com/stretchr/testify/assert"
)

func Test_GetDetail(t *testing.T) {
	patterns := []struct {
		description             string
		start                   int
		pageSize                int
		notCache                string
		level                   string
		pwd                     string
		newMockDictService      func(ctrl *gomock.Controller) services.DictionaryService
		newMockTranslateService func(ctrl *gomock.Controller) services.TranslateService
		expect                  []models.Word
		err                     error
	}{
		{
			description: "success",
			start:       0,
			pageSize:    1,
			notCache:    "true",
			level:       "n1",
			pwd:         "sync_pass",
			newMockDictService: func(ctrl *gomock.Controller) services.DictionaryService {
				mock := mock_services.NewMockDictionaryService(ctrl)
				mock.EXPECT().GetDictionary(gomock.Any(), gomock.Eq("https://japanesetest4you.com/jlpt-n1-vocabulary-list/")).Return([]models.Word{
					{}, {}, {}, {},
				}, nil)
				return mock
			},
			newMockTranslateService: func(ctrl *gomock.Controller) services.TranslateService {
				mock := mock_services.NewMockTranslateService(ctrl)
				mock.EXPECT().TranslateData(gomock.Any(), gomock.Any()).Return([]models.Word{
					{}, {}, {}, {},
				})
				return mock
			},
			expect: []models.Word{
				{},
			},
			err: nil,
		},
		{
			description: "success over size",
			start:       0,
			pageSize:    8,
			notCache:    "true",
			level:       "n1",
			pwd:         "sync_pass",
			newMockDictService: func(ctrl *gomock.Controller) services.DictionaryService {
				mock := mock_services.NewMockDictionaryService(ctrl)
				mock.EXPECT().GetDictionary(gomock.Any(), gomock.Eq("https://japanesetest4you.com/jlpt-n1-vocabulary-list/")).Return([]models.Word{
					{}, {}, {}, {},
				}, nil)
				return mock
			},
			newMockTranslateService: func(ctrl *gomock.Controller) services.TranslateService {
				mock := mock_services.NewMockTranslateService(ctrl)
				mock.EXPECT().TranslateData(gomock.Any(), gomock.Any()).Return([]models.Word{
					{}, {}, {}, {},
				})
				return mock
			},
			expect: []models.Word{
				{}, {}, {}, {},
			},
			err: nil,
		},
		{
			description: "permission denied",
			start:       0,
			pageSize:    8,
			notCache:    "true",
			level:       "n1",
			pwd:         "1111",
			expect:      nil,
			err:         PermissionDeniedErr,
		},
	}

	for i, p := range patterns {
		t.Run(fmt.Sprintf("%d:%s", i, p.description), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var mockDict services.DictionaryService
			if p.newMockDictService != nil {
				mockDict = p.newMockDictService(ctrl)
			}
			var mockTrans services.TranslateService
			if p.newMockDictService != nil {
				mockTrans = p.newMockTranslateService(ctrl)
			}
			uc := dictUseCase{
				dictionaryService: mockDict,
				translateService:  mockTrans,
			}
			ctx := context.Background()
			os.Setenv("SYNC_PASS", "sync_pass")
			actual, err := uc.GetDict(ctx, p.start, p.pageSize, p.notCache, p.level, p.pwd)
			assert.Equal(t, p.expect, actual)
			assert.Equal(t, p.err, err)
		})
	}
}