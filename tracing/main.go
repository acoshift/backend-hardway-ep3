package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/opentracing-go"
)

/*
 docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.10
*/
// http://localhost:16686
func main() {
	rand.Seed(time.Now().UnixNano())

	tracer, closer, err := config.Configuration{
		ServiceName: "service-1",
	}.NewTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	http.HandleFunc("/h1", h1)
	http.HandleFunc("/h2", h2)

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}

func h1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := opentracing.SpanFromContext(ctx)
	defer span.Finish()

	time.Sleep(100 * time.Millisecond)
	do1(ctx)
	w.Write([]byte("ok"))
}

func do1(ctx context.Context) {
	span := opentracing.SpanFromContext(ctx)
	defer span.Finish()

	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
}

func h2(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := opentracing.SpanFromContext(ctx)
	defer span.Finish()

	time.Sleep(600 * time.Millisecond)
	w.Write([]byte("ok"))
}
