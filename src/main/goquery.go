package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"sync"
  "os"
)

type Result struct {
	CalTitle string
	Title    string
	Url      string
}

func GetPage(url string) []Result {
	results := []Result{}
	doc, _ := goquery.NewDocument(url)
	doc.Find("a.calendar-name").Each(func(_ int, s *goquery.Selection) {
		url, exists := s.Attr("href")
		if exists {
			caltitle := s.Text()
			entryPage, _ := goquery.NewDocument("http://qiita.com" + url)
			entryPage.Find("div.body h1>a").Each(func(_ int, s *goquery.Selection) {
				url, exists := s.Attr("href")
				if exists {
					result := Result{caltitle, s.Text(), url}
					results = append(results, result)
				}
			})
		}
	})
	return results
}

func GoGet(urls []string) <-chan []Result {
	var wg sync.WaitGroup
	ch := make(chan []Result)
	go func() {
		for _, url := range urls {
			wg.Add(1)
			go func(url string) {
				ch <- GetPage(url)
				wg.Done()
			}(url)
		}
		wg.Wait()
		close(ch)
	}()
	return ch
}

func main() {
	args := os.Args
  if (len(args) < 2) {
    panic("usage : goquery <url>")
  }

  urls := []string{}
  for index, arg := range args {
    if index != 0 {
      urls = append(urls, arg)
    }
  }

	ch := GoGet(urls)
	for {
		results, ok := <-ch
		if !ok {
			return
		}
		calTitle := ""
		for _, result := range results {
			if calTitle != result.CalTitle {
				fmt.Println("#" + result.CalTitle)
				calTitle = result.CalTitle
			}
			fmt.Println("[" + result.Title + "](" + result.Url + ")")
		}
	}
}
