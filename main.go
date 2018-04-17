package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
)

const (
	root     = "https://www.asahi.com/articles/ASL4J669JL4JUEHF016.html"
	selector = "td.link a[href$='.pdf']"
	saveDir  = "./nippo"
)

func getPDFURLs() ([]string, error) {
	var URLs []string

	res, err := http.Get(root)
	if err != nil {
		return URLs, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return URLs, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return URLs, err
	}

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		URLs = append(URLs, href)
	})

	return URLs, nil
}

func downloadPDF(u string) error {
	saveTo := fmt.Sprintf("%s/%s", saveDir, filepath.Base(u))
	if _, err := os.Stat(saveTo); os.IsNotExist(err) {
		res, err := http.Get(u)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			return err
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(saveTo, body, 0644); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	URLs, err := getPDFURLs()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		if err := os.Mkdir(saveDir, 0755); err != nil {
			log.Fatal(err)
		}
	}

	for _, u := range URLs {
		log.Println("Download", u)
		if err := downloadPDF(u); err != nil {
			log.Println(err)
		}
	}
}
