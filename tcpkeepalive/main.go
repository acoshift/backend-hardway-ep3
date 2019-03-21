package main

import (
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()

	server := http.Server{
		Handler: http.HandlerFunc(handler),
	}
	server.Serve(&keepAliveListener{lis.(*net.TCPListener)})
}

type keepAliveListener struct {
	*net.TCPListener
}

func (lis *keepAliveListener) Accept() (net.Conn, error) {
	conn, err := lis.TCPListener.AcceptTCP()
	if err != nil {
		return nil, err
	}
	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(10 * time.Second)
	return conn, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
