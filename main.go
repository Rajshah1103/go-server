package main

import (
	"bufio"
	"fmt"
	"github.com/Rajshah1103/go-server/router"
    "github.com/Rajshah1103/go-server/utils"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request:", err)
		return
	}

	method, path, _ := utils.ParseRequestLine(requestLine)
	if method == "" || path == "" {
		conn.Write([]byte(utils.BuildHTTPResponse("400 Bad Request", 400)))
		return
	}

	headers := utils.ReadHeaders(reader)
	fmt.Printf("📩 %s %s\n", method, path)

	if !router.HandleRoute(conn, method, path, headers) {
		conn.Write([]byte(utils.BuildHTTPResponse("404 Not Found", 404)))
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting the server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("🚀 Server listening on port 8080...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}
