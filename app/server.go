package main

import (
	"bufio"
	"fmt"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

const GET = "GET"

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

		scanner := bufio.NewScanner(conn)
		scanner.Split(bufio.ScanLines)

		scanner.Scan()
		head := scanner.Text()

		if err != nil {
			fmt.Println("Error parsing connection: ", err.Error())
			continue
		}

		parts := strings.Split(head, " ")
		verb, path := parts[0], parts[1]

		if verb == GET && strings.Contains(path, "/user-agent") {
			var ua string

			for scanner.Scan() {
				line := scanner.Text()
				if line == "" {
					break
				}
				if strings.Contains(line, "User-Agent:") {
					ua = strings.Split(line, "User-Agent: ")[1]
					break
				}
			}

			fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(ua), ua)
			conn.Close()
			continue
		}

		if verb == GET && strings.Contains(path, "/echo/") {
			splitPath := strings.Split(path, "echo/")
			text := splitPath[len(splitPath)-1:][0]
			fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(text), text)
			conn.Close()
			continue
		}

		if path != "/" {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			conn.Close()
			continue
		}

		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		conn.Close()

	}

}
