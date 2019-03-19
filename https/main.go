package main

import (
	"net/http"
)

// bundle cert and ca
// cat server.crt > server.bundle-crt && cat ca.crt >> server.bundle-crt
func main() {
	http.ListenAndServeTLS(":8443", "server.bundle-crt", "server.key", http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello TLS"))
}
