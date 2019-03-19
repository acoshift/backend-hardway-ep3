package main

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/get", get)
	http.HandleFunc("/incr", incr)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	// language=HTML
	w.Write([]byte(`
<!doctype html>
<title>SSE</title>
<h1><span id="current"></span></h1>
<button onclick="incr()">Incr</button>
<script>
	const current = document.getElementById('current')
	
	function incr () {
		fetch('/incr')
	}

	const source = new EventSource('/get')
	source.onmessage = (event) => {
		current.innerText = event.data
	}
</script>
`))
}

var (
	current uint64
	noti    notifier
)

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	f := w.(http.Flusher)
	fmt.Fprintf(w, "data: %d\n\n", atomic.LoadUint64(&current))
	f.Flush()

	ch := make(chan struct{}, 1)
	noti.Notify(ch)
	for {
		select {
		case <-r.Context().Done():
			noti.Cancel(ch)
		case <-ch:
			fmt.Fprintf(w, "data: %d\n\n", atomic.LoadUint64(&current))
			f.Flush()
		}
	}
}

func incr(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint64(&current, 1)
	go noti.Send()
}

type notifier struct {
	sync.Mutex
	cs []chan<- struct{}
}

func (n *notifier) Notify(c chan<- struct{}) {
	n.Lock()
	n.cs = append(n.cs, c)
	n.Unlock()
}

func (n *notifier) Send() {
	n.Lock()
	for _, c := range n.cs {
		c <- struct{}{}
	}
	n.Unlock()
}

func (n *notifier) Cancel(c chan<- struct{}) {
	n.Lock()
	for i, p := range n.cs {
		if p == c {
			close(c)
			n.cs = append(n.cs[:i], n.cs[i+1:]...)
		}
	}
	n.Unlock()
}
