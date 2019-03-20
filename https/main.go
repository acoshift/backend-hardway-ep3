package main

import (
	"net/http"
)

// bundle cert and ca
// cat server.crt > server.bundle-crt && cat ca.crt >> server.bundle-crt
func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/style.css", style)
	http.ListenAndServeTLS(":8443", "server.bundle-crt", "server.key", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if p, ok := w.(http.Pusher); ok {
		p.Push("/style.css", &http.PushOptions{})
		p.Push("/img.jpg", &http.PushOptions{})
	}

	w.Header().Set("Content-Type", "text/html")
	// language=HTML
	w.Write([]byte(`
		<!doctype html>
		<link rel="stylesheet" href="/style.css">
		<h1>Hello</h1>
	`))
}

func style(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	// language=CSS
	w.Write([]byte(`
h1 {
	color: red;
}
	`))
}
