package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

func GetPage(url string) {
	doc, _ := goquery.NewDocument(url)
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		fmt.Println(url)
	})
}

func main() {
	url := "http://qiita.com/advent-calendar/2013/"
	GetPage(url)
}
