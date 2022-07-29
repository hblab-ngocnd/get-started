package controller

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/hblab-ngocnd/get-started/helpers"
	"github.com/hblab-ngocnd/get-started/models"
	"github.com/hblab-ngocnd/get-started/services"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/html"
)

type dictHandler struct {
	translateService services.TranslateService
}

type Result struct {
	data []models.Word
	mu   sync.Mutex
}

func NewDictHandler() *dictHandler {
	return &dictHandler{
		translateService: services.NewTranslate(),
	}
}

func (f *dictHandler) Dict(c echo.Context) error {
	return c.Render(http.StatusOK, "dictionary.html", map[string]interface{}{"router": "dictionary"})
}

func (f *dictHandler) ApiDict(c echo.Context) error {
	ctx := context.Background()
	data, err := getData(ctx, "https://japanesetest4you.com/jlpt-n1-vocabulary-list/")
	if err != nil {
		log.Fatal(err)
	}
	return c.JSON(http.StatusOK, f.translateService.TranslateData(ctx, data))
}

func getData(ctx context.Context, url string) ([]models.Word, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	tag := helpers.GetElementByClass(doc, "entry clearfix")
	targets := helpers.GetListElementByTag(tag, "p")
	var wg sync.WaitGroup
	var result Result
	for i, target := range targets {
		id := i
		if os.Getenv("DEBUG") == "true" && i == 70 {
			break
		}
		tar := target
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := tar.FirstChild
			if c == nil {
				return
			}
			var detail string
			var errDetail error
			if detailURL, ok := helpers.GetAttribute(c, "href"); ok {
				detail, errDetail = getDetail(detailURL, id)
				if errDetail != nil {
					log.Println(errDetail)
				}
			} else {
				return
			}
			w := makeWord(c, detail, id)
			if w == nil {
				return
			}
			result.mu.Lock()
			result.data = append(result.data, *w)
			result.mu.Unlock()
		}()
	}
	wg.Wait()
	log.Println("clone done")
	return result.data, nil
}

func makeWord(c *html.Node, detail string, index int) *models.Word {
	if c.FirstChild == nil {
		return nil
	}
	idx := strings.Index(c.FirstChild.Data, ":")
	mean := c.FirstChild.Data[idx+1:]
	arr := strings.Split(c.FirstChild.Data[:idx], " ")
	text := arr[0]
	var alphabet string
	if len(arr) > 1 {
		alphabet = strings.TrimRight(strings.TrimLeft(strings.Join(arr[1:], " "), "("), ")")
	}
	return &models.Word{
		Index:    index,
		Text:     text,
		Alphabet: alphabet,
		MeanEng:  mean,
		Detail:   detail,
	}
}

func getDetail(url string, i int) (string, error) {
	log.Println("start with goroutine ", i)
	defer log.Println("end with goroutine ", i)
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return "", err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	var data []string
	tag := helpers.GetElementByClass(doc, "entry clearfix")
	if tag != nil {
		nodes := helpers.GetListElementByTag(tag, "p")
		data = []string{helpers.RenderNode(nodes[1])}
		for _, node := range nodes[3:] {
			if node.FirstChild != nil && node.FirstChild.Data == "img" {
				continue
			}
			data = append(data, helpers.RenderNode(node))
		}
	}
	return strings.Join(data, ""), nil
}
