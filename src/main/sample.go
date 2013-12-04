package main

import (
	"code.google.com/p/go.net/html"
	"fmt"
	"io"
	"net/http"
)

type Result struct {
	Url string
}

func ParseItem(r io.Reader) []Result {
	results := []Result{}
	doc, err := html.Parse(r)
	if err != nil {
		fmt.Println(err)
	}

	var result Result
	var f func(*html.Node)
	f = func(n *html.Node) {
		// n.Typeでノードの型をチェックできる、ElementNodeでHTMLタグのNode。
		// n.Dataでノートの値をチェックする、aタグをチェックしている
		if n.Type == html.ElementNode && n.Data == "a" {
			// n.Attrで属性を一覧する
			for _, a := range n.Attr {
				if a.Key == "href" {
					result.Url = a.Val
					results = append(results, result)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return results
}

func GetPage(url string) []Result {
	//http.GetでGetリクエストを発行する
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	// deferでやるとReaderを関数の終わりで必ずCloseしてくれる。便利!!
	defer res.Body.Close()
	results := ParseItem(res.Body)
	return results
}

func main() {
	url := "http://qiita.com/advent-calendar/2013/"
	results := GetPage(url)
	for _, result := range results {
		fmt.Println(result.Url)
	}
}
