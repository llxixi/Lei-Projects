package parser

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mongodb/mongo-go-driver/bson"
	_ "github.com/mongodb/mongo-go-driver/bson"
	_ "github.com/mongodb/mongo-go-driver/mongo"
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

func GetPageItems(content string) map[string]string {
	m := make(map[string]string)
	r := bytes.NewReader([]byte(content))
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatal(err)
	}
	sections := doc.Find(".gl-item")
	sections.Each(func(i int, s *goquery.Selection) {
		divsection := s.Find(".gl-i-wrap")
		asection := s.Find("a")
		href := asection.AttrOr("href", "")
		sku := divsection.AttrOr("data-sku", "")
		_, ok := m[sku]
		if !ok {
			m[sku] = href
		}
	})
	return m
}

func GetBookDetail(content string) interface{} {
	r := bytes.NewReader([]byte(content))
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatal(err)
	}
	name := doc.Find(".sku-name").Text()
	bookDetails := doc.Find("#parameter2 li")
	publisher := ""
	isbn := ""
	publishcount := ""
	jdsku := ""
	brand := ""
	cover := ""
	kaiben := ""
	publishDatetime := ""
	pagecount := ""
	workcount := ""

	bookDetails.Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.Contains(text, "出版社") {
			publisher = strings.Split(text, ":")[1]
		}
		if strings.Contains(text, "ISBN") {
			isbn = strings.Split(text, ":")[1]
		}
	})
	ds := bson.D{
		{"title", name},
		{"publisher", publisher},
		{"isbn", isbn},
	}
	return ds
}
