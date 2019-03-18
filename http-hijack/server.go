package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/sudo", sudo)

	fmt.Println("listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome"))
}

func sudo(w http.ResponseWriter, r *http.Request) {
	// check is w implement hijacker
	// some middleware may wrap original response writer without implement hijacker
	hijack, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "hijack not support!", http.StatusInternalServerError)
		return
	}

	// hijack the tcp connection from http library
	conn, wr, err := hijack.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// http library won't do anything to connection anymore
	// so we have to close the connection by our self
	defer conn.Close()

	// tell client that hijack completed
	if r.ProtoMajor <= 1 {
		fmt.Fprintf(wr, "%s 200 OK\n", r.Proto)
	} else if r.ProtoMajor == 2 {
		fmt.Fprintf(wr, "%s 200\n", r.Proto)
	} else {
		// unsupported proto
		return
	}
	fmt.Fprintln(wr)

	fmt.Println("hijacked connection")
	defer fmt.Println("disconnected")

	wr.WriteString("Welcome to sudo console\n\n")
	wr.WriteString("type help to list all commands\n\n")
	wr.Flush()

	for {
		// extend deadline after process a command
		conn.SetDeadline(time.Now().Add(30 * time.Minute))
		lineBytes, _, err := wr.ReadLine()
		if err != nil {
			// client may close the connection
			return
		}

		line := string(lineBytes)
		fmt.Println("received:", line)
		switch line {
		case "":
		case "help":
			help(wr)
		case "exit":
			fmt.Fprintln(wr, "bye!")
			fmt.Fprintln(wr, ">0")
			wr.Flush()
			return
		case "now":
			now(wr)
		case "echo":
			echo(wr)
		default:
			fmt.Fprintln(wr, "unknown command")
		}
		wr.Flush()
	}
}

func help(wr io.ReadWriter) {
	fmt.Fprintf(wr, "- now\n\tprint server time\n")
	fmt.Fprintf(wr, "- exit\n\texit console\n")
	fmt.Fprintf(wr, "- echo\\n<message>\n\techo message\n")
	fmt.Fprintf(wr, "\n")
}

func now(wr io.ReadWriter) {
	fmt.Fprintln(wr, time.Now().Format(time.RFC3339))
}

func echo(wr *bufio.ReadWriter) {
	message, _, _ := wr.ReadLine()
	fmt.Fprintln(wr, string(message))
}
