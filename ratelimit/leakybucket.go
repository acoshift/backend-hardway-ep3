package main

import (
	"net"
	"net/http"
	"sync"
	"time"
)

func main() {
	http.ListenAndServe(":8080", leakyBucketLimiter(100*time.Millisecond, 10)(http.NotFoundHandler()))
}

func leakyBucketLimiter(perRequest time.Duration, size int) func(h http.Handler) http.Handler {
	type leakyItem struct {
		Last  time.Time // last request time
		Count int       // requests in queue
	}

	var (
		mu      sync.RWMutex
		storage map[string]*leakyItem
	)

	// cleanup loop
	{
		maxDuration := perRequest + time.Second
		if maxDuration < time.Minute {
			maxDuration = time.Minute
		}

		cleanup := func() {
			deleteBefore := time.Now().Add(-maxDuration)

			mu.Lock()
			defer mu.Unlock()

			for k, t := range storage {
				if t.Count <= 0 && t.Last.Before(deleteBefore) {
					delete(storage, k)
				}
			}
		}

		go func() {
			for {
				time.Sleep(maxDuration)
				cleanup()
			}
		}()
	}

	take := func(key string) bool {
		mu.Lock()
		defer mu.Unlock()

		if storage == nil {
			storage = make(map[string]*leakyItem)
		}

		if storage[key] == nil {
			storage[key] = new(leakyItem)
		}

		t := storage[key]

		now := time.Now()

		// first request ?
		if t.Last.IsZero() {
			t.Last = now
			return true
		}

		next := t.Last.Add(perRequest)
		sleep := next.Sub(now)
		if sleep <= 0 {
			t.Last = now
			return true
		}

		if t.Count >= size {
			// queue full, drop the request
			return false
		}

		t.Last = next

		t.Count++
		mu.Unlock()

		time.Sleep(sleep)

		mu.Lock()
		t.Count--

		return true
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			addr, _, _ := net.SplitHostPort(r.RemoteAddr)
			if !take(addr) {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
