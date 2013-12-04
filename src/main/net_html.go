package main

import (
  "io"
	"fmt"
	"net/http"
  "sync"
  "code.google.com/p/go.net/html"
  "strings"
)

type Result struct {
  Title string
  Url string
}

func ParseItem(r io.Reader) []Result {
  results := []Result{}
  doc, err := html.Parse(r)
  if err != nil {
    fmt.Println(err)
  }

  found := false
  var result Result
  var f func(*html.Node)
  f = func(n *html.Node) {
    if n.Type == html.ElementNode && n.Data == "a" {
      for _, a := range n.Attr {
        if a.Key == "href" && strings.Contains(a.Val, "items") {
          result.Url = "http://qiita.com" + a.Val
          found = true
        }
      }
    }
    if n.Type == html.TextNode && found {
      found = false
      result.Title = n.Data
      results = append(results, result)
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
      f(c)
    }
  }
  f(doc)
  return results
}

func ParseCalendarTitle(r io.Reader) []Result {
  results := []Result{}
  doc, err := html.Parse(r)
  if err != nil {
    fmt.Println(err)
  }

  found := false
  var result Result
  var f func(*html.Node)
  f = func(n *html.Node) {
    if n.Type == html.ElementNode && n.Data == "a" {
      for _, a := range n.Attr {
        if a.Key == "class" && a.Val == "calendar-name" {
          found = true
        }
        if a.Key == "href" && found {
          result.Url = "http://qiita.com" + a.Val
          break
        }
      }
    }
    if n.Type == html.TextNode && found {
      found = false
      result.Title = n.Data
      results = append(results, result)
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
      f(c)
    }
  }
  f(doc)
  return results
}

func GetPage(url string) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
  results := ParseCalendarTitle(res.Body)
  for _, result := range results {
    fmt.Println("<h1>" + result.Title + "</h1>")
	  itemlistRes, _ := http.Get(result.Url)
    defer itemlistRes.Body.Close()
    items := ParseItem(itemlistRes.Body)
    for _, item := range items {
      fmt.Println(item.Title);
      fmt.Println(item.Url);
    }
  }
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
	urls := []string{"http://qiita.com/advent-calendar/2013/"}
	ch := GoGet(urls)
	for {
    _, ok := <-ch
    if !ok {
      return
    }
	}
}
