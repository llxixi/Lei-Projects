package downloader

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Downloader interface {
	Download(string) (string, error)
	Save(string, string, string) error
}

func Download(urlString string) (string, error) {
	res, err := http.Get(urlString)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Println("status code error: %d %s", res.StatusCode, res.Status)
	}
	//load the html document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	return doc.Html()
}

func Save(content string, directory string, name string) error {
	path := directory + name
	if IsDirectoryExist(directory) == false {
		os.Mkdir(directory, os.ModePerm)
	}
	data := []byte(content)
	err := ioutil.WriteFile(path, data, 0644)
	if err == nil {
		fmt.Println("写入文件成功:", len(content))
	} else {
		log.Println(err)
	}
	return err
}

func ReadFromPath(path string) (string, error) {

	btcontents, err := ioutil.ReadFile(path)
	//var builder = strings.Builder{}
	//builder.Write(btcontents)
	//contents := builder.String()
	data, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader(btcontents), simplifiedchinese.GBK.NewDecoder()))
	var builder2 = strings.Builder{}
	builder2.Write(data)
	contents2 := builder2.String()
	return contents2, err
}

// IsExist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func IsDirectoryExist(dir string) bool {
	_, err := os.Stat(dir)
	return err == nil || os.IsExist(err)
}

// IsFile checks whether the path is a file,
// it returns false when it's a directory or does not exist.
func IsFile(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return !fi.IsDir()
}
