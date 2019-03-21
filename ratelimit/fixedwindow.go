package main

import (
	"net"
	"net/http"
	"sync"
	"time"
)

func main() {
	http.ListenAndServe(":8080", fixedWindowRateLimiter(2, time.Second)(http.NotFoundHandler()))
}

func fixedWindowRateLimiter(rate int, unit time.Duration) func(h http.Handler) http.Handler {
	var (
		mu         sync.RWMutex
		lastWindow int64
		storage    map[string]int
	)

	take := func(key string) bool {
		currentWindow := time.Now().UnixNano() / int64(unit)

		mu.Lock()
		defer mu.Unlock()

		// is window outdated ?
		if lastWindow != currentWindow {
			// window outdated, create new window
			lastWindow = currentWindow
			if len(storage) > 0 || storage == nil {
				storage = make(map[string]int)
			}
		}

		// get available token
		available, ok := storage[key]
		if !ok {
			// bucket of given key not exists
			// set available to max
			available = rate
		}

		// can we take a token ?
		if available <= 0 {
			// token not available
			return false
		}

		// take a token
		storage[key] = available - 1

		return true
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
			addr, _, _ := net.SplitHostPort(r.RemoteAddr)
			if !take(addr) {
				// w.Header().Set("Retry-After", "1")
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
