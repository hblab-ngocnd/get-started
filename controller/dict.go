package controller

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/translate"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/html"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

type Word struct {
	Index    int    `json:"index"`
	Text     string `json:"text"`
	Alphabet string `json:"alphabet"`
	MeanEng  string `json:"mean_eng"`
	MeanVN   string `json:"mean_vn"`
	Detail   string `json:"detail"`
}

type dictHandler struct {
	dictService interface{}
}

type Result struct {
	data []Word
	mu   sync.Mutex
}

func NewDictHandler() *dictHandler {
	return &dictHandler{
		dictService: nil,
	}
}

func (f *dictHandler) Dict(c echo.Context) error {
	return c.Render(http.StatusOK, "dictionary.html", map[string]interface{}{"router": "dictionary"})
}

func (f *dictHandler) ApiDict(c echo.Context) error {
	data, err := getData("https://japanesetest4you.com/jlpt-n1-vocabulary-list/")
	if err != nil {
		log.Fatal(err)
	}
	return c.JSON(http.StatusOK, data)
}

func getData(url string) ([]Word, error) {
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
	tag := getElementByClass(doc, "entry clearfix")
	targets := getListElementByTag(tag, "p")
	var wg sync.WaitGroup
	var result Result
	for i, target := range targets {
		id := i
		if os.Getenv("DEBUG") == "true" && i == 10 {
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
			if detailURL, ok := getAttribute(c, "href"); ok {
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
	return translateData(result.data), nil
}

func makeWord(c *html.Node, detail string, index int) *Word {
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
	return &Word{
		Index:    index,
		Text:     text,
		Alphabet: alphabet,
		MeanEng:  mean,
		Detail:   detail,
	}
}

var bucketSize = 100

func translateData(data []Word) []Word {
	mapData := make(map[int]Word, len(data))
	maxIdx := 0
	for _, w := range data {
		if w.Index > maxIdx {
			maxIdx = w.Index
		}
		mapData[w.Index] = w
	}
	trans := make([]string, maxIdx+1)
	for i := 0; i <= maxIdx; i++ {
		if d, ok := mapData[i]; ok {
			trans[i] = d.MeanEng
		}
	}
	translated := make([]string, 0, len(trans))
	var wg sync.WaitGroup
	var mu sync.Mutex
	transMap := make(map[int][]string)
	for i := 0; i < len(trans); {
		e := i + bucketSize
		if e > len(trans) {
			e = len(trans)
		}
		start := i
		end := e
		wg.Add(1)
		go func() {
			defer wg.Done()
			bulk := translateToVN(trans[start:end])
			mu.Lock()
			transMap[start] = bulk
			mu.Unlock()
		}()
		i = e
	}
	wg.Wait()
	for i := 0; i < len(trans); i = i + bucketSize {
		if arr, ok := transMap[i]; ok {
			translated = append(translated, arr...)
		}
	}
	for i, vn := range translated {
		if v, ok := mapData[i]; ok {
			v.MeanVN = vn
			mapData[i] = v
		}
	}
	result := make([]Word, 0, len(mapData))
	for _, m := range mapData {
		result = append(result, m)
	}
	return result
}

func translateToVN(text []string) []string {
	log.Println("start translate")
	defer log.Println("end translate")
	apiKey := os.Getenv("GOOGLE_APPLICATION_API_KEY")
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	lang, _ := language.Parse("vi")
	client, err := translate.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Println(err)
	}
	defer client.Close()
	resp, err := client.Translate(ctx, text, lang, nil)
	if err != nil {
		log.Println(fmt.Errorf("Translate: %v", err))
		return []string{""}
	}
	if len(resp) == 0 {
		log.Println(fmt.Errorf("Translate returned empty response to text: %s", text))
		return []string{""}
	}
	result := make([]string, len(resp))
	for i, res := range resp {
		result[i] = res.Text
	}
	return result
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
	tag := getElementByClass(doc, "entry clearfix")
	if tag != nil {
		nodes := getListElementByTag(tag, "p")
		data = []string{renderNode(nodes[1])}
		for _, node := range nodes[3:] {
			if node.FirstChild != nil && node.FirstChild.Data == "img" {
				continue
			}
			data = append(data, renderNode(node))
		}
	}
	return strings.Join(data, ""), nil
}

func getListElementByTag(n *html.Node, tag string) []*html.Node {
	var result []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Data == tag {
			result = append(result, c)
		}
	}
	return result
}

func getAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)

	err := html.Render(w, n)

	if err != nil {
		return ""
	}
	return buf.String()
}

// nolint:unused // This function used next turn
func checkId(n *html.Node, id string) bool {
	if n.Type == html.ElementNode {
		s, ok := getAttribute(n, "id")
		if ok && s == id {
			return true
		}
	}
	return false
}

func checkClass(n *html.Node, class string) bool {
	if n.Type == html.ElementNode {
		s, ok := getAttribute(n, "class")
		if ok && strings.Contains(s, class) {
			return true
		}
	}
	return false
}

func traverse(n *html.Node, id string, fn func(node *html.Node, id string) bool) *html.Node {
	if fn(n, id) {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		res := traverse(c, id, fn)
		if res != nil {
			return res
		}
	}
	return nil
}

// nolint:deadcode,unused // This function used next turn
func getElementById(n *html.Node, id string) *html.Node {
	return traverse(n, id, checkId)
}

func getElementByClass(n *html.Node, class string) *html.Node {
	return traverse(n, class, checkClass)
}
