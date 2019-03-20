package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// curl http://localhost:8080 -X PUT -d '{"name":"a"}'
func main() {
	rand.Seed(time.Now().UnixNano())
	http.ListenAndServe(":8080", retry(2, 50*time.Millisecond)(http.HandlerFunc(handler)))
}

func handler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if rand.Int()%2 == 0 {
		http.Error(w, "random error", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}

func retry(retries int, backoffFactor time.Duration) func(http.Handler) http.Handler {
	if retries < 0 {
		retries = 0
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !canRetry(r) {
				h.ServeHTTP(w, r)
				return
			}

			buf := bytes.Buffer{}
			buf.ReadFrom(r.Body)
			r.Body = ioutil.NopCloser(&buf)

			ctx := r.Context()
			var nw *bufferedResponseWriter
			for i := 0; i <= retries; i++ {
				w.Header().Set("X-Retries", strconv.Itoa(i))
				r.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
				nw = newBufferedResponseWriter(w)
				h.ServeHTTP(nw, r)

				if nw.code < 500 {
					break
				}

				select {
				case <-time.After(backoffFactor * time.Duration(1<<uint(i))):
				case <-ctx.Done():
					return
				}
			}
			nw.Flush(w)
		})
	}
}

func newBufferedResponseWriter(w http.ResponseWriter) *bufferedResponseWriter {
	return &bufferedResponseWriter{
		header: cloneHeader(w.Header()),
	}
}

func cloneHeader(h http.Header) http.Header {
	r := make(http.Header)
	for k, v := range h {
		copy(r[k], v)
	}
	return r
}

type bufferedResponseWriter struct {
	buf         bytes.Buffer
	header      http.Header
	code        int
	wroteHeader bool
}

func (w *bufferedResponseWriter) Header() http.Header {
	return w.header
}

func (w *bufferedResponseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.code = code
}

func (w *bufferedResponseWriter) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.buf.Write(p)
}

func (w *bufferedResponseWriter) Flush(p http.ResponseWriter) {
	for k, v := range w.header {
		p.Header()[k] = v
	}
	p.WriteHeader(w.code)
	w.buf.WriteTo(p)
}

func canRetry(r *http.Request) bool {
	if !isIdempotent(r.Method) {
		return false
	}
	return true
}

func isIdempotent(method string) bool {
	switch method {
	case
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPut,
		http.MethodDelete:
		return true
	default:
		return false
	}
}
