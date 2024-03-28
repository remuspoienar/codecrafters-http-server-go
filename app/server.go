package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

const GET = "GET"
const POST = "POST"

var directory string

func handleConn(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 2048)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Could not read conn", err.Error())
		return
	}

	request, err := http.ReadRequest(bufio.NewReader(bytes.NewBuffer(buffer[:n])))
	if err != nil {
		fmt.Println("Could not parse request", err.Error())
		return
	}

	verb, pathObj := request.Method, request.URL
	path := pathObj.Path

	if strings.Contains(path, "/files/") {
		splitPath := strings.Split(path, "files/")
		filename := splitPath[1]
		fullPath := strings.TrimRight(directory, "/") + "/" + filename

		if verb == GET {
			text, err := os.ReadFile(fullPath)

			if err == nil {
				fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(text), text)
				return
			}
		} else if verb == POST {
			content, _ := io.ReadAll(request.Body)
			os.WriteFile(fullPath, content, 0644)

			fmt.Fprintf(conn, "HTTP/1.1 201 Created\r\n\r\n")
			return
		}
	}

	if verb == GET && strings.Contains(path, "/user-agent") {
		ua := request.UserAgent()

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
	defer l.Close()

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
