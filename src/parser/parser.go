package parser

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
)

type parser interface {
	GetDropList(content string) map[string]string
}

func GetDropList(content string) map[string]string {
	m := make(map[string]string)
	r := bytes.NewReader([]byte(content))
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatal(err)
	}
	sections := doc.Find(".menu-drop-list")
	count := sections.Length()
	fmt.Println(count)
	sections.Each(func(i int, s *goquery.Selection) {
		if i == 1 {
			// For each item found, get the band and title
			s.Find("a").Each(func(p int, selection *goquery.Selection) {
				category := selection.Text()
				category = strings.Replace(category, "/", "", -1)
				href := selection.AttrOr("href", "")
				m[category] = href
				fmt.Printf("童书分类 %d: %s - %s\n", p, category, href)
			})
		}
	})
	return m
}

func GetNextPageUrl(content string) string {
	r := bytes.NewReader([]byte(content))
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatal(err)
	}
	sections := doc.Find(".pn-next")
	href := sections.AttrOr("href", "")
	return href
}
