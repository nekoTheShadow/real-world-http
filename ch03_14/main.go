package main

import (
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {
	request, err := http.NewRequest("GET", "http://locahost:18888", nil)
	if err != nil {
		panic(err)
	}

	request.Header.Add("Content-Type", "image/jpeg")
	request.SetBasicAuth("ユーザー名", "パスワード")
	request.AddCookie(&http.Cookie{Name: "test", Value: "value"})

	dump, err := httputil.DumpRequest(request, true)
	if err != nil {
		panic(err)
	}
	log.Println(string(dump))
}
