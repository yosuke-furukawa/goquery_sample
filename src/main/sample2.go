package main

import (
	"code.google.com/p/go-html-transform/html/transform"
	"fmt"
	"io"
	"net/http"
)

type Result struct {
	Url string
}

func ParseItem(r io.Reader) {
	// transformのインスタンスを作る
	t, _ := transform.NewFromReader(r)
	// Applyメソッドで自分のDOMに反映する。
	// Applyメソッド内では、TransformFuncを受け付けるようになっており、
	// Replaceの他にもDOMを追加するAppendChildrenやPrependChildrenなどもある。
	// ページを加工するならこっちのが便利。
	// ちなみに以下の処理で不要なページを削っている
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
	//http.GetでGetリクエストを発行する
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	// deferでやるとReaderを関数の終わりで必ずCloseしてくれる。便利!!
	defer res.Body.Close()
	ParseItem(res.Body)
}

func main() {
	url := "http://qiita.com/advent-calendar/2013/"
	GetPage(url)
}
