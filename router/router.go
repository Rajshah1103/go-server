package router

import (
	"net"

	"github.com/Rajshah1103/go-server/handler"
)

type HandlerFunc func(conn net.Conn, method string, path string, headers map[string]string)

var routes = map[string]HandlerFunc{
	"/":        handler.Index,
	"/hello":   handler.Hello,
	"/healthz": handler.Health,
}

func HandleRoute(conn net.Conn, method, path string, headers map[string]string) bool {
	if method != "GET" {
		return false
	}
	if h, ok := routes[path]; ok {
		h(conn, method, path, headers)
		return true
	}
	return false
}
