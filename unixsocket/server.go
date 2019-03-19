package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	lis, err := net.Listen("unix", "/tmp/go-backendhardway.sock")
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()

	go http.Serve(lis, http.HandlerFunc(handler))

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
