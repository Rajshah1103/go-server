package router

import (
	"net"
)

type HandlerFunc func(conn net.Conn, method string, path string, headers map[string]string)

type Middleware func(HandlerFunc) HandlerFunc

var (
	routes     = map[string]HandlerFunc{}
	middleware []Middleware
)

// Register routes
func Register(path string, handler HandlerFunc) {
	routes[path] = handler
}

// Add global middleware
func Use(mw Middleware) {
	middleware = append(middleware, mw)
}

// Apply middleware in reverse
func applyMiddleware(h HandlerFunc) HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

func HandleRoute(conn net.Conn, method, path string, headers map[string]string) bool {
	if method != "GET" {
		return false
	}
	if h, ok := routes[path]; ok {
		handlerWithMiddleware := applyMiddleware(h)
		handlerWithMiddleware(conn, method, path, headers)
		return true
	}
	return false
}
