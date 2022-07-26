package services

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hblab-ngocnd/get-started/helpers"
	"github.com/hblab-ngocnd/get-started/models"
	"golang.org/x/net/html"
)

type dictionaryService struct {
}

func NewDictionary() *dictionaryService {
	return &dictionaryService{}
}

type DictionaryService interface {
	GetDictionary(context.Context, string) ([]models.Word, error)
	GetDetail(context.Context, string, int) (string, error)
}

func (d *dictionaryService) GetDictionary(ctx context.Context, url string) ([]models.Word, error) {
	ctx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	tag := helpers.GetElementByClass(doc, "entry clearfix")
	targets := helpers.GetListElementByTag(tag, "p")
	if len(targets) > 2 {
		targets = targets[2:]
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	mapWords := make(map[int]models.Word, len(targets))
	for i, target := range targets {
		id := i
		if os.Getenv("DEBUG") == "true" && i == 40 {
			break
		}
		tar := target
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := tar.FirstChild
			if c == nil {
				c = tar
			}
			var detail string
			var errDetail error
			detailURL, ok := helpers.GetAttribute(c, "href")
			if ok {
				detail, errDetail = d.getDetail(ctx, detailURL, id)
				if errDetail != nil {
					log.Println(errDetail)
				}
			}
			w := models.MakeWord(c, detailURL, detail, id)
			if w == nil {
				return
			}
			mu.Lock()
			mapWords[id] = *w
			mu.Unlock()
		}()
	}
	wg.Wait()
	log.Println("clone done")
	data := make([]models.Word, 0, len(mapWords))
	for i := 0; i < len(mapWords); i++ {
		if w, ok := mapWords[i]; ok {
			data = append(data, w)
		}
	}
	return data, nil
}
func (d *dictionaryService) GetDetail(ctx context.Context, url string, i int) (string, error) {
	return d.getDetail(ctx, url, i)
}
func (d *dictionaryService) getDetail(ctx context.Context, url string, i int) (string, error) {
	log.Println("start with goroutine ", i)
	defer log.Println("end with goroutine ", i)
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return "", err
	}
	client := http.DefaultClient

	res, err := client.Do(req)
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
