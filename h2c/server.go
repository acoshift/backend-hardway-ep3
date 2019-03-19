package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	h := http.HandlerFunc(handler)
	http.ListenAndServe(":8080", h2c.NewHandler(h, &http2.Server{}))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Proto)
	w.Write([]byte("Hello\n"))
}
