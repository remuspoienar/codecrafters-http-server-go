package main

import (
	"bufio"
	"fmt"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		reader := bufio.NewReader(conn)
		head, err := reader.ReadString('\r')

		if err != nil {
			fmt.Println("Error parsing connection: ", err.Error())
			continue
		}

		path := strings.Split(head, " ")[1]

		if path != "/" {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			conn.Close()
			continue
		}

		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		conn.Close()

	}

}
