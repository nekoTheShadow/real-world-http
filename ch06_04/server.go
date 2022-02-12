package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func handlerChunkedResponse(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("expected http.ResponseWriter to be an http.Flusher")
	}
	for i := 1; i <= 10; i++ {
		fmt.Fprintf(w, "Chunk $%d\n", i)
		flusher.Flush()
		time.Sleep(500 * time.Millisecond)
	}
	flusher.Flush()
}

func main() {
	http.HandleFunc("/chunked", handlerChunkedResponse)
	err := http.ListenAndServe(":18888", nil)
	log.Println(err)
}
