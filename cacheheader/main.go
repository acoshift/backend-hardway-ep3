package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"time"
)

func main() {
	h := http.HandlerFunc(api)
	http.Handle("/api", h)
	http.Handle("/cache-control", cacheControl(h))
	http.Handle("/last-modified", lastModified(h))
	http.Handle("/etag", etag(h))
	http.ListenAndServe(":8080", nil)
}

func api(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// language=JSON
	w.Write([]byte(`{"name": "acoshift"}`))
}

func cacheControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(25*time.Millisecond)
		w.Header().Set("Cache-Control", "public, max-age=600")
		h.ServeHTTP(w, r)
	})
}

func lastModified(h http.Handler) http.Handler {
	t := time.Now().Format(http.TimeFormat)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(25*time.Millisecond)
		w.Header().Set("Last-Modified", t)
		if r.Header.Get("If-Modified-Since") == t {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func etag(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(25*time.Millisecond)
		nw := etagResponseWriter{
			ResponseWriter: w,
		}
		defer nw.Flush(r.Header.Get("If-None-Match"))
		h.ServeHTTP(&nw, r)
	})
}

type etagResponseWriter struct {
	http.ResponseWriter
	buf         bytes.Buffer
	code        int
	wroteHeader bool
}

func (w *etagResponseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.code = code
}

func (w *etagResponseWriter) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.buf.Write(p)
}

func (w *etagResponseWriter) Flush(ifNoneMatch string) {
	rawDigest := sha256.Sum256(w.buf.Bytes())
	digest := "\"" + base64.RawStdEncoding.EncodeToString(rawDigest[:]) + "\""
	w.Header().Set("ETag", digest)
	if digest == ifNoneMatch {
		w.ResponseWriter.WriteHeader(http.StatusNotModified)
		return
	}
	w.ResponseWriter.WriteHeader(w.code)
	w.buf.WriteTo(w.ResponseWriter)
}
