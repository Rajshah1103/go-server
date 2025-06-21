package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"path/filepath"
)

const publicDir = "./public"

// function used to handle connections
func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// read the request line
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request:", err)
		return
	}

	// parse request line
	method, path, _ := parseRequestLine(requestLine)
	if method == "" || path == "" {
		conn.Write([]byte(buildHTTPResponse("400 Bad Request", 400)))
		return
	}

	// log the request
	fmt.Printf("📩 Received %s request for %s\n", method, path)

	// read and log headers
	headers := readHeaders(reader)

	fmt.Println("🧠 Parsed Headers:")
	for k, v := range headers {
		fmt.Printf("  %s: %s\n", k, v)
	}

	// routing logic
	var response string
	switch {
	case method == "GET" && fileExists(publicDir+path):
		response = serveStatic(publicDir + path)
	case method == "GET" && (path == "/" || path == ""):
		response = serveStatic(publicDir + "/index.html")
	case method == "GET" && path == "/hello":
		response = buildHTTPResponse("Hello, Raj! 👋")
	case method == "GET" && path == "/healthz":
		response = buildHTTPResponse("OK ✅")
	default:
		response = buildHTTPResponse("404 Not Found", 404)
	}

	conn.Write([]byte(response))
}

// func to check whether the file exists or not & check if directory
func fileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	return err == nil && !info.IsDir()
}

// serve static files with cotent tyoe
func serveStatic(filepath string) string {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return buildHTTPResponse("500 Internal Server Error", 500)
	}
	contentType := guessContentType(filepath)
	header := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: %s\r\nConnection: close\r\n\r\n",
		len(data), contentType,
	)
	return header + string(data)
}

func guessContentType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	default:
		return "application/octet-stream"
	}
}

func readHeaders(reader *bufio.Reader) map[string]string {
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break // end of headers
		}
		fmt.Println(line)

		colonIndex := strings.Index(line, ":")
		if colonIndex != -1 {
			key := strings.TrimSpace(line[:colonIndex])
			value := strings.TrimSpace(line[colonIndex+1:])
			headers[key] = value
		}
	}
	return headers
}

func parseRequestLine(line string) (method, path, version string) {
	parts := strings.Fields(strings.TrimSpace(line))
	if len(parts) >= 3 {
		return parts[0], parts[1], parts[2]
	}
	return "", "", ""
}

// function to build HTTP response
func buildHTTPResponse(body string, statusCode ...int) string {
	code := 200
	statusText := "OK"

	if len(statusCode) > 0 {
		if statusCode[0] == 404 {
			code = 404
			statusText = "Not Found"
		} else if statusCode[0] == 400 {
			code = 400
			statusText = "Bad Request"
		}
	}

	return fmt.Sprintf(
		"HTTP/1.1 %d %s\r\nContent-Length: %d\r\nContent-Type: text/plain\r\nConnection: close\r\n\r\n%s",
		code, statusText, len(body), body)
}

// main function
func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting the server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("🚀 Server is listening on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}
