package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// function used to handle connections
func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// read the request line
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request: ", err)
		return
	}
	method, path, _ := parseRequestLine(requestLine)

	// log the request
	fmt.Printf("Received %s request for %s\n", method, path)

	// build response
	var response string
	if path == "/" {
		response = buildHTTPResponse("Welcome to Raj's HTTP Server")
	} else {
		response = buildHTTPResponse("404 Not Found", 404)
	}

	conn.Write([]byte(response))

}

func parseRequestLine(line string) (method, path, version string) {
	parts := strings.Split(strings.TrimSpace(line), " ")
	if len(parts) >= 3 {
		return parts[0], parts[1], parts[2]
	}
	return "", "", ""
}

// function to build http response

func buildHTTPResponse(body string, statusCode ...int) string {
	code := 200
	statusText := "OK"
	if len(statusCode) > 0 && statusCode[0] == 404 {
		code = 404
		statusText = "Not Found"
	}
	return fmt.Sprintf(
		"HTTP/1.1 %d %s\r\nContent-Length: %d\r\nContent-Type: text/plain\r\n\r\n%s",
		code, statusText, len(body), body)
}

// main function

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting the server: ", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connections: ", err)
		}
		// concurrent handling
		go handleConnection(conn)
	}
}
