package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/ws", websocket)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	// language=HTML
	w.Write([]byte(`
<!doctype html>
<title>Web Socket</title>
<form onsubmit="send(event)">
	<input id="input" autocomplete="off">
	<button>Send</button>
</form>
<div id="msg"></div>
<script>
	const ws = new WebSocket("ws://localhost:8080/ws")
	const input = document.getElementById('input')
	const message = document.getElementById('msg')
	
	function send (event) {
		event.preventDefault()
		ws.send(input.value)
		input.value = ''
		input.focus()
	}
	
	function appendMessage (msg) {
		const el = document.createElement('div')
		el.innerText = msg
		message.appendChild(el)
	}

	ws.onopen = () => { appendMessage('connected') }
	ws.onclose = () => { appendMessage('disconnected') }
	ws.onmessage = (event) => { appendMessage(event.data) }
	ws.onerror = () => { appendMessage('error') }
</script>
`))
}

func websocket(w http.ResponseWriter, r *http.Request) {
	upgraded := false
	defer func() {
		if !upgraded {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			return
		}
	}()
	if !contains(r.Header.Get("Connection"), "Upgrade") {
		return
	}
	if !contains(r.Header.Get("Upgrade"), "websocket") {
		return
	}
	if r.Method != http.MethodGet {
		return
	}
	if !contains(r.Header.Get("Sec-Websocket-Version"), "13") {
		return
	}
	key := r.Header.Get("Sec-Websocket-Key")
	if key == "" {
		return
	}

	conn, wr, err := w.(http.Hijacker).Hijack()
	if err != nil {
		return
	}
	defer conn.Close()

	upgraded = true

	fmt.Fprintln(wr, "HTTP/1.1 101 Switching Protocols")
	fmt.Fprintln(wr, "Connection: Upgrade")
	fmt.Fprintln(wr, "Upgrade: websocket")
	fmt.Fprintln(wr, "Sec-WebSocket-Accept: "+computeKey(key))
	fmt.Fprintln(wr)
	wr.Flush()

	conn.SetDeadline(time.Time{}) // no timeout

	io.Copy(os.Stdout, conn)
}

func contains(list string, item string) bool {
	xs := strings.Split(list, ",")
	for _, x := range xs {
		if strings.TrimSpace(x) == item {
			return true
		}
	}
	return false
}

func computeKey(key string) string {
	// https://tools.ietf.org/html/rfc6455#section-1.3
	h := sha1.New()
	h.Write([]byte(key))
	h.Write([]byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
