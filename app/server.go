package main

import (
	"bufio"
	"flag"
	"fmt"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

const GET = "GET"

var directory string

func handleConn(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	head := scanner.Text()

	if head == "" {
		fmt.Println("Error parsing connection: Empty request data")
		return
	}

	parts := strings.Split(head, " ")
	verb, path := parts[0], parts[1]

	if verb == GET && strings.Contains(path, "/files/") {
		splitPath := strings.Split(path, "files/")
		filename := splitPath[1]
		fullPath := strings.TrimRight(directory, "/") + "/" + filename

		text, err := os.ReadFile(fullPath)

		if err == nil {
			fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(text), text)
		}
	}

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
		return
	}

	if verb == GET && strings.Contains(path, "/echo/") {
		splitPath := strings.Split(path, "echo/")
		text := splitPath[len(splitPath)-1:][0]
		fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(text), text)

		return
	}

	if path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return
	}

	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))

}
func main() {
	flag.StringVar(&directory, "directory", "", "Directory to scan for files")
	flag.Parse()

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

		go handleConn(conn)
	}

}
