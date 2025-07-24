package handler

import (
	"net"

	"github.com/Rajshah1103/go-server/utils"
)

func Index(conn net.Conn, method, path string, headers map[string]string) {
	response := utils.ServeStatic("./public/index.html")
	conn.Write([]byte(response))
}
