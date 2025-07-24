package handler

import (
	"github.com/Rajshah1103/go-server/utils"
	"net"
)

func Hello(conn net.Conn, method, path string, headers map[string]string) {
	conn.Write([]byte(utils.BuildHTTPResponse("Hello, Raj! 👋")))
}
