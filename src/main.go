package main

import (
	"JDSpider/src/downloader"
	"JDSpider/src/parser"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	for {
		fmt.Println("1# Enter '1' to Download JD Main Starter Page.")
		fmt.Println("2# Enter '2' to Download Each Category Main Page.")
		fmt.Println("3# Enter '3' to Download List Page for Each Category.")
		fmt.Println("4# Enter '4' to Download Book Detail Page for Each Category.")
		fmt.Println("5# Enter '5' to Parse Book Detail Page for Each Category and save into mongodb.")
		fmt.Println("4# Enter 'exit' to exit the Program.")
		fmt.Print("Enter text: ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		fmt.Println(text)
		text = strings.Replace(text, "\n", "", -1)
		switch text {
		case "1":
			fmt.Println("1")
		case "4":
			fmt.Println("4")
			//searching each category folder for all list pages
			rd, err := ioutil.ReadDir("E:\\JD\\童书")
			if err != nil {
				log.Fatal(err)
			}
			//parsing each list page to get all the books url and download book detail page into detail folder in category folder
			for _, fi := range rd {
				fmt.Println("")
				fmt.Println(fi.Name())
				fmt.Println(fi.IsDir())
				fmt.Println(fi.Size())
				fmt.Println(fi.ModTime())
				fmt.Println(fi.Mode())
				if fi.IsDir() {
					sBuilder := strings.Builder{}
					sBuilder.WriteString("E:\\JD\\童书\\")
					sBuilder.WriteString(fi.Name())
					rd2, err := ioutil.ReadDir(sBuilder.String())
					if err != nil {
						log.Fatal(err)
					}
					for _, ab := range rd2 {
						if ab.IsDir() == false {
							sBuilder2 := strings.Builder{}
							sBuilder2.WriteString("E:\\JD\\童书\\")
							sBuilder2.WriteString(fi.Name())
							sBuilder2.WriteString("\\" + ab.Name())
							contents, _ := downloader.ReadFromPath(sBuilder2.String())
							//println(contents)
							m := parser.GetPageItems(contents)
							for key, value := range m {
								if downloader.IsFile("E:\\JD\\童书\\" + fi.Name() + "\\details\\" + key + ".html") {
									continue
								}
								contents, _ = downloader.Download("https:" + value)
								downloader.Save(contents, "E:\\JD\\童书\\"+fi.Name()+"\\details\\", key+".html")
							}
						}
						//rand.Seed(time.Now().UnixNano())
						//num := rand.Int31n(10)
						//fmt.Println(num)
						//time.Sleep(time.Duration(num) * time.Second)
					}
				}
				rand.Seed(time.Now().UnixNano())
				num := rand.Int31n(10)
				fmt.Println(num)
				time.Sleep(time.Duration(num) * time.Second)
			}
		case "5":
			fmt.Println("5")
			rand.Seed(time.Now().UnixNano())
			num := rand.Int31n(10)
			fmt.Println(num)
			time.Sleep(time.Duration(num) * time.Second)
		default:
			fmt.Println("unknown command")
		}
		if strings.Compare(text, "exit") == 0 {
			break
		}
	}
	return
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
