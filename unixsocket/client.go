package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	client := http.Client{
		Transport: &unixTransport{},
	}
	resp, err := client.Get("unix://tmp/go-backendhardway.sock")
	if err != nil {
		log.Fatal(err)
	}
	resp.Write(os.Stdout)
}

type unixTransport struct {
	once          sync.Once
	httpTransport *http.Transport
}

func (t *unixTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	t.once.Do(func() {
		t.httpTransport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				addr = strings.TrimSuffix(addr, ":80")
				return net.Dial("unix", addr)
			},
		}
	})
	r.URL.Scheme = "http"
	r.URL.Host = "/" + r.Host + r.URL.Path
	return t.httpTransport.RoundTrip(r)
}
