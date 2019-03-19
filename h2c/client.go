package main

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"golang.org/x/net/http2"
)

func main() {
	client := http.Client{
		Transport: &h2cTransport{},
	}
	resp, err := client.Get("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	resp.Write(os.Stdout)
}

type h2cTransport struct {
	once          sync.Once
	httpTransport *http2.Transport
}

func (t *h2cTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	t.once.Do(func() {
		t.httpTransport = &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (conn net.Conn, e error) {
				return net.Dial(network, addr)
			},
		}
	})
	r.URL.Scheme = "http"
	return t.httpTransport.RoundTrip(r)
}
