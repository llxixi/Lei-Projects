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
	name := doc.Find(".p-name").Text()
	bookDetails := doc.Find(".Ptable tr")
	publisher := ""
	isbn := ""
	publishcount := ""
	jdsku := ""
	//brand := ""
	cover := ""
	congshu := ""
	kaiben := ""
	publishDatetime := ""
	pagecount := ""
	paperMode := ""

	bookDetails.Each(func(i int, s *goquery.Selection) {
		text := s.Find("td:nth-child(1)").Text()
		text2 := s.Find("td:nth-child(2)").Text()

		if strings.Contains(text, "出版社") {
			publisher = text2
		}
		if strings.Contains(text, "ISBN") {
			isbn = text2
		}
		if strings.Contains(text, "版次") {
			publishcount = text2
		}
		if strings.Contains(text, "商品编码") {
			jdsku = text2
		}
		if strings.Contains(text, "包装") {
			cover = text2
		}

		if strings.Contains(text, "丛书名") {
			congshu = text2
		}
		if strings.Contains(text, "开本") {
			kaiben = text2
		}
		if strings.Contains(text, "出版时间") {
			publishDatetime = text2
		}
		if strings.Contains(text, "用纸") {
			paperMode = text2
		}
		if strings.Contains(text, "页数") {
			pagecount = text2
		}
	})
	ds := bson.D{
		{"title", name},
		{"publisher", publisher},
		{"isbn", isbn},
		{"publishcount", publishcount},
		{"jdsku", jdsku},
		{"cover", cover},
		{"congshu", congshu},
		{"publisher", publisher},
		{"kaiben", kaiben},
		{"publishDatetime", publishDatetime},
		{"paperMode", paperMode},
		{"pagecount", pagecount},
	}
	return ds
}
