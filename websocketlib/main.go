package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

func main() {
	http.HandleFunc("/", index)
	http.Handle("/ws", websocket.Server{
		Config:  websocket.Config{},
		Handler: ws,
	})
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

func ws(conn *websocket.Conn) {
	buf := make([]byte, 1000)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		fmt.Println("received:", string(buf[:n]))
		conn.Write([]byte("server: received " + string(buf[:n])))
	}
}
