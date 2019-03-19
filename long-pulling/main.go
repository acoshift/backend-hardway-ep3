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
<title>Pulling</title>
<h1><span id="current"></span></h1>
<button onclick="incr()">Incr</button>
<script>
	const current = document.getElementById('current')

	function reload (long) {
		fetch('/get?long=' + (long || 0))
			.then((resp) => resp.text())
			.then((resp) => {
				current.innerText = resp
				reload(1)
			})
			.catch(() => {
				reload(1)
			})
	}
	
	function incr () {
		fetch('/incr')
	}

	reload()
</script>
`))
}

var (
	current uint64
	noti    notifier
)

func get(w http.ResponseWriter, r *http.Request) {
	long := r.FormValue("long") == "1"
	if !long {
		fmt.Fprintln(w, atomic.LoadUint64(&current))
		return
	}

	ch := make(chan struct{})
	noti.Notify(ch)

	select {
	case <-r.Context().Done():
		noti.Cancel(ch)
	case <-ch:
		fmt.Fprintln(w, atomic.LoadUint64(&current))
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
		close(c)
	}
	n.cs = nil
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
