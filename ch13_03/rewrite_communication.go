package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type RetryTransporter struct{}

func (RetryTransporter) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	for i := 0; i < 3; i++ {
		resp, err := http.DefaultTransport.RoundTrip(req)
		if err != nil {
			log.Println("fail")
			time.Sleep(time.Second)
			continue
		}
		return resp, nil
	}
	return nil, fmt.Errorf("failed to request to %s", req.URL.String())
}

func main() {
	target, _ := url.Parse("http://127.0.0.1:9001")
	rp := httputil.NewSingleHostReverseProxy(target)
	rp.Transport = &RetryTransporter{}
	server := http.Server{
		Addr:    "127.0.0.1.9000",
		Handler: rp,
	}
	log.Println("Start Listening at :9000")
	log.Fatalln(server.ListenAndServe())
}
