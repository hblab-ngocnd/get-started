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
}

func (d *dictionaryService) GetDictionary(ctx context.Context, url string) ([]models.Word, error) {
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
	var wg sync.WaitGroup
	var mu sync.Mutex
	mapWords := make(map[int]models.Word, len(targets)-5)
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
				detail, errDetail = d.getDetail(ctx, detailURL, id)
				if errDetail != nil {
					log.Println(errDetail)
				}
			} else {
				return
			}
			w := models.MakeWord(c, detail, id)
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
	for i := 0; i < len(mapWords)+5; i++ {
		if w, ok := mapWords[i]; ok {
			data = append(data, w)
		}
	}
	return data, nil
}

func (d *dictionaryService) getDetail(ctx context.Context, url string, i int) (string, error) {
	log.Println("start with goroutine ", i)
	defer log.Println("end with goroutine ", i)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
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