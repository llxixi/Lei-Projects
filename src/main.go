package main

import (
	"JDSpider/src/downloader"
	"JDSpider/src/parser"
	"log"
	"strconv"
)

func main() {
	baseUrl := "https://list.jd.com"
	jdStartUrl := "https://list.jd.com/list.html?cat=1713,3263,3394"
	//Get JD Starter Page contents and save
	path := "E://JD/童书/starter.html"
	var docContents = ""
	if downloader.IsFile(path) == false {
		contents, err := downloader.Download(jdStartUrl)
		if err != nil {
			log.Fatal(err)
		}
		docContents = contents
		downloader.Save(contents, "E://JD/童书/", "starter.html")
	} else {
		contents, err := downloader.ReadFromPath(path)
		if err != nil {
			log.Fatal(err)
		}
		docContents = contents
	}

	m := parser.GetDropList(docContents)

	for e := range m {
		var contents = ""
		path := "E://JD/童书/" + e + "/1.html"
		if downloader.IsFile(path) {
			GetRestPages(contents, path, e)
			continue
		}
		fullUrl := baseUrl + m[e]
		contents, _ = downloader.Download(fullUrl)

		downloader.Save(contents, "E://JD/童书/"+e, "/1.html")
		GetRestPages(contents, path, e)
	}
}

func GetRestPages(contents string, path string, cate string) {
	contents, _ = downloader.ReadFromPath(path)

	count := 1
	for {
		count++
		nextPageUrl := parser.GetNextPageUrl(contents)
		if nextPageUrl == "" {
			break
		}
		contents, _ = downloader.Download("https://list.jd.com" + nextPageUrl)

		downloader.Save(contents, "E://JD/童书/"+cate, "/"+strconv.Itoa(count)+".html")
	}
}
