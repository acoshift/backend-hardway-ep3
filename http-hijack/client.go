package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	fmt.Println("< connecting to http://localhost:8080/sudo")

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	{
		// write http request to connection
		req, err := http.NewRequest("GET", "http://localhost:8080/sudo", nil)
		if err != nil {
			log.Println(err)
			return
		}
		req.Write(conn)
	}

	r := bufio.NewReader(conn)
	r.ReadLine() // HTTP/1.1 200 OK
	r.ReadLine() // \n

	fmt.Println("< connected")
	fmt.Println("----------------")

	go func() {
		for {
			lineBytes, _, err := r.ReadLine()
			if err != nil {
				os.Exit(0)
				return
			}
			line := string(lineBytes)
			if line == ">0" {
				fmt.Println("< connection closed")
				os.Exit(0)
				return
			}
			fmt.Println("> " + line)
		}
	}()

	for {
		var line string
		fmt.Scanln(&line)

		_, err := fmt.Fprintln(conn, line)
		if err != nil {
			return
		}
	}
}
