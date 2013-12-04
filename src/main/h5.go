package main

import (
  "io"
	"fmt"
	"net/http"
  "sync"
  "code.google.com/p/go-html-transform/html/transform"
//  "strings"
)

type Result struct {
  Title string
  Url string
}

func ParseItem(r io.Reader) {
  t, _ := transform.NewFromReader(r);
  t.Apply(transform.Replace(), "script")
  t.Apply(transform.Replace(), "footer")
  t.Apply(transform.Replace(), "meta")
  t.Apply(transform.Replace(), "ul")
  t.Apply(transform.Replace(), "li")
  t.Apply(transform.Replace(), "link")
  t.Apply(transform.Replace(), "div.day")
  t.Apply(transform.Replace(), "div.advent-calendar-breadcrumb")
  t.Apply(transform.Replace(), "a.post-user-icon")
  fmt.Println(t.String())
}

func GetPage(url string) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
  ParseItem(res.Body)
}

func GoGet(urls []string) (<-chan string) {
	var wg sync.WaitGroup
	ch := make(chan string)
  go func() {
  for _, url := range urls {
    wg.Add(1)
      go func(url string) {
        GetPage(url)
          ch <- url
          wg.Done()
      }(url)
  }
  wg.Wait()
  close(ch)
  }()
	return ch
}

func main() {
	urls := []string{"http://qiita.com/advent-calendar/2013/one-minute"}
	ch := GoGet(urls)
	for {
    _, ok := <-ch
    if !ok {
      return
    }
	}
}
