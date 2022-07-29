package models

import (
	"strings"

	"golang.org/x/net/html"
)

type Word struct {
	Index    int    `json:"index"`
	Text     string `json:"text"`
	Alphabet string `json:"alphabet"`
	MeanEng  string `json:"mean_eng"`
	MeanVN   string `json:"mean_vn"`
	Detail   string `json:"detail"`
}

func MakeWord(c *html.Node, detail string, index int) *Word {
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
