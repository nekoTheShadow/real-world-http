package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/slow_page", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	dump, err := httputil.DumpResponse(res, true)
	if err != nil {
		panic(err)
	}
	log.Println(string(dump))
}
