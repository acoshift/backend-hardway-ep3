package main

import (
	"fmt"
	"net/http"
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

	function reload () {
		fetch('/get')
			.then((resp) => resp.text())
			.then((resp) => {
				current.innerText = resp
			})
	}
	
	function incr () {
		fetch('/incr')
	}

	setInterval(reload, 1000)
</script>
`))
}

var current uint64

func get(w http.ResponseWriter, r *http.Request) {
	x := atomic.LoadUint64(&current)
	fmt.Fprintln(w, x)
}

func incr(w http.ResponseWriter, r *http.Request) {
	atomic.AddUint64(&current, 1)
}
