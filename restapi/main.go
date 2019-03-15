package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/nocontent", noContentHandler)
	http.HandleFunc("/json", jsonHandler)
	http.HandleFunc("/fizzbuzz", fizzBuzzHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	//http.ResponseWriterに書き込むとレスポンスのbodyになる
	fmt.Fprint(w, "1")

	//ステータスコード設定
	w.WriteHeader(200)
}

func noContentHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(204)
}

func jsonHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
	//json形式で返却
	fmt.Fprint(w, `{"year": "2019", "status": 200}`)
}

func fizzBuzzHandler(w http.ResponseWriter, r *http.Request) {

	//map[string][]string形式で取得
	v := r.URL.Query()
	if v == nil {
		return
	}

	//テストコード側でkeyは固定値nで、valueは一つしかaddしていないので0番目を取得
	n := v["n"][0]
	//intに変換
	m, err := strconv.Atoi(n)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, err)
		return
	}

	switch {
	case m%15 == 0:
		fmt.Fprintf(w, `{"Value": "%s"}`, "FizzBuzz")
	case m%3 == 0:
		fmt.Fprintf(w, `{"Value": "%s"}`, "Fizz")
	case m%5 == 0:
		fmt.Fprintf(w, `{"Value": "%s"}`, "Buzz")
	default:
		fmt.Fprintf(w, `{"Value": "%d"}`, m)
	}

	return
}
